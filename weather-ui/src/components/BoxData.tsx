import { ReactNode } from "react";
import { useNavigate } from "react-router-dom";
import "./BoxData.css"
import { decode } from 'html-entities';
import * as CSS from 'csstype';

export interface Props {
    icon: string
    title: string
    unit: string
    style?: CSS.Properties
    navigate: string
    children: ReactNode
}
const BoxData = (props:Props) => {
    const navigate = useNavigate();
    const routeChange = () =>{
        const path = "/almanac/" + props.navigate;
        navigate(path);
    }
   function setEntity(icon: string) {
       return decode(icon, {level: 'html5'});
   }

    return (
        <>
        <div className="box-container" style={props.style} onClick={() => routeChange()}>
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
