package questions

import (
	"log"
	"math/rand"
	"strings"

	"github.com/benrhyshoward/hitthespot/model"
	"github.com/google/uuid"
)

func generatePortmanteauQuestions(user model.User, noMoreQuestionsChannel chan (struct{})) chan (model.Question) {
	output := make(chan (model.Question))
	go func() {
		log.Print("Starting portmanteau generation")

		allArtists, err := getTopArtists(user)
		if err != nil {
			return
		}

		rand.Shuffle(len(allArtists), func(i, j int) {
			allArtists[i], allArtists[j] = allArtists[j], allArtists[i]
		})

		used := make(map[string]bool)

		for _, artist1 := range allArtists {
			if used[artist1.Name] {
				continue
			}

			for i := 1; i < (len([]rune(artist1.Name)) - 1); i++ {
				if used[artist1.Name] {
					break
				}
				possibleJoin := substr(artist1.Name, i, len([]rune(artist1.Name)))

				for _, artist2 := range allArtists {
					if used[artist2.Name] || artist1.Name == artist2.Name {
						continue
					}
					if strings.HasPrefix(strings.ToLower(artist2.Name), strings.ToLower(possibleJoin)) &&
						len(artist1.Images) > 0 &&
						len(artist2.Images) > 0 {
						portmanteau := artist1.Name + substr(artist2.Name, len([]rune(possibleJoin)), len([]rune(artist2.Name)))
						images := []string{artist1.Images[0].URL, artist2.Images[0].URL}
						question := model.Question{
							Id:          uuid.New().String(),
							Type:        model.Portmanteau,
							Description: "Find the artist portmanteau",
							Images:      images,
							Options:     []string{},
							Answer: model.Answer{
								Value: portmanteau,
							},
							Guesses: []model.Guess{},
						}
						used[artist1.Name] = true
						used[artist2.Name] = true
						select {
						case output <- question:
							log.Print("Sending portmanteau question to channel")
						case <-noMoreQuestionsChannel:
							log.Print("Stopping portmanteau generation")
							return
						}
						break
					}
				}
			}
		}
	}()
	return output
}

func substr(input string, start int, end int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if end > len(asRunes) {
		end = len(asRunes)
	}

	return string(asRunes[start:end])
}
