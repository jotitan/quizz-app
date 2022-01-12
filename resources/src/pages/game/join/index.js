import React, {useContext, useState} from 'react';
import GameContext from "../../../context/GameContext";
import {Input, notification} from "antd";
import RoundButton from "../../../components/RoundButton";
import './join.css';
import {useSearchParams} from "react-router-dom";
import getBaseFront from "../../../services/httpHelper";

export default function JoinGame(){
    const [query] = useSearchParams();
    const [name,setName] = useState('');
    const [game,setGame] = useState(query.get("game") || '');
    const {joinGame} = useContext(GameContext)

    const connect = ()=>joinGame(game,name)
        .then(d=>window.location.href=`${getBaseFront()}/player/play/${game}/${d.id}?name=${name}`)
        .catch(e=>notification["error"]({message:'Impossible to connect',description:e}))

    return (
        <div className={"join-block"}>
            <div>
                <Input placeholder={"Pseudo"} value={name} size={"large"} onChange={e=>setName(e.target.value)}/>
            </div>
            <div>
                <Input placeholder={"Code partie"} value={game} size={"large"} onChange={e=>setGame(e.target.value)}/>
            </div>
            <RoundButton title={"Rejoindre"} action={connect}/>
        </div>)
}
