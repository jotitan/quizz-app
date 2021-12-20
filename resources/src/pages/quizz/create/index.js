import React, {useState,useContext} from "react";
import QuizzContext from "../../../context/QuizzContext";
import './create.css';
import {Input, notification} from "antd";

export default function CreateQuizz(){
    const {createQuizz} = useContext(QuizzContext)
    const [title,setTitle] = useState('');
    const [description,setDescription] = useState('');

    const save = ()=>{
        createQuizz(title,description)
            .then(d=>{
                notification['success']({message:'Quizz created'})
                window.location.href=`/quizz/${d.id}`;
            })
            .catch(()=>notification["error"]({message:"Impossible to create quizz"}))
    }

    return <div>
        <h1>Create quizz</h1>

        <div style={{padding:20}}>
            <div>
                Titre
            </div>
            <div className={"field"}>
                <Input value={title} onChange={e=>setTitle(e.target.value)}/>
            </div>
            <div>Description</div>
            <div className={"field"}>
                <Input value={description} onChange={e=>setDescription(e.target.value)}/>
            </div>
            <button onClick={save}>Sauvegarder</button>
        </div>

    </div>
}
