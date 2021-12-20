import React, {useEffect, useState} from "react";
import {Divider, InputNumber, Space, Tooltip} from "antd";
import {
    DeleteOutlined,
    DislikeTwoTone,
    LikeTwoTone,
    MinusCircleTwoTone,
    PlusCircleTwoTone,
    SoundOutlined
} from '@ant-design/icons';
import './EditQuestion.css'
import AudioUpload from "../AudioUpload";

export default function EditQuestion({question,updateQuestion}){

    const [title,setTitle] = useState(question.title)
    const [maxTime,setMaxTime] = useState(question.maxTime || 30)
    const [music,setMusic] = useState({})
    useEffect(()=>setTitle(question.title),[question])
    const update = () => {
        updateQuestion(q=>{
            let copy = {...q};
            copy.title = title;
            copy.time = maxTime;
            if(music.music != null){
                const replace = copy.music != null && (copy.music.delete || copy.music.replace);
                copy.music = music;
                copy.music.replace = replace;
            }
            return copy
        })
    }
    useEffect(()=>{
        update()
        // eslint-disable-next-line
    },[music])

    const wrapSetMusic = m => {
        setMusic(previous=> {
            if(previous.delete){
                m.delete = previous.delete;
            }
            return m;
        })
    }

    const deleteMusic = ()=> {
        updateQuestion(q=>{
            let copy = {...q};
            copy.music = {delete:true};
            copy.filename = null;
            return copy;
        })
    }

    const addAnswer = ()=>{
        updateQuestion(q=>{
            let copy = {...q};
            let answer = {text:'',good:false,maxTime:30}
            if(copy.answers == null){
                copy.answers = [answer]
            } else{
                copy.answers.push(answer);
            }
            return copy
        })
    }
    const removeAnswer = ()=>{
        updateQuestion(q=>{
            let copy = {...q};
            copy.answers = copy.answers == null ? []:copy.answers.slice(0,-1)
            return copy;
        })
    }
    const showMusic = ()=> {
        return <div style={{paddingTop:10}}>
            <SoundOutlined /> {question.filename}
            <Tooltip title={"Supprimer"}>
                <DeleteOutlined onClick={deleteMusic}/>
            </Tooltip>
        </div>;
    }

    const displayAnswer = (a,pos) => {
        return <div key={`div_${pos}`}><Space>
            <input style={{width:400}} value={a.text} onChange={e=>{
                updateQuestion(q=>{
                    let copy = {...q}
                    copy.answers[pos].text = e.target.value
                    return copy
                })
            }
            }/>
            <span onClick={()=>updateQuestion(q=>{
                let copy = {...q}
                copy.answers[pos].good = !copy.answers[pos].good;
                return copy
            })}>
                {a.good ? <LikeTwoTone twoToneColor={"#52c41a"} className={"button-action"}/>:<DislikeTwoTone twoToneColor={"#ca1e1e"} className={"button-action"}/>}
            </span>
        </Space></div>
    }
    return (
        <div>
            <div className={"question-label"}>Titre</div>
            <div style={{paddingLeft:10}}>
                <input value={title}
                       style={{width:'100%'}}
                       onChange={e=>setTitle(e.target.value)}
                       onBlur={()=>update()} key={'input_title'}/>
            </div>
            <div className={"question-label"}>Durée (en secondes)</div>
            <div style={{paddingLeft:10}}>
                <InputNumber value={maxTime}
                             style={{width:'100%'}}
                             step={10}
                             max={120}
                             min={10}
                             onChange={setMaxTime}
                             onBlur={()=>update()} key={'input_max_time'}/>
            </div>
            {question.filename ? showMusic():<AudioUpload setMusic={wrapSetMusic}/>}

            <Divider orientation={"left"}>
                <Space>
                    Réponses
                    <Tooltip title={"Ajouter réponse"}>
                        <PlusCircleTwoTone onClick={addAnswer} className={"button-action"}/>
                    </Tooltip>
                    <Tooltip title={"Supprimer réponse"}>
                        <MinusCircleTwoTone onClick={removeAnswer} className={"button-action"}/>
                    </Tooltip>
                </Space>
            </Divider>
            {question.answers != null ? question.answers.map(displayAnswer):''}
        </div>);
}
