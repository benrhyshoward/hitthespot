import React, { useEffect } from 'react';
import { useParams, Link } from 'react-router-dom'
import {
    fetchRound,
    selectRound,
    selectScore,
    requestNewQuestion,
    abandonRound,
} from './roundSlice' 
import { useSelector, useDispatch } from 'react-redux';
import { Typography, Grid, Button, CircularProgress, makeStyles } from '@material-ui/core';
import { Logout } from '../auth/Logout'
import { FreeTextQuestion } from './FreeTextQuestion';
import { MultipleChoiceQuestion } from './MutipleChoiceQuestion';
import { PortmanteauQuestion } from './PortmanteauQuestion';
import { useHistory } from 'react-router-dom'

const useStyles = makeStyles({
    root: {
        margin: '0px',
        textAlign: 'center'
    },
    question: {
      whiteSpace: 'pre-line'
    }
  });

export function Round() {
    
    const classes = useStyles();

    const history = useHistory();
    const { id } = useParams();
    const dispatch = useDispatch();

    const round = useSelector(selectRound);
    const score = useSelector(selectScore)

    useEffect(() => {
        dispatch(fetchRound(id));
    }, [])

    function Question(props){
        const question = props.question
        if (question.Type === "FreeText") {
            return (
                <FreeTextQuestion roundId={props.roundId} question={question}/>
            )
        } else if (question.Type === "MultipleChoice") {
            return (
                <MultipleChoiceQuestion roundId={props.roundId} question={question}/>
            )
        } else if (question.Type === "Portmanteau") {
            return (
                <PortmanteauQuestion roundId={props.roundId} question={question}/>
            )
        }
        return null
    }

    function MoveOnInput(props){
        if (endOfRound(props.round)){
            return ScoresButton();
        } else {
            return NextQuestionButton();
        }
    }

    function NextQuestionButton() {
        return (
            <Button variant="contained" color="primary" onClick={() => dispatch(requestNewQuestion(round.Id))}>Next question</Button>
        )
    }

    function ScoresButton() {
        return (
            <Button variant="contained" color="primary" component={Link} to={"/"+id+"/s"}>Final scores</Button>
        )
    }

    function currentQuestion (round) {
        if (!round || !round.Questions || round.Abandoned){
            return null
        }
        //Find the first unanswered and unabandoned question
        const unansweredQuestions = round.Questions.filter(q => !questionOver(q));
        if (unansweredQuestions.length > 0) {
            return unansweredQuestions[0]
        } 
        //Otherwise find the most recently answered or abandoned question
        const answeredQuestions = round.Questions.filter(q => questionOver(q));
        if (answeredQuestions.length >0) {
            return answeredQuestions.sort((q1, q2) => {
                        var q1CorrectGuesses = q1.Guesses.filter(g => g.Correct)
                        var q2CorrectGuesses = q2.Guesses.filter(g => g.Correct)

                        var q1AnsweredAt = q1CorrectGuesses.length > 0 ? q1CorrectGuesses[0].Created : null
                        var q2AnsweredAt = q2CorrectGuesses.length > 0 ? q2CorrectGuesses[0].Created : null
                        return Math.max( q2.AbandonedAt, q2AnsweredAt) - Math.max(q1.AbandonedAt, q1AnsweredAt);
                    })[answeredQuestions.length-1]
        }
        return null
    }

    function endOfRound (round) {
        return (
            round &&                                             //round exists
            (round.Abandoned || (                                //round hasn't been abandoned
            round.Questions &&                                   //questions exist for round
            round.Questions.length === round.TotalQuestions &&   //have reached maximum number of questions
            round.Questions.every(q => questionOver(q))))        //every question has been answered
        )                   
    }

    function questionOver (question) {
        return (
            question &&                                                 //question exists
            (question.Abandoned || (                                    //question has been abandoned                             
            question.Guesses &&                                         //guesses exist for question                                   
            question.Guesses.some((guess) => guess.Correct === true)))  //at least one of the guesses is correct
        )
    }

    function questionNumber(round){
        if (!round || !round.Questions) {
            return 1
        }
        //counting the number of questions which have been abandoned or have a correct guess
        const finishedQuestions = round.Questions.filter(q => questionOver(q)).length

        if (questionOver(currentQuestion(round))){
            return finishedQuestions
        }
        return finishedQuestions + 1
    }

    return (
        <Grid container className={classes.root} direction="row" justify="center" alignItems="center" spacing={4}>
            <Grid item xs={12}>
                <Typography variant="h6">
                Score: {score} | Question {questionNumber(round)} of {round.TotalQuestions}
                </Typography>
            </Grid>
            { !currentQuestion(round) &&
                <CircularProgress/>
            }
            { currentQuestion(round) &&
                <Question roundId={round.Id} question={currentQuestion(round)}/>
            }
            { (endOfRound(round) || questionOver(currentQuestion(round))) &&
                <Grid item xs={12}>
                    <MoveOnInput round={round}/>
                </Grid>
            } 
            { !endOfRound(round) &&
                <Grid item xs={12}>
                    <Button variant="contained" color="secondary" onClick={() => dispatch(abandonRound(round.Id, history))}>End Round</Button>
                </Grid>
            }
            <Grid item xs={12}>
                <Logout/>        
            </Grid>
        </Grid>
    );
}


