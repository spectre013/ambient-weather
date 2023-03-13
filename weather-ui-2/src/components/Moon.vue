<template>
  <div class="weather-item">
    <div class="chartforecast">
      <span class="yearpopup"></span>
    </div>
    <span class="moduletitle">Moon Phase</span>
    <div id="dldata" v-if="props.astro">
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
            stroke-width="6.25%">
            <path d="M16 14 L16 23 M16 8 L16 10"></path>
            <circle cx="16" cy="16" r="14"></circle>
          </svg>
          {{ now() }}</span
        >
      </div>
      <div class="moonphasemoduleposition">
        <div class="moonrise1">
          <svg
            id="weather34 moon rise"
            viewBox="0 0 32 32"
            width="6"
            height="6"
            fill="none"
            stroke="#ff9350"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="10%">
            <path d="M30 20 L16 8 2 20"></path>
          </svg>
          Moon <br />
          <blueu>Rise &nbsp;<moonrisecolor>
              {{ props.astro.moonrise }}
              <div class="weather34moonmodulepos">
                <div id="weather34moonphases"></div>
                <div class="weather34moonmodule">
                  <svg id="weather34 simple moonphase">
                    <circle cx="50" cy="50" r="49.5" fill="rgba(86, 95, 103, 0.8)"></circle>
                    <path id="weather34shape" fill="currentcolor" v-bind:d="weather34Moon()"></path>
                  </svg>
                </div>
              </div>
            </moonrisecolor>
          </blueu>
        </div>

        <div class="fullmoon1">
          <svg
            id="weather34 full moon"
            viewBox="0 0 32 32"
            width="6"
            height="6"
            fill="#aaa"
            stroke="#aaa"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="6.25%">
            <circle cx="16" cy="16" r="14"></circle>
            <path d="M6 6 L26 26"></path>
          </svg>
          Next Full Moon <br />
          <div class="nextfullmoon">
            <value><moonm>{{ format(props.astro.fullmoon) }}</moonm></value>
          </div>
        </div>
        <div class="newmoon1">
          <svg
            id="weather34 new moon"
            viewBox="0 0 32 32"
            width="6"
            height="6"
            fill="#777"
            stroke="#777"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="6.25%">
            <circle cx="16" cy="16" r="14"></circle>
            <path d="M6 6 L26 26"></path>
          </svg>
          Next New Moon
          <div class="nextnewmoon">
            <value><moonm>{{ format(props.astro.nextnewmoon) }}</moonm></value>
          </div>
        </div>
        <div class="moonset1">
          <svg
            id="weather34 moon set"
            viewBox="0 0 32 32"
            width="6"
            height="6"
            fill="none"
            stroke="#f26c4f"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="10%">
            <path d="M30 12 L16 24 2 12"></path>
          </svg>
          Moon
          <div class="nextnewmoon">
            Set&nbsp;<moonsetcolor> {{ props.astro.moonset }}</moonsetcolor>
          </div>
        </div>
        <div class="meteorshower">
          <svg xmlns="http://www.w3.org/2000/svg" width="10px" height="10px" viewBox="0 0 16 16">
            <path
              fill="currentcolor"
              d="M0 0l14.527 13.615s.274.382-.081.764c-.355.382-.82.055-.82.055L0 0zm4.315 1.364l11.277
              10.368s.274.382-.081.764c-.355.382-.82.055-.82.055L4.315 1.364zm-3.032 2.92l11.278
              10.368s.273.382-.082.764c-.355.382-.819.054-.819.054L1.283 4.284zm6.679-1.747l7.88
              7.244s.19.267-.058.534-.572.038-.572.038l-7.25-7.816zm-5.68 5.13l7.88
              7.244s.19.266-.058.533-.572.038-.572.038l-7.25-7.815zm9.406-3.438l3.597
              3.285s.094.125-.029.25c-.122.125-.283.018-.283.018L11.688 4.23zm-7.592 7.04l3.597
              3.285s.095.125-.028.25-.283.018-.283.018l-3.286-3.553z"
            ></path>
          </svg>
          Perseids Aug 11th-13th
        </div>

        <div class="weather34moonphasem2">Moon Phase<br />{{ props.astro.phase }}</div>
        <div class="weather34luminancem2">Luminance<br />{{ props.astro.illuminated }}%</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import moment from 'moment';
const props = defineProps({
  astro: Object,
});

function format(date) {
  return moment.utc(date).format('MM-DD-YYYY');
}
function now() {
  return moment().format('HH:mm:ss');
}

function weather34Moon() {
  let day = Date.now() / 86400000;
  let referenceweather34Moon = Date.UTC(2018, 0, 17, 2, 17, 0, 0);
  let refweather34Day = referenceweather34Moon / 86400000;
  let phase = ((day - refweather34Day) % 29.530588853) * 1.008;
  let s = String;
  let weather34moonCurve;
  let lf = Math.min(3 - 4 * (phase / 30), 1);
  let lc = Math.abs(lf * 50);
  let lb = lf < 0 ? '0' : '1';
  let rf = Math.min(3 + 4 * ((phase - 30) / 30), 1);
  let rc = Math.abs(rf * 50);
  let rb = rf < 0 ? '0' : '1';
  weather34moonCurve =
      'M 50,0 ' +
      'a ' +
      s(lc) +
      ',50 0 0 ' +
      lb +
      ' 0,100 ' +
      'a ' +
      s(rc) +
      ',50 0 0 ' +
      rb +
      ' 0,-100';
  return weather34moonCurve;
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped></style>
