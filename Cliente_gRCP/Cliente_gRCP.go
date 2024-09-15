package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	pb "github.com/CristianPezcado/Tarea_SD/proto" // Reemplaza con la ruta correcta

	"google.golang.org/grpc"
)

const (
	address        = "localhost:50051"
	externalAPIURL = "http://localhost:8080/submit" // Reemplaza con la URL de tu API externa
)

func main() {
	// Conectar al servidor gRPC
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewExampleServiceClient(conn)

	// Obtener datos de la API externa
	data, err := fetchExternalData()
	if err != nil {
		log.Fatalf("Error fetching external data: %v", err)
	}

	// Enviar los datos obtenidos al servidor gRPC
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.MessageRequest{Message: data}
	res, err := c.SendMessage(ctx, req)
	if err != nil {
		log.Fatalf("Error while sending message: %v", err)
	}

	log.Printf("Server response: %s", res.GetConfirmation())
}

// fetchExternalData realiza una solicitud HTTP a la API externa y devuelve el contenido como string.
func fetchExternalData() (string, error) {
	resp, err := http.Get(externalAPIURL)
	if err != nil {
		return "", fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(body), nil
}
