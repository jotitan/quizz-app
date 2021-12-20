import React, {useContext, useEffect, useState} from 'react';
import GameContext from "../../context/GameContext";
import './Score.css';

export default function Score({id,sid,isEnded}){
    const {getScore} = useContext(GameContext);
    const [scores,setScores] = useState([])
    const [top3,setTop3] = useState([]);
    const lengthTop3 = 3;
    const format = values => {
        let computeScore = Object.entries(values)
            .sort((e1, e2) => e2[1] - e1[1])
            .map(e => {
                return {player: e[0], score: e[1]}
            });
        setScores(computeScore);
        setTop3(computeScore.length > lengthTop3 ? computeScore.filter((_,i)=>i<lengthTop3) : computeScore)
    }

    useEffect(()=>getScore(id,sid).then(format),[getScore,id,sid])
    return (
        <div>
            <div style={{textAlign:'center'}}>
                <h1>Score {isEnded ? 'final':''}</h1>
                <div className={'group-top-3'}>
                {top3.map((s,rank)=>
                    <div key={`score_${s.player}`} className={'score-player top3'}>
                        <span className={`rank ${isEnded && rank === 0?'winner':''}`}>{rank+1}</span> -
                        <span className={isEnded && rank === 0?'winner':''}>{s.player}</span> : {s.score}
                    </div>)}
                </div>
                {scores.filter((_,i)=>i>=lengthTop3 && i < 20).map((s,rank)=>
                    <div key={`score_${s.player}`} className={'score-player'}>
                        <span className={'rank'}>{rank+1+lengthTop3}</span> - {s.player} : {s.score}
                    </div>)}
            </div>
        </div>
    )
}
