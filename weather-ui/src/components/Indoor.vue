<template>
    <div class="weather-item">
        <div class="chartforecast">
            <span class="yearpopup"> <a alt="almanac temperature" title="almanac temperature" href="w34tempalmanac.php" data-lity="">
                <svg viewBox="0 0 32 32" width="8" height="8" fill="none" stroke="#777" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M14 9 L3 9 3 29 23 29 23 18 M18 4 L28 4 28 14 M28 4 L14 18"></path></svg> Almanac </a></span>
            <span class="yearpopup"> <a alt="Temperature" title="Temperature"  href="w34highcharts/dark-charts.html?chart='tempallplot'&amp;span='weekly'&amp;temp='C'&amp;pressure='hPa'&amp;wind='mph'&amp;rain='mm" data-lity="">
                <svg version="1.1" width="8pt" height="8pt" x="0px" y="0px" viewBox="0 0 496 496" style="enable-background:new 0 0 496 496;" xml:space="preserve"><path style="fill:#719FA3;" d="M496,428c0,6.4-5.6,12-12,12H12c-6.4,0-12-5.6-12-12l0,0c0-6.4,5.6-12,12-12h472 C490.4,416,496,421.6,496,428L496,428z"></path><rect x="24" y="56" style="fill:#1589AD;" width="88" height="344"></rect><polyline style="fill:#04567F;" points="24,56 112,56 112,400 "></polyline><rect x="144" y="128" style="fill:#24966A;" width="88" height="272"></rect><polyline style="fill:#007763;" points="144,128 232,128 232,400 "></polyline><rect x="264" y="208" style="fill:#E8961F;" width="88" height="192"></rect><polyline style="fill:#E57520;" points="264,208 352,208 352,400 "></polyline><rect x="384" y="272" style="fill:#D32A0F;" width="88" height="128"></rect><polyline style="fill:#AF1909;" points="384,272 472,272 472,400 "></polyline><g></g></svg> Temp Dew</a></span>

        </div>
        <span class="moduletitle" v-if="title"> {{title }} (<valuetitleunit>&deg;F</valuetitleunit>) </span><br>
        <div id="temperature"  v-if="live && minmax && trend">
            <div class="updatedtime"><span><svg id="i-info" viewBox="0 0 32 32" width="6" height="6" fill="#9aba2f" stroke="#9aba2f" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M16 14 L16 23 M16 8 L16 10"></path><circle cx="16" cy="16" r="14"></circle></svg>
                {{live.date | now }} </span></div>
            <div class="tempcontainer">
                <div class="maxdata">{{ minmax.max.day.value }}&deg; | {{ minmax.min.day.value }}&deg;</div>
                <div class="maxdata"></div>
                <div v-bind:class="temp| tempClass">{{ temp }}<smalltempunit>&deg;F</smalltempunit></div>
                <div class="temptrendx"></div>
            </div>
            <div class="heatcircle">
                <div class="heatcircle-content">
                    <valuetextheading1>Humidity</valuetextheading1><br>
                    <div class="tempconverter1">
                        <div v-bind:class="humidity | smallTempClass">{{ humidity }}<smalltempunit2>&deg;F<smalltempunit2></smalltempunit2>
                        </smalltempunit2>
                        </div>
                    </div>
                </div>
                <div class="heatcircle2">
                    <div class="heatcircle-content">
                        <valuetextheading1>Avg Today</valuetextheading1>
                        <div class="tempconverter1">
                            <div v-bind:class="minmax.avg.day.value | smallTempClass">{{ minmax.avg.day.value }}<smalltempunit2>&deg;F</smalltempunit2>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="heatcircle3">
                    <div class="heatcircle-content">
                        <valuetextheading1>High</valuetextheading1>
                        <div class="tempconverter1">
                            <div v-bind:class="minmax.max.day.value | smallTempClass">{{ minmax.max.day.value }}<smalltempunit2>&deg;F</smalltempunit2>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="heatcircle4">
                    <div class="heatcircle-content">
                        <valuetextheading1>Low</valuetextheading1>
                        <div class="tempconverter1">
                            <div v-bind:class="minmax.min.day.value | smallTempClass">&nbsp;{{ minmax.min.day.value }}<smalltempunit2>&deg;F</smalltempunit2>&nbsp;</div>
                        </div>
                    </div>
                </div>
            </div>
            <div class="tempconverter2">
                <div v-bind:class="temp | tempCircle">{{ temp | toCelcius }}<smalltempunit2>&deg;C<smalltempunit2></smalltempunit2>
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
    name: 'indoor',
        props: {
            title: String,
            loc: String,
            live: Object
    },
    data () {
        return {
            trend: null,
            minmax: null,
        }
    },
    mounted() {
        axios.get('/api/trend/temp').then(response => (this.trend = response.data))
        axios.get('/api/minmax/temp'+this.loc+'f').then(response => (this.minmax = response.data));
    },
    methods: {

    },
    filters: {
        now: function(date) {
            return moment(date).format("HH:mm:ss");
        },
        toCelcius: function(temp) {
            return ((temp -32) * 5/9).toFixed(2);
        },
    },
    computed: {
        temp: function() {
            let key = 'temp'+this.loc.toLowerCase()+"f";
            return this.live[key];
        },
        humidity: function() {
            let key = 'humidity'+this.loc;
            return this.live[key];
        },
        tempClass: function() {
            let temp =  this.live['temp'+this.loc+"f"];
            if(temp < 14) {
              return "outsideminus10";
            }
            else if(temp <= 23) {
              return "outsideminus5";
            }
            else if(temp <= 32) {
              return "outsidezero";
            }
            else if(temp < 41) {
              return "outside0-5";
            }
            else if(temp < 50) {
              return "outside6-10";
            }
            else if(temp < 59) {
              return "outside11-15";
            }
            else if(temp < 68) {
              return "outside16-20";
            }
            else if(temp < 77) {
              return "outside21-25";
            }
            else if(temp < 86) {
              return "outside26-30";
            }
            else if(temp < 95) {
              return "outside31-35";
            }
            else if(temp < 104) {
              return " outside36-40";
            }
            else if(temp < 113) {
              return " outside41-45";
            }
            else if(temp < 150) {
              return " outside50";
            }
            return 'outsideminus10';
            }
        },
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
    #home {
        position: relative;
        top: 10px;
        left:5px;
    }
</style>
