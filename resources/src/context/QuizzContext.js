import React from 'react';

const notImplemented = () => new Error('not implemented yet');

export default React.createContext({
    getQuizz:notImplemented,
    getQuizzes:notImplemented,
    createOrUpdateQuizz:notImplemented,
    addQuestion:notImplemented,
    deleteQuestion:notImplemented,
    deleteQuizz:notImplemented,
})
