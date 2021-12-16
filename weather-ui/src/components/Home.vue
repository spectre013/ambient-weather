<template>
    <div id="app" v-if="$store.getters.theme">
        <Menu />
        <div class="weather2-container">
            <Time/>
            <MinMax :temp="temp" />
            <Rainfall :current="current"/>
            <Alert :forecast="forecast"/>
        </div>
        <div class="weather-container">
            <Temperature :current="live" :temp="temp" v-on:openModal="showModal"/>
            <Wind :current="live" :wind="wind" />
            <RainfallDetails :current="current"/>
        </div>
        <div class="weather-container">
            <Barometer :current="live" />
            <Moon :astro="astro" />
            <Daylight :astro="astro"/>
        </div>
        <div class="weather-container">
            <Uv :current="live" :astro="astro"/>
            <Lightining :current="live" :stats="lightning"/>
            <Indoor :live="live" :loc="'in'" :title="'Indoor'"/>
        </div>
        <div class="weather-container">
            <Indoor :live="live" :loc="'2'" :title="'Master'"/>
            <Indoor :live="live" :loc="'3'" :title="'Office'"/>
          <Indoor :live="live" :loc="'1'" :title="'Basement'"/>
        </div>
        <div class="weatherfooter-container">
            <Footer />
        </div>
        <ChartModal v-if="models.chart" @close="closeModal" :options="modalOptions"/>
        <AlertModal v-if="models.alert" :forecast="forecast" @close="closeModal" :options="modalOptions"/>
        <AlmanacModel v-if="models.almanac" @close="closeModal" :options="modalOptions"/>
        <Radar v-if="models.radar" @close="closeModal"/>
        <MetarModal v-if="models.metar" :astro="astro" @close="closeModal"/>
        <RainfallAlmanac v-if="models.rainfallalmanac" :current=current @close="closeModal"/>
    </div>

</template>

<script>

    import Time from './Time';
    import MinMax from './MinMax';
    import Rainfall from './Rainfall';
    import Alert from './Alert';
    import Temperature from './Temperature';
    import RainfallDetails from './RainfallDetails';
    import Wind from './Wind';
    import Barometer from './Barometer';
    import Daylight from './Daylight';
    import Moon from './Moon';
    import Indoor from './Indoor';
    import Lightining from "./Lightining";
    import Uv from './Uv';
    import ChartModal from './ChartModal';
    import AlertModal from './AlertModal';
    import AlmanacModel from './AlmancModel';
    import Footer from './Footer';
    import Radar from './Radar';
    import MetarModal from "./MetarModal";
    import RainfallAlmanac from "./RainfallAlmanac";
    import Menu from "./Menu";
    import axios from 'axios';


    export default {
        name: 'app',
        components: {
            Time,
            MinMax,
            Rainfall,
            Alert,
            Temperature,
            RainfallDetails,
            Wind,
            Barometer,
            Daylight,
            Moon,
            Indoor,
            Uv,
            Footer,
            Radar,
            ChartModal,
            AlertModal,
            AlmanacModel,
            MetarModal,
            Menu,
            RainfallAlmanac,
            Lightining

        },
        data () {
            return {
                loaded: false,
                live: null,
                current: null,
                forecast: null,
                temp: null,
                lightning: null,
                astro: null,
                wind: null,
                theme: 'dark',
                models :{
                    chart:false,
                    alert:false,
                    almanac: false,
                    rainfallalmanac: false,
                    radar: false,
                    metar: false,
                    forecasthourly: false,
                    forecastsummary: false
                },
                modalOptions:null,
            }
        },
        methods: {
            showModal(type,options) {
                this.modalOptions = options;
                this.models[type] = true;
            },
            closeModal(type) {
                this.models[type] = false;
            },
        },
        computed: {

        },
        watch: {

        },
        beforeCreate() {
            this.$store.dispatch('getSettings');
        },
        beforeMount() {
            this.setStyle();
            this.loaded = true;
        },
        mounted () {
            axios.get('/api/luna').then(response => (this.astro = response.data));
            function updateData(self){
                axios.get('/api/wind').then(response => (self.wind = response.data))
                axios.get('/api/minmax/tempf').then(response => (self.temp = response.data));
                axios.get('/api/minmax/lightning').then(response => (self.lightning = response.data));
                axios.get('/api/current').then(response => (self.current = response.data));
                setTimeout(function() { updateData(self); }, 60000);
            }
            updateData(this);
            this.$options.sockets.onmessage = (msg) => {
                let data = JSON.parse(msg.data);
                this.live = data
            };
            window.addEventListener("keyup", e => {
                if(e.key === 'Escape') {
                    Object.keys(this.models).forEach((modal) => {
                        this.closeModal(modal);
                    });
                }
            },this);

        },
        created() {
            this.$on('unitsChange', units => {
                this.units = units;
            });
        }
    }
</script>

<style>
</style>
