import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import VueNativeSock from 'vue-native-websocket-vue3';

import './assets/main.css'

let wsurl = 'wss://' + window.location.host;
if (window.location.protocol === 'http:') {
    wsurl = 'ws://' + window.location.host;
}

wsurl += '/api/ws';
const app = createApp(App)
app.use(VueNativeSock, wsurl, {
    reconnection: true,
});

app.use(router)

app.config.productionTip = false;
app.config.ignoredElements = [
    /^darksky/,
    /^temp/,
    /^weather34/,
    /^span/,
    'oorange',
    'orange',
    'ogreen',
    'spanelightning',
    'ored',
    'heatindex',
    'spanmaxwind',
    'rainf',
    'spanwindtitle',
    'spanmaxwind',
    'oblue',
    'yesterdaytimemax',
    'spaneindoortemp',
    'suptemp',
    'trendmovementrising',
    'trendmovementfalling',
    'rainblue',
    'grey',
    'spancalm',
    'supmb',
    'max',
    'supmb',
    'windrun',
    'darkgrey',
    'convtext',
    'steady',
    'strongnumbers',
    'period',
    'min',
    'minutes',
    'hrs',
    'spanindoortempfalling',
    'indoororange1',
    'spanindoortempfalling',
    'indoororange1',
    'spanindoortempfalling',
    'luminance1',
    'uppercase',
    'span1',
    'uviforecasthourgreen',
    'value',
    'valuetext',
    'smalluvunit',
    'alertvalue',
    'noalert',
    'spanyellow',
    'valuetitleunit',
    'topblue1',
    'smallwindunit',
    'smallrainunita',
    'valuetextheading1',
    'raiblue',
    'smallrainunit2',
    'rainratetextheading',
    'maxred',
    'hours',
    'blueu',
    'moonrisecolor',
    'moonm',
    'moonsetcolor',
    'tgreen',
    'topgreen1',
    'smallwindunit',
    'toporange1',
    'minblue',
    'smalltempunit',
    'smalltempunit2',
    'smallrainunit',
    'valuetext1',
    'uviforecasthouryellow',
    'trendmovementfallingx',
    'redu',
    'blue',
    'articlegraph',
    'uviforecasthourred',
    'uviforecasthourorange',
    'stationid',
    'yellow',
    'green',
    'darkskytempwindhome',
    'darkskyiconcurrent',
    'darkskyrainhome1',
    'unit',
    'blue1',
    'windunit',
    'gust',
    'purpleu',
    'orange1',
    'uv',
    'topbarmetricc',
    'topbarimperialf',
    'topbarimperial',
    'topbarmetric',
    'strikeicon',
    'laststrike',
];

app.mount('#app')
