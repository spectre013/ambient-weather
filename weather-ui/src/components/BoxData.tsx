import { ReactNode } from "react";
import "./BoxData.css"
import * as CSS from 'csstype';

export interface Props {
    icon: string
    title: string
    unit: string
    style?: CSS.Properties
    children: ReactNode
}
const BoxData = (props:Props) => {

    return (
        <>
        <div className="box-container" style={props.style}>
            <div className="box-grid">
                <div className="icon">
                    <div className="header-icon">
                        <i className={`fa-solid ${props.icon}`}></i> &nbsp;
                        <span className="header-text">{props.title}</span>
                    </div>
                </div>
                <div className="scale">{ props.unit}</div>
                <div className="content">
                    {props.children}
                </div>
            </div>
        </div>
        </>
    )
}
export default BoxData
