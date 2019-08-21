<template>
    <div class="weather-item">
        <div class="chartforecast">
            <!-- HOURLY & Outlook for homeweather station-->
            <span class="yearpopup"> <a alt="nearby metar station" title="nearby metar station" href="metarnearby.php" data-lity=""><svg viewBox="0 0 32 32" width="8" height="8" fill="none" stroke="#777" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M14 9 L3 9 3 29 23 29 23 18 M18 4 L28 4 28 14 M28 4 L14 18"></path></svg> Nearby Metar </a></span><span class="monthpopup"><a href="windyradar.php" title="Windy.com Radar" alt="Windy.com Radar" data-lity=""><svg viewBox="0 0 32 32" width="8" height="8" fill="none" stroke="#777" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M14 9 L3 9 3 29 23 29 23 18 M18 4 L28 4 28 14 M28 4 L14 18"></path></svg> Radar</a></span>
        </div>
            <span class="moduletitle">Current Conditions</span><br>
        <div id="currentsky" v-if="forecast">
            <div class="updatedtimecurrent">
                <svg id="i-info" viewBox="0 0 32 32" width="6" height="6" fill="#9aba2f" stroke="#9aba2f" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M16 14 L16 23 M16 8 L16 10"></path><circle cx="16" cy="16" r="14"></circle></svg>
                {{forecast.currently.time | now }}</div>
            <div class="darkskyiconcurrent">
                <span1><img rel="prefetch" :src="'css/darkskyicons/'+forecast.currently.icon+'.svg'" width="70px" height="60px" alt="weather34 icon"></span1>
            </div>
            <div class="darkskysummary"><span>
                <uppercase>{{ forecast.currently.summary}}<br>Conditions</uppercase>
            </span></div>
            <!-- Darksky Hourly Next 3 hours HOMEWEATHER STATION-->
            <div class="darkskynexthours">

                Hourly Forecast <br> Temperature <oorange v-if="hourly.temperature >=65">{{ hourly.temperature }}&deg;F</oorange>
                <ogreen v-if="this.hourly.temperature >39 && this.hourly.temperature < 64">{{ hourly.temperature }}&deg;F</ogreen>
                <oblue v-if="hourly.temperature < 39">{{ hourly.temperature }}&deg;F</oblue> <ogreen> {{ hourly.summary }} </ogreen> Conditions.
                <oorange><br>Wind Gust</oorange>
                <ogreen v-if="this.hourly.windGust >=0 && this.hourly.windGust < 30"> {{ hourly.windGust }}</ogreen>
                <oorange v-if="this.hourly.windGust >=30"> {{ hourly.windGust }}</oorange> mph. Rainfall <svg id="weather34 raindrop" x="0px" y="0px" viewBox="0 0 512 512" width="8px" fill="#01a4b5" stroke="#01a4b5" stroke-width="3%"><g><g><path d="M348.242,124.971C306.633,58.176,264.434,4.423,264.013,3.889C262.08,1.433,259.125,0,256,0	c-3.126,0-6.079,1.433-8.013,3.889c-0.422,0.535-42.621,54.287-84.229,121.083c-56.485,90.679-85.127,161.219-85.127,209.66 C78.632,432.433,158.199,512,256,512c97.802,0,177.368-79.567,177.368-177.369C433.368,286.19,404.728,215.65,348.242,124.971z M256,491.602c-86.554,0-156.97-70.416-156.97-156.97c0-93.472,123.907-263.861,156.971-307.658 C289.065,70.762,412.97,241.122,412.97,334.632C412.97,421.185,342.554,491.602,256,491.602z"></path></g></g><g><g><path d="M275.451,86.98c-1.961-2.815-3.884-5.555-5.758-8.21c-3.249-4.601-9.612-5.698-14.215-2.45 c-4.601,3.249-5.698,9.613-2.45,14.215c1.852,2.623,3.75,5.328,5.688,8.108c1.982,2.846,5.154,4.369,8.377,4.369 c2.012,0,4.046-0.595,5.822-1.833C277.536,97.959,278.672,91.602,275.451,86.98z"></path></g></g><g><g><path d="M362.724,231.075c-16.546-33.415-38.187-70.496-64.319-110.213c-3.095-4.704-9.42-6.01-14.126-2.914 c-4.706,3.096-6.01,9.421-2.914,14.127c25.677,39.025,46.9,75.379,63.081,108.052c1.779,3.592,5.391,5.675,9.148,5.675 c1.521,0,3.064-0.342,4.517-1.062C363.159,242.241,365.224,236.123,362.724,231.075z"></path></g></g></svg>
                {{ hourly.precipProbability }}%<oblue> {{ hourly.precipIntensity.toFixed(2) }}</oblue> in<br>
                <oorange>UVI</oorange> Forecast <uviforecasthourred v-if="this.hourly.uvIndex >=8"> {{ hourly.uvIndex }}</uviforecasthourred>
                                                <uviforecasthourorange v-if="this.hourly.uvIndex >=6 && this.hourly.uvIndex <= 7"> {{ hourly.uvIndex }}</uviforecasthourorange>
                                                <uviforecasthouryellow v-if="this.hourly.uvIndex >=3 && this.hourly.uvIndex <= 6"> {{ hourly.uvIndex }}</uviforecasthouryellow>
                                                <uviforecasthourgreen v-if="this.hourly.uvIndex <3"> {{ hourly.uvIndex }}</uviforecasthourgreen>
            </div>
        </div>
    </div>
