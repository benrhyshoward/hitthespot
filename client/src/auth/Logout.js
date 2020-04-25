import React from 'react';
import { useSelector } from 'react-redux';
import {
    selectUserId
  } from './authSlice';
import Button from '@material-ui/core/Button'
import Grid from '@material-ui/core/Grid'
import { Typography, makeStyles } from '@material-ui/core';

const useStyles = makeStyles({
  userId: {
    color: '#1DB954',
    fontWeight: 'bold'
  }, 
  message: {
    textAlign: 'center',
  }
});

export function Logout() {
    const userId = useSelector(selectUserId);
    const classes = useStyles()

    let url;

    if (process.env.NODE_ENV !== 'production') {
      //If react app is being served separately from backend then need the full path
      url = process.env.REACT_APP_GO_SERVER_EXTERNAL_URL+"/auth/logout"
    } else {
      url = "/auth/logout"
    }

  return (
    <Grid container direction="column" justify="center" alignItems="center" spacing={2}>
      <Grid item>
        <Typography className={classes.message} variant="body2">
          Logged in as <span className={classes.userId}>{userId}</span>
        </Typography>
      </Grid>
      <Grid item>
        <Button variant="contained" target="_self" href={url}> 
          Log Out
        </Button>
      </Grid>
    </Grid>
  );
}
