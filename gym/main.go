package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math"
	mathrand "math/rand"
	"net"
	"os"
	"time"

	pb "gym/proto/grpc-server/proto"

	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

type GimnasioServer struct {
	pb.UnimplementedComunicacionServer
	rabbitMQConn *amqp.Connection
	claveAES     string
	region       string
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

func NewGimnasioServer(region string) *GimnasioServer {
	return &GimnasioServer{
		region:   region,
		claveAES: os.Getenv("CLAVE_AES"), // Clave desde variable de entorno
	}
}

func (gs *GimnasioServer) connectRabbitMQ() error {
	var err error
	gs.rabbitMQConn, err = amqp.Dial(os.Getenv("RABBITMQ_URL"))
	return err
}

func (gs *GimnasioServer) SimularCombate(ent1, ent2 *pb.Entrenador) (int32, float64) {
	diff := float64(ent1.Ranking - ent2.Ranking)
	k := 100.0
	prob := 1.0 / (1.0 + math.Exp(-diff/k))
	if mathrand.Float64() <= prob {
		return ent1.Id, prob * 100
	}
	return ent2.Id, (1 - prob) * 100
}

func (gs *GimnasioServer) cifrarResultado(resultado []byte) (string, error) {
	block, err := aes.NewCipher([]byte(gs.claveAES))
	if err != nil {
		return "", err
	}

	// IV aleatorio para cada mensaje
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	// Padding PKCS7
	padding := aes.BlockSize - len(resultado)%aes.BlockSize
	for i := 0; i < padding; i++ {
		resultado = append(resultado, byte(padding))
	}

	ciphertext := make([]byte, len(resultado))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, resultado)

	// Concatenar IV + texto cifrado
	encrypted := append(iv, ciphertext...)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func (gs *GimnasioServer) publicarResultado(combateID int64, resultado []byte) error {
	if gs.rabbitMQConn == nil {
		if err := gs.connectRabbitMQ(); err != nil {
			return err
		}
	}

	ch, err := gs.rabbitMQConn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"resultados_combate",
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	// Reintentar 3 veces en caso de error
	for i := 0; i < 3; i++ {
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         resultado,
				DeliveryMode: amqp.Persistent,
				Headers: amqp.Table{
					"region": gs.region,
				},
			})
		if err == nil {
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	return err
}

func (gs *GimnasioServer) AsignarCombate(ctx context.Context, req *pb.Combate) (*pb.Respuesta, error) {
	log.Printf("Nuevo combate asignado en %s: %s vs %s (Torneo %d)",
		gs.region, req.Entrenador1.Nombre, req.Entrenador2.Nombre, req.IdTorneo.IdTorneo)

	// Simular combate
	ganadorID, probabilidad := gs.SimularCombate(req.Entrenador1, req.Entrenador2)
	ganador := req.Entrenador1
	if ganadorID == req.Entrenador2.Id {
		ganador = req.Entrenador2
	}

	// Crear estructura de resultado
	resultado := ResultadoCombate{
		CombateID:      req.IdCombate,
		TorneoID:       req.IdTorneo.IdTorneo,
		Entrenador1ID:  req.Entrenador1.Id,
		Entrenador1Nom: req.Entrenador1.Nombre,
		Entrenador2ID:  req.Entrenador2.Id,
		Entrenador2Nom: req.Entrenador2.Nombre,
		GanadorID:      ganadorID,
		GanadorNom:     ganador.Nombre,
		Fecha:          time.Now().Format("2006-01-02"),
		TipoMensaje:    "resultado_combate",
	}

	// Convertir a JSON
	jsonResultado, err := json.Marshal(resultado)
	if err != nil {
		return nil, fmt.Errorf("error al serializar resultado: %v", err)
	}

	// Cifrar resultado
	resultadoCifrado, err := gs.cifrarResultado(jsonResultado)
	if err != nil {
		return nil, fmt.Errorf("error al cifrar resultado: %v", err)
	}

	// Publicar en RabbitMQ
	if err := gs.publicarResultado(req.IdCombate, []byte(resultadoCifrado)); err != nil {
		return nil, fmt.Errorf("error al publicar resultado: %v", err)
	}

	log.Printf("Combate %d completado. Ganador: %s (Probabilidad: %.1f%%)",
		req.IdCombate, ganador.Nombre, probabilidad)

	return &pb.Respuesta{
		Mensaje: fmt.Sprintf("Combate procesado. Ganador: %s", ganador.Nombre),
	}, nil
}

func main() {
	region := os.Getenv("REGION")
	if region == "" {
		region = "Kanto"
	}

	gs := NewGimnasioServer(region)
	defer func() {
		if gs.rabbitMQConn != nil {
			gs.rabbitMQConn.Close()
		}
	}()

	// Intentar conectar a RabbitMQ al inicio
	if err := gs.connectRabbitMQ(); err != nil {
		log.Printf("Advertencia: no se pudo conectar a RabbitMQ al inicio: %v", err)
	}

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Error al escuchar: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterComunicacionServer(srv, gs)

	log.Printf("Gimnasio de %s escuchando en :50052", region)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Error al servir: %v", err)
	}
}
