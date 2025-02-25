import React, {useContext, useEffect, useState} from "react";
import QuizzContext from "../../../context/QuizzContext";
import {Card, Tooltip} from "antd";
import {PlusCircleOutlined, ThunderboltOutlined} from "@ant-design/icons";
import './list.css';
import {getBase, getBaseFront} from "../../../services/httpHelper";

export default function ListQuizzes(){
    const [quizzes,setQuizzes] = useState([]);
    const {getQuizzes} = useContext(QuizzContext)

    useEffect(()=>getQuizzes().then(setQuizzes),[getQuizzes])

    const getBackgroundStyle = q => {
        return q.img ? {
            backgroundImage:`url(${getBase()}/quizz/${q.id}/cover)`,
            height:150,
            backgroundRepeat:'no-repeat',
            backgroundSize:'contain',
            backgroundPosition:'center'
        }:{backgroundImage: 'red'}
    }

    const computeExtra = quizz => {
        return quizz.nb > 0 ?
        <Tooltip title={"Jouer"}>
            <ThunderboltOutlined className={"icon-action"} onClick={()=>window.location.href = `${getBaseFront()}/game/host/${quizz.id}`}/>
        </Tooltip> :''
    }

    const showQuizz = quizz => {
        return <Card title={<a href={`${getBaseFront()}/quizz/${quizz.id}`}>{quizz.name} ({quizz.nb})</a>}
                     key={`div_${quizz.id}`}
                     headStyle={{height:40}}
                     className={`quizz-card ${quizz.img?'card-image':''}`}
                     extra={computeExtra(quizz)}>
            <p style={getBackgroundStyle(quizz)} className={'description'}>
               <span>{quizz.description}</span>
            </p>
        </Card>
    }

    return <div>
        <h1>Tous les quizz
            <Tooltip title={"Créer un quizz"}><a href={`${getBaseFront()}/quizz/create`} style={{cursor:'pointer', marginLeft:15}}><PlusCircleOutlined /></a></Tooltip></h1>
        {quizzes.sort((a,b)=>a.name.localeCompare(b.name)).map(showQuizz)}
    </div>
}
