import React from 'react';
import {Badge, Divider} from "antd";
import RoundButton from "../RoundButton";

export default function ShowPlayers({users,start}){

    return <div>
        <Divider orientation={"left"}>
            Players - <Badge count={users.length} style={{ backgroundColor: '#1a5ec4' }}/>
        </Divider>
        <div>
            {users.map(u=><span key={`player_${u}`} className={"player"}>{u}</span>)}
        </div>
        <RoundButton title={"Démarrer"} action={start}/>
    </div>
}
