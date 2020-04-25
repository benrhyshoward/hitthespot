package questions

import (
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/benrhyshoward/hitthespot/model"
	"github.com/google/uuid"
)

func generateAnagramQuestions(user model.User, noMoreQuestionsChannel chan (struct{})) chan (model.Question) {
	output := make(chan (model.Question))
	go func() {
		log.Print("Starting anagram generation")
		topArtists, err := getTopArtists(user)
		if err != nil {
			return
		}

		if len(topArtists) == 0 {
			close(output)
			return
		}

		for {
			rand.Shuffle(len(topArtists), func(i, j int) {
				topArtists[i], topArtists[j] = topArtists[j], topArtists[i]
			})

			for _, artist := range topArtists {
				anagram := getAnagram(artist.Name)
				hint := getAnagramHint(artist.Name)

				question := model.Question{
					Id:          uuid.New().String(),
					Type:        model.FreeText,
					Description: "Unscramble the artist name",
					Content:     anagram + "\n" + hint,
					Options:     []string{},
					Answer: model.Answer{
						Value: strings.ToUpper(artist.Name),
					},
					Guesses: []model.Guess{},
				}
				select {
				case output <- question:
					log.Print("Sending anagram question to channel")
				case <-noMoreQuestionsChannel:
					log.Print("Stopping anagram generation")
					return
				}
			}
		}
	}()
	return output
}

func getAnagram(word string) string {
	//removing spaces, uppercasing, and splitting into rune array
	letters := []rune(strings.ToUpper(strings.Replace(word, " ", "", -1)))
	rand.Shuffle(len(letters), func(i, j int) {
		letters[i], letters[j] = letters[j], letters[i]
	})
	return strings.TrimSpace(string(letters))
}

//Generates '(x,y,z,...)' where x y z are lengths of individual words
func getAnagramHint(word string) string {
	hint := "("
	words := strings.Split(word, " ")
	for i, word := range words {
		hint = hint + strconv.Itoa(len(word))
		if i != len(words)-1 {
			hint = hint + ", "
		}
	}
	return hint + ")"
}
