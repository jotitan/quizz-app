import React, {useContext, useEffect, useState} from "react";
import {useParams} from "react-router-dom";
import GameContext from "../../../context/GameContext";
import './play.css';
import '../../../App.css';
import ShowPlayers from "../../../components/ShowPlayers";
import Score from "../../../components/Score";
import QuizzContext from "../../../context/QuizzContext";
import RoundButton from "../../../components/RoundButton";
import ShowQuestion from "../../../components/ShowQuestion";
import ComputeScore from "../../../components/ComputeScore";
import QRCode from 'qrcode.react';
import {getBaseFront} from "../../../services/httpHelper";

const getInitStatus = (id,idGame,idSecure) => {
    if(id !=null){
        return 'summary';
    }
    if(idGame !== '' && idSecure !==''){
        return 'connecting';
    }
    return 'shame';
}

export default function HostGame(){
    let id = useParams().id;
    let idGame = useParams().id_game;
    let idSecure = useParams().id_secure;

    const [status,setStatus] = useState('waiting')
    const {createGame,connectMasterGame, startGame, nextQuestion,getQuizzFromGame,endQuestion} = useContext(GameContext)
    const {getQuizz} = useContext(QuizzContext)
    const [question,setQuestion] = useState({})
    const [quizz,setQuizz] = useState({})
    const [users,setUsers] = useState([]);
    const [currentQuestion,setCurrentQuestion] = useState(-1);
    const [isEndQuestion,setIsEndQuestion] = useState(false);

    const getQuizzMethod = ()=>id != null ? ()=>getQuizz(id) : ()=>getQuizzFromGame(idGame,idSecure);
    const showScore = ()=> setStatus("score")

    const doEndQuestion = ()=> {
        endQuestion(idGame,idSecure)
            .then(()=>setStatus("compute_score"))
    }

    useEffect(()=>{
        getQuizzMethod()().then(q=>{
            setQuizz(q)
            setStatus(getInitStatus(id,idGame,idSecure))
        })
        // eslint-disable-next-line
    },[])

    const launchGame = ()=> createGame(id).then(g=>window.location.href = `${getBaseFront()}/game/host/${g.id}/${g.secureId}`)

    const connectGame = ()=>{
        //Start SSE
        let sse = connectMasterGame(idGame,idSecure)
        sse.addEventListener("welcome",event=>{
            let data = JSON.parse(event.data);
            switch(data.status){
                case "waiting":
                    setStatus('waiting_players')
                    setUsers(data.users)
                    break;
                case "answer":
                    setQuestion(data.question)
                    setUsers(data.users)
                    setStatus('waiting_answers')
                    break;
                case "compute_score":
                    setStatus('compute_score')
                    break;
                case "score":
                    setCurrentQuestion(data.current);
                    setStatus("score");
                    break;
                default:setStatus("score")
            }
        })
        sse.addEventListener("join",event=>{
            let data = JSON.parse(event.data);

            setUsers(u=>{
                // Check if already exist in list
                if(u.some(v=>v===data.player)){
                    return u;
                }
                let copy = [...u]
                copy.push(data.player);
                return copy;
            });
        })
        sse.addEventListener("end-answers",event=>{
            setIsEndQuestion(true);
        })
        sse.addEventListener("answer",event=>{
            let data = JSON.parse(event.data);
            setUsers(u=>{
                let copy = [...u]
                copy.push(data.player);
                return copy;
            });
        })
    }

    useEffect(()=>{
        if(status === 'connecting'){
            connectGame();
        }
        // eslint-disable-next-line
    },[status])

    const start = ()=> startGame(idGame,idSecure).then(()=>setStatus("score"))

    const goNextQuestion = ()=> {
        setUsers([])
        nextQuestion(idGame,idSecure).then(q=> {
            setQuestion(q)
            setCurrentQuestion(q.position);
            setIsEndQuestion(false);
            setStatus("waiting_answers")
        })
    }

    const endGame = ()=> window.location.href=getBaseFront();

    const isEnded = ()=> currentQuestion  >= quizz.questions.length-1;

    const summary = ()=>{
        return (<div style={{textAlign:'center'}}><h1>Play game '{quizz.name}' ({quizz.questions.length})</h1>
            <RoundButton title={"Lancer la partie"} action={launchGame}/>
        </div>)
    }
    const url = `${window.location.origin}/${getBaseFront()}/player?game=${idGame}`
    switch(status){
        case 'summary':return summary()
        case 'connecting':return 'Connecting'
        case 'waiting_players':return (
            <div>
                <div style={{marginTop:20}}>
                    Pour vous connecter à la partie : {url} avec le code : {idGame}
                </div>
                <div style={{textAlign:'center',marginTop:50}}>
                    <QRCode value={url} level={'M'} fgColor={"#0034ad"} size={256}/>
                </div>
                <ShowPlayers users={users} start={start}/>
            </div>)
        case 'compute_score':return <div>
            <ComputeScore gameId={idGame} secureId={idSecure}/>
            <RoundButton title={"Score"} action={showScore}/>
        </div>
        case 'score':
            return <div>
                <Score id={idGame} sid={idSecure} isEnded={isEnded()}/>
                {isEnded()?
                    <RoundButton title={"Retour accueil"} action={endGame}/>:
                    <RoundButton title={"Question suivante"} action={goNextQuestion}/>}
            </div>
        case 'waiting_answers':return <ShowQuestion question={question} total={quizz.questions.length} users={users} isEndQuestion={isEndQuestion} doEndQuestion={doEndQuestion}/>
        default :return ''
    }
}