<!--    <div class="weather-item">-->
<!--        <div class="chartforecast">-->
<!--            &lt;!&ndash; HOURLY & Outlook for homeweather station&ndash;&gt;-->
<!--            <span class="yearpopup"> <a href="metarnearby.php" data-featherlight="iframe">-->
<!--                <svg viewBox="0 0 32 32" width="10" height="10" fill="none" stroke="currentcolor" stroke-linecap="round"-->
<!--                     stroke-linejoin="round" stroke-width="6.25%">-->
<!--                    <path d="M14 9 L3 9 3 29 23 29 23 18 M18 4 L28 4 28 14 M28 4 L14 18"></path>-->
<!--                </svg> Nearby Metar</a>-->
<!--        </span>-->

<!--        </div>-->
<!--        <span class="moduletitle">Hourly Forecast</span><br>-->
<!--        <div id="currentsky" v-if="forecast">-->
<!--            <div class="updatedtimecurrent">-->
<!--                <svg id="i-info" viewBox="0 0 32 32" width="7" height="7" fill="#9aba2f" stroke="#9aba2f"-->
<!--                     stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%">-->
<!--                    <path d="M16 14 L16 23 M16 8 L16 10"></path>-->
<!--                    <circle cx="16" cy="16" r="14"></circle>-->
<!--                </svg> {{forecast.currently.time | now }}</div>-->
<!--            <div class="darkskyiconcurrent">-->
<!--                <span1>-->
<!--                    <img rel="prefetch" :src="'css/darkskyicons/'+forecast.currently.icon+'.svg'" width="75px"></span1>-->
<!--            </div>-->
<!--            <div class="darkskysummary"><span>-->
<!--                <uppercase>{{ forecast.currently.summary}}</uppercase>-->
<!--            </span></div>-->
<!--            &lt;!&ndash; Darksky Hourly Next 3 hours HOMEWEATHER STATION&ndash;&gt;-->
<!--            <div class="darkskynexthours" v-if="hourly">-->
<!--                <br>-->
<!--                Hourly Forecast <br> Temperature <span v-bind:class="temp">{{ hourly.temperature }}</span>&deg;F<ogreen> {{ hourly.summary }} </ogreen>Conditions.-->
<!--                <oorange><br>Wind Gust</oorange>-->
<!--                <span v-bind:class="wind">{{ hourly.windGust }}</span> mph. Rainfall <svg id="weather34 raindrop" x="0px" y="0px" viewBox="0 0 512 512" width="8px" fill="#01a4b5" stroke="#01a4b5" stroke-width="3%"><g><g><path d="M348.242,124.971C306.633,58.176,264.434,4.423,264.013,3.889C262.08,1.433,259.125,0,256,0	c-3.126,0-6.079,1.433-8.013,3.889c-0.422,0.535-42.621,54.287-84.229,121.083c-56.485,90.679-85.127,161.219-85.127,209.66 C78.632,432.433,158.199,512,256,512c97.802,0,177.368-79.567,177.368-177.369C433.368,286.19,404.728,215.65,348.242,124.971z M256,491.602c-86.554,0-156.97-70.416-156.97-156.97c0-93.472,123.907-263.861,156.971-307.658 C289.065,70.762,412.97,241.122,412.97,334.632C412.97,421.185,342.554,491.602,256,491.602z"></path></g></g><g><g><path d="M275.451,86.98c-1.961-2.815-3.884-5.555-5.758-8.21c-3.249-4.601-9.612-5.698-14.215-2.45 c-4.601,3.249-5.698,9.613-2.45,14.215c1.852,2.623,3.75,5.328,5.688,8.108c1.982,2.846,5.154,4.369,8.377,4.369 c2.012,0,4.046-0.595,5.822-1.833C277.536,97.959,278.672,91.602,275.451,86.98z"></path></g></g><g><g><path d="M362.724,231.075c-16.546-33.415-38.187-70.496-64.319-110.213c-3.095-4.704-9.42-6.01-14.126-2.914 c-4.706,3.096-6.01,9.421-2.914,14.127c25.677,39.025,46.9,75.379,63.081,108.052c1.779,3.592,5.391,5.675,9.148,5.675 c1.521,0,3.064-0.342,4.517-1.062C363.159,242.241,365.224,236.123,362.724,231.075z"></path></g></g>-->
<!--                </svg> {{ hourly.precipProbability }}%<oblue> {{ hourly.precipIntensity }}</oblue> in<br>-->
<!--                <oorange>UVI</oorange> Forecast <span v-bind:class='uvClass'> {{ hourly.uvIndex }} </span>-->
<!--                <svg id="weather34 uvdark" width="14px" height="14px" viewBox="0 0 200 200" xmlns="http://www.w3.org/2000/svg">-->
<!--                    <path fill="#ccc" stroke="#ccc" stroke-width=".1"-->
<!--                          d="M147.5 5h8v29.2h-8V5zM96.6 34.5l6-6 16.8 17c-2 1.8-3.8 4.2-6.2 5.7-5.4-5.7-11-11-16.6-16.7zM184.6 45.4l17-16.8 5.8 6-17 16.8-5.8-6zM143.6 46.8c8.3-1.5 17-.7 24.6 3 11 5.2 19 15.8 21.2 27.7 1.3 8 .4 16.5-3.2 23.8-5 10.6-15.3 18.6-27 21-8.3 1.5-17.4.6-25-3.5-10.8-5.3-18.7-16-20.6-27.8-1.2-7.8-.3-16 3.2-23.3 5-10.5 15.3-18.6 26.8-21zM72 80.5h29.2v8H72v-8zM201.8 80.5H231v8h-29.2v-8zM96.6 133.5l17-17 5.8 6-17 17-5.8-6zM184.6 122.5l6-6 16.8 17-6 6c-5.5-5.8-11.2-11.4-16.8-17zM147.5 134.8h8V164h-8v-29.2z">-->
<!--                    </path>-->
<!--                </svg>-->
<!--            </div>-->
<!--        </div>-->
<!--    </div>-->
</template>

