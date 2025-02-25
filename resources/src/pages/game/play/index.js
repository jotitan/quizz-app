import  React,{useContext, useEffect, useState} from 'react';
import {useParams, useSearchParams} from "react-router-dom";
import GameContext from "../../../context/GameContext";
import PlayerAnswer from "../../../components/PlayerAnswer";
import {Col, Row} from "antd";

export default function PlayGame(){
    let params = useParams();
    let [query] = useSearchParams();
    const [status,setStatus] = useState('connecting')
    const {connectPlayerGame, getPlayerDetail} = useContext(GameContext)
    const [question,setQuestion] = useState({})
    const [score,setScore] = useState({rank:0,score:0});
    const [position, setPosition] = useState(0);

    const getPosition = () => {
        getPlayerDetail(params.id_game, params.id_player).then(p=>setPosition(p.position))
    }

    const createSSEConnection = ()=>{
        const sse = connectPlayerGame(params.id_game,params.id_player,query.get("name"));
        sse.addEventListener('welcome',event=>{
            const data = JSON.parse(event.data)
            switch(data.status){
                case 'waiting':setStatus('waiting');break;
                case 'score':
                    setScore({score:data.score,rank:data.rank,end:data.end})
                    if(data.end){
                        setStatus('end_game')
                    }else {
                        setStatus('waiting');
                    }
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

    const showEnd = ()=> {
        return <div style={{textAlign:'center',fontSize:24}}>
            <div>Fin de partie</div>
            <div style={{fontWeight:'bold',marginTop:20}}>{score.rank === 0 ? 'Bravo, tu as gagné le quizz':''}</div>
        </div>
    }

    useEffect(()=>{
        createSSEConnection()
        getPosition()
        // eslint-disable-next-line
    },[])

    const switchStatus = ()=> {
        switch(status){
            case "connecting":return "En attente de connexion";
            case "waiting":return "Vous êtes connecté, en attente de démarrage";
            case "answer_done":return "Réponse envoyée, en attente...";
            case "end_game":return showEnd();
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
                <Col span={8} style={{textAlign:'center'}}>Score : <b>{score.score}</b></Col>
                <Col span={8} style={{textAlign:'center'}}><img alt="icon" style={{height:28}} src={`/icons/icon_${position < 10 ? '0':''}${position}.svg`}/></Col>
                <Col span={8} style={{textAlign:'center'}}>Place : <b>{score.rank+1}</b></Col>
            </Row>
        </div>
    }

    return <div>
        {showScore()}
        {switchStatus()}
    </div>
}
