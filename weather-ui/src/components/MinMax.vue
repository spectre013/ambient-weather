<script setup>
import * as weather from '@/weather'
import moment from 'moment';
import { useStore } from 'vuex'

const store = useStore()
const props = defineProps({
  temp: {}
});

// function at(date) {
//   return moment(date).format('HH:mm');
// }
function mFormat(date) {
  return moment(date).format('MMM');
}
function dayFormat(date) {
  return moment(date).format('MMM Do');
}
function year(date) {
  return moment(date).format('YYYY');
}


</script>

<template>
  <div class="weather34box indoor">
    <div class="title">
      <svg
          id="i-info"
          viewBox="0 0 32 32"
          width="9"
          height="9"
          fill="none"
          stroke="currentcolor"
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="6.25%"
      >
        <path d="M16 14 L16 23 M16 8 L16 10"></path>
        <circle cx="16" cy="16" r="14"></circle>
      </svg>
      &nbsp;<oblue>Min</oblue> | <ored>Max</ored> Temperatures
    </div>
    <div class="value" v-if="temp">
      <div id="position2">
        <div class="topmin">
          <topblue1>
            {{ weather.tempDisplay(props.temp.min.year.value, store.units) }}
            <smallwindunit>&deg;{{ weather.tempLabel(store) }}</smallwindunit>
          </topblue1>
        </div>
        <div class="minword">{{ mFormat(props.temp.min.year.date) }}</div>
        <div class="mintimedate">{{ dayFormat(props.temp.min.year.date) }}</div>
        <div class="yearwordbig">{{ year(props.temp.min.year.date) }}</div>
        <div class="topmax">
          <toporange1>
            {{ weather.tempDisplay(props.temp.max.year.value, store.units) }}
            <smallwindunit>&deg;{{ weather.tempLabel(store) }}</smallwindunit>
          </toporange1>
        </div>
        <div class="maxword">{{ year(props.temp.max.year.date)}}</div>
        <div class="maxtimedate">{{  dayFormat(props.temp.max.year.date) }}</div>
      </div>
    </div>
  </div>
</template>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped></style>
