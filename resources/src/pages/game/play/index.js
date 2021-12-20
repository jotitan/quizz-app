import React from 'react';
import  {useContext, useEffect, useState} from 'react';
import {useParams, useSearchParams} from "react-router-dom";
import GameContext from "../../../context/GameContext";
import PlayerAnswer from "../../../components/PlayerAnswer";
import {Col, Row} from "antd";

export default function PlayGame(){
    let params = useParams();
    let [query] = useSearchParams();
    const [status,setStatus] = useState('connecting')
    const {connectPlayerGame} = useContext(GameContext)
    const [question,setQuestion] = useState({})
    const [score,setScore] = useState({rank:0,score:0});

    const createSSEConnection = ()=>{
        const sse = connectPlayerGame(params.id_game,params.id_player,query.get("name"));
        sse.addEventListener('welcome',event=>{
            const data = JSON.parse(event.data)
            switch(data.status){
                case 'waiting':setStatus('waiting');break;
                case 'score':
                    setScore({score:data.score,rank:data.rank})
                    setStatus('waiting');
                    break;
                case 'question':
                    setQuestion(data.question);
                    setStatus('answer');
                    setScore({score:data.score,rank:data.rank})
                    break;
                default:
            }
        })
        sse.addEventListener('score',event=>{
            setScore(JSON.parse(event.data));
        })
        sse.addEventListener('question',event=>{
            setStatus('answer')
            setQuestion(JSON.parse(event.data));
        })
    }

    useEffect(()=>{
        createSSEConnection()
        // eslint-disable-next-line
    },[])

    const switchStatus = ()=> {
        switch(status){
            case "connecting":return "En attente de connexion";
            case "waiting":return "Vous êtes connecté, en attente de démarrage";
            case "answer_done":return "Réponse envoyée, en attente...";
            case "answer":return <PlayerAnswer question={question} gameId={params.id_game} playerId={params.id_player} setStatus={setStatus}/>
            default:return ""
        }
    }

    const showScore = ()=> {
        if(status === 'connecting'){
            return "";
        }
        return <div>
            <Row style={{borderBottom:'solid 1px gray',padding:'10px 0px',fontVariant:'small-caps'}}>
                <Col span={12} style={{textAlign:'center'}}>Score : {score.score}</Col>
                <Col span={12} style={{textAlign:'center'}}>Place : {score.rank+1}</Col>
            </Row>
        </div>
    }

    return <div>
        {showScore()}
        {switchStatus()}
    </div>
}
