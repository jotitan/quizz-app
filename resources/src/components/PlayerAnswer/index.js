import React, {useContext} from 'react';
import '../ShowQuestion/showQuestion.css'
import './PlayerAnswer.css';
import ShowAnswers from "../ShowAnswers";
import GameContext from "../../context/GameContext";

export default function PlayerAnswer({question,gameId,playerId,setStatus}){
    const {sendAnswer} = useContext(GameContext)


    const actionAnswer = pos=> {
        sendAnswer(gameId,playerId,pos)
        setStatus('answer_done')
    }

    const showAnswer = nb=> <ShowAnswers action={actionAnswer} answers={Array.from({length: nb}, (v, i) => i).map(()=>'')}/>;

    return <div style={{height:'calc(100% - 45px)'}}>
        <div className={'title-player'}>
            Question {question.position+1}
        </div>
        <div style={{height:'calc(100% - 50px)'}}>
            {showAnswer(question.nb)}
        </div>
    </div>
}
