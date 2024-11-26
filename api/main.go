package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/gorilla/mux"
)

type Erro struct {
	Erro string `json:"erro"`
}

type Usuario struct {
	ID   uint64 `json:"id,omitempty"`
	Nome string `json:"nome,omitempty"`
}

func conectarDB() (*sql.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	stringConexao := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, erro := sql.Open("postgres", stringConexao)
	if erro != nil {
		return nil, erro
	}

	if erro = db.Ping(); erro != nil {
		return nil, erro
	}

	return db, nil
}

func respostaDeErro(w http.ResponseWriter, statusCode int, erro error) {
	erroStruct := Erro{erro.Error()}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if erro := json.NewEncoder(w).Encode(erroStruct); erro != nil {
		http.Error(w, erro.Error(), http.StatusInternalServerError)
	}
}

func respostaDeSucesso(w http.ResponseWriter, statusCode int, dados interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if erro := json.NewEncoder(w).Encode(dados); erro != nil {
		http.Error(w, erro.Error(), http.StatusInternalServerError)
	}
}

func inserir(w http.ResponseWriter, r *http.Request) {
	corpoReq, erro := io.ReadAll(r.Body)
	if erro != nil {
		respostaDeErro(w, http.StatusUnprocessableEntity, erro)
		return
	}
	defer r.Body.Close()

	var usuario Usuario
	if erro = json.Unmarshal(corpoReq, &usuario); erro != nil {
		respostaDeErro(w, http.StatusBadRequest, erro)
		return
	}

	db, erro := conectarDB()
	if erro != nil {
		respostaDeErro(w, http.StatusInternalServerError, erro)
		return
	}
	defer db.Close()

	sqlStatement := `INSERT INTO usuarios (nome) VALUES ($1) RETURNING id`
	if erro = db.QueryRow(sqlStatement, usuario.Nome).Scan(&usuario.ID); erro != nil {
		respostaDeErro(w, http.StatusInternalServerError, erro)
		return
	}

	respostaDeSucesso(w, http.StatusCreated, usuario)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", inserir).Methods(http.MethodPost)

	fmt.Println("Rodando na porta 5000")
	log.Fatal(http.ListenAndServe(":5000", r))
}
