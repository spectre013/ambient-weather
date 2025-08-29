import {useNavigate} from "react-router-dom";
import {WeatherContext} from "../Context.ts";
import {useContext} from "react";

const Header = () => {
    const navigate = useNavigate();
    const ctx = useContext(WeatherContext);

    const handleClick = () => {
        navigate('/'); // Navigate to the details page for the specific stat
    }

    return (
        <header className="details-header">
            <h1><span className="material-symbols-sharp">thermostat</span> {ctx.shortname} -  Temperature</h1>
            <div className="hasclick" onClick={handleClick}><span className="material-symbols-sharp">home</span></div>
        </header>
    );
}

export default Header;