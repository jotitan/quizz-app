import React, {useEffect, useState} from "react";
import './App.css';
import {BrowserRouter as Router, Route, Routes} from "react-router-dom";
import ListQuizzes from "./pages/quizz/list";
import ShowQuizz from "./pages/quizz/show";
import CreateQuizz from "./pages/quizz/create";
import QuizzContext from "./context/QuizzContext";
import UserApi from "./services/userApi";
import GameContext from "./context/GameContext";
import QuizzApi from "./services/quizzApi"
import GameApi from "./services/gameApi"
import HostGame from "./pages/game/host";
import JoinGame from "./pages/game/join";
import PlayGame from "./pages/game/play";
import Header from "./components/Header";
import {getBaseFront} from "./services/httpHelper";
import useLocalStorage from "./services/local-storage.hook";
import NoAccess from "./pages/noaccess";

function App() {
    const [loading, setLoading] = useState(true)
    const [isAdmin, setIsAdmin] = useLocalStorage("is_admin")

    const gameApi = {
        createGame: GameApi.createGame,
        startGame: GameApi.startGame,
        connectMasterGame: GameApi.connectMasterGame,
        getScore: GameApi.getScore,
        getAnswersRepartition: GameApi.getAnswersRepartition,
        joinGame: GameApi.joinGame,
        connectPlayerGame: GameApi.connectPlayerGame,
        nextQuestion: GameApi.nextQuestion,
        sendAnswer: GameApi.sendAnswer,
        getQuizzFromGame: GameApi.getQuizzFromGame,
        endQuestion: GameApi.endQuestion,
        computeScores: GameApi.computeScores,
        getMusic: GameApi.getMusic,
        getPlayerDetail: GameApi.getPlayerDetail,
    }
    const quizzApi = {
        createOrUpdateQuizz: QuizzApi.createOrUpdateQuizz,
        getQuizz: QuizzApi.getQuizz,
        addQuestion: QuizzApi.addQuestion,
        getQuizzes: QuizzApi.getQuizzes,
        deleteQuestion: QuizzApi.deleteQuestion,
        deleteQuizz: QuizzApi.deleteQuizz,
    }
    const userApi = {
        isAdmin: UserApi.isAdmin,
    }

    useEffect(()=> {
        userApi.isAdmin().then(details => {
            setIsAdmin(details.is_admin)
            setLoading(false);
        })
        //eslint-disable-next-line
    },[setIsAdmin])

    function createSecureRoute(path, c){
        return  <Route path={path} element={isAdmin ? React.createElement(c):<NoAccess/>}/>;
    }

    function getContent(){
        return <Header>
            <div>
                <QuizzContext.Provider value={quizzApi}>
                    <GameContext.Provider value={gameApi}>
                            <Router basename={getBaseFront()}>
                                <Routes>
                                    {createSecureRoute("/",ListQuizzes)}
                                    {createSecureRoute("/quizz/create",CreateQuizz)}
                                    {createSecureRoute("/quizzes",ListQuizzes)}
                                    {createSecureRoute("/game/host/:id",HostGame)}
                                    {createSecureRoute("/game/host/:id_game/:id_secure",HostGame)}
                                    <Route path="/player/play/:id_game/:id_player" element={<PlayGame/>}/>
                                    <Route path="/quizz/:id" element={<ShowQuizz/>}/>
                                    <Route path="/player" element={<JoinGame/>}/>
                                </Routes>
                            </Router>
                    </GameContext.Provider>
                </QuizzContext.Provider>
            </div>
        </Header>
    }

    return loading ? 'loading' : getContent();

}

export default App;
