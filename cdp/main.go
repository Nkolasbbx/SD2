package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"

	pb "cdp/proto/grpc-server/proto"

	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

type CDP struct {
	pb.UnimplementedComunicacionServer
	mu                 sync.RWMutex
	combatesProcesados map[int64]bool
	rabbitMQConn       *amqp.Connection
	lcpClient          pb.ComunicacionClient
	snpClient          pb.ComunicacionClient
	claveAES           string
}

type ResultadoCombate struct {
	CombateID      int64  `json:"combate_id"`
	TorneoID       int64  `json:"torneo_id"`
	Entrenador1ID  int32  `json:"id_entrenador_1"`
	Entrenador1Nom string `json:"nombre_entrenador_1"`
	Entrenador2ID  int32  `json:"id_entrenador_2"`
	Entrenador2Nom string `json:"nombre_entrenador_2"`
	GanadorID      int32  `json:"id_ganador"`
	GanadorNom     string `json:"nombre_ganador"`
	Fecha          string `json:"fecha"`
	TipoMensaje    string `json:"tipo_mensaje"`
}

func NewCDP() *CDP {
	return &CDP{
		combatesProcesados: make(map[int64]bool),
		claveAES:           os.Getenv("CLAVE_AES"),
	}
}

func (c *CDP) connectRabbitMQ() error {
	var err error
	c.rabbitMQConn, err = amqp.Dial(os.Getenv("RABBITMQ_URL"))
	return err
}

func (c *CDP) connectLCP() error {
	conn, err := grpc.Dial(
		os.Getenv("LCP_ADDR"),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return err
	}
	c.lcpClient = pb.NewComunicacionClient(conn)
	return nil
}

func (c *CDP) descifrarResultado(cifradoB64 string) ([]byte, error) {
	key := sha256.Sum256([]byte(c.claveAES))
	ciphertext, err := base64.StdEncoding.DecodeString(cifradoB64)
	if err != nil {
		return nil, fmt.Errorf("error decodificando base64: %v", err)
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("texto cifrado demasiado corto")
	}

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("error creando cipher: %v", err)
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	plain := make([]byte, len(ciphertext))
	mode.CryptBlocks(plain, ciphertext)

	// Remover padding PKCS7
	padding := int(plain[len(plain)-1])
	if padding < 1 || padding > aes.BlockSize {
		return nil, fmt.Errorf("padding inválido")
	}

	return plain[:len(plain)-padding], nil
}

func (c *CDP) validarEstructura(resultado *ResultadoCombate) error {
	if resultado.CombateID == 0 {
		return fmt.Errorf("combateID no puede ser cero")
	}
	if resultado.TorneoID == 0 {
		return fmt.Errorf("torneoID no puede ser cero")
	}
	if resultado.GanadorID != resultado.Entrenador1ID && resultado.GanadorID != resultado.Entrenador2ID {
		return fmt.Errorf("ganador no coincide con los participantes")
	}
	if _, err := time.Parse("2006-01-02", resultado.Fecha); err != nil {
		return fmt.Errorf("formato de fecha inválido")
	}
	return nil
}

func (c *CDP) validarEntrenadores(ctx context.Context, id1, id2 int32) (bool, error) {
	ent1, err := c.lcpClient.ConsultarEntrenador(ctx, &pb.EntrenadorID{Id: id1})
	if err != nil {
		return false, fmt.Errorf("error consultando entrenador 1: %v", err)
	}

	ent2, err := c.lcpClient.ConsultarEntrenador(ctx, &pb.EntrenadorID{Id: id2})
	if err != nil {
		return false, fmt.Errorf("error consultando entrenador 2: %v", err)
	}

	if ent1.Estado != "Activo" || ent2.Estado != "Activo" {
		return false, nil
	}
	return true, nil
}

func (c *CDP) enviarResultadoLCP(resultado *ResultadoCombate) error {
	if c.rabbitMQConn == nil {
		if err := c.connectRabbitMQ(); err != nil {
			return err
		}
	}

	ch, err := c.rabbitMQConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"resultados_validados",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	body, err := json.Marshal(resultado)
	if err != nil {
		return err
	}

	return ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		})
}

func (c *CDP) procesarMensaje(body []byte) {
	var msg struct {
		Region    string `json:"region"`
		Resultado string `json:"resultado"`
	}
	if err := json.Unmarshal(body, &msg); err != nil {
		log.Printf("Error decodificando mensaje: %v", err)
		return
	}

	// Descifrar resultado
	plain, err := c.descifrarResultado(msg.Resultado)
	if err != nil {
		log.Printf("Error descifrando resultado: %v", err)
		return
	}

	var resultado ResultadoCombate
	if err := json.Unmarshal(plain, &resultado); err != nil {
		log.Printf("Error decodificando resultado: %v", err)
		return
	}

	// Verificar duplicados
	c.mu.Lock()
	if c.combatesProcesados[resultado.CombateID] {
		c.mu.Unlock()
		log.Printf("Combate duplicado (ID: %d), descartado", resultado.CombateID)
		return
	}
	c.combatesProcesados[resultado.CombateID] = true
	c.mu.Unlock()

	// Validar estructura
	if err := c.validarEstructura(&resultado); err != nil {
		log.Printf("Resultado inválido (ID: %d): %v", resultado.CombateID, err)
		return
	}

	// Validar entrenadores con LCP
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	valid, err := c.validarEntrenadores(ctx, resultado.Entrenador1ID, resultado.Entrenador2ID)
	if err != nil {
		log.Printf("Error validando entrenadores: %v", err)
		return
	}
	if !valid {
		log.Printf("Entrenadores inválidos o inactivos (Combate ID: %d)", resultado.CombateID)
		return
	}

	// Enviar a LCP
	if err := c.enviarResultadoLCP(&resultado); err != nil {
		log.Printf("Error enviando resultado a LCP: %v", err)
		return
	}

	log.Printf("Resultado validado y enviado (Combate ID: %d, Ganador: %s)",
		resultado.CombateID, resultado.GanadorNom)
}

func (c *CDP) consumirResultados() {
	for {
		if c.rabbitMQConn == nil {
			if err := c.connectRabbitMQ(); err != nil {
				log.Printf("Error conectando a RabbitMQ: %v. Reintentando...", err)
				time.Sleep(5 * time.Second)
				continue
			}
		}

		ch, err := c.rabbitMQConn.Channel()
		if err != nil {
			log.Printf("Error abriendo canal: %v", err)
			c.rabbitMQConn.Close()
			c.rabbitMQConn = nil
			time.Sleep(5 * time.Second)
			continue
		}

		q, err := ch.QueueDeclare(
			"resultados_combate",
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			log.Printf("Error declarando cola: %v", err)
			ch.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			false,  // auto-ack (manualmente)
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		if err != nil {
			log.Printf("Error consumiendo mensajes: %v", err)
			ch.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		log.Println("CDP listo para procesar resultados...")
		for d := range msgs {
			go func(d amqp.Delivery) {
				c.procesarMensaje(d.Body)
				d.Ack(false)
			}(d)
		}
	}
}

func main() {
	cdp := NewCDP()

	// Conectar a LCP
	if err := cdp.connectLCP(); err != nil {
		log.Fatalf("Error conectando a LCP: %v", err)
	}

	// Iniciar consumidor de RabbitMQ en goroutine
	go cdp.consumirResultados()

	// Servidor gRPC para posibles consultas
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}
	srv := grpc.NewServer()
	pb.RegisterComunicacionServer(srv, cdp)

	log.Println("CDP iniciado, escuchando en :50053")
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
