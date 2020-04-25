package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/benrhyshoward/hitthespot/db"
	"github.com/benrhyshoward/hitthespot/model"
	"github.com/benrhyshoward/hitthespot/questions"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type PatchRoundRequest struct {
	Abandoned *bool
}

func GetRound(w http.ResponseWriter, r *http.Request) {

	user, err := getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	roundID := mux.Vars(r)["round_id"]

	round, err := db.GetRoundById(user.Id, roundID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(round)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func GetRounds(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	roundFilter := func(round model.Round) bool { return true }

	active := r.URL.Query().Get("active")
	if active != "" {
		activeBool, err := strconv.ParseBool(active)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		roundFilter = func(round model.Round) bool { return round.InProgress() == activeBool }
	}

	//Fetch rounds from DB
	rounds, err := db.GetRounds(user.Id, roundFilter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json, err := json.Marshal(rounds)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

func PostRound(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	rounds, err := db.GetRounds(user.Id, func(round model.Round) bool { return round.InProgress() == true })
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(rounds) > 0 {
		http.Error(w, "There is already a round with questions left to answer", http.StatusConflict)
		return
	}

	roundID := uuid.New().String()

	time := time.Now()

	newRound := model.Round{
		Id:             roundID,
		Abandoned:      false,
		Created:        &time,
		TotalQuestions: 10,
		Questions:      []model.Question{},
	}

	err = db.AddRound(user.Id, newRound)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	questions.RegisterActiveRound(newRound, user)

	json, err := json.Marshal(newRound)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}

//Only allows updating the 'Abandoned' field
func PatchRound(w http.ResponseWriter, r *http.Request) {
	user, err := getUserFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	roundID := vars["round_id"]

	var req PatchRoundRequest
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err = dec.Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	round, err := db.GetRoundById(user.Id, roundID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !round.InProgress() {
		http.Error(w, "Round is already over", http.StatusConflict)
		return
	}

	if round.Abandoned == true && *req.Abandoned == false {
		http.Error(w, "Rounds cannot be un-abandoned", http.StatusUnauthorized)
		return
	}

	if req.Abandoned != nil {
		round.Abandoned = *req.Abandoned
	}

	if *req.Abandoned == true {
		activeRound, err := questions.GetActiveRound(round)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//If the rounds still has questions remaining, then signal to question generators to stop
		if round.TotalQuestions-len(round.Questions) > 0 {
			activeRound.AbandonedChannel <- struct{}{}
		}

		questions.DeregisterActiveRound(round)
	}

	err = db.UpdateRound(user.Id, round)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
