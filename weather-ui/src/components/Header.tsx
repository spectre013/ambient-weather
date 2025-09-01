import {useNavigate} from "react-router-dom";
import {WeatherContext} from "../Context.ts";
import {useContext} from "react";

interface HeaderProps {
    name: string
    icon: string
}

const Header = (props: HeaderProps) => {
    const navigate = useNavigate();
    const ctx = useContext(WeatherContext);

    const handleClick = () => {
        navigate('/'); // Navigate to the details page for the specific stat
    }

    return (
        <header className="details-header">
            <h1><span className="material-symbols-sharp">{props.icon}</span> {ctx.shortname} -  {props.name}</h1>
            <div className="hasclick" onClick={handleClick}><span className="material-symbols-sharp">home</span></div>
        </header>
    );
}

export default Header;