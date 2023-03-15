<template>
  <div class="weather-item">
    <div class="chartforecast">
      <span class="yearpopup">
        <a
          alt="almanac temperature"
          title="almanac temperature"
          href="#"
          v-on:click="
            openModal('almanac', { tempfield: 'temp' + props.loc + 'f', humidfield: 'humidity' + props.loc })
          "
          data-lity=""
        >
          <svg
            viewBox="0 0 32 32"
            width="8"
            height="8"
            fill="none"
            stroke="#777"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="6.25%"
          >
            <path d="M14 9 L3 9 3 29 23 29 23 18 M18 4 L28 4 28 14 M28 4 L14 18"></path>
          </svg>
          Almanac
        </a></span
      >
      <span class="yearpopup">
        <a
          v-on:click="openModal('chart', { time: 'year', type: 'temp', field: 'temp' + props.loc + 'f' })"
          alt="Temperature"
          title="Temperature"
          href="#"
        >
          <svg
            version="1.1"
            width="8pt"
            height="8pt"
            x="0px"
            y="0px"
            viewBox="0 0 496 496"
            style="enable-background: new 0 0 496 496"
            xml:space="preserve"
          >
            <path
              style="fill: #719fa3"
              d="M496,428c0,6.4-5.6,12-12,12H12c-6.4,0-12-5.6-12-12l0,0c0-6.4,5.6-12,12-12h472 C490.4,416,496,421.6,
              496,428L496,428z"
            ></path>
            <rect x="24" y="56" style="fill: #1589ad" width="88" height="344"></rect>
            <polyline style="fill: #04567f" points="24,56 112,56 112,400 "></polyline>
            <rect x="144" y="128" style="fill: #24966a" width="88" height="272"></rect>
            <polyline style="fill: #007763" points="144,128 232,128 232,400 "></polyline>
            <rect x="264" y="208" style="fill: #e8961f" width="88" height="192"></rect>
            <polyline style="fill: #e57520" points="264,208 352,208 352,400 "></polyline>
            <rect x="384" y="272" style="fill: #d32a0f" width="88" height="128"></rect>
            <polyline style="fill: #af1909" points="384,272 472,272 472,400 "></polyline>
            <g></g>
          </svg>
          {{ chartDates('Y') }}
        </a></span
      >
      <span class="todaypopup">
        <a
          alt="Feels Like"
          title="Feels Like"
          href="#"
          data-lity=""
          v-on:click="
            openModal('chart', { time: 'month', type: 'temp', field: 'temp' + props.loc + 'f' })
          "
        >
          <svg
            version="1.1"
            width="8pt"
            height="8pt"
            x="0px"
            y="0px"
            viewBox="0 0 496 496"
            style="enable-background: new 0 0 496 496"
            xml:space="preserve"
          >
            <path
              style="fill: #719fa3"
              d="M496,428c0,6.4-5.6,12-12,12H12c-6.4,0-12-5.6-12-12l0,0c0-6.4,5.6-12,12-12h472 C490.4,416,496,421.6,
              496,428L496,428z"
            ></path>
            <rect x="24" y="56" style="fill: #1589ad" width="88" height="344"></rect>
            <polyline style="fill: #04567f" points="24,56 112,56 112,400 "></polyline>
            <rect x="144" y="128" style="fill: #24966a" width="88" height="272"></rect>
            <polyline style="fill: #007763" points="144,128 232,128 232,400 "></polyline>
            <rect x="264" y="208" style="fill: #e8961f" width="88" height="192"></rect>
            <polyline style="fill: #e57520" points="264,208 352,208 352,400 "></polyline>
            <rect x="384" y="272" style="fill: #d32a0f" width="88" height="128"></rect>
            <polyline style="fill: #af1909" points="384,272 472,272 472,400 "></polyline>
            <g></g>
          </svg>
          {{ chartDates('MMM') }}
        </a></span
      >
      <span class="todaypopup">
        <a
          alt="Greenhouse"
          title="Greenhouse"
          href="#"
          data-lity=""
          v-on:click="openModal('chart', { time: 'day', type: 'temp', field: 'temp' + props.loc + 'f' })"
        >
          <svg
            width="8pt"
            height="8pt"
            x="0px"
            y="0px"
            viewBox="0 0 496 496"
            style="enable-background: new 0 0 496 496"
            xml:space="preserve"
          >
            <path
              style="fill: #719fa3"
              d="M496,428c0,6.4-5.6,12-12,12H12c-6.4,0-12-5.6-12-12l0,0c0-6.4,5.6-12,12-12h472 C490.4,416,496,421.6,
              496,428L496,428z"
            ></path>
            <rect x="24" y="56" style="fill: #1589ad" width="88" height="344"></rect>
            <polyline style="fill: #04567f" points="24,56 112,56 112,400 "></polyline>
            <rect x="144" y="128" style="fill: #24966a" width="88" height="272"></rect>
            <polyline style="fill: #007763" points="144,128 232,128 232,400 "></polyline>
            <rect x="264" y="208" style="fill: #e8961f" width="88" height="192"></rect>
            <polyline style="fill: #e57520" points="264,208 352,208 352,400 "></polyline>
            <rect x="384" y="272" style="fill: #d32a0f" width="88" height="128"></rect>
            <polyline style="fill: #af1909" points="384,272 472,272 472,400 "></polyline>
            <g></g>
          </svg>
          Today</a
        ></span
      >
    </div>
    <span class="moduletitle" v-if="title">
      {{ props.title }} (<valuetitleunit>&deg;{{ weather.tempLabel(store) }}</valuetitleunit
      >) </span
    ><br />
    <div id="temperature" v-if="live && minmax && trend">
      <div class="updatedtime">
        <span
          ><svg
            id="i-info"
            viewBox="0 0 32 32"
            width="6"
            height="6"
            fill="#9aba2f"
            stroke="#9aba2f"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="6.25%"
          >
            <path d="M16 14 L16 23 M16 8 L16 10"></path>
            <circle cx="16" cy="16" r="14"></circle>
          </svg>
          {{ now(props.live.date) }}
        </span>
      </div>
      <div class="tempcontainer">
        <div class="maxdata">
          {{ weather.tempDisplay(minmax.max.day.value, store.getters.units) }}&deg; |
          {{ weather.tempDisplay(minmax.min.day.value, store.getters.units) }}&deg;
        </div>
        <div class="maxdata"></div>
        <div v-bind:class="weather.tempClass(temp())">
          {{ weather.tempDisplay(temp(), store.getters.units)
          }}<smalltempunit>&deg;{{ weather.tempLabel(store) }}</smalltempunit>
        </div>
        <div class="temptrendx"></div>
      </div>
      <div class="heatcircle">
        <div class="heatcircle-content">
          <valuetextheading1>Humidity</valuetextheading1><br />
          <div class="tempconverter1">
            <div v-bind:class="weather.smallTempClass(humidity())">
              {{ humidity() }}<smalltempunit2>%</smalltempunit2>
            </div>
          </div>
        </div>
        <div class="heatcircle2">
          <div class="heatcircle-content">
            <valuetextheading1>Avg Today</valuetextheading1>
            <div class="tempconverter1">
              <div v-bind:class="weather.smallTempClass(minmax.avg.day.value)">
                {{ weather.tempDisplay(minmax.avg.day.value, store.getters.units) }}
                <smalltempunit2>&deg;{{ weather.tempLabel(store) }}</smalltempunit2>
              </div>
            </div>
          </div>
        </div>
        <div class="heatcircle3">
          <div class="heatcircle-content">
            <valuetextheading1>High</valuetextheading1>
            <div class="tempconverter1">
              <div v-bind:class="weather.smallTempClass(minmax.max.day.value)">
                {{ weather.tempDisplay(minmax.max.day.value, store.getters.units) }}
                <smalltempunit2>&deg;{{ weather.tempLabel(store) }}</smalltempunit2>
              </div>
            </div>
          </div>
        </div>
        <div class="heatcircle4">
          <div class="heatcircle-content">
            <valuetextheading1>Low</valuetextheading1>
            <div class="tempconverter1">
              <div v-bind:class="weather.smallTempClass(minmax.min.day.value)">
                &nbsp;{{ weather.tempDisplay(minmax.min.day.value, store.getters.units) }}
                <smalltempunit2>&deg;{{ weather.tempLabel(store) }}</smalltempunit2>&nbsp;
              </div>
            </div>
          </div>
        </div>
      </div>
      <div class="tempconverter2">
        <div v-bind:class="weather.tempCircle(temp())">
          {{ toCelcius(temp(), store.getters.units) }}
          <smalltempunit2>&nbsp;&deg;{{ weather.tempLabelAlt(store) }}</smalltempunit2>
        </div>
      </div>
    </div>
    <br />
  </div>
