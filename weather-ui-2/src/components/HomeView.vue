<script setup>
import {ref, onMounted, onBeforeMount} from 'vue';
import axios from 'axios';
import StationTime from './StationTime.vue'
import MinMax from './MinMax.vue';
import Rainfall from './Rainfall.vue';
import Alert from './Alert.vue';
import Temperature from './Temperature.vue';
import RainfallDetails from './RainfallDetails.vue';
import Wind from './Wind.vue';
import Barometer from './Barometer.vue';
import Daylight from './Daylight.vue';
import Moon from './Moon.vue';
import Indoor from './Indoor.vue';
import Lightining from './Lightining.vue';
import Uv from './Uv.vue';
import ChartModal from './ChartModal.vue';
// import AlertModal from './AlertModal';
import AlmanacModel from './AlmancModel.vue';
import Footer from './Footer.vue';
// import Radar from './Radar';
// import MetarModal from './MetarModal';
import RainfallAlmanac from './RainfallAlmanac.vue';
// import Menu from './Menu';
import {useStore} from "vuex";


let loaded = ref(false);
let live = ref(null);
let current = ref(null);
let alerts = ref(null);
let temp = ref(null);
let lightning = ref(null);
let astro = ref(null);
let wind = ref(null);
let theme = ref('dark');
let models = ref({
  chart: false,
  alert: false,
  almanac: false,
  rainfall: false,
  radar: false,
  metar: false,
});
const store = useStore()
let modalOptions = ref(null);
let connection = null;

function show(values) {
  console.log(values);
  modalOptions.value = values.options;
  models.value[values.type] = true;
}
function close(event) {
  models.value[event] = false;
}

onBeforeMount(() => {
  //setStyle();
  store.dispatch('getSettings');
  loaded.value = true;
});

onMounted(() => {
    axios.get('/api/luna').then((response) => (astro.value = response.data));
    function updateData() {
      axios.get('/api/wind').then((response) => (wind.value = response.data));
      axios.get('/api/minmax/tempf').then((response) => (temp.value = response.data));
      axios.get('/api/minmax/lightning').then((response) => (lightning.value = response.data));
      axios.get('/api/current').then((response) => (current.value = response.data));
      axios.get('/api/alerts').then((response) => (alerts.value = response.data));

      setTimeout(function () {
        updateData(self);
      }, 60000);
    }

    updateData();
    let wsurl = 'wss://' + window.location.host;
    if (window.location.protocol === 'http:') {
      wsurl = 'ws://' + window.location.host;
    }
    wsurl += '/api/ws';
    connection = new WebSocket(wsurl);

    connection.addEventListener('open', () => {
      console.log('Connection Open!');
    });

    // Listen for messages
    connection.addEventListener('message', (event) => {
      live.value = JSON.parse(event.data);
    });

    connection.onerror = function(error) {
      console.log(`[error]`,error);
    };

    window.addEventListener(
      'keyup',
      (e) => {
        if (e.key === 'Escape') {
          Object.keys(models.value).forEach((modal) => {
            close(modal);
          });
        }
      },
      this,
    );
  });
  // created() {
  //   this.$on('unitsChange', (units) => {
  //     this.units = units;
  //   });
  // },

</script>
<template>
  <div id="app">
<!--    <Menu />-->
    <div class="weather2-container">
      <StationTime />
      <MinMax :temp="temp" />
      <Rainfall :current="current" />
      <Alert :alerts="alerts" />
    </div>
    <div class="weather-container">
      <Temperature :current="live" :temp="temp" @openModal="show"  @closeModal="close"/>
      <Wind :current="live" :wind="wind" />
      <RainfallDetails :current="current" />
    </div>
        <div class="weather-container">
          <Barometer :current="live" />
          <Moon :astro="astro" />
          <Daylight :astro="astro" />
        </div>
        <div class="weather-container">
          <Uv :current="live" :astro="astro" />
          <Lightining :current="live" :stats="lightning" />
          <Indoor :live="live" :loc="'in'" :title="'Indoor'" />
        </div>
        <div class="weather-container">
          <Indoor :live="live" :loc="'2'" :title="'Master'" @openModal="show"  @closeModal="close"/>
          <Indoor :live="live" :loc="'3'" :title="'Office'" @openModal="show"  @closeModal="close"/>
          <Indoor :live="live" :loc="'1'" :title="'Basement'" @openModal="show"  @closeModal="close"/>
        </div>
        <div class="weatherfooter-container">
          <Footer />
        </div>
        <ChartModal v-if="models.chart" @closeModal="close" :options="modalOptions" />
    <!--    <AlertModal v-if="models.alert" :alerts="alerts" @close="closeModal" :options="modalOptions" />-->
        <AlmanacModel v-if="models.almanac" @closeModal="close" :options="modalOptions" />
    <!--    <Radar v-if="models.radar" @close="closeModal" />-->
    <!--    <MetarModal v-if="models.metar" :astro="astro" @close="closeModal" />-->
<!--        <RainfallAlmanac v-if="models.rainfallalmanac" :current="current" @closeModal="close" />-->
  </div>
</template>

<style></style>
