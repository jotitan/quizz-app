import React from 'react';
import '../../App.css';

export default function RoundButton({title,action}){

    return <div style={{textAlign:'center',position:'absolute',bottom:30,width:'100%'}}>
        <div style={{display:'inline-block'}}>
            <div className={"circle-button"}><span onClick={action}>{title}</span></div>
        </div>
    </div>
}
