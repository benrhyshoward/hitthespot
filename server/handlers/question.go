package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/benrhyshoward/hitthespot/db"
	"github.com/benrhyshoward/hitthespot/model"
	"github.com/benrhyshoward/hitthespot/questions"
	"github.com/gorilla/mux"
)

type PatchQuestionRequest struct {
	Abandoned *bool
}

func GetQuestion(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	roundID := vars["round_id"]
	questionID := vars["question_id"]

	question, err := db.GetQuestionById(user.Id, roundID, questionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func GetQuestions(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)

	roundID := vars["round_id"]

	questionFilter := func(q model.Question) bool { return true }

	answered := r.URL.Query().Get("answered")
	if answered != "" {
		answeredBool, err := strconv.ParseBool(answered)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		questionFilter = func(q model.Question) bool { return q.Over() == answeredBool }
	}

	questions, err := db.GetQuestions(user.Id, roundID, questionFilter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(questions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func PostQuestion(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	roundID := vars["round_id"]

	unansweredQuestions, err := db.GetQuestions(user.Id, roundID, func(q model.Question) bool { return q.Over() == false })
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(unansweredQuestions) > 0 {
		http.Error(w, "There is still an unanswered question for this round", http.StatusConflict)
		return
	}

	round, err := db.GetRoundById(user.Id, roundID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(round.Questions) >= round.TotalQuestions {
		http.Error(w, "No more questions for round", http.StatusConflict)
		return
	}

	activeRound, err := questions.GetActiveRound(round)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newQuestion, ok := <-activeRound.QuestionChannel
	if !ok {
		http.Error(w, "No more questions for round", http.StatusInternalServerError)
		return
	}
	time := time.Now()
	newQuestion.Created = &time

	err = db.AddQuestionToRound(user.Id, roundID, newQuestion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(newQuestion)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

//Only allows updating the 'Abandoned' field
func PatchQuestion(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	roundID := vars["round_id"]
	questionID := vars["question_id"]

	var req PatchQuestionRequest
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	question, err := db.GetQuestionById(user.Id, roundID, questionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if question.Over() {
		http.Error(w, "Question is already over", http.StatusConflict)
		return
	}

	if question.Abandoned == true && *req.Abandoned == false {
		http.Error(w, "Questions cannot be un-abandoned", http.StatusUnauthorized)
		return
	}

	if question.Abandoned == false && *req.Abandoned == true {
		time := time.Now()
		question.AbandonedAt = &time
	}

	if req.Abandoned != nil {
		question.Abandoned = *req.Abandoned
	}

	round, err := db.GetRoundById(user.Id, roundID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !round.InProgress() {
		questions.DeregisterActiveRound(round)
	}

	err = db.UpdateQuestion(user.Id, round.Id, question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
