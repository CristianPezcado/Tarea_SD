package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Definimos una estructura para los datos que esperamos recibir
type RequestData struct {
	ID      int    `json:"id"`
	Dominio string `json:"dominio"`
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	// Asegurarnos de que la solicitud sea POST
	if r.Method != http.MethodPost {
		log.Printf("Invalid request method: %s", r.Method)
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	log.Println("Received POST request")

	// Decodificamos el cuerpo de la solicitud JSON
	var data RequestData
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		log.Printf("Failed to parse JSON: %v", err)
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	// Log the received data
	log.Printf("Received data: ID=%s, Dominio=%s", data.ID, data.Dominio)

	// Respondemos con un mensaje de éxito
	response := fmt.Sprintf("RData recividaaa: ID=%s, Dominio=%s", data.ID, data.Dominio)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}

func main() {
	// Configuramos el manejador para la ruta "/submit"
	http.HandleFunc("/submit", submitHandler)

	// Iniciamos el servidor HTTP en el puerto 8080
	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
