import getBase from './httpHelper'

function getQuizz(id){
    return fetch(`${getBase()}/quizz/${id}`).then(d=>d.json())
}

function getQuizzes(){
    return fetch(`${getBase()}/quizzes`).then(d=>d.json())
}

function createQuizz(name,description){
    return fetch(`${getBase()}/quizz?name=${name}&description=${description}`,{method:'POST'}).then(d=>d.json())
}

function addQuestion(id,question){
    const form = new FormData();
    if(question.music != null && !question.music.delete){
        // Question with music already exists
        if(typeof question.music === 'string'){
            question.music = {keepExisting:true};
        }else{
            // Add new music
            question.range = question.music.range
            if(question.music.replace){
                question.music.delete = true;
            }
            form.append("music",question.music.music)
            question.music.music = null;
        }
    }else{
        if(question.filename){
            question.music = {keepExisting:true};
        }
    }
    form.append("question",JSON.stringify(question));
    const headers = {
        //'Content-Type': 'multipart/form-data'
    }
    return fetch(`${getBase()}/quizz/${id}/question`,{headers:headers,method:'POST',body:form})
}

function deleteQuestion(quizzId,questionId){
    return fetch(`${getBase()}/quizz/${quizzId}/question/${questionId}`,{method:'DELETE'})
}


function deleteQuizz(quizzId){
    return fetch(`${getBase()}/quizz/${quizzId}`,{method:'DELETE'})
}

const quizzApi = {
    getQuizz,
    getQuizzes,
    createQuizz,
    addQuestion,
    deleteQuestion,
    deleteQuizz,
}

export default quizzApi;
