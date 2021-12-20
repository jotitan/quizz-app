import React, {useEffect, useState} from 'react';
import {Button, Col, Row, Slider, Upload} from "antd";
import {PauseOutlined, RightOutlined, ScissorOutlined} from "@ant-design/icons";

export default function AudioUpload({setMusic}){

    const [musicSize,setMusicSize] = useState(0);
    const [musicFile,setMusicFile] = useState();
    const [musicBuffer,setMusicBuffer] = useState();
    const [source,setSource] = useState();
    const [startTime,setStartTime] = useState();
    const [readingTime,setReadingTime] = useState(0);
    const [positionMusic,setPositionMusic] = useState(0);
    const [range,setRange] = useState([0,0])
    const [inter,setInter] = useState();
    const [running,setRunning] = useState(false);

    useEffect(()=>{
        // Update object
        setMusic({music:musicFile,from:range[0],to:range[1]})
        // eslint-disable-next-line
    },[range,musicFile])

    const readSize = file => {
        let fileReader = new FileReader();
        fileReader.readAsArrayBuffer(file);
        fileReader.onload = function (e) {
            let ac = new window.AudioContext();
            ac.decodeAudioData(e.target.result).then(buffer => {
                setMusicBuffer(buffer);
                setReadingTime(0)
                const size = Math.round(buffer.duration);
                setMusicSize(size);
                setRange([0,size]);
            })
        }
    }

    const uploadDone = e=>{
        if(e.file.percent === 100) {
            setMusicFile(e.file.originFileObj);
            readSize(e.file.originFileObj);
        }
    }

    const stop = ()=>{
        source.stop();
        setRunning(false);
        clearInterval(inter)
        setInter(i=>{
            clearInterval(i)
            return null;
        })
       setReadingTime(time => time + (new Date() - startTime))
    }

    const play = (forcePosition=null) =>{
        if(musicFile == null){
            return;
        }
        if(running){
            stop();
        }
        const ac = new window.AudioContext();
        const src = ac.createBufferSource();
        src.buffer = musicBuffer;
        src.loop = false;
        src.connect(ac.destination);
        src.start(0,forcePosition != null ?forcePosition:(readingTime/1000));
        setInter(setInterval(()=>{
            setPositionMusic(pos=>pos+1)
        },1000))
        setStartTime(new Date())
        setRunning(true);
        setSource(src);
    }

    const [isSliding,setisSliding] = useState(false);

    const updatePositionMusic = pos => {
        setPositionMusic(pos)
        setisSliding(true);
    }

    const changePosition = pos => {
        if(source == null && isSliding){
            return
        }
        setPositionMusic(pos)
        play(pos)
    }

    const stopRequest = ({ onSuccess }) => {
        setTimeout(() => onSuccess("ok"), 0);
    };

    return   <div style={{marginTop:10}}>
        <Upload customRequest={stopRequest} maxCount={1}
                itemRender={(v,file)=>
                    <span>{file.name} ({musicSize}")</span>
                }
                onChange={uploadDone}>
            <Button>🎵 Ajouter musique</Button>
        </Upload>
            <Row>
                <Col flex={"30px"} >
                    {running?<PauseOutlined style={{marginTop:6,fontSize:20}} onClick={stop}/>:<RightOutlined style={{marginTop:8}} onClick={()=>play()}/>}
                </Col>
                <Col flex={"auto"}>
                    <Slider defaultValue={0}
                            value={positionMusic}
                            max={musicSize}
                            onChange={updatePositionMusic}
                            onAfterChange={changePosition}
                            tipFormatter={v=>`${v}"`}
                    />
                </Col>
            </Row>
        <Row>
            <Col flex={"30px"}>
                <ScissorOutlined style={{marginTop:6,fontSize:20}} />
            </Col>
            <Col flex={"auto"}>
                <Slider defaultValue={[0,musicSize]} value={range} range  max={musicSize}  tipFormatter={v=>`${v}"`} onChange={setRange}/>
            </Col>
        </Row>

    </div>
}
