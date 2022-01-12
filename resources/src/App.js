import React from "react";
import './App.css';
import {BrowserRouter as Router, Route, Routes} from "react-router-dom";
import ListQuizzes from "./pages/quizz/list";
import ShowQuizz from "./pages/quizz/show";
import CreateQuizz from "./pages/quizz/create";
import QuizzContext from "./context/QuizzContext";
import GameContext from "./context/GameContext";
import QuizzApi from "./services/quizzApi"
import GameApi from "./services/gameApi"
import HostGame from "./pages/game/host";
import JoinGame from "./pages/game/join";
import PlayGame from "./pages/game/play";
import Header from "./components/Header";
import {getBaseFront} from "./services/httpHelper";

function App() {
    const gameApi = {
        createGame:GameApi.createGame,
        startGame:GameApi.startGame,
        connectMasterGame:GameApi.connectMasterGame,
        getScore:GameApi.getScore,
        joinGame:GameApi.joinGame,
        connectPlayerGame:GameApi.connectPlayerGame,
        nextQuestion:GameApi.nextQuestion,
        sendAnswer:GameApi.sendAnswer,
        getQuizzFromGame:GameApi.getQuizzFromGame,
        endQuestion:GameApi.endQuestion,
        computeScores:GameApi.computeScores,
        getMusic:GameApi.getMusic,
    }
    const quizzApi = {
        createOrUpdateQuizz:QuizzApi.createOrUpdateQuizz,
        getQuizz:QuizzApi.getQuizz,
        addQuestion:QuizzApi.addQuestion,
        getQuizzes:QuizzApi.getQuizzes,
        deleteQuestion:QuizzApi.deleteQuestion,
        deleteQuizz:QuizzApi.deleteQuizz,
    }
    return (
        <Header>
        <div>
            <QuizzContext.Provider value={quizzApi}>
            <GameContext.Provider value={gameApi}>
                <Router basename={getBaseFront()}>
                    <Routes>
                        <Route path="/"  element={<ListQuizzes/>}/>
                        <Route path="/quizz/create" element={<CreateQuizz/>}/>
                        <Route path="/quizz/:id" element={<ShowQuizz/>}/>
                        <Route path="/quizzes"  element={<ListQuizzes/>}/>
                        <Route path="/game/host/:id"  element={<HostGame/>}/>
                        <Route path="/game/host/:id_game/:id_secure"  element={<HostGame/>}/>
                        <Route path="/player/play/:id_game/:id_player"  element={<PlayGame/>}/>
                        <Route path="/player"  element={<JoinGame/>}/>
                    </Routes>
                </Router>
            </GameContext.Provider>
            </QuizzContext.Provider>
        </div>
        </Header>

    );
}

export default App;
