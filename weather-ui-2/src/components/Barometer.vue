<template>
  <div class="weather-item">
    <div class="chartforecast"></div>
    <span class="moduletitle">Barometer (<valuetitleunit>{{ weather.baroLabel(store) }}</valuetitleunit>)</span><br />
    <div id="barometer" v-if="minmax && props.current && trend">
      <div class="updatedtime">
        <span><svg
            id="i-info"
            viewBox="0 0 32 32"
            width="6"
            height="6"
            fill="#9aba2f"
            stroke="#9aba2f"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="6.25%">
            <path d="M16 14 L16 23 M16 8 L16 10"></path>
            <circle cx="16" cy="16" r="14"></circle>
          </svg>
          {{ now(props.current.date) }}</span
        >
      </div>
      <div class="barometermax">
        <div class="barometerorange">
          <valuetext>Max ({{ timeFormat(minmax.max.day.date) }})<br />
            <maxred>
              <value>{{ weather.baroDisplay(minmax.max.day.value, store.getters.units) }}</value>
            </maxred>
            {{ weather.baroLabel(store) }}
          </valuetext>
        </div>
      </div>
      <div class="barometermin">
        <div class="barometerblue">
          <valuetext>Min ({{ timeFormat(minmax.min.day.date) }})<br />
            <minblue>
              <value>{{ weather.baroDisplay(minmax.min.day.value, store.getters.units) }}</value> </minblue>&nbsp;
            {{ weather.baroLabel(store) }}
          </valuetext>
        </div>
      </div>

      <div class="barometertrend2">
        <valuetext
          >&nbsp;&nbsp;Trend<ogreen>
            <svg
              id="steadybarometer"
              width="10"
              height="10"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentcolor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
            >
              <polyline points="9 18 15 12 9 6"></polyline></svg
          ></ogreen>
          <steady>
            <ogreen>
              <value>Steady</value>
            </ogreen>
          </steady>
        </valuetext>
      </div>

      <div class="homeweathercompass2">
        <div class="homeweathercompass-line2">
          <div
            class="weather34barometerarrowactual"
            v-bind:style="{ transform: 'rotate(' + props.current.baromrelin * 13.083 + 'deg)' }"
          ></div>
          <div
            class="weather34barometerarrowmin"
            v-bind:style="{ transform: 'rotate(' + minmax.min.day.value * 13.083 + 'deg)' }"
          ></div>
          <div
            class="weather34barometerarrowmax"
            v-bind:style="{ transform: 'rotate(' + minmax.max.day.value * 13.083 + 'deg)' }"
          ></div>
        </div>
        <div class="text2">
          <div class="pressuretext">
            <ogreen>{{ trend.trend }}</ogreen>
          </div>
          <darkgrey>{{ weather.baroDisplay(props.current.baromrelin, store.getters.units) }}
            <span>{{ weather.baroLabel(store) }}</span></darkgrey>
        </div>
      </div>

      <div class="barometerconverter">
        <div class="barometerconvertercircleblue">
          {{ toFixed(inToMMHG()) }}<smallrainunit>{{ weather.baroLabelAlt(store) }}</smallrainunit>
        </div>
      </div>

      <div class="barometerlimits">
        <div class="weather34-barometerruler">
          <weather34-barometerlimitmin>
            <value>28</value>
          </weather34-barometerlimitmin>
          <weather34-barometerlimitmax>
            <value>32</value>
          </weather34-barometerlimitmax>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import * as weather from '@/weather'
import moment from 'moment';
import {useStore} from "vuex";
import axios from 'axios';
import {onMounted, ref} from "vue";

const store = useStore()
const props = defineProps({
  current: Object
});

let trend = ref(null);
let minmax = ref(null);

onMounted(() => {
  function updateData() {
    axios.get('/api/trend/baromrelin').then((response) => (trend.value = response.data));
    axios.get('/api/minmax/baromrelin').then((response) => (minmax.value = response.data));
    setTimeout(function () {
      updateData();
    }, 60000);
  }
  updateData();
});


function toFixed(bar) {
  return bar.toFixed(2);
}
function now(date) {
  return moment(date).format('HH:mm:ss');
}
function timeFormat(date) {
  return moment(date).format('HH:mm');
}

function inToMMHG() {
  if (store.getters.units === 'imperial') {
    return props.current.baromrelin * 33.86;
  }
  return props.current.baromrelin;
}

</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped></style>
