<template>
  <div class="weather-item">
    <div class="chartforecast"></div>
    <span class="moduletitle">Direction | Wind (<valuetitleunit>{{ weather.windLabel(store) }}</valuetitleunit>)</span>
    <br />
    <div id="windspeed" v-if="wind && current">
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
          {{ now(props.current.date) }}
        </span>
      </div>
      <br />
      <div class="windspeedvalues">
        <div class="windspeedvalue">
          {{ weather.windDisplay(props.current.windspeedmph, store.getters.speed) }}
          <div class="windunitidspeed">Currently</div>
          <div class="windunitspeed">{{ weather.windLabel(store) }}</div>
        </div>
        <div class="windgustvalue">
          {{ weather.windDisplay(props.current.windgustmph, store.getters.speed) }}
          <div class="windunitgust">{{ weather.windLabel(store) }}</div>
          <div class="windunitidgust">Gust</div>
        </div>
      </div>
      <div class="windspeedtrend1">
        <valuetext>Max
          <max>
            <value>
              <maxred>{{ weather.windDisplay(props.wind.gust.value, store.getters.speed) }}</maxred>
            </value>
          </max>
          <supmb>&nbsp;{{ weather.windLabel(store) }}</supmb
          ><br />
          Gust ({{ timeFormat(props.wind.gust.date) }})
        </valuetext>
      </div>
      <div class="windconverter">
        <div :class="windConvertCircle">
          {{ toKTS(props.current.windgustmph) }}
          <smallrainunit>&nbsp;kts</smallrainunit>
        </div>
      </div>
      <div class="homeweathercompass1">
        <div class="homeweathercompass-line1">
          <div
            class="thearrow2"
            v-bind:style="{ transform: 'rotate(' + props.current.winddir + 'deg)' }"
          ></div>
          <!-- <div class="thearrow1" v-bind:style="{transform:'rotate('+current.winddir+'deg)'}"></div>-->
        </div>
        <div class="text1">
          <div class="windvalue1" id="windvalue">{{ props.current.winddir }}&deg;</div>
        </div>
        <div class="windirectiontext1">
          <span> {{ degToCompass(props.current.winddir) }}<br /></span>
        </div>
      </div>
      <div class="weather34windrun">
        <svg
          id="weather34 windrun symbol"
          viewBox="0 0 24 24"
          width="9pt"
          height="9pt"
          fill="#3b9cac"
        >
          <path
            d="M13,16H3a1,1,0,0,0,0,2H13a1,1,0,0,0,0-2ZM3,8H21a1,1,0,0,0,0-2H3A1,
            1,0,0,0,3,8Zm18,3H3a1,1,0,0,0,0,2H21a1,1,0,0,0,0-2Z"
          ></path></svg><grey>
          <valuetext1>&nbsp;{{ weather.windDisplay(windRun(), store.getters.speed) }}
            <grey><weather34windrunspan>mi</weather34windrunspan></grey></valuetext1></grey
        >
      </div>
      <br />
      <div class="windrun1">Wind Run</div>
      <div :class="beaufortClass()">
        <span v-html="beaufortSvg()"></span>&nbsp;{{ getBeaufort() }}
        <weather34bftspan>BFT</weather34bftspan>
      </div>
      <div class="beaufort1">{{ beaufortText() }}</div>
    </div>
  </div>
</template>

<script setup>
import * as weather from '@/weather'
import moment from 'moment';
import {useStore} from "vuex";

const store = useStore()
const props = defineProps({
  wind: Object,
  current: Object
});
//let beaufort = ref(null);

function windRun() {
  let run = props.wind.avg.value * moment().hours();
  return run.toFixed(2);
}

function getBeaufort() {
  let windkts = props.current.windspeedmph;
  if (props.current.windspeedmph > 0) {
    windkts = props.current.windspeedmph / 1.151;
  }
  let beaufort = 1;
  if (windkts >= 64) {
    beaufort = 12;
  } else if (windkts >= 56) {
    beaufort = 11;
  } else if (windkts >= 48) {
    beaufort = 10;
  } else if (windkts >= 41) {
    beaufort = 9;
  } else if (windkts >= 34) {
    beaufort = 8;
  } else if (windkts >= 28) {
    beaufort = 7;
  } else if (windkts >= 22) {
    beaufort = 6;
  } else if (windkts >= 17) {
    beaufort = 5;
  } else if (windkts >= 11) {
    beaufort = 4;
  } else if (windkts >= 7) {
    beaufort = 3;
  } else if (windkts >= 4) {
    beaufort = 2;
  } else if (windkts >= 1) {
    beaufort = 1;
  } else if (windkts < 1) {
    beaufort = 0;
  }

  return beaufort;
}



