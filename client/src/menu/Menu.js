import React, { useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import {
  fetchActiveRound,
  selectActiveRoundId,
  startNewRound
} from './menuSlice'
import { useHistory, Link } from 'react-router-dom'
import { Button, Typography, Grid, makeStyles } from '@material-ui/core';

const useStyles = makeStyles({
  message: {
    textAlign: 'center',
  },
});

export function Menu() {

  const dispatch = useDispatch();
  const history = useHistory();
  const classes = useStyles();

  const activeRoundId = useSelector(selectActiveRoundId);

  useEffect(() => {
      dispatch(fetchActiveRound());
  }, [])

  function UserInput() {
    if (activeRoundId){
      return ContinueRound()
    } else {
      return StartNewRound()
    }
  }

  function StartNewRound() {
    return (
      <Grid item>
        <Button variant="contained" color="primary" onClick={() => dispatch(startNewRound(history))}>Start new round</Button>
      </Grid>
    )
  }

  function ContinueRound() {
    return (
      <>
        <Grid item >
          <Typography className={classes.message} variant="body1">You already have a round in progress</Typography>
        </Grid>
        <Grid item>
          <Button variant="contained" color="primary" component={Link} to={"/"+activeRoundId}>Continue round</Button>
        </Grid>
      </>
    )
  }

  return (
    <div>   
      <Grid container direction="column" justify="center" alignItems="center" spacing={5}>
        <UserInput/>
      </Grid>
    </div>
  );
}
