
export interface Props {
    children: JSX.Element
}
const BoxAlmanac = (props:Props) => {

    return (
        <div>
            {props.children}
        </div>
    )
}
export default BoxAlmanac
