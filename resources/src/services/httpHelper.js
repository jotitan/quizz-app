
export function getBase(){
    if(window.document.location.origin.indexOf("localhost")!==-1) {
        return 'http://localhost:9001/api'
    }
    return `${window.document.location.origin}/api`
}

export function getBaseFront(){
    if(window.document.location.origin.indexOf("localhost")!==-1) {
        return "";
    }
    return "/quizz_app";
}

export default function httpHelper(){}
