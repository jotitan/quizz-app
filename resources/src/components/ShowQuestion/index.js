import React, {useEffect, useState} from 'react';
import './showQuestion.css';
import Countdown from "../Countdown";
import {CodepenCircleOutlined, GithubOutlined, GitlabOutlined, QqOutlined, TeamOutlined} from "@ant-design/icons";
import {Button, Col, Row, Space} from "antd";
import ReadMusic from "../ReadMusic";
import {useParams} from "react-router-dom";

const icons = [<GithubOutlined />,<GitlabOutlined />,<CodepenCircleOutlined />,<QqOutlined />]

export default function ShowQuestion({question,total,totalUsers, users, isEndQuestion=false,doEndQuestion}){
    const secureId = useParams().id_secure;
    const gameId = useParams().id_game;
    const [showButton,setShowButton] = useState(false);

    const showAnswer = (answer,pos)=>{
        return <div key={pos} className={'answer'}>
            <div className={`color${pos}`}>
                {icons[pos]}
                {answer.text}
            </div>
        </div>
    }

    useEffect(()=>{
        if(isEndQuestion){
            setShowButton(true)
        }
    },[isEndQuestion])

    const showNextButton = ()=>setShowButton(true);

    return <div>
        <div style={{height:'60px',paddingLeft:20,backgroundColor:'black'}}>
            <Row style={{paddingTop:10}}>
                <Col span={18}>
                    <span className={'title-question'}>Question n°{question.position+1} / {total}</span>
                </Col>
                <Col span={6} style={{right:20,textAlign:'right'}} className={'title-question'}>
                    <TeamOutlined /> {users.length} / {totalUsers}
                </Col>
            </Row>

        </div>
        <Space className={"question"}>
            <div style={{display:"inline-block",width:question.filename !=='' ? 'calc(55vw - 10px)':'calc(100vw - 20px)'}}>
                <div><p>{question.title}</p></div>
            </div>
            {question.filename !=='' ?<div style={{width:'calc(45vw - 15px)',display:'inline-block'}}>
                    <ReadMusic question={question} secureId={secureId} gameId={gameId}/>
            </div>:''}
        </Space>
        <div style={{height:'50vh'}}>
            {question.answers.map(showAnswer)}
        </div>
        <div style={{height:'60px'}}>
            <Countdown duration={question.time} style={{display:'inline-block',width:'80%',paddingLeft:20,paddingTop:16}} endAction={showNextButton}/>
            <div style={{display:'inline-block',float:'right',paddingTop:10,paddingRight:20}}>
                <Button disabled={!showButton} type={"primary"} onClick={doEndQuestion}>Suite</Button>
            </div>
        </div>
    </div>
}