function beaufortScale() {
    let beaufortValue = getBeaufort(props.current.windspeedmph);

    if (beaufortValue >= 12) {
      return {
        svg:
          '<svg id="weather34 bft12" width="12pt" height="12pt" viewBox="0 0 96 96">' +
          '<path fill="currentcolor" stroke="currentcolor" stroke-width="0.09375" ' +
          'opacity="1.00" ' +
          'd=" M 0.00 29.96 C 5.55 36.68 11.21 43.31 16.80 50.00 C 18.26 49.99 19.73 49.99 21.19 49.99 C ' +
          '18.93 47.26 16.67 44.53 14.40 41.79 C 15.94 40.54 17.47 39.27 19.00 38.00 C 22.34 42.00 25.66 ' +
          '46.01 29.01 50.00 C 42.72 49.98 56.43 50.03 70.14 49.98 C 71.17 47.82 72.07 45.50 73.83 43.81 ' +
          'C 77.91 39.62 84.85 39.15 89.85 41.94 C 93.15 43.97 95.29 47.56 96.00 51.33 L 96.00 54.56 C ' +
          '95.35 58.38 93.17 62.01 89.84 64.06 C 85.44 66.52 79.67 66.42 75.46 63.60 C 72.81 61.81 71.37 ' +
          '58.87 70.15 56.02 C 46.76 55.98 23.38 56.01 0.00 56.00 L 0.00 29.96 Z" /></svg>',
        text: 'Hurricane',
        class: 'weather34beaufort6',
      };
    } else if (beaufortValue >= 11) {
      return {
        svg:
          '<svg id="weather34 bft11" width="12pt" height="12pt" viewBox="0 0 96 96">' +
          '<path fill="currentcolor" stroke="currentcolor" stroke-width="0.09375" opacity="1.00" ' +
          'd=" M 0.00 29.97 C 5.57 36.68 11.20 43.33 16.81 50.00 C 34.60 49.99 52.38 50.02 70.16 ' +
          '49.99 C 71.98 43.63 78.44 39.00 85.10 40.36 C 90.77 40.90 95.07 45.87 96.00 51.29 L ' +
          '96.00 54.67 C 95.15 59.33 91.95 63.89 87.21 65.17 C 82.45 66.67 76.62 65.56 73.32 61.64 ' +
          'C 71.87 60.01 71.03 57.98 70.16 56.01 C 46.77 55.99 23.39 56.01 0.00 56.00 L 0.00 29.97 Z" />' +
          '</svg>',
        text: 'Violent Storm',
        class: 'weather34beaufort6',
      };
    } else if (beaufortValue >= 10) {
      return {
        svg:
          '<svg id="weather34 bft10"  width="12pt" height="12pt" viewBox="0 0 96 96" version="1.1" >' +
          '<path fill="currentcolor" stroke="currentcolor" stroke-width="0.09375" opacity="1.00" ' +
          'd=" M 0.00 29.97 C 5.30 36.33 10.66 42.65 15.99 48.98 C 16.01 42.66 15.99 36.34 16.00 30.02 C 21.62 ' +
          '36.67 27.19 43.35 32.81 50.00 C 34.20 50.00 35.60 49.99 36.99 50.00 C 33.74 45.99 30.46 42.01 27.21 ' +
          '38.00 C 28.66 36.67 30.12 35.34 31.58 34.01 C 36.02 39.32 40.38 44.69 44.81 50.00 C 53.27 49.99 ' +
          '61.72 50.02 70.18 49.99 C 71.39 46.85 73.14 43.69 76.15 41.96 C 80.11 39.71 85.11 39.63 89.20 ' +
          '41.59 C 92.87 43.50 95.27 47.34 96.00 51.35 L 96.00 54.56 C 95.18 60.08 90.75 65.14 85.02 65.65 C ' +
          '78.40 66.97 71.95 62.35 70.18 56.01 C 46.79 55.98 23.39 56.01 0.00 56.00 L 0.00 29.97 Z" /></svg>',
        text: 'Storm',
        class: 'weather34beaufort6',
      };
    } else if (beaufortValue >= 9) {
      return {
        svg:
          '<svg id="weather34 bft9" width="12pt" height="12pt" viewBox="0 0 96 96">' +
          '<path fill="currentcolor" stroke="currentcolor" stroke-width="0.09375" opacity="1.00" ' +
          'd=" M 0.00 29.97 C 5.29 36.34 10.66 42.65 15.99 48.99 C 16.01 42.66 15.99 36.34 16.00 ' +
          '30.01 C 21.61 36.67 27.19 43.34 32.80 50.00 C 45.26 49.99 57.71 50.02 70.16 49.98 C 71.97 ' +
          '43.66 78.38 39.03 85.02 40.35 C 90.73 40.87 95.12 45.87 96.00 51.36 L 96.00 54.55 C 95.18 ' +
          '60.08 90.75 65.14 85.00 65.66 C 78.37 66.96 71.98 62.34 70.16 56.02 C 46.77 55.98 23.39 56.01 ' +
          '0.00 56.00 L 0.00 29.97 Z" />' +
          '</svg>',
        text: 'Strong Gale',
        class: 'weather34beaufort6',
      };
    } else if (beaufortValue >= 8) {
      return {
        svg:
          '<svg id="weather34 bft8" width="12pt" height="12pt" viewBox="0 0 96 96">' +
          '<path fill="currentcolor" stroke="currentcolor" stroke-width="0.09375" opacity="1.00" ' +
          'd=" M 4.64 30.07 C 10.05 36.70 15.41 43.37 20.82 50.01 C 22.21 50.00 23.60 50.00 25.00 49.99 ' +
          'C 20.66 44.66 16.33 39.33 12.00 34.00 C 13.54 32.67 15.07 31.34 16.60 30.00 C 22.01 36.67 27.40 ' +
          '43.35 32.81 50.01 C 34.21 50.00 35.60 49.99 37.00 49.99 C 32.66 44.66 28.33 39.33 24.00 34.00 ' +
          'C 25.54 32.67 27.07 31.34 28.60 30.00 C 34.01 36.67 39.40 43.35 44.82 50.01 C 46.21 50.00 47.60 ' +
          '50.00 49.00 49.99 C 44.66 44.66 40.33 39.33 36.00 34.00 C 37.54 32.67 39.07 31.34 40.60 30.00 C ' +
          '46.01 36.67 51.40 43.35 56.81 50.01 C 58.34 50.00 59.86 50.00 61.39 49.99 C 58.60 46.59 55.80 ' +
          '43.20 53.00 39.80 C 54.54 38.53 56.07 37.27 57.61 36.01 C 61.73 40.79 65.44 45.94 69.89 50.42 C ' +
          '71.21 47.70 72.41 44.73 74.89 42.83 C 79.11 39.58 85.30 39.36 89.89 41.99 C 93.19 43.96 95.20 ' +
          '47.55 96.00 51.23 L 96.00 54.77 C 95.21 58.43 93.21 62.00 89.94 63.98 C 85.52 66.55 79.63 66.43 ' +
          '75.40 63.55 C 72.77 61.77 71.38 58.81 70.11 56.01 C 52.74 55.99 35.38 56.01 18.01 56.00 C 11.92 ' +
          '48.95 6.57 41.23 0.00 34.64 L 0.00 33.40 C 1.68 32.49 3.18 31.30 4.64 30.07 Z" />' +
          '</svg>',
        text: 'Gale',
        class: 'weather34beaufort6',
      };
    } else if (beaufortValue >= 7) {
      return {
        svg:
          '<svg id="weather34 bft7" width="12pt" height="12pt" viewBox="0 0 96 96">' +
          '<path fill="currentcolor" stroke="currentcolor" stroke-width="0.09375" opacity="1.00" ' +
          'd=" M 0.00 34.01 C 1.53 32.68 3.03 31.30 4.61 30.03 C 10.06 36.65 15.35 43.40 20.85 49.98 ' +
          'C 22.23 50.02 23.60 50.02 24.98 49.97 C 20.67 44.64 16.31 39.36 12.04 34.00 C 13.53 32.64 ' +
          '15.05 31.31 16.61 30.03 C 22.05 36.65 27.35 43.39 32.84 49.98 C 34.22 50.02 35.60 50.02 36.98 ' +
          '49.98 C 32.69 44.64 28.30 39.37 24.04 34.00 C 25.53 32.64 27.05 31.31 28.61 30.03 C 34.05 36.65 ' +
          '39.36 43.39 44.83 49.98 C 46.35 50.02 47.86 50.02 49.38 49.99 C 46.62 46.57 43.78 43.22 41.03 ' +
          '39.80 C 42.53 38.52 44.05 37.24 45.61 36.03 C 49.51 40.65 53.29 45.38 57.22 49.98 C 61.55 50.03 ' +
          '65.88 50.00 70.21 49.99 C 71.17 47.29 72.62 44.67 74.86 42.84 C 78.91 39.72 84.66 39.43 89.20 ' +
          '41.60 C 92.85 43.49 95.26 47.32 96.00 51.30 L 96.00 54.66 C 95.11 60.04 90.82 65.13 85.16 65.58 C ' +
          '78.59 67.06 71.90 62.43 70.21 56.01 C 52.82 55.97 35.43 56.04 18.04 55.98 C 11.96 48.71 6.04 41.31 ' +
          '0.00 34.01 L 0.00 34.01 Z" />' +
          '</svg>',
        text: 'Near Gale',
        class: 'weather34beaufort6',
      };
    } else if (beaufortValue >= 6) {
      return {
        svg:
          '<svg id="weather34 bft6" width="12pt" height="12pt" viewBox="0 0 96 96">' +
          '<path fill="currentcolor" stroke="currentcolor" stroke-width="0.09375" opacity="1.00" ' +
          'd=" M 4.55 30.01 C 10.03 36.62 15.37 43.35 20.81 50.00 C 22.20 50.00 23.60 50.00 24.99 49.99 C ' +
          '20.67 44.65 16.33 39.34 12.01 34.00 C 13.53 32.66 15.07 31.33 16.60 30.00 C 22.02 36.67 27.39 ' +
          '43.38 32.84 50.02 C 34.22 50.01 35.60 49.99 36.98 49.98 C 32.67 44.64 28.31 39.34 24.01 33.99 C ' +
          '25.54 32.66 27.07 31.33 28.60 30.01 C 34.01 36.67 39.39 43.35 44.81 50.00 C 53.26 49.99 61.71 ' +
          '50.01 70.15 49.99 C 71.04 48.00 71.89 45.95 73.36 44.31 C 76.67 40.43 82.45 39.34 87.19 40.83 C ' +
          '91.91 42.08 95.07 46.60 96.00 51.22 L 96.00 54.75 C 95.20 58.73 92.83 62.57 89.13 64.44 C 84.81 ' +
          '66.48 79.42 66.27 75.43 63.58 C 72.80 61.79 71.34 58.86 70.15 56.01 C 52.77 55.99 35.39 56.01 18.01 ' +
          '56.00 C 11.92 48.94 6.51 41.22 0.00 34.56 L 0.00 33.45 C 1.83 32.80 3.11 31.23 4.55 30.01 Z" />' +
          '</svg>',
        text: 'Strong Breeze',
        class: 'weather34beaufort4-5',
      };
    } else if (beaufortValue >= 5) {
      return {
        svg:
          '<svg id="weather34 bft5" width="12pt" height="12pt" viewBox="0 0 96 96">' +
          '<path fill="currentcolor" stroke="currentcolor" stroke-width="0.09375" opacity="1.00" ' +
          'd=" M 4.55 30.01 C 10.04 36.62 15.37 43.37 20.82 50.01 C 22.21 50.00 23.60 49.99 25.00 49.99 C ' +
          '20.67 44.66 16.33 39.33 12.00 34.00 C 13.53 32.67 15.07 31.34 16.60 30.01 C 22.01 36.67 27.39 43.35 ' +
          '32.82 50.01 C 45.26 49.98 57.71 50.02 70.15 49.99 C 71.41 46.91 73.07 43.77 76.03 42.02 C 79.40 ' +
          '40.12 83.56 39.63 87.24 40.85 C 91.95 42.11 95.08 46.63 96.00 51.23 L 96.00 55.03 C 95.11 58.56 ' +
          '93.16 61.97 90.02 63.95 C 85.60 66.53 79.71 66.45 75.44 63.58 C 72.80 61.79 71.34 58.86 70.15 56.01' +
          ' C 52.77 55.99 35.39 56.00 18.02 56.00 C 11.93 48.90 6.44 41.24 0.00 34.48 L 0.00 33.53 C 1.72 32.64' +
          ' 3.15 31.32 4.55 30.01 Z" />' +
          '</svg>',
        text: 'Fresh Breeze',
        class: 'weather34beaufort4-5',
      };
    } else if (beaufortValue >= 4) {
      return {
        svg:
          '<svg id="weather34 bft4" width="12pt" height="12pt" viewBox="0 0 96 96" >' +
          '<path fill="currentcolor" stroke="currentcolor" stroke-width="0.09375" opacity="1.00" ' +
          'd=" M 0.00 33.39 C 1.62 32.38 3.17 31.27 4.69 30.10 C 10.05 36.74 15.43 43.37 20.80 50.01 C 22.27 ' +
          '49.99 23.73 49.99 25.20 49.99 C 22.39 46.60 19.61 43.19 16.80 39.80 C 18.34 38.53 19.87 37.27 21.40 ' +
          '36.00 C 25.26 40.67 29.13 45.33 33.00 50.00 C 45.36 49.99 57.72 50.02 70.08 49.98 C 71.35 47.43 ' +
          '72.52 44.67 74.84 42.87 C 79.08 39.57 85.34 39.34 89.94 42.02 C 93.23 44.01 95.21 47.59 96.00 51.27' +
          ' L 96.00 54.84 C 95.16 58.45 93.23 61.98 89.99 63.95 C 85.38 66.65 79.11 66.44 74.86 63.15 C 72.54' +
          ' 61.35 71.34 58.58 70.08 56.02 C 52.72 55.98 35.37 56.01 18.01 56.00 C 11.92 48.97 6.60 41.23 0.00 ' +
          '34.67 L 0.00 33.39 Z" />' +
          '</svg>',
        text: 'Moderate Breeze',
        class: 'weather34beaufort3-4',
      };
    } else if (beaufortValue >= 3) {
      return {
        svg:
          '<svg id="weather34 bft3" width="12pt" height="12pt" viewBox="0 0 96 96">' +
          '<path fill="currentcolor" stroke="currentcolor" stroke-width="0.09375" opacity="1.00" ' +
          'd=" M 0.00 33.44 C 1.67 32.50 3.17 31.28 4.64 30.06 C 10.04 36.70 15.41 43.36 20.80 50.00 C 37.24 ' +
          '49.99 53.68 50.02 70.12 49.98 C 71.39 47.19 72.76 44.24 75.38 42.46 C 79.66 39.55 85.61 39.46 90.05 ' +
          '42.08 C 93.25 44.09 95.22 47.60 96.00 51.23 L 96.00 54.90 C 95.16 58.48 93.20 61.96 90.01 63.95 C ' +
          '85.59 66.53 79.71 66.44 75.44 63.58 C 72.79 61.80 71.39 58.83 70.12 56.02 C 52.75 55.98 35.38 56.01 ' +
          '18.01 56.00 C 11.92 48.94 6.53 41.24 0.00 34.58 L 0.00 33.44 Z" />' +
          '</svg>',
        text: 'Gentle Breeze',
        class: 'weather34beaufort1-3',
      };
    } else if (beaufortValue >= 2) {
      return {
        svg:
          '<svg id="weather34 bft2" width="12pt" height="12pt" viewBox="0 0 96 96">' +
          '<path fill="currentcolor" stroke="currentcolor" stroke-width="0.09375" opacity="1.00" ' +
          'd=" M 0.00 33.38 C 1.68 32.46 3.19 31.28 4.67 30.09 C 10.04 36.72 15.42 43.36 20.80 50.00 C 37.23 ' +
          '49.99 53.66 50.03 70.09 49.98 C 71.41 47.21 72.76 44.23 75.39 42.45 C 79.66 39.54 85.60 39.45 90.03 ' +
          '42.07 C 93.26 44.08 95.26 47.64 96.00 51.31 L 96.00 54.79 C 95.15 58.68 92.92 62.47 89.30 64.34 C ' +
          '84.74 66.62 78.83 66.29 74.79 63.09 C 72.52 61.29 71.31 58.57 70.10 56.02 C 46.73 55.97 23.37 56.02 ' +
          '0.00 56.00 L 0.00 49.94 C 4.33 50.04 8.66 50.00 13.00 49.99 C 8.62 44.92 4.88 39.28 0.00 34.68 L ' +
          '0.00 33.38 Z" />' +
          '</svg>',
        text: 'Light Breeze',
        class: 'weather34beaufort1-3',
      };
    } else if (beaufortValue >= 1) {
      return {
        svg:
          '<svg id="weather34 bft1" width="12pt" height="12pt" viewBox="0 0 96 96" >' +
          '<path fill="currentcolor" stroke="currentcolor" stroke-width="0.09375" opacity="1.00" ' +
          'd=" M 73.92 43.89 C 77.12 40.10 82.80 39.45 87.34 40.81 C 91.48 42.01 93.99 45.85 96.00 49.39 L 96.00 ' +
          '56.58 C 94.00 60.14 91.49 63.99 87.34 65.19 C 82.80 66.55 77.13 65.90 73.92 62.11 C 72.32 60.28 71.03 ' +
          '58.19 69.69 56.16 C 46.47 55.76 23.23 56.12 0.00 56.00 L 0.00 50.00 C 23.23 49.88 46.47 50.24 69.69 ' +
          '49.84 C 71.03 47.81 72.31 45.73 73.92 43.89 Z" />' +
          '</svg>',
        text: 'Light Air',
        class: 'weather34beaufort1-3',
      };
    } else if (beaufortValue >= 0) {
      return {
        svg:
          '<svg id="weather34 bft0 calm" width="12pt" height="12pt" viewBox="0 0 96 96">' +
          '<path fill="currentcolor" stroke="currentcolor" stroke-width="0.09375" opacity="1.00" ' +
          'd=" M 42.39 20.62 C 51.40 18.80 61.26 21.58 67.72 28.18 C 72.37 33.32 76.08 39.84 76.03 46.95 C ' +
          '76.90 61.75 63.98 75.82 49.07 75.88 C 34.47 76.94 20.93 64.49 20.01 50.07 C 18.76 36.54 29.04 23.08 ' +
          '42.39 20.62 M 42.50 28.79 C 33.81 31.19 27.33 40.09 27.98 49.12 C 28.45 59.28 37.72 68.28 48.01 ' +
          '67.94 C 57.97 68.25 66.94 59.86 67.96 50.09 C 68.99 41.50 63.52 32.72 55.56 29.51 C 51.45 27.80 46.77' +
          ' 27.67 42.50 28.79 Z" />' +
          '<path fill="#fff" stroke="#fff" stroke-width="0.09375" opacity="1.00" d=" M ' +
          '47.33 36.16 C 54.36 35.51 60.80 42.06 59.99 49.07 C 59.51 55.87 52.46 61.15 45.80 59.61 C 38.41 58.73 ' +
          '33.64 49.68 37.08 43.09 C 38.85 39.03 42.93 36.37 47.33 36.16 Z" />' +
          '</svg>',
        text: 'Calm',
        class: 'weather34beaufort1-3',
      };
    }
}

