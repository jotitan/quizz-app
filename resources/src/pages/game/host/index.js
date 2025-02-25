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
import {Switch} from "antd";

const getInitStatus = (id,idGame,idSecure) => {
    if(id !=null){
        return 'summary';
    }
    if(idGame !== '' && idSecure !==''){
        return 'connecting';
    }
    return 'shame';
}

//const testUsers = ["Robert","Marine","Roger","Philippe","Sandrine","Aline","Alice","Emma","Jeayon","Thomas","Henri","René","Florent","Fatima","Yohann","Brice","Jean-François", "Francis","Jean-Philippe','Laura","Diane","Pauline","Cécile","Céline"].map((u,i) => {return {player:u, position:i}})

export default function HostGame(){
    let id = useParams().id;
    let idGame = useParams().id_game;
    let idSecure = useParams().id_secure;

    const [status,setStatus] = useState('waiting')
    const {createGame,connectMasterGame, startGame, nextQuestion,getQuizzFromGame,endQuestion} = useContext(GameContext)
    const {getQuizz} = useContext(QuizzContext)
    const [totalPlayers,setTotalPlayers] = useState(0)
    const [question,setQuestion] = useState({})
    const [quizz,setQuizz] = useState({})
    const [users,setUsers] = useState([]);
    const [currentQuestion,setCurrentQuestion] = useState(-1);
    const [isEndQuestion,setIsEndQuestion] = useState(false);
    const [speedGame,setSpeedGame] = useState(false);

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

    const launchGame = ()=> createGame(id, speedGame).then(g=>window.location.href = `${getBaseFront()}/game/host/${g.id}/${g.secureId}`)

    const connectGame = ()=>{
        //Start SSE
        let sse = connectMasterGame(idGame,idSecure)
        sse.addEventListener("welcome",event=>{
            let data = JSON.parse(event.data);
            switch(data.status){
                case "waiting":
                    setStatus('waiting_players')
                    //setUsers(data.users)
                    break;
                case "answer":
                    setQuestion(data.question)
                    setUsers(data.users)
                    setTotalPlayers(data.total_players)
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
                if(u.some(v=>v.position===data.position)){
                    return u;
                }
                let copy = [...u]
                copy.push({player:data.player, position:data.position});
                setTotalPlayers(copy.length)
                return copy;
            });
        })
        sse.addEventListener("end-answers",event=>{
            setIsEndQuestion(true);
        })
        sse.addEventListener("answer",event=>{
            let data = JSON.parse(event.data);
            if(data.total_players != null) {
                setTotalPlayers(data.total_players)
            }
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

    const endGame = ()=> window.location.href=`${getBaseFront()}/`;

    const isEnded = ()=> currentQuestion  >= quizz.questions.length-1;

    const summary = ()=>{
        return (<div style={{textAlign:'center'}}>
            <h1>Play game <b>{quizz.name}</b></h1>
            <div> {quizz.questions.length} question(s)</div>
            <Switch onChange={setSpeedGame} /> Mode course
            <div>
                <button onClick={launchGame}>Lancer la partie</button>
            </div>
        </div>)
    }
    const url = `${window.location.origin}${getBaseFront()}/player?game=${idGame}`

    switch(status){
        case 'summary':return summary()
        case 'connecting':return 'Connecting'
        case 'waiting_players':return (
            <div>
                <div style={{textAlign:'center', width:'49%', display:'inline-block', marginTop:50, borderRight:"solid"}}>
                    <QRCode value={url} level={'M'} fgColor={"#0034ad"} size={384}/>
                </div>
                <div style={{textAlign:'center', width:'49%', display:'inline-block', verticalAlign:"top", paddingTop:20, marginTop:50}}>
                    <h2>Se connecter</h2>
                    <p>{url}</p>
                    <p>Code : <b>{idGame}</b></p>
                    <button onClick={start} className={"launch-game"}>Lancer la partie</button>
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
        case 'waiting_answers':return <ShowQuestion question={question} total={quizz.questions.length} users={users} totalUsers={totalPlayers} isEndQuestion={isEndQuestion} doEndQuestion={doEndQuestion}/>
        default :return ''
    }
}
