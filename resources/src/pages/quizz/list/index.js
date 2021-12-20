import React, {useContext, useEffect, useState} from "react";
import QuizzContext from "../../../context/QuizzContext";
import {Card, notification, Popconfirm, Tooltip} from "antd";
import {DeleteOutlined, ThunderboltOutlined} from "@ant-design/icons";
import './list.css';
import RoundButton from "../../../components/RoundButton";

export default function ListQuizzes(){
    const [quizzes,setQuizzes] = useState([]);
    const [forceRefresh,setForceRefresh] = useState(false);
    const {getQuizzes,deleteQuizz} = useContext(QuizzContext)


    useEffect(()=>getQuizzes().then(setQuizzes),[getQuizzes,forceRefresh])

    const doDeleteQuizz = id => {
        deleteQuizz(id).then(()=>{
            notification["success"]({message:"Quizz supprimé"})
            setForceRefresh(f=>!f)
        })
    }

    const showQuizz = quizz => {
        return <Card title={<a href={`/quizz/${quizz.id}`}>{quizz.name}</a>} key={`div_${quizz.id}`} style={{ width: 300,marginLeft:20 }}
                     extra={<Tooltip title={"Jouer"}>
                         <ThunderboltOutlined className={"icon-action"} onClick={()=>window.location.href = `/game/host/${quizz.id}`}/>
                     </Tooltip>}>
            <p>{quizz.nb} question(s)</p>
            <Popconfirm okText={"Supprimer"} cancelText={"Annuler"} onConfirm={()=>doDeleteQuizz(quizz.id)} title={"Etes vous sur de vouloir supprimer ce quizz"}>
                <Tooltip title={"Supprimer"}>
                    <DeleteOutlined className={"icon-action"}/>
                </Tooltip>
            </Popconfirm>

        </Card>
    }

    return <div>
        <h1>Quizz application</h1>
        {quizzes.map(showQuizz)}

        <RoundButton title={"Créer un quizz"} action={()=>window.location.href='/quizz/create'}/>
    </div>
}