function now(date) {
  return moment(date).format('HH:mm:ss');
}
function degToCompass(num) {
  let val = Math.floor(num / 22.5 + 0.5);
  let arr = [
    'North',
    'NNE',
    'NE',
    'ENE',
    'East',
    'ESE',
    'SE',
    'SSE',
    'South',
    'SSW',
    'SW',
    'WSW',
    'West',
    'WNW',
    'NW',
    'NNW',
  ];
  return arr[val % 16];
}
function timeFormat(date) {
  return moment(date).format('HH:mm');
}
function toKTS(speed) {
  return (speed / 1.15).toFixed(2);
}

function beaufortSvg() {
  return this.beaufortScale().svg;
}
function beaufortText() {
  return this.beaufortScale().text;
}
function beaufortClass() {
  return this.beaufortScale().class;
}
function windConvertCircle() {
  let windkts = this.current.windspeedmph / 1.151;
  let c = 'windconvertercirclegreen1 tgreen';
  if (windkts >= 26.9978) {
    c = 'windconvertercirclered1 tred';
  } else if (windkts >= 21.5983) {
    c = 'windconvertercircleorange1 torange';
  } else if (windkts >= 16.1987) {
    c = 'windconvertercircleblue1';
  } else if (windkts < 16.1987) {
    c = 'windconvertercirclegreen1 tgreen';
  }
  return c;
}


