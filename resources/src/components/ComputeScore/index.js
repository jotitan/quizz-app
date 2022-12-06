import React,{useContext, useEffect, useState} from "react";
import GameContext from "../../context/GameContext";
import './ComputeScore.css';
import ChartResults from "../ChartResults";

export default function ComputeScore({gameId,secureId}){
    const {computeScores} = useContext(GameContext)
    const [scoreQuestion,setScoreQuestion] = useState({winners:[],good:[],repartition:[]})

    useEffect(()=>{
        computeScores(gameId,secureId).then(scores => {
            setScoreQuestion(scores)
        })
        // eslint-disable-next-line
    },[])

    const showPlayer = p=>{
        return <span key={`good_${p}`} className={'good-player'}>{p}</span>
    }

    const showPlayers = players => {
        switch(players.length){
            case 0 : return <div>Personne n'a trouvé la réponse :(</div>;
            case 1 : return <div>
                Bravo {showPlayer(players[0])}, tu es le seul à avoir trouvé
            </div>
            default:return <div>
                Bravo aux {players.length} qui ont trouvés :
                <div style={{marginTop:30}}>
                    {scoreQuestion.winners.map(showPlayer)}
                </div>
            </div>;
        }
    }

    return <div>
        <div style={{textAlign:'center',marginTop:20}}>
            {scoreQuestion.good.length >= 2 ? <div className={"good-answer"}>
                Les bonnes réponses sont : {scoreQuestion.good.map((v,i)=><span>{i>0?",":""} {v}</span>)}
            </div>:<div className={"good-answer"}>
                La bonne réponse est : <span>{scoreQuestion.good[0]}</span>
            </div>}
        </div>
        <div style={{textAlign:'center',marginTop:30,fontSize:18}}>
            {showPlayers(scoreQuestion.winners)}
        </div>
        <div>
            <ChartResults data={scoreQuestion.repartition} answers={scoreQuestion.good}/>
        </div>

    </div>
}
