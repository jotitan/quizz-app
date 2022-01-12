import React, {useState, useContext, useEffect} from "react";

import './EditQuizz.css';
import {Button, Input, notification, Upload} from "antd";
import QuizzContext from "../../context/QuizzContext";
import {DeleteOutlined} from "@ant-design/icons";
import getBase from "../../services/httpHelper";
import getBaseFront from "../../services/httpHelper";

export default function EditQuizz({quizz,runEdit = false}){
    const {createOrUpdateQuizz} = useContext(QuizzContext)
    const [title,setTitle] = useState(quizz.name);
    const [image,setImage] = useState(null);
    const [removeImage,setRemoveImage] = useState(false);
    const [description,setDescription] = useState(quizz.description);
    const {TextArea} = Input;

    const save = ()=>{
        createOrUpdateQuizz(title,description,quizz.id,image,removeImage)
            .then(d=>{
                notification['success']({message:`Quizz ${quizz.id == null ? 'created':'updated'}`})
                window.location.href=`${getBaseFront()}/quizz/${d.id}`;
            })
            .catch(()=>notification["error"]({message:"Impossible to create quizz"}))
    }

    useEffect(()=>{
        if(runEdit) {
            save();
        }
        // eslint-disable-next-line
    },[runEdit])

    const loadIllustration = e => {
        if(e.file.status === 'done'){
            setImage(e.file)
        }
    }

    const removeIllustration = ()=> {
        setImage(null);
        if(quizz.image != null){
            setRemoveImage(true);
        }
    }

    const stopRequest = ({ onSuccess }) => {
        setTimeout(() => onSuccess("ok"), 0);
    };

    return <div>
        <div style={{padding:20}}>
            {image == null && quizz.image ? <div style={{textAlign:'center'}}>
                <img alt={'illustration'} height={100} src={`${getBase()}/quizz/${quizz.id}/cover`}/>
            </div>:''}
            <div>
                Titre
            </div>
            <div className={"field"}>
                <Input value={title} onChange={e=>setTitle(e.target.value)} maxLength={26}/>
            </div>
            <div>Description</div>
            <div className={"field"}>
                <TextArea value={description} onChange={e=>setDescription(e.target.value)} rows={3}/>
            </div>
            <div>Illustration</div>
            <div className={"field"}>
                <Upload customRequest={stopRequest}
                        maxCount={1}
                        listType="picture"
                        itemRender = {() => {
                            if(image != null){
                                return <p>
                                    <img src={image.thumbUrl} height={100} alt={'preview'}/>
                                    {image.name}
                                    <DeleteOutlined onClick={removeIllustration}/>
                                </p>
                            }
                            return '';
                        }}
                        onRemove={removeIllustration}
                        onChange={loadIllustration}

                >
                    <Button>{quizz.image ? 'Changer':'Choisir'}</Button>

                </Upload>
            </div>
        </div>
    </div>
}
