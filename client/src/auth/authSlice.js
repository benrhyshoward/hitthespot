import { createSlice } from '@reduxjs/toolkit';

export const slice = createSlice({
  name: 'auth',
  initialState: {
    userId: "",
    authenticated: false,
  },
  reducers: {
    setAuthenticated: (state, action) => {
      state.authenticated = action.payload;
    },
    setUserId: (state, action) => {
      state.userId = action.payload;
    }
  },
});

export const { setAuthenticated } = slice.actions;
export const { setUserId } = slice.actions;

export const fetchUser = () => dispatch => {
    fetch("/api/user", {credentials: "include"})
    .then(res => {
      //If we can successfully query the user endpoint then user is authenticated
      dispatch(setAuthenticated(res.status===200))
      return res.json()
    })
    .then(json => {
      dispatch(setUserId(json.Id))
    })
};

export const selectAuthenticated = state => state.auth.authenticated
export const selectUserId = state => state.auth.userId

export default slice.reducer;
