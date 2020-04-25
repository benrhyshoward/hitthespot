import { createSlice } from '@reduxjs/toolkit';

export const slice = createSlice({
  name: 'round',
  initialState: {
    round:{},
    score:0,
    correctNotifcation: null,
  },
  reducers: {
    setRound: (state, action) => {
      state.round = action.payload;
    },
    setScore: (state, action) => {
      state.score = action.payload;
    },
    setCorrectNotifcation: (state, action) => {
      state.correctNotifcation = action.payload;
    }
  },
});

export const { setRound } = slice.actions;
export const { setScore } = slice.actions;
export const { setCorrectNotifcation } = slice.actions;

export const fetchRound = (id) => dispatch => {
    fetch("/api/rounds/"+id, {
      credentials: "include"})
    .then(res => res.json())
    .then(json => {
        dispatch(setRound(json))
        dispatch(setScore(calculateScoreForRound(json)))
    })
};

export const requestNewQuestion = (roundId) => dispatch => {
  dispatch(setCorrectNotifcation(null))

  fetch("/api/rounds/"+roundId+"/questions", {
    credentials: "include", 
    method:"POST"})
  .then(res => 
    res.json())
  .then(json => {
    dispatch(fetchRound(roundId))
  })
};

export const abandonRound = (roundId, history) => dispatch => {
  fetch("/api/rounds/"+roundId, {
    credentials: "include", 
    method:"PATCH",
    body:JSON.stringify({Abandoned:true})
  })
  .then(res => {
    dispatch(fetchRound(roundId))
    history.push("/"+roundId+"/s")
  })
};

export const abandonQuestion = (roundId, questionId) => dispatch => {
  dispatch(setCorrectNotifcation(null))

  fetch("/api/rounds/"+roundId+"/questions/"+questionId, {
    credentials: "include", 
    method:"PATCH",
    body:JSON.stringify({Abandoned:true})
  })
  .then(res => {
    dispatch(fetchRound(roundId))
  })
};

export const sendGuess = (roundId, questionId, guess) => dispatch => {
  fetch("/api/rounds/"+roundId+"/questions/"+questionId+"/guesses", {
    credentials: "include", 
    method:"POST",
    body:JSON.stringify({Guess:guess})
  })
  .then(res => res.json())
  .then(json => {
    dispatch(setCorrectNotifcation(json.Correct))
    dispatch(fetchRound(roundId))
  })
};

//Sending a guess, then if incorrect, requesting the correct answer
export const sendSingleGuess = (roundId, questionId, guess) => dispatch => {
  fetch("/api/rounds/"+roundId+"/questions/"+questionId+"/guesses", {
    credentials: "include", 
    method:"POST",
    body:JSON.stringify({Guess:guess})
  })
  .then(res => res.json())
  .then(json => {
    if (!json.Correct){
      dispatch(abandonQuestion(roundId, questionId))
    }
   dispatch(fetchRound(roundId))
  })
};

const calculateScoreForRound = (round) => {
    //Summing all scores from all guesses in the given round
    return round.Questions.reduce((t, c) => {
        return t + c.Guesses.reduce((t2, c2) => {
          return t2 += c2.Score
        },0)
    },0)
}

export const selectRound = state => state.round.round
export const selectScore = state => state.round.score
export const selectCorrectNotification = state => state.round.correctNotifcation

export default slice.reducer;