</template>
<script setup>
import * as weather from '@/weather';
import moment from 'moment';
import axios from 'axios';
import {onMounted, ref} from "vue";
import {useStore} from "vuex";

const store = useStore()
const props = defineProps({
  title: String,
  loc: String,
  live: Object,
});

const emit = defineEmits(['openModal']);
let trend=ref(null);
let minmax=ref(null);

onMounted(() => {
  function updateData() {
    axios
        .get('/api/trend/temp' + props.loc + 'f')
        .then((response) => (trend.value = response.data));
    axios
        .get('/api/minmax/temp' + props.loc + 'f')
        .then((response) => (minmax.value = response.data));
    setTimeout(function () {
      updateData(self);
    }, 60000);
  }
  updateData(this);
});

function openModal(type, options) {
  emit('openModal',{type:type,options:options})
}


// function closeModal() {
//   emit('closeModal','chart');
// }

function chartDates(format) {
  return moment().format(format);
}
function now(date) {
  return moment(date).format('HH:mm:ss');
}
function toCelcius(temp, units) {
  if (units === 'imperial') {
    return (((temp - 32) * 5) / 9).toFixed(2);
  }
  return temp;
}
function temp() {
  let key = 'temp' + props.loc.toLowerCase() + 'f';
  return props.live[key];
}
function humidity() {
  let key = 'humidity' + props.loc;
  return props.live[key];
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
#home {
  position: relative;
  top: 10px;
  left: 5px;
}
</style>
