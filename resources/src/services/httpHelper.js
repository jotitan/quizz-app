
export default function getBase(){
    if(window.document.location.origin.indexOf("location")!==-1) {
        return 'http://localhost:9001/api'
    }
    return `${window.document.location.origin}/api`
}
