import React, {useContext, useEffect, useState} from 'react';
import {SoundOutlined, SoundTwoTone} from "@ant-design/icons";
import {Tooltip} from "antd";
import GameContext from "../../context/GameContext";


export default function ReadMusic({gameId,secureId,question}){
    const {getMusic} = useContext(GameContext);
    const [source,setSource] = useState();
    const [buffer,setBuffer] = useState();
    const [playing,setPlaying] = useState(false);

    const play = ()=> {
        if(playing){
            source.stop();
            setPlaying(false);
            return
        }
        const ac = new AudioContext();
        const src = ac.createBufferSource();
        src.buffer = buffer;
        src.connect(ac.destination);
        src.start(0);
        setSource(src);
        setPlaying(true);
    }

    const load = ()=> {
        getMusic(gameId,secureId,question.id).then(music=> {
            new AudioContext().decodeAudioData(music).then(setBuffer);
        });
    }

    useEffect(()=>{
        load()
        // eslint-disable-next-line
    },[]);


    return <div>
        <Tooltip title={"Lecture"}>
            {playing?<SoundTwoTone  style={{fontSize:30,cursor:"pointer"}} onClick={play}/>:
                <SoundOutlined style={{fontSize:30,cursor:"pointer"}} onClick={play}/>}
        </Tooltip>
    </div>
}
