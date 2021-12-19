<template>
  <div class="weather-item">
    <div class="chartforecast"></div>
    <span class="moduletitle"
      >Barometer (<valuetitleunit>{{ baroLabel }}</valuetitleunit
      >)</span
    ><br />
    <div id="barometer" v-if="minmax && current">
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
          {{ current.date | now }}</span
        >
      </div>
      <div class="barometermax">
        <div class="barometerorange">
          <valuetext
            >Max ({{ minmax.max.day.date | timeFormat }})<br />
            <maxred>
              <value>{{ minmax.max.day.value | baroDisplay($store.getters.units) }}</value>
            </maxred>
            {{ baroLabel }}
          </valuetext>
        </div>
      </div>
      <div class="barometermin">
        <div class="barometerblue">
          <valuetext
            >Min ({{ minmax.min.day.date | timeFormat }})<br />
            <minblue>
              <value>{{
                minmax.min.day.value | baroDisplay($store.getters.units)
              }}</value> </minblue
            >&nbsp;{{ baroLabel }}
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
            v-bind:style="{ transform: 'rotate(' + current.baromrelin * 13.083 + 'deg)' }"
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
          <darkgrey
            >{{ current.baromrelin | baroDisplay($store.getters.units) }}
            <span>{{ baroLabel }}</span></darkgrey
          >
        </div>
      </div>

      <div class="barometerconverter">
        <div class="barometerconvertercircleblue">
          {{ inToMMHG.toFixed(2) }}<smallrainunit>{{ baroLabelAlt }}</smallrainunit>
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
import moment from 'moment';
import axios from 'axios';

export default {
  name: 'barometer',
  props: {
    current: Object,
  },
  data() {
    return {
      trend: Object,
      minmax: Object,
    };
  },
  mounted() {
    function updateData(self) {
      axios.get('/api/trend/baromrelin').then((response) => (self.trend = response.data));
      axios.get('/api/minmax/baromrelin').then((response) => (self.minmax = response.data));
      setTimeout(function () {
        updateData(self);
      }, 60000);
    }
    updateData(this);
  },
  methods: {},
  filters: {
    toFixed: function (bar) {
      return bar.toFixed(2);
    },
    now: function (date) {
      return moment(date).format('HH:mm:ss');
    },
    timeFormat: function (date) {
      return moment(date).format('HH:mm');
    },
  },
  computed: {
    inToMMHG: function () {
      if (this.$store.getters.units === 'imperial') {
        return this.current.baromrelin * 33.86;
      }
      return this.current.baromrelin;
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped></style>
