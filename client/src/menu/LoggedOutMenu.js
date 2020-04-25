import React from 'react';
import { Login } from '../auth/Login'
import { Title } from './Title'
import {Grid, makeStyles} from '@material-ui/core'

const useStyles = makeStyles({
  root: {
      margin: '0px'
  }
});

export function LoggedOutMenu() {

  const classes = useStyles();

    return (
      <Grid container className={classes.root} direction="column" justify="center" alignItems="center" spacing={6}>
        <Grid item>
            <Title/>        
        </Grid>
        <Grid item>
            <Login/>        
        </Grid>
      </Grid>
    );
  }
  