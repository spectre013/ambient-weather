import Vue from 'vue';
import computed from './mixins/temperature';
import VueNativeSock from 'vue-native-websocket';
import store from './store/index';
import VueRouter from 'vue-router';

import App from './App.vue';
import Home from './components/Home.vue';
import Forecast from './components/Forecast';

Vue.config.productionTip = false;
Vue.config.ignoredElements = [
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

Vue.use(VueRouter);
Vue.mixin(computed);
let wsurl = 'wss://' + window.location.host;
if (window.location.protocol === 'http:') {
  wsurl = 'ws://' + window.location.host;
}

wsurl += '/api/ws';

Vue.use(VueNativeSock, wsurl, {
  reconnection: true,
});

const routes = [
  { path: '/', component: Home },
  { path: '/forecast', component: Forecast },
];

export default routes;

const router = new VueRouter({ mode: 'history', routes });

new Vue({
  router,
  store,
  render: (h) => h(App),
}).$mount('#app');
