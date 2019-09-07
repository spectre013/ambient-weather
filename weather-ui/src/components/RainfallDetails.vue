<template>
    <div class="weather-item">
        <div class="chartforecast">
            <span class="yearpopup"> <a alt="almanac rainfall" title="almanac rainfall" href="w34rainfallalmanac.php"><svg viewBox="0 0 32 32" width="8" height="8" fill="none" stroke="#777" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M14 9 L3 9 3 29 23 29 23 18 M18 4 L28 4 28 14 M28 4 L14 18"></path></svg> Almanac </a></span>
            <span class="todaypopup"> <a alt="Rainfall" title="Rainfall" href="w34highcharts/dark-charts.html?chart='rainplot'&amp;span='weekly'&amp;temp='C'&amp;pressure='hPa'&amp;wind='mph'&amp;rain='mm" data-lity=""><svg version="1.1" width="8pt" height="8pt" x="0px" y="0px" viewBox="0 0 496 496" style="enable-background:new 0 0 496 496;" xml:space="preserve"><path style="fill:#719FA3;" d="M496,428c0,6.4-5.6,12-12,12H12c-6.4,0-12-5.6-12-12l0,0c0-6.4,5.6-12,12-12h472 C490.4,416,496,421.6,496,428L496,428z"></path><rect x="24" y="56" style="fill:#1589AD;" width="88" height="344"></rect><polyline style="fill:#04567F;" points="24,56 112,56 112,400 "></polyline><rect x="144" y="128" style="fill:#24966A;" width="88" height="272"></rect><polyline style="fill:#007763;" points="144,128 232,128 232,400 "></polyline><rect x="264" y="208" style="fill:#E8961F;" width="88" height="192"></rect><polyline style="fill:#E57520;" points="264,208 352,208 352,400 "></polyline><rect x="384" y="272" style="fill:#D32A0F;" width="88" height="128"></rect><polyline style="fill:#AF1909;" points="384,272 472,272 472,400 "></polyline><g></g></svg> Rainfall </a></span>
        </div>
        <span class="moduletitle">Rainfall Today (<valuetitleunit>{{ rainLabel }}</valuetitleunit>)</span><br>
        <div id="rainfall" v-if="current">
            <div class="updatedtime"><span><svg id="i-info" viewBox="0 0 32 32" width="6" height="6" fill="#9aba2f" stroke="#9aba2f" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M16 14 L16 23 M16 8 L16 10"></path><circle cx="16" cy="16" r="14"></circle></svg>
                {{current.date | now }}</span>
            </div>
            <div class="weather34i-rairate-bar">
                <div id="raincontainer">
                    <div id="weather34rainbeaker">
                        <div id="weather34rainwater" v-bind:style="rainLevel">
                        </div>
                    </div>
                </div>
            </div>
            <div class="raincontainer1">
                <div class="raintoday1">{{ current.dailyrainin | rainDisplay($store.getters.units)}}<sup><smallrainunita>{{ rainLabel }}</smallrainunita></sup></div>
            </div>
            <div class="heatcircle">
                <div class="heatcircle-content">
                    &nbsp;<valuetextheading1>{{ current.date | year }}</valuetextheading1> <br>
                    <div class="tempconverter1">
                        <div class="rainmodulehome">
                            <raiblue>{{ current.yearlyrainin | rainDisplay($store.getters.units)}}</raiblue>
                            <smallrainunit2>{{ rainLabel }}<smallrainunit2></smallrainunit2>
                            </smallrainunit2>
                        </div>
                    </div>
                </div>
                <div class="heatcircle2">
                    <div class="heatcircle-content">
                        &nbsp;&nbsp;&nbsp;<valuetextheading1>{{ current.date | month }}</valuetextheading1> <br>
                        <div class="tempconverter1">
                            <div class="rainmodulehome">
                                <raiblue>{{ current.monthlyrainin | rainDisplay($store.getters.units)}}</raiblue>
                                <smallrainunit2>{{ rainLabel }}</smallrainunit2>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="heatcircle3">
                    <div class="heatcircle-content">
                        &nbsp;&nbsp;<valuetextheading1>Last Hour</valuetextheading1><br>
                        <div class="tempconverter1">
                            <div class="rainmodulehome">
                                <raiblue>{{ current.hourlyrain | rainDisplay($store.getters.units)}}</raiblue>
                                <smallrainunit2>{{ rainLabel }}</smallrainunit2>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="heatcircle4">
                    <div class="heatcircle-content">
                        &nbsp;&nbsp;<valuetextheading1>Last 24hr</valuetextheading1><br>
                        <div class="tempconverter1">
                            <div class="rainmodulehome">
                                <raiblue>{{ current.dailyrainin | rainDisplay($store.getters.units)}}</raiblue>
                                <smallrainunit2>{{ rainLabel }}</smallrainunit2>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            <div class="rainconverter">
                <div class="rainconvertercircle"><raiblue>Last Rain:</raiblue> {{ current.lastrain | full }}<smallrainunit>
                </smallrainunit>
                </div>
            </div>
            <div class="rainrateextra">
                <div class="rainratemodulehome">
                    <rainratetextheading>&nbsp;Rate&nbsp;</rainratetextheading>
                    <raiblue>{{ current.hourlyrainin | rainDisplay($store.getters.units)}}</raiblue>
                    <smallrainunit2>{{ rainLabel }}</smallrainunit2>
                </div>
            </div>
        </div>
    </div>
</template>


<script>
import moment from 'moment';

export default {
  name: 'raindetails',
  props: {
    current: Object
  },
  mounted() {
   
  },
  methods: {
   
  },
  filters: {
     year: function (date) {
      return moment(date).format('YYYY');
    },
    month: function (date) {
      return moment(date).format('MMM');
    },
    full: function(date) {
        return moment(date).format('YYYY-MM-DD HH:mm:ss');
    },
    now: function (date) {
        return moment(date).format('HH:mm:ss');
    }
  },
  computed: {
    rainLevel: function() {
        let rain = this.current.dailyrainin;
        if(rain > 0) {
            return {height:rain * 50.8+'px'};
        }
        return {height:0+'px'};
    }
  },
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
