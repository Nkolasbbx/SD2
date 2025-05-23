package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	pb "entrenador/proto/grpc-server/proto"

	"github.com/streadway/amqp"
	"google.golang.org/grpc"
)

type EntrenadorJSON struct {
	ID         string `json:"id"`
	Nombre     string `json:"nombre"`
	Region     string `json:"region"`
	Ranking    int    `json:"ranking"`
	Estado     string `json:"estado"`
	Suspension int    `json:"suspension"` // Corregido de "suspencion"
}

var (
	entrenadores []*pb.Entrenador
	jugador      pb.Entrenador
	jugadorMutex sync.Mutex
	notifChan    = make(chan string)
)

func EnviarEntrenadores(client pb.ComunicacionClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	jugadorMutex.Lock()
	defer jugadorMutex.Unlock()

	respuesta, err := client.EnviarEntrenadores(ctx, &pb.ListaEntrenadores{
		Entrenadores: entrenadores,
	})
	if err != nil {
		log.Printf("Error al enviar entrenadores: %v", err)
		return
	}
	log.Println("Entrenadores registrados en LCP:", respuesta.Mensaje)
}

func leerEntrenadores() error {
	archivo, err := os.Open("entrenadores_pequeno.json")
	if err != nil {
		return fmt.Errorf("error al abrir archivo: %v", err)
	}
	defer archivo.Close()

	var entrenadoresJSON []EntrenadorJSON
	if err := json.NewDecoder(archivo).Decode(&entrenadoresJSON); err != nil {
		return fmt.Errorf("error decodificando JSON: %v", err)
	}

	for _, e := range entrenadoresJSON {
		idInt, err := strconv.Atoi(e.ID)
		if err != nil {
			log.Printf("ID inválido para entrenador %s: %v", e.Nombre, err)
			continue
		}
		entrenadores = append(entrenadores, &pb.Entrenador{
			Id:         int32(idInt),
			Nombre:     e.Nombre,
			Region:     e.Region,
			Ranking:    int32(e.Ranking),
			Estado:     e.Estado,
			Suspension: int32(e.Suspension),
		})
	}
	return nil
}

func inicializarJugador(nombre string) {
	jugadorMutex.Lock()
	defer jugadorMutex.Unlock()

	jugador = pb.Entrenador{
		Id:         generarIDUnico(),
		Nombre:     nombre,
		Region:     "Kanto",
		Ranking:    800,
		Estado:     "Activo",
		Suspension: 0,
	}
	entrenadores = append(entrenadores, &jugador)
}

func generarIDUnico() int32 {
	maxID := int32(0)
	for _, e := range entrenadores {
		if e.Id > maxID {
			maxID = e.Id
		}
	}
	return maxID + 1
}

func conectarLCP() (pb.ComunicacionClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(
		"localhost:50051", //"lcp:50051",  Usando nombre del servicio Docker
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("no se pudo conectar a LCP: %v", err)
	}
	return pb.NewComunicacionClient(conn), conn, nil
}

