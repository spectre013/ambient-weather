<script setup>
import {ref, onMounted, onBeforeMount, getCurrentInstance} from 'vue';
import axios from 'axios';
import StationTime from './StationTime.vue'
import MinMax from './MinMax.vue';
import Rainfall from './Rainfall.vue';
import Alert from './Alert.vue';
import Temperature from './Temperature.vue';
// import RainfallDetails from './RainfallDetails';
// import Wind from './Wind';
// import Barometer from './Barometer';
// import Daylight from './Daylight';
// import Moon from './Moon';
// import Indoor from './Indoor';
// import Lightining from './Lightining';
// import Uv from './Uv';
// import ChartModal from './ChartModal';
// import AlertModal from './AlertModal';
// import AlmanacModel from './AlmancModel';
// import Footer from './Footer';
// import Radar from './Radar';
// import MetarModal from './MetarModal';
// import RainfallAlmanac from './RainfallAlmanac';
// import Menu from './Menu';

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
  rainfallalmanac: false,
      radar: false,
      metar: false,
});

let modalOptions = ref(null)

const { proxy } = getCurrentInstance();

function showModal(type, options) {
  modalOptions.value = options;
  models[type] = true;
}
function closeModal(type) {
  models[type] = false;
}

// beforeCreate() {
//   this.$store.dispatch('getSettings');


onBeforeMount(() => {
  //this.setStyle();
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

    proxy.$socket.onmessage = ({msg}) => {
      console.log(msg);
      live.value = JSON.parse(msg.data);
    }
    window.addEventListener(
      'keyup',
      (e) => {
        if (e.key === 'Escape') {
          Object.keys(models).forEach((modal) => {
            closeModal(modal);
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
      <Temperature :current="live" :temp="temp" v-on:openModal="showModal" />
<!--      <Wind :current="live" :wind="wind" />-->
<!--      <RainfallDetails :current="current" />-->
    </div>
    <!--    <div class="weather-container">-->
    <!--      <Barometer :current="live" />-->
    <!--      <Moon :astro="astro" />-->
    <!--      <Daylight :astro="astro" />-->
    <!--    </div>-->
    <!--    <div class="weather-container">-->
    <!--      <Uv :current="live" :astro="astro" />-->
    <!--      <Lightining :current="live" :stats="lightning" />-->
    <!--      <Indoor :live="live" :loc="'in'" :title="'Indoor'" />-->
    <!--    </div>-->
    <!--    <div class="weather-container">-->
    <!--      <Indoor :live="live" :loc="'2'" :title="'Master'" />-->
    <!--      <Indoor :live="live" :loc="'3'" :title="'Office'" />-->
    <!--      <Indoor :live="live" :loc="'1'" :title="'Basement'" />-->
    <!--    </div>-->
    <!--    <div class="weatherfooter-container">-->
    <!--      <Footer />-->
    <!--    </div>-->
    <!--    <ChartModal v-if="models.chart" @close="closeModal" :options="modalOptions" />-->
    <!--    <AlertModal v-if="models.alert" :alerts="alerts" @close="closeModal" :options="modalOptions" />-->
    <!--    <AlmanacModel v-if="models.almanac" @close="closeModal" :options="modalOptions" />-->
    <!--    <Radar v-if="models.radar" @close="closeModal" />-->
    <!--    <MetarModal v-if="models.metar" :astro="astro" @close="closeModal" />-->
    <!--    <RainfallAlmanac v-if="models.rainfallalmanac" :current="current" @close="closeModal" />-->
  </div>
</template>

<style></style>
