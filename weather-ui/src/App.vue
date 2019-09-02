<template>
  <div id="app">
    <div class="weather2-container">
      <Time/>
      <MinMax :temp="temp" />
      <Rainfall :current="current"/>
      <Alert :forecast="forecast"/>
    </div>
    <div class="weather-container">
      <Temperature :current="live" :temp="temp" v-on:openModal="showModal"/>
      <Forecast :forecast="forecast"/>
      <RainfallDetails :current="current"/>
    </div>
    <div class="weather-container">
      <Wind :current="live" :wind="wind" />
      <Barometer :current="live" />
      <Daylight :astro="astro"/>
    </div>
    <div class="weather-container">
      <Moon :astro="astro" />
      <Current :forecast="forecast"/>
      <Uv :current="live" :astro="astro"/>
    </div>
    <div class="weather-container">
      <Indoor :live="live" :loc="'in'" :title="'Indoor'"/>
      <Indoor :live="live" :loc="'1'" :title="'Basement'"/>
      <Indoor :live="live" :loc="'2'" :title="'Master Bedroom'"/>

    </div>
    <div class="weatherfooter-container">
      <Footer />
    </div>
    <ChartModal v-if="models.chart" @close="closeModal" :options="modalOptions"/>
    <AlertModal v-if="models.alert" :forecast="forecast" @close="closeModal" :options="modalOptions"/>
    <AlmanacModel v-if="models.almanac" @close="closeModal"/>
    <Radar v-if="models.radar" @close="closeModal"/>
    <MetarModal v-if="models.metar" :astro="astro" @close="closeModal"/>
    <ForecastSummary v-if="models.forecastsummary" :forecast="forecast" @close="closeModal"/>
    <ForecastHourly v-if="models.forecasthourly" :forecast="forecast" @close="closeModal"/>
  </div>

</template>

<script>
import io from 'socket.io-client';
import Forecast from './components/Forecast';
import Time from './components/Time';
import MinMax from './components/MinMax';
import Rainfall from './components/Rainfall';
import Alert from './components/Alert';
import Temperature from './components/Temperature';
import RainfallDetails from './components/RainfallDetails';
import Wind from './components/Wind';
import Barometer from './components/Barometer';
import Daylight from './components/Daylight';
import Moon from './components/Moon';
import Indoor from './components/Indoor';
import Current from './components/Current'
import Uv from './components/Uv';
import ChartModal from './components/ChartModal';
import AlertModal from './components/AlertModal';
import AlmanacModel from './components/AlmancModel';
import Footer from './components/Footer';
import Radar from './components/Radar';
import MetarModal from "./components/MetarModal";
import ForecastSummary from "./components/ForecastSummary";
import ForecastHourly from "./components/ForecastHourly";
import axios from 'axios';

export default {
  name: 'app',
  components: {
    Forecast,
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
    Current,
    Uv,
    Footer,
    Radar,
    ChartModal,
    AlertModal,
    AlmanacModel,
    MetarModal,
    ForecastSummary,
    ForecastHourly

  },
  data () {
    return {
      live: null,
      current: null,
      forecast: null,
      temp: null,
      astro: null,
      wind: null,
      models :{
        chart:false,
        alert:false,
        almanac: false,
        radar: false,
        metar: false,
        forecasthourly: false,
        forecastsummary: false
      },
      modalOptions:null,
      socket : io('/', {path:'/api/ws'})
    }
  },
  methods: {
    showModal(type,options) {
      this.modalOptions = options;
      this.models[type] = true;
      console.log(type, this.models[type]);
    },
    closeModal(type) {
      this.models[type] = false;
    }
  },
  computed: {

  },
  mounted () {
    this.socket.on('data', (data) => {
        this.live = data;
    });
    axios.get('/api/forecast').then(response => (this.forecast = response.data));
    axios.get('/api/luna').then(response => (this.astro = response.data));
    function updateData(self){
      axios.get('/api/wind').then(response => (self.wind = response.data))
      axios.get('/api/minmax/tempf').then(response => (self.temp = response.data));
      axios.get('/api/current').then(response => (self.current = response.data));
        setTimeout(function() { updateData(self); }, 60000);
    }
    updateData(this);

    window.addEventListener("keyup", e => {
        if(e.key === 'Escape') {
          Object.keys(this.models).forEach((modal) => {
            this.closeModal(modal);
          });
        }
    },this);

  },
}
</script>

<style>
</style>
