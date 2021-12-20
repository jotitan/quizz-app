import React from "react";
import {Card, Row} from "antd";
import {PlusOutlined,BarChartOutlined} from "@ant-design/icons";

export default function Home(){
    return <div>

        <Row>
            <Card title={<><BarChartOutlined/>Quizz</>} extra={<a href="/quizzes">Accéder</a>} style={{ width: 300,marginLeft:100 }}>
                <p>Lister tous les quizz du jeu</p>
            </Card>

            <Card title={<><PlusOutlined /> Créer</>} extra={<a href="/quizz/create">Accéder</a>} style={{ width: 300,marginLeft:100 }}>
                <p>Créer un nouveau quizz</p>
            </Card>
        </Row>
    </div>
}
