import React from 'react';
import classNames from 'classnames';
import {
    sendSingleGuess
} from './roundSlice' 
import { useDispatch } from 'react-redux';
import { Typography, Grid, Button, makeStyles } from '@material-ui/core';


const useStyles = makeStyles({
    question: {
      textAlign: 'center',
      whiteSpace: 'pre-line'
    },
    description: {
        textAlign: 'center',
        whiteSpace: 'pre-line'
      },
    button: {
        maxWidth: '500px'
    },
    buttonCorrect: {
        backgroundColor: '#1fce45a8'
    },
    buttonIncorrect: {
        backgroundColor: '#ff0000ab'
    },
    disabled: {
        pointerEvents: 'none'
    },
  });

export function MultipleChoiceQuestion(props) {

    const roundId = props.roundId
    const question = props.question
    const classes = useStyles();

    const dispatch = useDispatch();

    function OptionGrid() {
        let disabledClass;
        if (question.Guesses && question.Guesses.length >0) {
            disabledClass = classes.disabled //only allowed one guess for multiple choice questions
        }
        return (
            <Grid container justify="center" alignItems="center" spacing={1}>
                {question.Options.map(option => {
                    return (
                        <Grid item align="center" xs={12} key={option}>
                            <Button 
                            className={classNames(classes.button, getOptionClass(option), disabledClass)}variant="outlined" 
                            onClick={() => dispatch(sendSingleGuess(roundId, question.Id, option))}
                            fullWidth>
                                {option}
                            </Button>
                        </Grid>
                    )
                })
                }
            </Grid>
        )
    }

    function getOptionClass(option){
        if (question.Answer && question.Answer.Value === option) {
            return classes.buttonCorrect
        }
        if (!question.Guesses || question.Guesses.length ===0 ) {
            return null
        }
        for (let i = 0; i<question.Guesses.length; i++){
            const guess = question.Guesses[i]
            if (guess.Content === option) {
                if (guess.Correct) {
                    return classes.buttonCorrect
                } else {
                    return classes.buttonIncorrect
                }
            }
        }
        return null
    }

    return (
        <Grid container direction="row" item xs={12} spacing={3}>
            <Grid item xs={12}>
                <Typography variant="h6" className={classes.description}>{question.Description}</Typography>
            </Grid>
            <Grid item xs={12}>
                <Typography variant="body1" className={classes.question}>{question.Content}</Typography>
            </Grid>
            <Grid item xs={12}>
                <OptionGrid/>
            </Grid>
            {question.Answer && question.Answer.ExtraInfo &&
                <Grid item xs={12}> 
                    <Typography variant="body1" className={classes.question}>{question.Answer.ExtraInfo}</Typography>
                </Grid>
            }
        </Grid>
    );
}


