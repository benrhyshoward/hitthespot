import React from 'react';
import {
    selectCorrectNotification
} from './roundSlice' 
import { useSelector } from 'react-redux';
import { Typography, makeStyles } from '@material-ui/core';

const useStyles = makeStyles({
    correctColor: {
        color: '#66bb6a'
      },
    incorrectColor: {
        color: '#ef5350'
      },
  });

export function CorrectNotification(){

    const correctNotification = useSelector(selectCorrectNotification);
    const classes = useStyles();

    function getCorrectMessage(){
        if (correctNotification === true) {
            return "Correct!"
        } else if (correctNotification === false) {
            return "Nope"
        } 
        return ""
    }
    
    function getCorrectClass(){
        if (correctNotification === true) {
            return classes.correctColor
        } else if (correctNotification === false) {
            return classes.incorrectColor
        } 
        return null
    }

    return (
        <Typography className={getCorrectClass()} variant="subtitle1">{getCorrectMessage()}&nbsp;</Typography>
    )
}