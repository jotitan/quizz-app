import React, {useEffect, useState} from 'react';
import {Progress} from "antd";

export default function Countdown({duration,endAction=()=>{},style}){
    const [percent,setPercent] = useState(0);
    const [leftTime,setLeftTime] = useState(duration);
    const [iv,setIv] = useState(null);
    useEffect(()=>{
        setIv(setInterval(()=>{
            setLeftTime(left=>{
                left-=0.2
                setPercent(((duration-left)/duration)*100)
                return left
            })
        },200))
        // eslint-disable-next-line
    },[])

    useEffect(()=>{
        if(leftTime <=0){
            clearInterval(iv)
            endAction();
        }
        // eslint-disable-next-line
    },[leftTime,iv])

    return <Progress style={style}
            format={()=><span style={{fontWeight:leftTime<=5?'bold':'normal'}}>{Math.round(leftTime)}"</span>}
              strokeColor={{
                  from: '#FFC800',
                  to: '#F46036',
              }}
              strokeWidth={28}
              percent={percent}
              status="active"
    />
}
