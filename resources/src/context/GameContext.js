import React from 'react';

const notImplemented = () => new Error('not implemented yet');

export default React.createContext({
    createGame:notImplemented,
    connectMasterGame:notImplemented,
    startGame:notImplemented,
    getScore:notImplemented,
    getAnswersRepartition:notImplemented,
    joinGame:notImplemented,
    connectPlayerGame:notImplemented,
    nextQuestion:notImplemented,
    sendAnswer:notImplemented,
    getQuizzFromGame:notImplemented,
    endQuestion:notImplemented,
    computeScores:notImplemented,
    getMusic:notImplemented,
})
