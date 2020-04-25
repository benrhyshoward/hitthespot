package handlers

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/benrhyshoward/hitthespot/db"
	"github.com/benrhyshoward/hitthespot/model"
	"github.com/benrhyshoward/hitthespot/questions"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"github.com/zmb3/spotify"
	"go.mongodb.org/mongo-driver/mongo"
)

var auth spotify.Authenticator

//Cache of active user sessions
var sessions *cache.Cache

type UserSession struct {
	UserID string
	Client spotify.Client
}

func init() {
	auth = spotify.NewAuthenticator(
		os.Getenv("GO_SERVER_EXTERNAL_URL")+"/auth/callback",
		spotify.ScopeUserTopRead,
		spotify.ScopeUserReadRecentlyPlayed)

	sessions = cache.New(24*time.Hour, 10*time.Minute)

}

func Login(w http.ResponseWriter, r *http.Request) {

	sessionID := uuid.New().String()

	loginURL := getURL(sha256Hash(sessionID))

	expire := time.Now().Add(24 * time.Hour)
	cookie := http.Cookie{
		Name:     "htssess",
		Value:    sessionID,
		Path:     "/",
		Expires:  expire,
		Secure:   os.Getenv("HTTPS") == "true",
		SameSite: 4, //SameSite=None
	}

	http.SetCookie(w, &cookie)
	http.Redirect(w, r, loginURL, http.StatusSeeOther)
}

func Callback(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("htssess")
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	sessionID := sessionCookie.Value
	sessionHash := sha256Hash(sessionID)

	token, err := auth.Token(sessionHash, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	client := auth.NewClient(token)

	spotifyUser, err := client.CurrentUser()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	spotifyUserID := spotifyUser.ID

	user, err := db.GetUserById(spotifyUserID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			//If no user currently exists then create a new one
			user = model.User{
				Id:      spotifyUserID,
				Created: time.Now(),
				Rounds:  []model.Round{},
			}
			err := db.AddUser(user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	user.Client = client

	sessions.Set(sessionID,
		UserSession{
			UserID: user.Id,
			Client: client,
		},
		cache.DefaultExpiration)

	//Find the rounds which the user has in progress and register them as active
	inProgressRounds, err := db.GetRounds(user.Id, func(round model.Round) bool { return round.InProgress() == true })
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, inProgressRound := range inProgressRounds {
		questions.RegisterActiveRound(inProgressRound, user)
	}

	http.Redirect(w, r, os.Getenv("FRONTEND_SERVER_EXTERNAL_URL"), http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	sessionID, err := r.Cookie("htssess")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	sessions.Delete(sessionID.Value)

	//Setting an expired cookie to remove from browsers
	expire := time.Now().Add(-24 * time.Hour)
	cookie := http.Cookie{
		Name:     "htssess",
		Value:    sessionID.Value,
		Secure:   os.Getenv("HTTPS") == "true",
		Path:     "/",
		Expires:  expire,
		SameSite: 4, //SameSite=None
	}

	http.SetCookie(w, &cookie)

	http.Redirect(w, r, os.Getenv("FRONTEND_SERVER_EXTERNAL_URL"), http.StatusSeeOther)
}

func getURL(state string) string {
	return auth.AuthURL(state)
}

func getUserFromRequest(r *http.Request) (model.User, error) {
	var user model.User
	sessionID, err := r.Cookie("htssess")
	if err != nil {
		return user, err
	}

	userSession, found := sessions.Get(sessionID.Value)
	if found {
		userSessionCast := userSession.(UserSession)
		user, err = db.GetUserById(userSessionCast.UserID)
		if err != nil {
			return user, err
		}
		user.Client = userSessionCast.Client
		return user, nil
	}

	return user, errors.New("User not found")
}

func sha256Hash(input string) string {
	h := sha256.New()
	h.Write([]byte(input))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
