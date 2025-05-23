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
	mu           sync.RWMutex
	entrenadores map[int32]*pb.Entrenador // ID -> Entrenador
	torneos      map[int64]*pb.Torneo     // ID -> Torneo
	combates     map[int64]*pb.Combate    // ID -> Combate
	maxTorneos   int
}

func NewLCP() *LCP {
	return &LCP{
		entrenadores: make(map[int32]*pb.Entrenador),
		torneos:      make(map[int64]*pb.Torneo),
		combates:     make(map[int64]*pb.Combate),
		maxTorneos:   8,
	}
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

	entrenador := req.Entrenador
	storedEntrenador, exists := s.entrenadores[entrenador.Id]
	if !exists {
		return nil, fmt.Errorf("entrenador no registrado")
	}

	// Validar estado del entrenador
	switch storedEntrenador.Estado {
	case "Expulsado":
		return &pb.InscribirResponse{
			Mensaje: "Entrenador expulsado permanentemente",
			En:      storedEntrenador,
		}, nil
	case "Suspendido":
		if storedEntrenador.Suspension > 0 {
			storedEntrenador.Suspension--
			if storedEntrenador.Suspension == 0 {
				storedEntrenador.Estado = "Activo"
			}
			return &pb.InscribirResponse{
				Mensaje: fmt.Sprintf("Entrenador suspendido (%d torneos restantes)", storedEntrenador.Suspension),
				En:      storedEntrenador,
			}, nil
		}
	}

	// Verificar si ya está inscrito
	for _, e := range torneo.Jugadores {
		if e.Id == entrenador.Id {
			return &pb.InscribirResponse{
				Mensaje: "Ya estás inscrito en este torneo",
				En:      entrenador,
			}, nil
		}
	}

	// Inscribir al torneo
	torneo.Jugadores = append(torneo.Jugadores, entrenador)
	torneo.Inscritos++

	// Si el torneo está lleno, asignar combate
	if torneo.Inscritos >= 2 {
		torneo.Estado = "En curso"
		go s.asignarCombate(torneo)
	}

	return &pb.InscribirResponse{
		Mensaje: "Inscripción exitosa",
		En:      entrenador,
	}, nil
}

func (s *LCP) asignarCombate(torneo *pb.Torneo) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(torneo.Jugadores) < 2 {
		return
	}

	// Crear combate
	combateID := time.Now().UnixNano()
	combate := &pb.Combate{
		IdCombate:   combateID,
		IdTorneo:    torneo,
		Entrenador1: torneo.Jugadores[0],
		Entrenador2: torneo.Jugadores[1],
	}

	// Aquí deberías implementar la llamada gRPC al gimnasio de la región
	// Por ejemplo: gimnasioClient.AsignarCombate(context.Background(), combate)

	s.combates[combateID] = combate
	log.Printf("Combate asignado en %s entre %s y %s",
		torneo.Region, torneo.Jugadores[0].Nombre, torneo.Jugadores[1].Nombre)
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
		log.Fatalf("failed to listen: %v", err)
	}

	lcp := NewLCP()

	go lcp.generarTorneos()

	s := grpc.NewServer()
	pb.RegisterComunicacionServer(s, lcp)

	log.Printf("LCP escuchando en %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
