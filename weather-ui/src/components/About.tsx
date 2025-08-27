import './About.css'
import {useNavigate} from 'react-router-dom';
import {useEffect, useState} from "react";
import {AboutModel} from "../models/current.ts";
import {tempLabel, windLabel} from "../util/weather.ts";
import moment from "moment";

const About = () => {
    const navigate = useNavigate();
    const [fLoaded, setFLoaded] = useState(false);
    const [about, setAbout] = useState<AboutModel>({} as AboutModel);
    const urls = [
        '/api/about',
    ];

    useEffect(() => {
        const fetchPromises = urls.map(url => fetch(url));
        Promise.all(fetchPromises)
            .then(responses => {
                // 'responses' will be an array of Response objects
                // Process each response to extract JSON data
                return Promise.all(responses.map(response => response.json()));
            })
            .then(data => {
                const aboutData: AboutModel = data[0] as unknown as AboutModel;
                setAbout(aboutData);
                setFLoaded(true);
            })
            .catch(error => {
                // Handle any errors that occurred during fetching or JSON parsing
                console.error('Error fetching data:', error);
            });
    },[]);

    if (!fLoaded) {
        return 'loading';
    }

    const handleClick = () => {
        navigate('/'); // Navigate to the /dashboard route
    };

    return (
        <>
            <div className="about-dashboard">
                <div className="content">
                    <main className="main-content">
                        <header className="details-header">
                            <h1><span className="material-symbols-sharp">weather_hail</span> Losron Ranch -  About</h1>
                            <div className="hasclick" onClick={handleClick}><span className="material-symbols-sharp">home</span></div>
                        </header>
                        <div>
                            <p>Lorson Ranch weather is run off an <a
                                href="https://ambientweather.com/ws-2902-smart-weather-station">Ambient Weather WS-2902</a> located about 10 meters off the ground.</p>
                            <p>The first began around 2018 and has been through several iterations, the first was VueJS with a Golang back end. In 2020 the site was switched
                            to React and Golang. The third iteration was a complete UI rework and done with Golang nad HTMX. The current version has moved back to React with a
                                Golang back end.</p>
                            <p>

                                Currently there is { about.records.toLocaleString('en-US')} with the first record being recorded on 2020-07-14. A few records that have
                                been recorded are a high temperature of { about.maxtemp}&deg;{tempLabel('imperial')}.
                                A low temperature of { about.mintemp}&deg;{tempLabel('imperial')} and finally a max wind gust of { about.maxgust } {windLabel('imperial')}
                            </p>

                            <div className="list-container">
                                <h2 className="list-header">Station Sensors</h2>
                                <div className="list-wrapper">
                                    <ul className="sensor-list">
                                        <li className="sensor-item">Dew Point</li>
                                        <li className="sensor-item">Forecast (Pressure Based)</li>
                                        <li className="sensor-item">Heat Index</li>
                                        <li className="sensor-item">Indoor Humidity</li>
                                        <li className="sensor-item">Indoor Temp</li>
                                        <li className="sensor-item">Lightning Detection</li>
                                        <li className="sensor-item">Moon Phase</li>
                                        <li className="sensor-item">Moonrise & Moonset</li>
                                        <li className="sensor-item">Outdoor Humidity</li>
                                        <li className="sensor-item">Outdoor Temperature</li>
                                    </ul>
                                    <ul className="sensor-list">
                                        <li className="sensor-item">Rainfall</li>
                                        <li className="sensor-item">Relative Pressure</li>
                                        <li className="sensor-item">Solar Radiation</li>
                                        <li className="sensor-item">Sunrise & Sunset</li>
                                        <li className="sensor-item">Ultrasonic Wind</li>
                                        <li className="sensor-item">UV</li>
                                        <li className="sensor-item">Wind Direction</li>
                                        <li className="sensor-item">Wind Speed</li>
                                        <li className="sensor-item">Wind Chill</li>
                                        <li className="sensor-item">Air Quality</li>
                                    </ul>
                                </div>
                            </div>
                            <p>If you would like to host this site your self you can get the code from <a href="https://github.com/spectre013/ambient-weather">ambient-weather</a> repository
                            on Github</p>
                        </div>
                    </main>
                </div>
            </div>
            <div className='copyright'>&copy; 2018-{moment().format('YYYY')} : zoms.net</div>
        </>
    )
}
export default About