</script>
<style scoped>
.thearrow2 {
  position: absolute;
  z-index: 200;
  top: 0;
  left: 50%;
  margin-left: -5px;
  width: 10px;
  height: 50%;
  -webkit-transform-origin: 50% 100%;
  -moz-transform-origin: 50% 100%;
  -o-transform-origin: 50% 100%;
  -ms-transform-origin: 50% 100%;
  transform-origin: 50% 100%;
  -webkit-transition-duration: 3s;
  -moz-transition-duration: 3s;
  -o-transition-duration: 3s;
  -ms-transition-duration: 3s;
  transition-duration: 3s;
}

.thearrow2:after {
  content: '';
  position: absolute;
  left: 50%;
  top: 0;
  width: 0;
  height: 0;
  border-style: solid;
  border-width: 14px 9px 0 9px;
  border-color: RGBA(255, 121, 58, 1) transparent transparent transparent;
  -webkit-transform: translate(-50%, -50%);
  -moz-transform: translate(-50%, -50%);
  -o-transform: translate(-50%, -50%);
  -ms-transform: translate(-50%, -50%);
  transform: translate(-50%, -50%);
  -webkit-transition-duration: 3s;
  -moz-transition-duration: 3s;
  -o-transition-duration: 3s;
  -ms-transition-duration: 3s;
  transition-duration: 3s;
}

