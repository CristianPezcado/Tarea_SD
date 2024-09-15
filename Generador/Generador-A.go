package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dbpool *pgxpool.Pool

func getRandomRecord(ctx context.Context) {
	// Consulta para contar el número total de registros en la tabla
	var count int
	err := dbpool.QueryRow(ctx, "SELECT COUNT(*) FROM dominios").Scan(&count)
	if err != nil {
		log.Fatalf("Error al contar los registros: %v", err)
	}

	if count == 0 {
		log.Println("No hay registros en la tabla.")
		return
	}

	// Generar un índice aleatorio
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(count)

	// Consultar el registro en el índice aleatorio
	var id int
	var dominio string
	query := "SELECT id, dominio FROM dominios LIMIT 1 OFFSET $1"
	err = dbpool.QueryRow(ctx, query, randomIndex).Scan(&id, &dominio)
	if err != nil {
		log.Fatalf("Error al obtener el registro aleatorio: %v", err)
	}

	// Imprimir el registro aleatorio en la consola
	fmt.Printf("Registro Aleatorio:\nID: %d\nDominio: %s\n", id, dominio)
}

type Record struct {
	ID      int    `json:"id"`
	Dominio string `json:"dominio"`
}

func getRandomRecord_Js(ctx context.Context) {
	// Consulta para contar el número total de registros en la tabla
	var count int
	err := dbpool.QueryRow(ctx, "SELECT COUNT(*) FROM dominios").Scan(&count)
	if err != nil {
		log.Fatalf("Error al contar los registros: %v", err)
	}

	if count == 0 {
		log.Println("No hay registros en la tabla.")
		return
	}

	// Generar un índice aleatorio
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(count)

	// Consultar el registro en el índice aleatorio
	var id int
	var dominio string
	query := "SELECT id, dominio FROM dominios LIMIT 1 OFFSET $1"
	err = dbpool.QueryRow(ctx, query, randomIndex).Scan(&id, &dominio)
	if err != nil {
		log.Fatalf("Error al obtener el registro aleatorio: %v", err)
	}

	// Crear un mapa con los datos del registro
	record := Record{
		ID:      id,
		Dominio: dominio,
	}

	// Convertir el mapa a JSON
	jsonData, err := json.Marshal(record)
	if err != nil {
		log.Fatalf("Error al convertir el registro a JSON: %v", err)
	}
	log.Printf("Datos JSON: %s", jsonData)

	// Crear un cliente Resty
	client := resty.New()

	// Configurar la solicitud HTTP POST
	apiURL := "http://localhost:8080/submit" // Reemplaza con la URL de tu API
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(jsonData).
		Post(apiURL)

	if err != nil {
		log.Fatalf("Error al enviar la solicitud HTTP: %v", err)
	}

	// Verificar el código de estado de la respuesta
	if resp.StatusCode() != 200 {
		body := resp.Body()
		log.Fatalf("Error en la respuesta de la API: %v, Cuerpo de la respuesta: %s", resp.Status(), string(body))
	}

	log.Println("Registro enviado correctamente a la API.")
}
func getRandomRecord_S(ctx context.Context) string {
	// Consulta para contar el número total de registros en la tabla
	var count int
	err := dbpool.QueryRow(ctx, "SELECT COUNT(*) FROM dominios").Scan(&count)
	if err != nil {
		return ""
	}

	if count == 0 {
		return ""
	}

	// Generar un índice aleatorio
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(count)

	// Consultar el dominio en el índice aleatorio
	var dominio string
	query := "SELECT dominio FROM dominios LIMIT 1 OFFSET $1"
	err = dbpool.QueryRow(ctx, query, randomIndex).Scan(&dominio)
	if err != nil {
		return ""
	}

	return dominio
}
func runDig(domain string) (string, error) {
	// Ejecutar el comando `dig`
	cmd := exec.Command("dig", domain)

	// Captura la salida del comando
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error al ejecutar dig: %v", err)
	}

	return string(output), nil
}
func queryDNSPeriodically(domain string) {
	ticker := time.NewTicker(10 * time.Second) // Configura el intervalo de tiempo para cada consulta
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Llamar a runDig para obtener el resultado de dig
			result, err := runDig(domain)
			if err != nil {
				log.Printf("Error al realizar la consulta DNS: %v", err)
				continue
			}

			// Mostrar el resultado en la consola
			fmt.Printf("Resultado de dig para %s:\n%s\n", domain, result)
		}
	}
}
func main() {
	databaseURL := "postgresql://user:user@172.20.0.5:5432/tarea1"
	var err error
	dbpool, err = pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("No se pudo establecer conexión con PostgreSQL: %v", err)
	}
	defer dbpool.Close()

	// Crear un ticker que emite cada segundo
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Llamar a la función para obtener y mostrar un registro aleatorio
			getRandomRecord(context.Background())
			//log.Printf("DNS:")
			//queryDNSPeriodically(getRandomRecord_S(context.Background()))
			getRandomRecord_Js(context.Background())
		}
	}
}
