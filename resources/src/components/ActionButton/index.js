import React from 'react';
import '../../App.css';

export default function ActionButton({title,action, icon}){

    return <div style={{textAlign:'center',position:'absolute',bottom:30,width:'calc(100% - 50px)'}}>
        <div style={{display:'inline-block'}}>
            <div className={"circle-button"}><span onClick={action}>{icon!=null ? icon:''} {title}</span></div>
        </div>
    </div>
}
