import React from "react";
import { Column } from '@ant-design/plots';

export default function ChartResults({data, answers}){
    if(Object.keys(data).length === 0){
        return '';
    }
    const config = {
        data,
        xField: 'answer',
        yField: 'nb',
        seriesField:'answer',

        color: ['#FFC800', '#086375', '#F46036','#14591D'],
        legend:false,
        columnStyle:({answer})=>{return answers.includes(answer) ? {lineDash: [6, 6],r:10,stroke:'#1cbd31',lineWidth:5}:{}},
        label: {
            position: 'middle',
            // 'top', 'bottom', 'middle',
            style: {
                fill: 'white',
                opacity: 0.8,
                fontSize:30,
                fontWeight:700
            },
        },
        xAxis: {
            label: {
                autoHide: true,
                autoRotate: false,
                style:{
                    fontSize:20,
                }
            },
        },
    };
    return <Column {...config} />;
}