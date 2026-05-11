let data = {}
let templateList = {
    "footer": {
        template: document.getElementById('footer-template').innerHTML,
        target: "footer",
        data: {year: moment().format('YYYY')}
    },
};

function buildTemplates() {
    for (let key in templateList) {
        Mustache.parse(templateList[key].template);
    }
}

function render(template,target, data) {
    console.log(template, target, data);
    document.getElementById(target).innerHTML = Mustache.render(template, data);
}

function renderTemplates() {
    for (let key in templateList) {
        const temp = templateList[key];
        render(temp.template,temp.target, temp.data);
    }
}

function wsConnect() {
    const socket = new WebSocket('ws://localhost:8100/api/ws');

// 2. Handle the "open" event to send data once connected
    socket.onopen = (event) => {
        console.log("Connected to the server!");
    };

// 3. Listen for incoming messages from the server
    socket.onmessage = (event) => {
        data = JSON.parse(event.data);
        renderTemplates();
    };

// 4. Handle errors and connection closures
    socket.onerror = (error) => console.error("WebSocket Error:", error);
    socket.onclose = () => console.log("Connection closed.");
}

buildTemplates();
wsConnect();


