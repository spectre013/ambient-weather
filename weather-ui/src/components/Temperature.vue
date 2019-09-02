<template>
    <div class="weather-item">
        <div class="chartforecast">
            <span class="yearpopup"> <a alt="almanac temperature" title="almanac temperature" href="#" v-on:click="openModal('almanac')">
                <svg viewBox="0 0 32 32" width="8" height="8" fill="none" stroke="#777" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M14 9 L3 9 3 29 23 29 23 18 M18 4 L28 4 28 14 M28 4 L14 18"></path></svg> Almanac </a></span>
            <span class="yearpopup"> <a v-on:click="openModal('chart',{time:'year',type:'temp',field:'tempf'})" alt="Temperature" title="Temperature"  href="#">
                <svg version="1.1" width="8pt" height="8pt" x="0px" y="0px" viewBox="0 0 496 496" style="enable-background:new 0 0 496 496;" xml:space="preserve"><path style="fill:#719FA3;" d="M496,428c0,6.4-5.6,12-12,12H12c-6.4,0-12-5.6-12-12l0,0c0-6.4,5.6-12,12-12h472 C490.4,416,496,421.6,496,428L496,428z"></path><rect x="24" y="56" style="fill:#1589AD;" width="88" height="344"></rect><polyline style="fill:#04567F;" points="24,56 112,56 112,400 "></polyline><rect x="144" y="128" style="fill:#24966A;" width="88" height="272"></rect><polyline style="fill:#007763;" points="144,128 232,128 232,400 "></polyline><rect x="264" y="208" style="fill:#E8961F;" width="88" height="192"></rect><polyline style="fill:#E57520;" points="264,208 352,208 352,400 "></polyline><rect x="384" y="272" style="fill:#D32A0F;" width="88" height="128"></rect><polyline style="fill:#AF1909;" points="384,272 472,272 472,400 "></polyline><g></g></svg>
                {{ 'Y' | chartDates }} </a></span>
            <span class="todaypopup">
                <a alt="Feels Like" title="Feels Like" href="#" data-lity="" v-on:click="openModal('chart',{time:'month',type:'temp',field:'tempf'})">
                <svg version="1.1" width="8pt" height="8pt" x="0px" y="0px" viewBox="0 0 496 496" style="enable-background:new 0 0 496 496;" xml:space="preserve"><path style="fill:#719FA3;" d="M496,428c0,6.4-5.6,12-12,12H12c-6.4,0-12-5.6-12-12l0,0c0-6.4,5.6-12,12-12h472 C490.4,416,496,421.6,496,428L496,428z"></path><rect x="24" y="56" style="fill:#1589AD;" width="88" height="344"></rect><polyline style="fill:#04567F;" points="24,56 112,56 112,400 "></polyline><rect x="144" y="128" style="fill:#24966A;" width="88" height="272"></rect><polyline style="fill:#007763;" points="144,128 232,128 232,400 "></polyline><rect x="264" y="208" style="fill:#E8961F;" width="88" height="192"></rect><polyline style="fill:#E57520;" points="264,208 352,208 352,400 "></polyline><rect x="384" y="272" style="fill:#D32A0F;" width="88" height="128"></rect><polyline style="fill:#AF1909;" points="384,272 472,272 472,400 "></polyline><g></g></svg>
                {{ 'MMM' | chartDates }} </a></span>
            <span class="todaypopup">
                <a alt="Greenhouse" title="Greenhouse" href="#" data-lity="" v-on:click="openModal('chart',{time:'day',type:'temp',field:'tempf'})">
                    <svg version="1.1" width="8pt" height="8pt" x="0px" y="0px" viewBox="0 0 496 496" style="enable-background:new 0 0 496 496;" xml:space="preserve"><path style="fill:#719FA3;" d="M496,428c0,6.4-5.6,12-12,12H12c-6.4,0-12-5.6-12-12l0,0c0-6.4,5.6-12,12-12h472 C490.4,416,496,421.6,496,428L496,428z"></path><rect x="24" y="56" style="fill:#1589AD;" width="88" height="344"></rect><polyline style="fill:#04567F;" points="24,56 112,56 112,400 "></polyline><rect x="144" y="128" style="fill:#24966A;" width="88" height="272"></rect><polyline style="fill:#007763;" points="144,128 232,128 232,400 "></polyline><rect x="264" y="208" style="fill:#E8961F;" width="88" height="192"></rect><polyline style="fill:#E57520;" points="264,208 352,208 352,400 "></polyline><rect x="384" y="272" style="fill:#D32A0F;" width="88" height="128"></rect><polyline style="fill:#AF1909;" points="384,272 472,272 472,400 "></polyline><g></g></svg>
                Today</a></span>

        </div>
        <span class="moduletitle"> Temperature (<valuetitleunit>&deg;F</valuetitleunit>) </span><br>
        <div id="temperature"  v-if="temp && current && trend">
            <div class="updatedtime"><span><svg id="i-info" viewBox="0 0 32 32" width="6" height="6" fill="#9aba2f" stroke="#9aba2f" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M16 14 L16 23 M16 8 L16 10"></path><circle cx="16" cy="16" r="14"></circle></svg>
                {{current.date | now }} </span></div>
            <div class="tempcontainer">
                <div class="maxdata">{{temp.max.day.value}}&deg; | {{ temp.min.day.value }}&deg;</div>
                <div v-bind:class="current.tempf | tempClass">{{ current.tempf }}<smalltempunit>&deg;F</smalltempunit></div>
                <div class="temptrendx">
                    <trendmovementfallingx v-if="trend.trend == 'down'">&nbsp;&nbsp;&nbsp;
                        <valuetext>Trend <svg id="falling" width="9" height="9" viewBox="0 0 24 24"><polyline points="23 18 13.5 8.5 8.5 13.5 1 6" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"></polyline><polyline points="17 18 23 18 23 12" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"></polyline></svg>
                            {{ trend.by }}&deg;</valuetext>
                    </trendmovementfallingx>
                    <trendmovementfallingx v-if="trend.trend == 'up'">&nbsp;&nbsp;&nbsp;
                        <valuetext>Trend <svg id="falling" width="9" height="9" viewBox="0 0 24 24"><polyline points="23 18 13.5 8.5 8.5 13.5 1 6" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"></polyline><polyline points="17 18 23 18 23 12" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"></polyline></svg>
                            {{ trend.by }}&deg;</valuetext>
                    </trendmovementfallingx>
                </div>
            </div>
            <div class="heatcircle">
                <div class="heatcircle-content">
                    <valuetextheading1>Feels</valuetextheading1><br>
                    <div class="tempconverter1">
                        <div v-bind:class="current.feelsLike | smallTempClass">{{ current.feelsLike }}<smalltempunit2>&deg;F<smalltempunit2></smalltempunit2>
                        </smalltempunit2>
                        </div>
                    </div>
                </div>
                <div class="heatcircle2">
                    <div class="heatcircle-content">
                        <valuetextheading1>Avg Today</valuetextheading1>
                        <div class="tempconverter1">
                            <div v-bind:class="temp.avg.day.value | smallTempClass">{{ temp.avg.day.value }}<smalltempunit2>&deg;F</smalltempunit2>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="heatcircle3">
                    <div class="heatcircle-content">
                        <valuetextheading1>Humidity</valuetextheading1>
                        <div class="tempconverter1">
                            <div v-bind:class="current.humidity | humidityClass">{{ current.humidity }}<smalltempunit2>%</smalltempunit2>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="heatcircle4">
                    <div class="heatcircle-content">
                        <valuetextheading1>Dewpoint</valuetextheading1>
                        <div class="tempconverter1">
                            <div v-bind:class="current.dewPoint | dewPointClass">&nbsp;{{ current.dewPoint }}<smalltempunit2>&deg;F</smalltempunit2></div>
                        </div>
                    </div>
                </div>
            </div>
            <div class="tempconverter2">
                <div v-bind:class="current.tempf | tempCircle">{{ current.tempf | toCelcius }}<smalltempunit2>&deg;C<smalltempunit2></smalltempunit2>
                </smalltempunit2>
                </div>
            </div>
        </div>
        <br>
    </div>
