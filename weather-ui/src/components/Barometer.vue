<template>
  <div class="weather-item">
    <div class="chartforecast">
        <span class="yearpopup"> <a alt="barometer almanac" title="barometer almanac" href="w34barometeralmanac.php" data-lity=""><svg viewBox="0 0 32 32" width="8" height="8" fill="none" stroke="#777" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M14 9 L3 9 3 29 23 29 23 18 M18 4 L28 4 28 14 M28 4 L14 18"></path></svg> Almanac</a></span>
        <span class="todaypopup"> <a alt="Pressure" title="Pressure" href="w34highcharts/dark-charts.html?chart='barometerplot'&amp;span='weekly'&amp;temp='C'&amp;pressure='hPa'&amp;wind='mph'&amp;rain='mm" data-lity=""><svg version="1.1" width="8pt" height="8pt" x="0px" y="0px" viewBox="0 0 496 496" style="enable-background:new 0 0 496 496;" xml:space="preserve"><path style="fill:#719FA3;" d="M496,428c0,6.4-5.6,12-12,12H12c-6.4,0-12-5.6-12-12l0,0c0-6.4,5.6-12,12-12h472 C490.4,416,496,421.6,496,428L496,428z"></path><rect x="24" y="56" style="fill:#1589AD;" width="88" height="344"></rect><polyline style="fill:#04567F;" points="24,56 112,56 112,400 "></polyline><rect x="144" y="128" style="fill:#24966A;" width="88" height="272"></rect><polyline style="fill:#007763;" points="144,128 232,128 232,400 "></polyline><rect x="264" y="208" style="fill:#E8961F;" width="88" height="192"></rect><polyline style="fill:#E57520;" points="264,208 352,208 352,400 "></polyline><rect x="384" y="272" style="fill:#D32A0F;" width="88" height="128"></rect><polyline style="fill:#AF1909;" points="384,272 472,272 472,400 "></polyline><g></g></svg> Pressure </a></span>
    </div>
    <span class="moduletitle">Barometer (<valuetitleunit>{{ baroLabel }}</valuetitleunit>)</span><br>
    <div id="barometer" v-if="minmax && current">
      <div class="updatedtime"><span><svg id="i-info" viewBox="0 0 32 32" width="6" height="6" fill="#9aba2f" stroke="#9aba2f" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M16 14 L16 23 M16 8 L16 10"></path><circle cx="16" cy="16" r="14"></circle></svg>
        {{current.date | now }}</span></div>
      <div class="barometermax">
        <div class="barometerorange">
          <valuetext>Max ({{ minmax.max.day.date | timeFormat}})<br>
            <maxred>
              <value>{{ minmax.max.day.value | baroDisplay($store.getters.units) }}</value>
            </maxred> {{ baroLabel }}
          </valuetext>
        </div>
      </div>
      <div class="barometermin">
        <div class="barometerblue">
          <valuetext>Min ({{ minmax.min.day.date | timeFormat}})<br>
            <minblue>
              <value>{{ minmax.min.day.value | baroDisplay($store.getters.units) }}</value>
            </minblue>&nbsp;{{baroLabel}}
          </valuetext>
        </div>
      </div>

      <div class="barometertrend2">
        <valuetext>&nbsp;&nbsp;Trend<ogreen> <svg id="steadybarometer" width="10" height="10" viewBox="0 0 24 24" fill="none" stroke="currentcolor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="9 18 15 12 9 6"></polyline></svg></ogreen>
          <steady>
            <ogreen>
              <value>Steady</value>
            </ogreen>
          </steady>
        </valuetext>
      </div>

      <div class="homeweathercompass2">
        <div class="homeweathercompass-line2">
          <div class="weather34barometerarrowactual" v-bind:style="{transform:'rotate('+((current.baromrelin) * 13.083)+'deg)'}"></div>
          <div class="weather34barometerarrowmin" v-bind:style="{transform:'rotate('+((minmax.min.day.value) * 13.083)+'deg)'}"></div>
          <div class="weather34barometerarrowmax" v-bind:style="{transform:'rotate('+((minmax.max.day.value) * 13.083)+'deg)'}"></div>
        </div>
        <div class="text2">
          <div class="pressuretext">
            <ogreen>{{ trend.trend}}</ogreen>
          </div>
          <darkgrey>{{ current.baromrelin | baroDisplay($store.getters.units) }} <span>{{ baroLabel }}</span></darkgrey>
        </div>
      </div>

      <div class="barometerconverter">
        <div class="barometerconvertercircleblue">{{ inToMMHG.toFixed(2) }}<smallrainunit>{{ baroLabelAlt }}</smallrainunit>
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

<script>
import moment from "moment";
import axios from 'axios';

export default {
  name: "barometer",
  props: {
    current: Object,
  },
  data () {
    return {
      trend: null,
      minmax: null
    }
  },
  mounted() {
    function updateData(self){
      axios.get('/api/trend/baromrelin').then(response => (self.trend = response.data))
      axios.get('/api/minmax/baromrelin').then(response => (self.minmax = response.data));
      setTimeout(function() { updateData(self); }, 60000);
    }
    updateData(this);
  },
  methods: {},
  filters: {
    toFixed: function(bar) {
      return bar.toFixed(2);
    },
    now: function(date) {
      return moment(date).format("HH:mm:ss");
    },
    timeFormat: function (date) {
      return moment(date).format('HH:mm');
    },
  },
  computed: {
      inToMMHG: function() {
          if(this.$store.getters.units === 'imperial') {
            return this.current.baromrelin * 33.86;
          }
          return this.current.baromrelin;
      }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
