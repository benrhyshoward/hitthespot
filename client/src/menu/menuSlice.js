import { createSlice } from '@reduxjs/toolkit';
import {
  requestNewQuestion,
} from '../round/roundSlice' 

export const slice = createSlice({
  name: 'menu',
  initialState: {
    activeRoundId:""
  },
  reducers: {
    setActiveRoundId: (state, action) => {
      state.activeRoundId = action.payload;
    }
  },
});

export const { setActiveRoundId } = slice.actions;

export const fetchActiveRound = () => dispatch => {
  dispatch(setActiveRoundId(""))
    fetch("/api/rounds?active=true", {credentials: "include"})
    .then(res => res.json())
    .then(json => {
      if (json.length > 0 ){
        dispatch(setActiveRoundId(json[0].Id))
      }
    })
};

export const startNewRound = (history) => dispatch => {
  fetch("/api/rounds", {credentials: "include", method:"POST"})
  .then(res => res.json())
  .then(json => {
    dispatch(setActiveRoundId(json.Id))
    
    //Moving the player straight into the round after its been created
    history.push('/' + json.Id)
    dispatch(requestNewQuestion(json.Id));
  })
};

export const selectActiveRoundId = state => state.menu.activeRoundId

export default slice.reducer;
