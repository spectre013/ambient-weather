import { ReactNode } from "react";
import "./BoxData.css"
import { decode } from 'html-entities';
import * as CSS from 'csstype';

export interface Props {
    icon: string
    title: string
    unit: string
    style?: CSS.Properties
    children: ReactNode
}
const BoxData = (props:Props) => {
   function setEntity(icon: string) {
       return decode(icon, {level: 'html5'});
   }

    return (
        <>
        <div className="box-container" style={props.style}>
            <div className="box-grid">
                <div className="icon">
                    <div className="header-icon">
                        { setEntity(props.icon) } &nbsp;
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
