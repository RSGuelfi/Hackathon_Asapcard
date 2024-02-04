package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Pessoa representa a estrutura dos dados do CSV
type Pessoa struct {
	IDTransacao   string    `json:"id_transacao"`
	DataTransacao time.Time `json:"data_transacao"`
	Documento     int       `json:"documento"`
	Nome          string    `json:"nome"`
	Idade         int       `json:"idade"`
	Valor         float64   `json:"valor"`
	NumParcelas   int       `json:"num_parcelas"`
}

func main() {
	// Abrir o arquivo CSV para leitura
	file, err := os.Open("input-data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Criar um leitor CSV
	reader := csv.NewReader(file)

	// Ler todas as linhas do arquivo
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Configurar a conexão com o PostgreSQL
	dsn := "host=localhost user=postgres password=mypassword dbname=postgres sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	// Automigrar o modelo Pessoa para criar a tabela no PostgreSQL
	err = db.AutoMigrate(&Pessoa{})
	if err != nil {
		log.Fatal(err)
	}

	// Iterar sobre as linhas e criar registros no PostgreSQL
	for _, line := range lines {
		// Seu critério de parse aqui
		if len(line) >= 7 {
			dataTransacao, err := time.Parse("2006-01-02T15:04:05Z", line[1])
			if err != nil {
				log.Println("Erro ao converter DataTransacao:", err)
				continue
			}

			documento, err := strconv.Atoi(line[2])
			if err != nil {
				log.Println("Erro ao converter Documento para int:", err)
				continue
			}

			idade, err := strconv.Atoi(line[4])
			if err != nil {
				log.Println("Erro ao converter Idade para int:", err)
				continue
			}

			valor, err := strconv.ParseFloat(line[5], 64)
			if err != nil {
				log.Println("Erro ao converter Valor para float64:", err)
				continue
			}

			numParcelas, err := strconv.Atoi(line[6])
			if err != nil {
				log.Println("Erro ao converter NumParcelas para int:", err)
				continue
			}

			pessoa := Pessoa{
				IDTransacao:   line[0],
				DataTransacao: dataTransacao,
				Documento:     documento,
				Nome:          line[3],
				Idade:         idade,
				Valor:         valor,
				NumParcelas:   numParcelas,
			}

			// Criar um registro no PostgreSQL
			result := db.Create(&pessoa)
			if result.Error != nil {
				log.Println("Erro ao adicionar pessoa ao PostgreSQL:", result.Error)
			}
		}
	}

	fmt.Println("Dados adicionados ao PostgreSQL com sucesso!")
}
