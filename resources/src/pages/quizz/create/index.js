import React, {useState} from "react";
import './create.css';
import EditQuizz from "../../../components/EditQuizz";
import {Button} from "antd";

export default function CreateQuizz(){
    const [runEdit,setRunEdit] = useState(false);
    return <div>
        <h1>Create quizz</h1>
        <EditQuizz quizz={{}} runEdit={runEdit}/>

        <Button onClick={()=>setRunEdit(true)}>Enregistrer</Button>
    </div>
}
