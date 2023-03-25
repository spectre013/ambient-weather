<script setup>
import {ref, onMounted, onBeforeMount} from 'vue';
import axios from 'axios';
import StationTime from './StationTime.vue'
import MinMax from './MinMax.vue';
import Rainfall from './RainfallBox.vue';
import Alert from './AlertBox.vue';
import Temperature from './TemperatureBox.vue';
import RainfallDetails from './RainfallDetails.vue';
import Wind from './WindBox.vue';
import Barometer from './BarometerBox.vue';
import Daylight from './DaylightBox.vue';
import Moon from './MoonBox.vue';
import Indoor from './IndoorBox.vue';
import Lightining from './LightiningBox.vue';
import Uv from './UvBox.vue';
import ChartModal from './ChartModal.vue';
import AlmanacModel from './AlmancModel.vue';
import Footer from './FooterBox.vue';
import AlertModal from './AlertModal.vue';
// import Radar from './RadarModal.vue';
// import MetarModal from './MetarModal';
//import RainfallAlmanac from './RainfallAlmanac.vue';
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
//let theme = ref('dark');
let models = ref({
  chart: false,
  indoor: false,
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
    connection.addEventListener('close', () => {
      console.log('Connection Close!');
      setTimeout(function() {
        connection = new WebSocket(wsurl);
      }, 1000);
    });
    // Listen for messages
    connection.addEventListener('message', (event) => {
      live.value = JSON.parse(event.data);
    });

    connection.onerror = function(error) {
      console.log(`[error]`,error);
      connection.close();
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
      <Alert :alerts="alerts" @openModal="show"  @closeModal="close" />
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
          <Lightining :current="live" :stats="lightning" @openModal="show"  @closeModal="close"/>
          <Indoor :live="live" :loc="'in'" :title="'Indo0or'" @openModal="show"  @closeModal="close" />
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
        <AlertModal v-if="models.alert" :alerts="alerts" @closeModal="close" :options="modalOptions" />
        <AlmanacModel v-if="models.almanac" @closeModal="close" :options="modalOptions" />
    <!--    <Radar v-if="models.radar" @close="closeModal" />-->
    <!--    <MetarModal v-if="models.metar" :astro="astro" @close="closeModal" />-->
<!--        <RainfallAlmanac v-if="models.rainfallalmanac" :current="current" @closeModal="close" />-->
  </div>
</template>

<style></style>
