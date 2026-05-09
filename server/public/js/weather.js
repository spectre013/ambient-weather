(function () {
    
    const state = {
        data: {},
        forecast: {},
        units: "imperial",
        isConnected: false,
        loading: true,
        templates: {} // Cached templates go here
    };

    function initTemplates() {
        state.templates = {
            footer: document.getElementById('footer-template').innerHTML,
            updated: document.getElementById('updated-template').innerHTML,
            temperature: document.getElementById('temperature-template').innerHTML,
            dailyForecast: document.getElementById('daily-template').innerHTML,
            hourlyForecast: document.getElementById('hourly-template').innerHTML,
            details: document.getElementById('details-template').innerHTML
        };

        for (let key in state.templates) {
            Mustache.parse(state.templates[key]);
        }
    }

    // 3. Setup Global Event Listeners (Event Delegation)
    function setupEventListeners() {
        // Listen for clicks on the document body to avoid duplicate listeners
        document.body.addEventListener('click', (event) => {
            // Handle Rain Modal Open
            const rainBtn = event.target.closest('#rain');
            if (rainBtn) {
                const modalId = rainBtn.parentNode.id + "Modal";
                const modal = document.getElementById(modalId);
                if (modal) modal.style.display = "block";
            }

            // Handle Modal Close (clicking outside)
            if (event.target.classList.contains('modal')) {
                event.target.style.display = "none";
            }
        });
    }

    // 4. Helper Methods
    function getConditions(conditions) {
        let tag = conditions;
        if (tag.includes(",")) {
            tag = tag.split(",")[0].trim();
        }
        if (tag.includes(" ")) {
            let words = tag.split(" ");
            words[0] = words[0][0] + ".";
            // Added a space to join so "Partly Cloudy" becomes "P. CLOUDY"
            tag = words.join(" ").toUpperCase();
        }
        return tag;
    }

    function renderDashboard(liveData) {
        state.data = liveData;

        // Render Single Elements
        document.getElementById('footer').innerHTML = Mustache.render(state.templates.footer, { year: moment().format('YYYY') });

        document.getElementById('updated').innerHTML = Mustache.render(state.templates.updated, {
            updated: moment(liveData.date).format('LLLL'),
            connected: state.isConnected
        });

        document.getElementById('temperature').innerHTML = Mustache.render(state.templates.temperature, {
            temp: liveData.temp.temp.toFixed(0),
            tempLabel: tempLabel(state.units),
            conditions: state.forecast[0].conditions,
            feelslike: liveData.temp.feelslike.toFixed(0),
            min: liveData.temp.minmax.min.day.value.toFixed(0),
            max: liveData.temp.minmax.max.day.value.toFixed(0),
        });

        let start = (moment().hour() > 12) ? 1 : 0;
        let end = start + 5;
        let dailyHtml = "";
        for (let i = start; i < end; i++) {
            dailyHtml += Mustache.render(state.templates.dailyForecast, {
                date: moment(state.forecast[i].datetime).tz("America/Denver").format('ddd'),
                tempmax: state.forecast[i].tempmax.toFixed(0),
                tempmin: state.forecast[i].tempmin.toFixed(0),
                conditions: getConditions(state.forecast[i].conditions),
                tempLabel: tempLabel(state.units),
            });
        }
        document.getElementById('daily-forecast').innerHTML = dailyHtml; // One DOM update

        let hours = JSON.parse(state.forecast[0].hours);
        let hourlyHtml = "";
        for (let i = 0; i < hours.length; i++) {
            hourlyHtml += Mustache.render(state.templates.hourlyForecast, {
                date: moment.unix(hours[i].datetimeEpoch).tz("America/Denver").format('h a'),
                id: moment.unix(hours[i].datetimeEpoch).tz("America/Denver").format('H'),
                temp: hours[i].temp.toFixed(0),
                tempLabel: tempLabel(state.units),
            });
        }
        document.getElementById('hourly-forecast').innerHTML = hourlyHtml; // One DOM update

        renderDetails('details-a', [
            { name: "RAIN", value: liveData.rain.daily, label: " " + rainLabel(state.units) },
            { name: "WIND SPEED", value: liveData.wind.windspeedmph, label: " " + windLabel(state.units) },
            { name: "WIND GUST", value: liveData.wind.windgustmph, label: " " + windLabel(state.units) }
        ]);

        renderDetails('details-b', [
            { name: "BAROMETER", value: liveData.barometer.baromrelin, label: " " + baroLabel(state.units) },
            { name: "HUMIDITY", value: liveData.temp.humidity, label: "%" },
            { name: "DEWPOINT", value: liveData.temp.dewpoint.toFixed(0), label: tempLabel(state.units) }
        ]);

        renderDetails('details-c', [
            { name: "LIGHTNING", value: liveData.lightning.day, label: "" },
            { name: "LIGHTNING DISTANCE", value: liveData.lightning.distance, label: " " + distanceLabel(state.units) },
            { name: "UV", value: liveData.uv.uv, label: "" },
            { name: "SOLAR RADIATION", value: liveData.uv.uv, label: "" }
        ]);

        renderDetails('details-d', [
            { name: "SUNRISE", value: moment(liveData.astro.sunrise).local().format('LT'), label: "" },
            { name: "SUNSET", value: moment(liveData.astro.sunset).local().format('LT'), label: "" }, // FIXED: Removed distanceLabel
            { name: "MOONPHASE", value: "WAXING", label: "" }
        ]);

        renderDetails('details-e', [
            { name: "LIVINGROOM", value: liveData.tempin.temp.toFixed(0), label: tempLabel(state.units) },
            { name: "MASTER", value: (liveData.temp1 ? liveData.temp1.temp : liveData.temp3.temp).toFixed(0), label: tempLabel(state.units) }, // FIXED: Using temp1 to avoid duplicating temp3
            { name: "HANNAH", value: liveData.temp2.temp.toFixed(0), label: tempLabel(state.units) },
            { name: "BASEMENT", value: liveData.temp3.temp.toFixed(0), label: tempLabel(state.units) },
            { name: "GARAGE", value: liveData.temp4.temp.toFixed(0), label: tempLabel(state.units) }
        ]);

        handleLoaded();
    }

    function renderDetails(targetId, dataArray) {
        let html = "";
        for (let i = 0; i < dataArray.length; i++) {
            dataArray[i].id = dataArray[i].name.replace(" ", "").toLowerCase();
            html += Mustache.render(state.templates.details, dataArray[i]);
        }
        document.getElementById(targetId).innerHTML = html;
    }

    function wsConnect() {
        const socket = new WebSocket('ws://localhost:8000/api/ws');

        socket.onopen = () => {
            console.log("Connected to the server!");
            state.isConnected = true;
        };

        socket.onmessage = (event) => {
            const liveData = JSON.parse(event.data);
            renderDashboard(liveData);
        };

        socket.onerror = (error) => console.error("WebSocket Error:", error);

        socket.onclose = () => {
            console.log("Connection closed.");
            state.isConnected = false;
        };
    }

    function handleLoaded() {
        if (state.loading === true) {
            state.loading = false;
            console.log("Loaded");
            document.getElementById("loading").style.display = "none";
            document.getElementById("dashboard").style.display = "block";

            let hour = moment().tz("America/Denver").hour();
            if (hour < 20) {
                hour = hour + 4;
            }
            const targetItem = document.getElementById('hour' + hour);
            if (targetItem) targetItem.scrollIntoView(false);
        }
    }

    async function fetchForecast() {
        try {
            const response = await fetch('/api/forecast');
            if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
            return await response.json();
        } catch (error) {
            console.error("Could not fetch forecast:", error);
            throw error;
        }
    }

    async function init() {
        initTemplates();
        setupEventListeners();

        try {
            // Await the HTTP request BEFORE opening the WebSocket
            state.forecast = await fetchForecast();
            wsConnect();
        } catch (error) {
            console.error("Initialization failed:", error);
        }
    }

    // Start App
    init();

})();