import React, {useContext, useEffect, useState} from "react";
import QuizzContext from "../../../context/QuizzContext";
import {Card, Tooltip} from "antd";
import {ThunderboltOutlined} from "@ant-design/icons";
import './list.css';
import RoundButton from "../../../components/RoundButton";
import getBase from "../../../services/httpHelper";

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

    const showQuizz = quizz => {
        return <Card title={<a href={`/quizz/${quizz.id}`}>{quizz.name} ({quizz.nb})</a>}
                     key={`div_${quizz.id}`}
                     headStyle={{height:40}}
                     className={`quizz-card ${quizz.img?'card-image':''}`}
                     extra={<Tooltip title={"Jouer"}>
                         <ThunderboltOutlined className={"icon-action"} onClick={()=>window.location.href = `/game/host/${quizz.id}`}/>
                     </Tooltip>}>
            <p style={getBackgroundStyle(quizz)} className={'description'}>
               <span>{quizz.description}</span>
            </p>
        </Card>
    }

    return <div>
        <h1>Tous les quizz</h1>
        {quizzes.sort((a,b)=>a.name.localeCompare(b.name)).map(showQuizz)}

        <RoundButton title={"Créer un quizz"} action={()=>window.location.href='/quizz/create'}/>
    </div>
}
