<template>
  <div class="weather-item">
    <div class="chartforecast">
<span class="yearpopup"> <a href="#" v-on:click="openModal('forecastsummary')" ><svg viewBox="0 0 32 32" width="10" height="10" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M14 9 L3 9 3 29 23 29 23 18 M18 4 L28 4 28 14 M28 4 L14 18"></path></svg> Forecast Summary</a></span>
<span class="yearpopup"> <a href="#" v-on:click="openModal('forecasthourly')"><svg viewBox="0 0 32 32" width="10" height="10" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M14 9 L3 9 3 29 23 29 23 18 M18 4 L28 4 28 14 M28 4 L14 18"></path></svg> Hourly Forecast </a></span>
      </div>
      <span class="moduletitle"> Forecast  </span>
        <br>
          <div id="currentfore" v-if="forecast">
            <div class="updatedtimecurrent">
                <svg id="i-info" viewBox="0 0 32 32" width="7" height="7" fill="#9aba2f" stroke="#9aba2f" stroke-linecap="round"
                    stroke-linejoin="round" stroke-width="6.25%">
                    <path d="M16 14 L16 23 M16 8 L16 10"></path>
                    <circle cx="16" cy="16" r="14"></circle>
                </svg> {{ now() }}
            </div>

            <div class="darkskyforecasthome">
              <div class="darkskydiv">
                <div class="darkskyforecastinghome" v-for="day in forecast.daily.data.slice(0,3)" v-bind:key="day.time">
                  <div class="darkskyweekdayhome">{{ day.time | shortDate}}<br>
                      <img :src="'css/darkskyicons/'+day.icon+'.svg'" width="45"><br>
                  </div>
                  <darkskytemphihome>
                      <span>{{ day.temperatureHigh | toNumber }}&deg;</span>
                  </darkskytemphihome>&nbsp;|&nbsp;
                  <darkskytemplohome>
                      <span>{{ day.temperatureLow | toNumber }}&deg;</span>
                  </darkskytemplohome><br>
                  <img src="css/windicons/avgw.svg" width="20" v-bind:style="{transform:'rotate('+day.windBearing+'deg)'}">
                  <darkskytempwindhome>
                      <span2>{{ day.windBearing | degToCompass }}</span2><br>
                      <span4>{{ day.windSpeed }} </span4>
                      <span2>
                          <oorange> mph</oorange><br>
                          <svg id="weather34 raindrop" x="0px" y="0px" viewBox="0 0 512 512" width="8px" fill="#01a4b5"
                              stroke="#01a4b5" stroke-width="3%">
                                                <g><g><path d="M348.242,124.971C306.633,58.176,264.434,4.423,264.013,3.889C262.08,1.433,259.125,0,256,0	c-3.126,0-6.079,1.433-8.013,3.889c-0.422,0.535-42.621,54.287-84.229,121.083c-56.485,90.679-85.127,161.219-85.127,209.66
                          C78.632,432.433,158.199,512,256,512c97.802,0,177.368-79.567,177.368-177.369C433.368,286.19,404.728,215.65,348.242,124.971z
                          M256,491.602c-86.554,0-156.97-70.416-156.97-156.97c0-93.472,123.907-263.861,156.971-307.658
                          C289.065,70.762,412.97,241.122,412.97,334.632C412.97,421.185,342.554,491.602,256,491.602z"></path></g></g><g><g><path d="M275.451,86.98c-1.961-2.815-3.884-5.555-5.758-8.21c-3.249-4.601-9.612-5.698-14.215-2.45
                          c-4.601,3.249-5.698,9.613-2.45,14.215c1.852,2.623,3.75,5.328,5.688,8.108c1.982,2.846,5.154,4.369,8.377,4.369
                          c2.012,0,4.046-0.595,5.822-1.833C277.536,97.959,278.672,91.602,275.451,86.98z"></path></g></g><g><g><path d="M362.724,231.075c-16.546-33.415-38.187-70.496-64.319-110.213c-3.095-4.704-9.42-6.01-14.126-2.914
                          c-4.706,3.096-6.01,9.421-2.914,14.127c25.677,39.025,46.9,75.379,63.081,108.052c1.779,3.592,5.391,5.675,9.148,5.675
                          c1.521,0,3.064-0.342,4.517-1.062C363.159,242.241,365.224,236.123,362.724,231.075z"></path></g></g></svg>
                          <darkskytempwindhome><span> {{ day.precipProbability }}</span> in</darkskytempwindhome>
                      </span2>
                  </darkskytempwindhome>
                </div>
              </div>
            </div>
          </div>
    </div>
</template>

<script>
import moment from 'moment';

export default {
  name: 'forecast',
  props: {
    forecast: Object
  },
  methods: {
    now: function () {
      return moment().format('HH:mm:ss');
    },
      openModal: function(type,options) {
          this.$parent.showModal(type,options);
      },
  },
  filters: {
    moment: function (date) {
      return moment.unix(date).format('MMMM Do YYYY');
    },
    shortDate: function (date) {
      return moment.unix(date).format('ddd MMM Do');
    },
    toNumber: function (num) {
      return Math.floor(num);
    },
    degToCompass: function(num) {
      var val = Math.floor((num / 22.5) + 0.5);
      var arr = ["North", "NNE", "NE", "ENE", "East", "ESE", "SE", "SSE", "South", "SSW", "SW", "WSW", "West", "WNW", "NW", "NNW"];
      return arr[(val % 16)];
    }
  },
  computed: {
    
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
