import React, { useEffect } from 'react';
import { useDispatch } from 'react-redux';
import {
    fetchUser
  } from './authSlice';
import Button from '@material-ui/core/Button'
import Grid from '@material-ui/core/Grid'
import { makeStyles } from '@material-ui/core';

const useStyles = makeStyles({
  loginButton: {
      textAlign: 'center'
  }
});

export function Login() {
    const dispatch = useDispatch();
    const classes = useStyles();

    let url;

    if (process.env.NODE_ENV !== 'production') {
      //If react app is being served separately from backend then need the full path
      url = process.env.REACT_APP_GO_SERVER_EXTERNAL_URL+"/auth/login"
    } else {
      url = "/auth/login"
    }

    useEffect(() => {
        dispatch(fetchUser());
    }, [])

  return (
    <Grid container direction="column" justify="center" alignItems="center" spacing={6}>
      <Grid item>
        <Button className={classes.loginButton} variant="contained" color="primary" target="_self" href={url}> 
          Login with Spotify
        </Button>
      </Grid>
    </Grid>
  );
}
