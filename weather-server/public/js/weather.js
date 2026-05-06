
let data = {}
let isConnected = false;
let units="imperial";
let forecast = {};
let loading = true;
fetchForecast().then(forecastData => {
    forecast = forecastData;
}).catch(error => {
    console.error("Error fetching forecast data:", error);
});
let templateList = {}

function setTemplates(data) {
    templateList = {
        footer: {
            template: document.getElementById('footer-template').innerHTML,
            target: "footer",
            data: {year: moment().format('YYYY')}
        },
        updated: {
            template: document.getElementById('updated-template').innerHTML,
            target: "updated",
            data: {updated: moment(data.date).format('MMM Do YYYY hh:mm:ss a'), connected: isConnected},
        },
        temperature: {
            template: document.getElementById('temperature-template').innerHTML,
            target: "temperature",
            data: {
                temperature: data.temp.temp,
                tempClass: tempColor(data.temp.temp),
                tempDisplay: tempDisplay(data.temp.temp, units),
                tempLabel: tempLabel(units),
                conditions: forecast[0].conditions,
                icon:forecast[0].icon
            }
        },
        forecast: {
            template: document.getElementById('forecast-template').innerHTML,
            target: "forecast",
        },
        feelslike: {
            template: document.getElementById('stat-template').innerHTML,
            target: "stats-cards",
            data: {
                name: "Feels Like",
                icon: "thermostat",
                className: tempColor(data.temp.feelslike),
                value: data.temp.feelslike,
                unit: tempLabel(units),
            }
        },
        windspeed: {
            template: document.getElementById('stat-template').innerHTML,
            target: "stats-cards",
            data: {
                name: "Wind Speed",
                icon: "air",
                className: beaufortClass(data.wind.windspeedmph),
                template: document.getElementById('stat-template').innerHTML,
                target: "stats-cards",
                value: data.wind.windspeedmph,
                unit: windLabel(units)
            }
        },
        windsgust: {
            template: document.getElementById('stat-template').innerHTML,
            target: "stats-cards",
            data: {
                name: "Wind Gust",
                icon: "air",
                className: beaufortClass(data.wind.windgustmph),
                value: data.wind.windgustmph,
                unit: windLabel(units),
            }
        },
        humidity: {
            template: document.getElementById('stat-template').innerHTML,
            target: "stats-cards",
            data: {
                name: "Humidity",
                icon: "humidity_mid",
                className: "",
                value: data.humidity.humidity,
                unit: "%",
            }
        },
        barometer: {
            template: document.getElementById('stat-template').innerHTML,
            target: "stats-cards",
            data: {
                name: "Barometer",
                icon: "water_drop",
                className: "",
                value: data.barometer.baromrelin,
                unit: baroLabel(units),
            }
        },
        dewpoint: {
            template: document.getElementById('stat-template').innerHTML,
            target: "stats-cards",
            data: {
                name: "Dewpoint",
                icon: "thermostat",
                className: tempColor(data.humidity.dewpoint),
                value: data.humidity.dewpoint,
                unit: tempLabel(units),
            }
        },
        precipitation: {
            template: document.getElementById('stat-template').innerHTML,
            target: "stats-cards",
            data: {
                name: "Precipitation",
                icon: "water_drop",
                className: "",
                value: data.rain.daily,
                unit: rainLabel(units),
            }
        },
        lightningcount: {
            template: document.getElementById('stat-template').innerHTML,
            target: "stats-cards",
            data: {
                name: "Lightning Today",
                icon: "bolt",
                className: "",
                value: data.lightning.day,
                unit: "",
            }
        },
        lightningdist: {
            template: document.getElementById('stat-template').innerHTML,
            target: "stats-cards",
            data: {
                name: "Lightning Distance",
                icon: "bolt",
                className: "",
                value: data.lightning.distance,
                unit: distanceLabel(units),
            }
        },
        internal: {
            template: document.getElementById('stat-template').innerHTML,
            target: "internal-cards",
        },
    };

}

function buildTemplates() {
    for (let key in templateList) {
        Mustache.parse(templateList[key].template);
    }
}

function render(template,target, data) {
    //console.log(template, target, data);
    if (target === "forecast" || target === "stats-cards") {
        document.getElementById(target).insertAdjacentHTML('beforeend',Mustache.render(template, data));
    } else {
        document.getElementById(target).innerHTML = Mustache.render(template, data);
    }
}

function renderTemplates() {
    document.getElementById("stats-cards").innerHTML = "";
    for (let key in templateList) {
        const temp = templateList[key];
        if(key === "forecast") {
            document.getElementById(temp.target).innerHTML = "";
            for (let i = 0; i < forecast.length; i++) {
                forecastData = {
                    day: "day" + i,
                    date: moment(forecast[i].datetime).format('ddd'),
                    tempmax: forecast[i].tempmax,
                    tempmin: forecast[i].tempmin,
                    conditions: forecast[i].conditions,
                    icon: forecast[i].icon,
                    description: forecast[i].description,
                    summary: forecast[i].summary,
                    gradient: getTemperatureGradient(forecast[i].tempmin, forecast[i].tempmax),
                }
                render(temp.template, temp.target, forecastData);
            }
        } else if (key === "internal") {
            document.getElementById(temp.target).innerHTML = "";
            for (let i = 0; i < 5; i++) {
                let node = "temp"+i;
                if(i === 0) {
                    node = "tempin";
                }
                let internaldata = {
                    name: "Feels Like",
                        icon: "thermostat",
                        className: tempColor(data[node].temp),
                        value: data[node].temp,
                        unit: tempLabel(units),
                }
                render(temp.template, temp.target, internaldata);
            }
        } else {
            render(temp.template, temp.target, temp.data);
        }
    }
}

function wsConnect() {
    const socket = new WebSocket('ws://localhost:8000/api/ws');

// 2. Handle the "open" event to send data once connected
    socket.onopen = (event) => {
        console.log("Connected to the server!");
        isConnected = true;
    };

// 3. Listen for incoming messages from the server
    socket.onmessage = (event) => {
        data = JSON.parse(event.data);
        setTemplates(data);
        renderTemplates();
        loading = false;
        isLoaded()

    };

// 4. Handle errors and connection closures
    socket.onerror = (error) => console.error("WebSocket Error:", error);
    socket.onclose = () => console.log("Connection closed.");
}

function isLoaded() {
    if(loading === false) {
        console.log("Loaded");
        document.getElementById("loading").style.display = "none";
        document.getElementById("root").style.display = "block";
    }
}

async function fetchForecast() {
    try {
        const response = await fetch('/api/forecast');

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        return  await response.json();;
    } catch (error) {
        console.error("Could not fetch forecast:", error);
        throw error;
    }
}

buildTemplates();
wsConnect();

function openModal(modalId) {
    document.getElementById(modalId).style.display = "block";
}

function closeModal(modalId) {
    document.getElementById(modalId).style.display = "none";
}

// Close modal if user clicks outside of it
window.onclick = function(event) {
    if (event.target.classList.contains('modal')) {
        event.target.style.display = "none";
    }
}

