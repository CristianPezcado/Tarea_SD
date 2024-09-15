package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/CristianPezcado/Tarea_SD/proto" // Reemplaza con la ruta correcta

	"google.golang.org/grpc"
)

// server es el tipo que implementa pb.ExampleServiceServer.
type server struct {
	pb.UnimplementedExampleServiceServer
}

// SendMessage maneja las solicitudes de los clientes y devuelve una respuesta.
func (s *server) SendMessage(ctx context.Context, req *pb.MessageRequest) (*pb.MessageResponse, error) {
	fmt.Printf("Received message: %s\n", req.GetMessage())
	return &pb.MessageResponse{Confirmation: "Message received successfully!"}, nil
}

func main() {
	// Configuración del servidor
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterExampleServiceServer(s, &server{})

	fmt.Println("Server is listening on port 50051...")
	if err := s.Serve(listen); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
