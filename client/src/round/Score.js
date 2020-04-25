import React, { useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import {
  fetchRound,
  selectScore,
  selectRound
} from './roundSlice'
import { useParams, Link } from 'react-router-dom'
import { makeStyles, Grid, Button, Typography} from '@material-ui/core'

const useStyles = makeStyles({
  root: {
    margin: '0px',
    textAlign: 'center'
  }
});

function totalQuestions(round){
  if (!round || !round.Questions) {
    return 0
  }
  return round.Questions.length
}

function correctAnswers(round){
  if (!round || !round.Questions) {
    return 0
  }
  return round.Questions.filter(q => q.Guesses.some((guess) => guess.Correct === true)).length
}


export function Score() {

    const classes = useStyles();

    const { id } = useParams();

    const dispatch = useDispatch();

    const round = useSelector(selectRound);
    const score = useSelector(selectScore);

    useEffect(() => {
        dispatch(fetchRound(id));
    }, [])

  return (
    <Grid container className={classes.root} direction="column" justify="center" alignItems="center" spacing={4}>
      <Grid item>
          <Typography variant="h3">Round Summary</Typography>
        </Grid>
      <Grid item>
          <Typography variant="h5">Total Questions : {totalQuestions(round)}</Typography>
      </Grid>
      <Grid item>
          <Typography variant="h5">Correct Answers : {correctAnswers(round)}</Typography>
      </Grid>
      <Grid item>
          <Typography variant="h4">Final score:  {score}</Typography>
        </Grid>
      <Grid item>
        <Button variant="contained" color="primary" component={Link} to={'/'}>Back to main menu</Button>
      </Grid>
    </Grid>
  );
}