.thearrow2:before {
  content: '  o o o';
  color: rgba(255, 124, 57, 1);
  font-family: Arial, Helvetica, sans-serif;
  font-size: 6px;
  width: 6px;
  height: 6px;
  position: absolute;
  z-index: 9;
  left: 2px;
  top: -5px;
  border: 2px solid RGBA(255, 255, 255, 0.8);
  -webkit-border-radius: 100%;
  -moz-border-radius: 100%;
  -o-border-radius: 100%;
  -ms-border-radius: 100%;
  border-radius: 100%;
}

.thearrow1 {
  position: absolute;
  z-index: 150;
  top: 0;
  left: 50%;
  margin-left: -5px;
  -webkit-transform-origin: 50% 100%;
  -moz-transform-origin: 50% 100%;
  -o-transform-origin: 50% 100%;
  -ms-transform-origin: 50% 100%;
  transform-origin: 50% 100%;
  -webkit-transition-duration: 0s;
  -moz-transition-duration: 0s;
  -o-transition-duration: 0s;
  -ms-transition-duration: 0s;
  transition-duration: 0s;
  background: 0;
}

.thearrow1:after {
  content: '';
  position: absolute;
  text-align: left;
  left: 50%;
  font-size: 8px;
  top: 0;
  width: 0;
  height: 0;
  -webkit-border-radius: 0;
  border-radius: 0;
  border-left: 6px solid transparent;
  border-right: 6px solid transparent;
  border-top: 9px solid rgb(144, 177, 42);
  border-bottom: 0;
  -webkit-transform: translate(-50%, -50%);
  -moz-transform: translate(-50%, -50%);
  -o-transform: translate(-50%, -50%);
  -ms-transform: translate(-50%, -50%);
  transform: translate(-50%, -50%);
  -webkit-transition-duration: 3s;
  -moz-transition-duration: 3s;
  -o-transition-duration: 3s;
  -ms-transition-duration: 3s;
  transition-duration: 3s;
  background: 0;
}

.thearrow1:before {
  content: '  o o o o';
  color: rgb(144, 177, 42);
  font-family: Arial, Helvetica, sans-serif;
  font-size: 4px;
  width: 1px;
  height: 1px;
  position: absolute;
  z-index: 1;
  left: 3px;
  top: -4px;
  border: 2px dotted RGBA(255, 255, 255, 0.8);
  -webkit-border-radius: 100%;
  -moz-border-radius: 100%;
  -o-border-radius: 100%;
  -ms-border-radius: 100%;
  border-radius: 100%;
}

.avgw {
  width: 27px;
  height: 27px;
  margin-left: 35px;
}

spancalm {
  position: relative;
  font-family: weathertext, Arial, serif;
  font-size: 26px;
}
</style>
