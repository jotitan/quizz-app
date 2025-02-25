import React from 'react';
import {Badge, Divider} from "antd";
import "./showplayers.css";

export default function ShowPlayers({users,start}){

    const buildImage = i => {
        return {backgroundImage:`url(/icons/icon_${i < 10 ? '0':''}${i}.svg)`,backgroundRepeat:'no-repeat',paddingLeft:'4vw',backgroundPositionX:5}
    }

    const showPlayer = p => <div key={`player_${p.player}_${p.position}`} className={"player"} style={buildImage(p.position)}>
        {p.player}
    </div>

    return <div>
        <Divider orientation={"left"}>
            Players - <Badge count={users.length} style={{ backgroundColor: '#1a5ec4' }}/>
        </Divider>
        <div className={"players"}>
            {users.map(showPlayer)}
        </div>
    </div>
}
