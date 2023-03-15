import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from './store/index';
import moment from "moment/moment";

const app = createApp(App)
app.use(router);
app.use(store);

app.config.productionTip = false;

app.directive("sundial",{
    mounted(canvas,data) {
        let d_crcl = (24 * 60) / 2;
        let now = moment().unix();
        let sunRise = moment().startOf('day').add(timeToDecimal(data.value.sunrise), 'hours').unix();
        let sunSet = moment().startOf('day').add(timeToDecimal(data.value.sunset), 'hours').unix();

        function timeToDecimal(t) {
            let arr = t.split(':');
            let dec = parseInt((arr[1] / 6) * 10, 10);
            return parseFloat(parseInt(arr[0], 10) + '.' + (dec < 10 ? '0' : '') + dec);
        }
        function clc_crcl(ts) {
            let h = parseInt(moment.unix(ts).format('H'));
            let m = parseInt(moment.unix(ts).format('m'));
            let calc = m + h * 60;
            calc = 0.5 + calc / d_crcl;
            if (calc > 2.0) {
                calc = calc - 2;
            }
            return calc.toFixed(5);
        }
        let start = parseFloat(clc_crcl(sunRise));
        let end = parseFloat(clc_crcl(sunSet));
        let pos = parseFloat(clc_crcl(moment().unix()));
        let sn_clr = '';
        if (now > sunSet || now < sunRise) {
            sn_clr = 'rgba(86,95,103,0)';
        } else {
            sn_clr = 'rgba(255, 112,50,1)';
        }

        let ctx = canvas.getContext('2d');
        ctx.imageSmoothingEnabled = false;
        ctx.beginPath();
        ctx.arc(63, 65, 55, 0, 2 * Math.PI);
        ctx.lineWidth = 0;
        ctx.strokeStyle = '#565f67';
        ctx.stroke();
        ctx.beginPath();
        ctx.arc(63, 65, 55, start * Math.PI, end * Math.PI);
        ctx.lineWidth = 2;
        ctx.strokeStyle = '#3b9cac';
        ctx.stroke();
        ctx.beginPath();
        ctx.arc(63, 65, 55, pos * Math.PI, pos * Math.PI);
        ctx.lineWidth = 0;
        ctx.strokeStyle = `"${sn_clr}"`;
        ctx.stroke();
    }
})


app.mount('#app')
