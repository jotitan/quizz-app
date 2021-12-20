import getBase from './httpHelper'

function createGame(id){
    return fetch(`${getBase()}/game/create/${id}`,{method:'POST'}).then(d=>d.json())
}

function connectMasterGame(id,secureId){
    return new EventSource(`${getBase()}/game/${id}/connect/${secureId}`)
}

function startGame(id,secureId){
    return fetch(`${getBase()}/game/${id}/start/${secureId}`,{method:'POST'})
}

function getScore(id,secureId){
    return fetch(`${getBase()}/game/${id}/score/${secureId}`).then(d=>d.json())
}

function joinGame(id,name){
    return fetch(`${getBase()}/player/join/${id}?name=${name}`,{method:'POST'}).then(d=>{
        return d.status === 400 ?
            new Promise((_,r)=>d.text().then(err=>r(err))):d.json()
    })
}

function connectPlayerGame(id, playerId, playerName){
    return new EventSource(`${getBase()}/player/connect/${id}/${playerId}?name=${playerName}`)
}

function nextQuestion(id,secureId){
    return fetch(`${getBase()}/game/${id}/playNextQuestion/${secureId}`,{method:'POST'}).then(d=>d.json())
}

function getQuizzFromGame(id,secureId){
    return fetch(`${getBase()}/game/${id}/quizz/${secureId}`).then(d=>d.json())
}

function sendAnswer(gameId,playerId,pos){
    return fetch(`${getBase()}/player/answer/${gameId}/${playerId}?answer=${pos}`,{method:'POST'})
}

function endQuestion(id,secureId){
    return fetch(`${getBase()}/game/${id}/forceEndQuestion/${secureId}`,{method:'POST'})
}

function computeScores(id,secureId){
    return fetch(`${getBase()}/game/${id}/computeScores/${secureId}`,{method:'POST'}).then(d=>d.json());
}

function getMusic(gameId,secureId,questionId){
    return fetch(`${getBase()}/game/${gameId}/music/${questionId}/${secureId}`).then(d=>d.arrayBuffer());
}

const gameApi = {
    createGame,
    startGame,
    connectMasterGame,
    getScore,
    joinGame,
    connectPlayerGame,
    nextQuestion,
    sendAnswer,
    getQuizzFromGame,
    endQuestion,
    computeScores,
    getMusic,
}

export default gameApi;
