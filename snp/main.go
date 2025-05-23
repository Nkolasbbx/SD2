package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	pb "snp/proto/grpc-server/proto" // Ajusta la ruta según tu estructura

	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

type SNP struct {
	rabbitConn *amqp.Connection
	lcpClient  pb.ComunicacionClient
}

type EventoLCP struct {
	Tipo         string `json:"tipo"`          // ranking_actualizado, penalizacion, nuevo_torneo
	EntrenadorID int32  `json:"entrenador_id"` // 0 = broadcast
	Mensaje      string `json:"mensaje"`
	Data         struct {
		NuevoRanking int32  `json:"nuevo_ranking,omitempty"`
		TorneoID     int64  `json:"torneo_id,omitempty"`
		Region       string `json:"region,omitempty"`
		Suspension   int32  `json:"suspension,omitempty"`
		Motivo       string `json:"motivo,omitempty"`
	} `json:"data,omitempty"`
}

func NewSNP() *SNP {
	return &SNP{}
}

func (s *SNP) connectRabbitMQ() error {
	var err error
	s.rabbitConn, err = amqp.Dial(os.Getenv("RABBITMQ_URL"))
	return err
}

func (s *SNP) connectLCP() error {
	conn, err := grpc.Dial(
		os.Getenv("LCP_ADDR"),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return err
	}
	s.lcpClient = pb.NewComunicacionClient(conn)
	return nil
}

func (s *SNP) procesarEvento(evento *EventoLCP) error {
	ch, err := s.rabbitConn.Channel()
	if err != nil {
		return fmt.Errorf("error abriendo canal: %v", err)
	}
	defer ch.Close()

	// Determinar cola(s) destino
	var colas []string
	if evento.EntrenadorID == 0 {
		// Mensaje broadcast a todos los entrenadores
		colas = []string{"notificaciones_general"}
	} else {
		// Mensaje específico para un entrenador
		colas = []string{fmt.Sprintf("notificaciones_entrenador_%d", evento.EntrenadorID)}
	}

	// Construir notificación según tipo
	var notificacion map[string]interface{}
	switch evento.Tipo {
	case "ranking_actualizado":
		notificacion = map[string]interface{}{
			"tipo":          "ranking",
			"entrenador_id": evento.EntrenadorID,
			"nuevo_ranking": evento.Data.NuevoRanking,
			"mensaje":       evento.Mensaje,
			"timestamp":     time.Now().Format(time.RFC3339),
		}
	case "penalizacion":
		notificacion = map[string]interface{}{
			"tipo":          "penalizacion",
			"entrenador_id": evento.EntrenadorID,
			"suspension":    evento.Data.Suspension,
			"motivo":        evento.Data.Motivo,
			"mensaje":       evento.Mensaje,
			"timestamp":     time.Now().Format(time.RFC3339),
		}
	case "nuevo_torneo":
		notificacion = map[string]interface{}{
			"tipo":      "nuevo_torneo",
			"torneo_id": evento.Data.TorneoID,
			"region":    evento.Data.Region,
			"mensaje":   evento.Mensaje,
			"timestamp": time.Now().Format(time.RFC3339),
		}
	default:
		return fmt.Errorf("tipo de evento no reconocido: %s", evento.Tipo)
	}

	// Publicar en cada cola destino
	body, err := json.Marshal(notificacion)
	if err != nil {
		return fmt.Errorf("error serializando notificación: %v", err)
	}

	for _, cola := range colas {
		err = ch.Publish(
			"",    // exchange
			cola,  // routing key
			false, // mandatory
			false, // immediate
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         body,
				DeliveryMode: amqp.Persistent, // Mensajes persistentes
				Headers: amqp.Table{
					"tipo_evento": evento.Tipo,
				},
			},
		)
		if err != nil {
			return fmt.Errorf("error publicando en cola %s: %v", cola, err)
		}
		log.Printf("Notificación enviada a %s: %s", cola, string(body))
	}

	return nil
}

func (s *SNP) consumirEventos() {
	for {
		if s.rabbitConn == nil {
			if err := s.connectRabbitMQ(); err != nil {
				log.Printf("Error conectando a RabbitMQ: %v. Reintentando...", err)
				time.Sleep(5 * time.Second)
				continue
			}
		}

		ch, err := s.rabbitConn.Channel()
		if err != nil {
			log.Printf("Error abriendo canal: %v", err)
			s.rabbitConn.Close()
			s.rabbitConn = nil
			time.Sleep(5 * time.Second)
			continue
		}

		q, err := ch.QueueDeclare(
			"eventos_lcp",
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
			false,  // auto-ack (manual)
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

		log.Println("SNP listo para procesar eventos...")
		for d := range msgs {
			go func(d amqp.Delivery) {
				var evento EventoLCP
				if err := json.Unmarshal(d.Body, &evento); err != nil {
					log.Printf("Error decodificando evento: %v", err)
					d.Nack(false, false) // Descartar mensaje mal formado
					return
				}

				if err := s.procesarEvento(&evento); err != nil {
					log.Printf("Error procesando evento: %v", err)
					d.Nack(false, true) // Reintentar más tarde
					return
				}

				d.Ack(false) // Confirmar procesamiento
			}(d)
		}
	}
}

func main() {
	snp := NewSNP()

	// Conectar a LCP para posibles verificaciones
	if err := snp.connectLCP(); err != nil {
		log.Fatalf("Error conectando a LCP: %v", err)
	}

	// Iniciar consumidor de eventos
	go snp.consumirEventos()

	// Mantener el servicio corriendo
	select {}
}
