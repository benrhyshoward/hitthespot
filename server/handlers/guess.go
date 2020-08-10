package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/benrhyshoward/hitthespot/db"
	"github.com/benrhyshoward/hitthespot/model"
	"github.com/benrhyshoward/hitthespot/questions"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type GuessRequest struct {
	Guess *string
}

//Mapping question types to point values
var points = map[string]int{
	model.FreeText:       6,
	model.MultipleChoice: 4,
	model.Portmanteau:    5,
}

func GetGuess(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	roundID := vars["round_id"]
	questionID := vars["question_id"]
	guessID := vars["guess_id"]

	guess, err := db.GetGuessById(user.Id, roundID, questionID, guessID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(guess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func GetGuesses(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	roundID := vars["round_id"]
	questionID := vars["question_id"]

	guessFilter := func(g model.Guess) bool { return true }

	correct := r.URL.Query().Get("correct")
	if correct != "" {
		correctBool, err := strconv.ParseBool(correct)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		guessFilter = func(q model.Guess) bool { return q.Correct == correctBool }
	}

	guesses, err := db.GetGuesses(user.Id, roundID, questionID, guessFilter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(guesses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func PostGuess(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	roundID := vars["round_id"]
	questionID := vars["question_id"]

	var req GuessRequest
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Guess == nil {
		http.Error(w, "Missing field : Guess", http.StatusBadRequest)
		return
	}

	question, err := db.GetQuestionById(user.Id, roundID, questionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if question.Over() {
		http.Error(w, "Question already over", http.StatusConflict)
		return
	}

	if question.Type == model.MultipleChoice && len(question.Guesses) > 0 {
		http.Error(w, "Only allowed one guess for multiple choice questions", http.StatusConflict)
		return
	}

	//Trimmed case insensitive comparison
	correct := strings.TrimSpace(strings.ToLower(question.Answer.Value)) == strings.TrimSpace(strings.ToLower(*req.Guess))

	var score int
	if correct {
		score = points[question.Type]
	} else {
		score = 0
	}

	time := time.Now()

	newGuess := model.Guess{
		Id:      uuid.New().String(),
		Created: &time,
		Content: *req.Guess,
		Correct: correct,
		Score:   score,
	}

	err = db.AddGuessToQuestion(user.Id, roundID, questionID, newGuess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	round, err := db.GetRoundById(user.Id, roundID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !round.InProgress() {
		questions.DeregisterActiveRound(round)
	}

	json, err := json.Marshal(newGuess)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
