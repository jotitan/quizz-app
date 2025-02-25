import React from 'react';
import './ShowAnswers.css';

import {CodepenCircleOutlined, GithubOutlined, GitlabOutlined, QqOutlined} from "@ant-design/icons";

const icons = [<GithubOutlined />,<GitlabOutlined />,<CodepenCircleOutlined />,<QqOutlined />]

export default function ShowAnswers({answers,action=null}){

    const nbLines = Math.ceil(answers.length/2)
    return answers.map((answer,pos)=>
        <div key={`answer${pos}`} style={{lineHeight:`calc(${100/nbLines - 2}vh - ${95/nbLines}px - 26px)`,height:`calc(${100/nbLines - 2}vh - ${95/nbLines}px)`}} className={`answers`}>
            <div
                className={`color${pos%4}`}
                onClick={()=>action != null ? action(pos):{}}
                style={{cursor:`${(action!=null ? 'pointer':'auto')}`}}>
                {icons[pos]} {answer}
            </div>
        </div>
    )
}
