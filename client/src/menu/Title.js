import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import { Grid, Typography } from '@material-ui/core';

const useStyles = makeStyles({
  header: {
    textAlign: 'center',
    '& span': {
      color: '#1DB954'
    }
  },
  description: {
    textAlign: 'center',
  },
});

export function Title() {

  const classes = useStyles();

  return (
      <Grid container direction="column" justify="center" alignItems="center" spacing={5}>
          <Grid item className={classes.header}>
            <Typography variant="h1">Hit The <span>Spot</span></Typography>
          </Grid>
          <Grid item className={classes.description}>
            <Typography variant="subtitle1">Personalised quiz based on your favourite music from Spotify</Typography>
          </Grid>
      </Grid>
  );
}
