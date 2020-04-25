package questions

import (
	"errors"
	"log"

	"github.com/benrhyshoward/hitthespot/model"
	"github.com/patrickmn/go-cache"
)

//Cache of active rounds mapped to channels containing questions for that round
var activeRoundQuestionChannels *cache.Cache

type ActiveRound struct {
	QuestionChannel        chan (model.Question)
	NoMoreQuestionsChannel chan (struct{})
	AbandonedChannel       chan (struct{})
}

func init() {
	activeRoundQuestionChannels = cache.New(cache.NoExpiration, cache.NoExpiration)
}

func RegisterActiveRound(round model.Round, user model.User) {
	_, found := activeRoundQuestionChannels.Get(round.Id)

	if !found {
		log.Print("Registering new active round " + round.Id)

		//Used by question merger to signal to generators to stop generating questions
		noMoreQuestionsChannel := make(chan (struct{}))

		//Used by handlers to signal to merger to stop providing questionss
		abandonedChannel := make(chan (struct{}))

		//Array of question channels, populated by Goroutines for each question type
		questionChannels := []chan (model.Question){
			generatePortmanteauQuestions(user, noMoreQuestionsChannel),
			generateAnagramQuestions(user, noMoreQuestionsChannel),
			generateLyricQuestions(user, noMoreQuestionsChannel),
		}

		//Merge all the question channels into a single channel to be read by handlers
		mergedChannel := mergeChannels(
			round,
			noMoreQuestionsChannel,
			abandonedChannel,
			questionChannels,
		)

		//Caching the question channel and abandoned channel for this round, both to be used by handlers
		activeRound := ActiveRound{
			QuestionChannel:  mergedChannel,
			AbandonedChannel: abandonedChannel,
		}
		activeRoundQuestionChannels.Set(round.Id, activeRound, cache.NoExpiration)
	}
}

func mergeChannels(round model.Round, noMoreQuestionsChannel chan (struct{}), abandonedChannel chan (struct{}), questionChannels []chan (model.Question)) chan (model.Question) {
	output := make(chan (model.Question))
	go func() {
		for questionsRemaining := round.TotalQuestions - len(round.Questions); questionsRemaining > 0; questionsRemaining-- {
			//if no more questions being provided by any channel then just wait for the round to be abandoned
			if len(questionChannels) == 0 {
				close(noMoreQuestionsChannel)
				close(output)
				<-abandonedChannel
				return
			}

			//Fetching questions from channels in a round robin
			channelIndex := questionsRemaining % len(questionChannels)

			question, ok := <-questionChannels[channelIndex]
			if !ok {
				//Channel is closed so can remove from list and continue
				questionChannels = append(questionChannels[:channelIndex], questionChannels[channelIndex+1:]...)
				questionsRemaining++
				continue
			}

			select {
			case output <- question:
			case <-abandonedChannel:
				questionsRemaining = 0
			}
		}
		//Close channel to broadcast to all generators to stop
		close(noMoreQuestionsChannel)
		close(output)
		return
	}()
	return output
}

func GetActiveRound(round model.Round) (ActiveRound, error) {
	channel, found := activeRoundQuestionChannels.Get(round.Id)

	if !found {
		return ActiveRound{}, errors.New("Round not active")
	}

	return channel.(ActiveRound), nil
}

func DeregisterActiveRound(round model.Round) {
	log.Print("Deregistering active round")
	activeRoundQuestionChannels.Delete(round.Id)
}
