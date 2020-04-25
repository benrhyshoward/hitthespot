import { configureStore } from '@reduxjs/toolkit';
import authReducer from '../auth/authSlice';
import menuReducer from '../menu/menuSlice';
import roundReducer from '../round/roundSlice';

export default configureStore({
  reducer: {
    auth: authReducer,
    menu: menuReducer,
    round: roundReducer
  },
});
