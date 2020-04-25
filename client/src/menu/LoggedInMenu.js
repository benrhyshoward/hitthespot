import React from 'react';
import { Menu } from './Menu'
import { Title } from './Title'
import { Logout } from '../auth/Logout'
import {Grid, makeStyles} from '@material-ui/core'

const useStyles = makeStyles({
  root: {
      margin: '0px'
  }
});

export function LoggedInMenu() {

  const classes = useStyles();
  
    return (
      <Grid container className={classes.root} direction="column" justify="center" alignItems="center" spacing={6}>
        <Grid item>
            <Title/>        
        </Grid>
        <Grid item>
            <Menu/>        
        </Grid>
        <Grid item>
            <Logout/>        
        </Grid>
      </Grid>
    );
  }
  