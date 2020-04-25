package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/benrhyshoward/hitthespot/handlers"
	"github.com/gorilla/mux"
)

func main() {

	//Seeding random number generation with current time
	rand.Seed(time.Now().UnixNano())

	r := mux.NewRouter()

	//Auth endpoints
	r.HandleFunc("/auth/login", handlers.Login).
		Methods("GET")
	r.HandleFunc("/auth/callback", handlers.Callback).
		Methods("GET")
	r.HandleFunc("/auth/logout", handlers.Logout).
		Methods("GET")

	//User endpoint
	r.HandleFunc("/api/user", handlers.GetUser).
		Methods("GET")

	//Round endpoints
	r.HandleFunc("/api/rounds/{round_id}", handlers.GetRound).
		Methods("GET")
	r.HandleFunc("/api/rounds/{round_id}", handlers.PatchRound).
		Methods("PATCH")
	r.HandleFunc("/api/rounds", handlers.GetRounds).
		Methods("GET")
	r.HandleFunc("/api/rounds", handlers.PostRound).
		Methods("POST")

	//Question endpoints
	r.HandleFunc("/api/rounds/{round_id}/questions/{question_id}", handlers.GetQuestion).
		Methods("GET")
	r.HandleFunc("/api/rounds/{round_id}/questions/{question_id}", handlers.PatchQuestion).
		Methods("PATCH")
	r.HandleFunc("/api/rounds/{round_id}/questions", handlers.GetQuestions).
		Methods("GET")
	r.HandleFunc("/api/rounds/{round_id}/questions", handlers.PostQuestion).
		Methods("POST")

	//Guess endpoints
	r.HandleFunc("/api/rounds/{round_id}/questions/{question_id}/guesses/{guess_id}", handlers.GetGuess).
		Methods("GET")
	r.HandleFunc("/api/rounds/{round_id}/questions/{question_id}/guesses", handlers.GetGuesses).
		Methods("GET")
	r.HandleFunc("/api/rounds/{round_id}/questions/{question_id}/guesses", handlers.PostGuess).
		Methods("POST")

	staticFilePath := os.Getenv("STATIC_FILE_PATH")
	if staticFilePath != "" {
		log.Print("Serving static files in " + staticFilePath)
		r.PathPrefix("/").Handler(http.FileServer(http.Dir(staticFilePath)))
	}

	//If HTTPS is enabled then expect a certificate and key in certs/certificate.pem and certs/key.pem
	if os.Getenv("HTTPS") == "true" {
		log.Print("Serving over https on port 8080")
		log.Fatal(http.ListenAndServeTLS(":8080", "certs/certificate.pem", "certs/key.pem", r))
	} else {
		log.Print("Serving over http on port 8080")
		log.Fatal(http.ListenAndServe(":8080", r))
	}

}