func initNotificaciones(entrenadorID int32) {
	for {
		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil {
			log.Printf("Error conectando a RabbitMQ, reintentando...: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		ch, err := conn.Channel()
		if err != nil {
			log.Printf("Error abriendo canal: %v", err)
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		cola := fmt.Sprintf("notificaciones_entrenador_%d", entrenadorID)
		_, err = ch.QueueDeclare(cola, true, false, false, false, nil)
		if err != nil {
			log.Printf("Error declarando cola: %v", err)
			ch.Close()
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		msgs, err := ch.Consume(cola, "", true, false, false, false, nil)
		if err != nil {
			log.Printf("Error consumiendo mensajes: %v", err)
			ch.Close()
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		log.Printf("Escuchando notificaciones para entrenador %d", entrenadorID)
		for d := range msgs {
			notifChan <- string(d.Body)
		}
	}
}

func listarTorneos(client pb.ComunicacionClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	respuesta, err := client.EnviarTorneos(ctx, &pb.Empty{})
	if err != nil {
		log.Printf("Error obteniendo torneos: %v", err)
		return
	}

	fmt.Println("\n=== Torneos Disponibles ===")
	for _, t := range respuesta.Torneos {
		fmt.Printf("ID: %d | Región: %s | Inscritos: %d/%d | Estado: %s\n",
			t.IdTorneo, t.Region, t.Inscritos, len(t.Jugadores), t.Estado)
	}
}

func inscribirEnTorneo(client pb.ComunicacionClient, torneoID int64) {
	jugadorMutex.Lock()
	defer jugadorMutex.Unlock()

	if !puedeInscribirse(jugador) {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	respuesta, err := client.InscribirEntrenador(ctx, &pb.InscribirRequest{
		IdTorneo:   torneoID,
		Entrenador: &jugador,
	})
	if err != nil {
		log.Printf("Error inscribiendo en torneo: %v", err)
		return
	}

	jugador = *respuesta.En
	fmt.Println("\n" + respuesta.Mensaje)
}

func puedeInscribirse(e pb.Entrenador) bool {
	switch e.Estado {
	case "Expulsado":
		fmt.Println("\nERROR: Estás expulsado permanentemente de los torneos.")
		return false
	case "Suspendido":
		if e.Suspension > 0 {
			fmt.Printf("\nERROR: Estás suspendido por %d torneos más.\n", e.Suspension)
			return false
		}
	}
	return true
}

func mostrarEstado() {
	jugadorMutex.Lock()
	defer jugadorMutex.Unlock()

	fmt.Println("\n=== Tu Estado ===")
	fmt.Printf("ID: %d\nNombre: %s\nRegión: %s\nRanking: %d\nEstado: %s\n",
		jugador.Id, jugador.Nombre, jugador.Region, jugador.Ranking, jugador.Estado)
	if jugador.Suspension > 0 {
		fmt.Printf("Suspensión: %d torneos restantes\n", jugador.Suspension)
	}
}

func main() {
	// Inicialización
	if err := leerEntrenadores(); err != nil {
		log.Fatalf("Error cargando entrenadores: %v", err)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Ingresa tu nombre como entrenador: ")
	nombre, _ := reader.ReadString('\n')
	nombre = strings.TrimSpace(nombre)

	inicializarJugador(nombre)

	clientLCP, connLCP, err := conectarLCP()
	if err != nil {
		log.Fatalf("Error conectando a LCP: %v", err)
	}
	defer connLCP.Close()

	EnviarEntrenadores(clientLCP)
	go initNotificaciones(jugador.Id)

	// Menú principal
	for {
		select {
		case notif := <-notifChan:
			fmt.Printf("\n¡Nueva notificación!\n%s\n\n", notif)
		default:
			// No bloquea si no hay notificaciones
		}

		fmt.Println("\n--- Menú Principal ---")
		fmt.Println("1. Ver torneos disponibles")
		fmt.Println("2. Inscribirse en un torneo")
		fmt.Println("3. Ver notificaciones")
		fmt.Println("4. Ver mi estado")
		fmt.Println("5. Salir")

		fmt.Print("Elige una opción: ")
		opcionStr, _ := reader.ReadString('\n')
		opcion, err := strconv.Atoi(strings.TrimSpace(opcionStr))
		if err != nil {
			fmt.Println("\nERROR: Ingresa un número válido")
			continue
		}

		switch opcion {
		case 1:
			listarTorneos(clientLCP)
		case 2:
			fmt.Print("\nIngresa el ID del torneo: ")
			idStr, _ := reader.ReadString('\n')
			id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64)
			if err != nil {
				fmt.Println("\nERROR: ID inválido")
				continue
			}
			inscribirEnTorneo(clientLCP, id)
		case 3:
			fmt.Println("\nMostrando notificaciones pendientes...")
			timeout := time.After(500 * time.Millisecond)
		notificacionesLoop:
			for {
				select {
				case notif := <-notifChan:
					fmt.Printf(" - %s\n", notif)
				case <-timeout:
					fmt.Println("No hay más notificaciones pendientes")
					break notificacionesLoop
				}
			}
		case 4:
			mostrarEstado()
		case 5:
			fmt.Println("\n¡Hasta pronto!")
			return
		default:
			fmt.Println("\nERROR: Opción no válida")
		}
	}
}
