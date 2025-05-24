package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"

	pb "lcp/proto/grpc-server/proto"

	"google.golang.org/grpc"
)

type LCP struct {
	pb.UnimplementedComunicacionServer
	mu             sync.RWMutex
	entrenadores   map[int32]*pb.Entrenador
	torneos        map[int64]*pb.Torneo
	combates       map[int64]*pb.Combate
	maxTorneos     int
	gimnasioClient pb.ComunicacionClient
}

func NewLCP() *LCP {
	return &LCP{
		entrenadores: make(map[int32]*pb.Entrenador),
		torneos:      make(map[int64]*pb.Torneo),
		combates:     make(map[int64]*pb.Combate),
		maxTorneos:   8,
	}
}

// Conecta con el gimnasio remoto (simulado aquí en "localhost:50052")
func (s *LCP) conectarAGimnasio() {
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error conectando con el gimnasio: %v", err)
	}
	s.gimnasioClient = pb.NewComunicacionClient(conn)
}

func (s *LCP) EnviarEntrenadores(ctx context.Context, req *pb.ListaEntrenadores) (*pb.Respuesta, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, e := range req.Entrenadores {
		if _, exists := s.entrenadores[e.Id]; !exists {
			s.entrenadores[e.Id] = e
		}
	}

	return &pb.Respuesta{Mensaje: fmt.Sprintf("%d entrenadores registrados", len(req.Entrenadores))}, nil
}

func (s *LCP) EnviarTorneos(ctx context.Context, req *pb.Empty) (*pb.ListaTorneos, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var torneos []*pb.Torneo
	for _, t := range s.torneos {
		if t.Estado == "En espera" {
			torneos = append(torneos, t)
		}
	}

	return &pb.ListaTorneos{Torneos: torneos}, nil
}

func (s *LCP) InscribirEntrenador(ctx context.Context, req *pb.InscribirRequest) (*pb.InscribirResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	torneo, exists := s.torneos[req.IdTorneo]
	if !exists {
		return nil, fmt.Errorf("torneo no encontrado")
	}

	entrenador, ok := s.entrenadores[req.Entrenador.Id]
	if !ok {
		return nil, fmt.Errorf("entrenador no registrado")
	}

	// Validación de estado
	switch entrenador.Estado {
	case "Expulsado":
		log.Printf("[SNP] Entrenador %s expulsado - intento de inscripción", entrenador.Nombre)
		return &pb.InscribirResponse{Mensaje: "Entrenador expulsado", En: entrenador}, nil
	case "Suspendido":
		if entrenador.Suspension > 0 {
			entrenador.Suspension--
			log.Printf("[SNP] Entrenador %s suspendido (%d torneos restantes)", entrenador.Nombre, entrenador.Suspension)
			if entrenador.Suspension == 0 {
				entrenador.Estado = "Activo"
			}
			return &pb.InscribirResponse{Mensaje: fmt.Sprintf("Suspendido (%d torneos restantes)", entrenador.Suspension), En: entrenador}, nil
		}
	}

	// Ya inscrito
	for _, e := range torneo.Jugadores {
		if e.Id == entrenador.Id {
			return &pb.InscribirResponse{Mensaje: "Ya estás inscrito", En: entrenador}, nil
		}
	}

	// Inscripción
	torneo.Jugadores = append(torneo.Jugadores, entrenador)
	torneo.Inscritos++

	log.Printf("[SNP] Confirmación de inscripción: %s al torneo %d", entrenador.Nombre, torneo.IdTorneo)

	if torneo.Inscritos >= 2 {
		torneo.Estado = "En curso"
		go s.asignarCombate(torneo)
	}

	return &pb.InscribirResponse{Mensaje: "Inscripción exitosa", En: entrenador}, nil
}

func (s *LCP) asignarCombate(torneo *pb.Torneo) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(torneo.Jugadores) < 2 {
		return
	}

	combateID := time.Now().UnixNano()
	combate := &pb.Combate{
		IdCombate:   combateID,
		IdTorneo:    torneo,
		Entrenador1: torneo.Jugadores[0],
		Entrenador2: torneo.Jugadores[1],
	}

	s.combates[combateID] = combate

	// Enviar al gimnasio
	if s.gimnasioClient != nil {
		_, err := s.gimnasioClient.AsignarCombate(context.Background(), combate)
		if err != nil {
			log.Printf("Error al enviar combate al gimnasio: %v", err)
		} else {
			log.Printf("Combate asignado a gimnasio: %s vs %s (ID %d)", combate.Entrenador1.Nombre, combate.Entrenador2.Nombre, combateID)
		}
	}
}

func (s *LCP) generarTorneos() {
	regiones := []string{"Kanto", "Johto", "Hoenn", "Sinnoh", "Unova", "Kalos", "Alola", "Galar", "Paldea"}

	for {
		s.mu.Lock()
		if len(s.torneos) < s.maxTorneos {
			region := regiones[rand.Intn(len(regiones))]
			torneoID := time.Now().UnixNano()

			s.torneos[torneoID] = &pb.Torneo{
				IdTorneo:  torneoID,
				Region:    region,
				Estado:    "En espera",
				Inscritos: 0,
			}
			log.Printf("Nuevo torneo generado en %s (ID: %d)", region, torneoID)
		}
		s.mu.Unlock()

		time.Sleep(30 * time.Second)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("No se pudo escuchar: %v", err)
	}

	lcp := NewLCP()
	lcp.conectarAGimnasio() // <-- Conexión gRPC al gimnasio

	go lcp.generarTorneos()

	s := grpc.NewServer()
	pb.RegisterComunicacionServer(s, lcp)

	log.Printf("Servidor LCP listo en %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Fallo al servir: %v", err)
	}
}