<script>
    import moment from 'moment';

    export default {
        name: 'current',
        props: {
            forecast: Object
        },
        data() {
            return {
                hourly: null,
            }
        },
        watch: {
          forecast: function() {
              this.hourly = this.forecast.hourly.data[0];
          }
        },
        mounted() {

        },
        methods: {

        },
        filters: {
            now: function (date) {
                return moment(date).format('HH:mm:ss');
            }
        },
        computed: {
            temp: function() {
                if(this.hourly.temperature >=65) {
                    return 'oorange';
                } else if(this.hourly.temperature <= 40) {
                    return 'oblue';
                } else if(this.hourly.temperature < 65){
                    return 'ogreen';
                }
                return "";
            },
            wind: function() {
                if(this.hourly.windGust >= 40) {
                    return 'oorange';
                } else if( this.hourly.windGust >=0 ){
                    return 'ogreen';
                }
                return "";
            },
            uvClass: function() {
                if(this.hourly.uvIndex >= 8) {
                    return 'uviforecasthourred';
                } else if(this.hourly.uvIndex >= 6) {
                    return 'uviforecasthourorange';
                } else if(this.hourly.uvIndex >= 3) {
                    return 'uviforecasthouryellow';
                } else if(this.hourly.uvIndex >= 0.1) {
                    return 'uviforecasthourgreen';
                }
                return "";
            }
        },
    }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
    uppercase{ text-transform:capitalize;}
</style>
