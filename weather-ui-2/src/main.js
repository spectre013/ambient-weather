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

app.config.globalProperties.$filters = {
    currencyUSD(value) {
        return '$' + value
    }
}

app.use(router)

app.mount('#app')
