import React from 'react';
import {PageHeader} from "antd";
import './Header.css';

const isHome = ()=> window.location.pathname === '/';

const showHeader = ()=> {
    return window.location.pathname.indexOf("game/host") === -1
    && window.location.pathname.indexOf("player") === -1
}

export default function Header({children}){

    return <div>
        {showHeader()?
            isHome() ? <PageHeader title={"Quizz app"} className="site-page-header"/>:
                    <PageHeader title={"Quizz app"} className="site-page-header" onBack={()=>window.location.href='/'}/>

            :''}
        {children}
    </div>
}
