<template>
  <div class="weather34box earthquake">
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
      Rainfall<oblue> Data </oblue>
    </div>
    <div class="value" v-if="current">
      <div id="position3">
        <div class="topmin">
          <topblue1>
            {{ weather.rainDisplay(props.current.monthlyrainin, store.getters.units) }}
            <smallwindunit>{{ weather.rainLabel(store) }}</smallwindunit></topblue1>
        </div>
        <div class="minword">{{ month(props.current.date) }}</div>
        <div class="mintimedate">Total</div>
        <div class="yearwordbig">{{ year(props.current.date) }}</div>
        <div class="topmax">
          <topblue1>
            {{ weather.rainDisplay(props.current.yearlyrainin, store.getters.units) }}
            <smallwindunit>{{ weather.rainLabel(store) }}</smallwindunit>
          </topblue1>
        </div>
        <div class="maxword">{{ year(props.current.date) }}</div>
        <div class="maxtimedate">Total</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import * as weather from '@/weather'
import moment from 'moment';
import { useStore } from 'vuex'

const store = useStore()
const props = defineProps({
  current: Object
});

function year(date) {
  return moment(date).format('YYYY');
}
function month(date) {
  return moment(date).format('MMM');
}


</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped></style>
