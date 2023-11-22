
export interface Props {
    icon: string,
    maintitle: string
    subtitle: string
    height: number
    width: number
    class: string
    max: string
    children: JSX.Element
}
const Box = (props:Props) => {

    function Max() {
        if(props.max == "") {
            return (
                <>
                </>
            )
        } else {
            return (
            <>
                | <span className={props.max}> MAX </span>
            </>
            )
        }
    }

    return (
        <div className="weatherbox">
            <div className="title">
                <img src={props.icon} height={props.height} width={props.width} alt="icon"/>
                {props.maintitle}  <span className={props.class}> {props.subtitle} </span> <Max/>
            </div>
            <div className="value">
                {props.children}
            </div>
        </div>
    )
}
export default Box
