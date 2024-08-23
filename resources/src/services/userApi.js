import {getBase} from './httpHelper'

function isAdmin(){
    return fetch(`${getBase()}/user/is_admin`).then(d=>d.json())
}



const userApi = {
    isAdmin
}

export default userApi;
