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
        "footer": {
            template: document.getElementById('footer-template').innerHTML,
            target: "footer",
            data: {year: moment().format('YYYY')}
        },
        "updated": {
            template: document.getElementById('updated-template').innerHTML,
            target: "updated",
            data: {updated: moment(data.date).format('MMM Do YYYY hh:mm:ss a'), connected: isConnected},
        },
        "temperature": {
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
        }
    };

}

function buildTemplates() {
    for (let key in templateList) {
        Mustache.parse(templateList[key].template);
    }
}

function render(template,target, data) {
    //console.log(template, target, data);
    document.getElementById(target).innerHTML = Mustache.render(template, data);
}

function renderTemplates() {
    for (let key in templateList) {
        const temp = templateList[key];
        render(temp.template,temp.target, temp.data);
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

function handleClick() {
    console.log("clicked");
}
buildTemplates();
wsConnect();