</template>

<script>
import moment from 'moment';
import axios from 'axios';

export default {
  name: 'temp',
  props: {
    current: Object,
      temp: Object,
  },
  data () {
    return {
      trend: null,
    }
  },
  mounted() {
      function updateData(self){
          axios.get('/api/trend/temp').then(response => (self.trend = response.data));
          setTimeout(function() { updateData(self); }, 60000);
      }
      updateData(this);

  },
  methods: {
    openModal: function(type,options) {
         this.$parent.showModal(type,options);
    },
      toCelcius: function(temp) {
          return ((temp -32) * 5/9).toFixed(2);
      },
    feelslike: function() {
          if(this.current.tempf < 35) {
              return "maxcircleblue";
          } else if(this.current.tempf < 50) {
              return "maxcirclegreen";
          } else if(this.current.tempf < 70) {
              return "maxcircleyellow";
          } else if(this.current.tempf < 80) {
              return "maxcircleorange";
          } else if(this.current.tempf > 80) {
              return "maxcirclered";
          }
      },
      svgfeels: function() {
          if(this.current.tempf < 35) {
              return "maxcircleblue";
          } else if(this.current.tempf < 50) {
              return "maxcirclegreen";
          } else if(this.current.tempf < 70) {
              return "maxcircleyellow";
          } else if(this.current.tempf < 80) {
              return "maxcircleorange";
          } else if(this.current.tempf > 80) {
              return "maxcirclered";
          }
      },
  },
 filters: {
    chartDates: function(format) {
      return moment().format(format);
    },
    now: function (date) {
        return moment(date).format('HH:mm:ss');
    },
    toCelcius: function(temp) {
        return ((temp -32) * 5/9).toFixed(2);
    },
  },
  computed: {

  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
