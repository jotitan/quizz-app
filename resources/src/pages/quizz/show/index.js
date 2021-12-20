import React, {useContext, useEffect, useState} from "react";
import {useParams} from "react-router-dom";
import QuizzContext from "../../../context/QuizzContext";
import {Divider, Modal, notification, Popconfirm, Space, Tooltip} from 'antd';
import {DeleteTwoTone, EditTwoTone, PlusCircleTwoTone} from '@ant-design/icons';
import 'antd/dist/antd.css';
import './show.css';
import '../../../App.css';
import EditQuestion from "../../../components/EditQuestion";
import RoundButton from "../../../components/RoundButton";


export default function ShowQuizz(){
    let id = useParams().id
    let [quizz,setQuizz] = useState({});
    let [refresh,setRefresh] = useState(false);
    let [showEditQuestion,setShowEditQuestion] = useState(false);
    const [editingQuestion,setEditingQuestion] = useState({});
    const {getQuizz,addQuestion,deleteQuestion} = useContext(QuizzContext)
    useEffect(()=>getQuizz(id).then(setQuizz),[getQuizz,setQuizz,refresh,id])

    const editQuestion = q => {
        setEditingQuestion(q);
        setShowEditQuestion(true);
    }

    const deleteQ = q => {
        deleteQuestion(id,q.id)
            .then(()=>setRefresh(r=>!r))
            .catch(()=>notification["error"]({message:"Impossible to delete"}))
    }

    const saveQuestion = ()=>{
        if(editingQuestion.id == null && editingQuestion.title == null){
            return;
        }
        setShowEditQuestion(false)
        addQuestion(id,editingQuestion)
            .then(()=>{
                notification["success"]({message:'Question added'})
                setRefresh(r=>!r)
            })
            .catch((e)=>notification["error"]({message:'Impossible to add question ' + e}))
    }

    const addNewQuestion = ()=>{
        addQuestion(id,{title:'',answers:[]}).then(()=>setRefresh(r=>!r))
    }

    const formatNbResponse = q=> q.answers != null ?`${q.answers.length} réponse(s)`:'0 réponse';

    const formatTime = q=> `⏰ ${q.time == null ?30:q.time}"`

    const showMusic = q => q.filename!==''?<span title={'Musique'}>🎵</span>:'';

    const showQuestion = (q,pos) => {
        return <div key={`q_${q.id}`}>
            <Space>
                <span style={{fontWeight:'bold'}}>{pos+1}</span> - {q.title} : {formatNbResponse(q)} {showMusic(q)} {formatTime(q)}
                <Tooltip title={"Modifier"}>
                    <EditTwoTone onClick={()=>editQuestion(q)} className={"button-action"}/>
                </Tooltip>
                <Tooltip title={"Supprimer"}>
                    <Popconfirm title={"Souhaitez vous supprimer cette question"} okText={"Supprimer"} cancelText={"Annuler"} onConfirm={()=>deleteQ(q)}>
                        <DeleteTwoTone style={{marginLeft:20}} className={"button-action"}/>
                    </Popconfirm>
                </Tooltip>
            </Space>
        </div>
    }
    const isQuizzPlayable = () => quizz.questions != null && quizz.questions.length > 0
            && quizz.questions.every(q=>q.answers != null && q.answers.length > 1)

    const showPlayButton = () => {
        return isQuizzPlayable() ?
            <RoundButton title={"Jouer"} action={()=>window.location.href=`/game/host/${quizz.id}`}/>
        :''
    }

    return (<div>
        <div style={{paddingLeft:50}}>
            {quizz.id != null ? <div>
                <h1>Quizz {quizz.name}</h1>
                {quizz.description}

                <Divider orientation={"left"}>
                    <Space>
                        Questions
                    <Tooltip title={"Ajouter question"}>
                        <PlusCircleTwoTone onClick={addNewQuestion} className={"button-action"}/>
                    </Tooltip>
                    </Space>
                </Divider>
                {quizz.questions != null ? quizz.questions.map(showQuestion):''}

                {showPlayButton()}
            </div>:<div>Quizz inconnu</div>}

        </div>
        <Modal title={"Edit question"} visible={showEditQuestion}
               className={"modal-question"}
               onOk={saveQuestion} onCancel={()=>setShowEditQuestion(false)}
               okText={"Sauvegarder"} cancelText={"Annuler"}>
            {editingQuestion == null || editingQuestion.id == null ? '':
                < EditQuestion question = {editingQuestion} updateQuestion={setEditingQuestion}/>
            }
        </Modal>
    </div>)
}
