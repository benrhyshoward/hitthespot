import React, { useState } from 'react';
import {
    sendGuess,
    abandonQuestion
} from './roundSlice' 
import { CorrectNotification } from './CorrectNotification'
import { useDispatch } from 'react-redux';
import { Typography, Grid, Button, makeStyles, TextField, FormGroup } from '@material-ui/core';

const useStyles = makeStyles({
    description: {
        textAlign: 'center',
        whiteSpace: 'pre-line'
      },
    image: {
        height: '100%',
        width: '100%',
        maxWidth: '300px',
        maxHeight: '300px',
        objectFit: 'cover'
      },
    image1 : {
        textAlign: 'right'
    },
    image2 : {
        textAlign: 'left'
    },
    form: {
        justifyContent: 'center'
      }
  });

export function PortmanteauQuestion(props) {

    const roundId = props.roundId
    const question = props.question
    const classes = useStyles();

    const [guess, setGuess] = useState("");

    const dispatch = useDispatch();

    function onGuess(e){
        e.preventDefault();
        dispatch(sendGuess(roundId, question.Id,guess));
        setGuess("")
    }
    return (
        <Grid container direction="row" item xs={12} spacing={2}>
            <Grid item xs={12}>
                <Typography variant="h6" className={classes.description}>{question.Description}</Typography>
            </Grid>
            <Grid item className={classes.image1} xs={6}> 
                <img alt=""className={classes.image} src={question.Images[0]}/>
            </Grid>
            <Grid item className={classes.image2} xs={6}> 
                <img alt="" className={classes.image} src={question.Images[1]}/>
            </Grid>
            {question.Answer.Value &&
                <Grid item xs={12}>
                    <Typography variant="h6" className={classes.question}>{question.Answer.Value}</Typography>
                </Grid>
            }
            {!question.Answer.Value &&
                <Grid item xs={12}>
                        <form noValidate autoComplete="off" onSubmit={onGuess}>
                                <FormGroup className={classes.form} row={true} >
                                    <TextField label="Answer" variant="outlined" value={guess} onChange={(event) => setGuess(event.target.value)} autoFocus/>
                                    <Button type="submit" variant="contained" color="primary">Submit</Button>
                                </FormGroup>    
                            </form>
                </Grid>
            }
            <Grid item xs={12}>
                <CorrectNotification/>
            </Grid>
            {!question.Answer.Value &&
                <Grid item xs={12}>
                    <Button variant="contained" color="secondary" onClick={() => dispatch(abandonQuestion(roundId, question.Id))}>Give Up</Button>
                </Grid>
            }
        </Grid>
    );
}


