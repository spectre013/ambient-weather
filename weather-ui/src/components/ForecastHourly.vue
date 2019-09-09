<template>
    <div class="modal-backdrop">
        <div class="modal">
            <section class="modal-body">
                <div id="body" v-if="forecast">
                    <div class="weather34darkbrowser" :url="title"></div>
                    <main class="grid">
                        <article>
                            <div class="darkskyforecasthome">
                                <div class="darkskydiv">
                                    <div class="darkskyforecastinghome" v-for="hour in forecast.hourly.data.slice(0,18)" v-bind:key="hour.time">
                                        <div class="darkskyweekdayhome">{{ hour.time | hourDisplay }}</div>
                                        <darkskytemphihome><img src="css/icons/temp34.svg" width="10px"><span>{{hour.temperature | tempDisplay($store.getters.units) }}°<grey></grey> &nbsp;</span>
                                            <grey>&nbsp;UVI:</grey> <orange>{{ hour.uvIndex }} </orange>
                                        </darkskytemphihome><br>
                                        <darkskyiconcurrent><img :src="'css/darkskyicons/'+hour.icon+'.svg'" width="50"></darkskyiconcurrent>
                                        <darkskytempwindhome>
                                            <span2 style="color:#ff8841;"><img src='css/windicons/avgw.svg' width='20' v-bind:style="{transform:'rotate('+hour.windBearing+'deg)'}"></span2>
                                            <span4> {{ hour.windSpeed | windDisplay($store.getters.speed) }} | <gust>{{ hour.windGust.toFixed(0) | windDisplay($store.getters.speed) }}</gust></span4>
                                            <windunit> {{ windLabel }}</windunit>
                                        </darkskytempwindhome><br>&nbsp;
                                        <darkskyrainhome>{{ hour.summary}}</darkskyrainhome><br>
                                            <darkskyrainhome1>{{hour.precipProbability.toFixed(2) }}% <blue1><svg id="weather34 raindrop" x="0px" y="0px" viewBox="0 0 512 512" width="8px" fill="#01a4b5" stroke="#01a4b5" stroke-width="3%"><g><g><path d="M348.242,124.971C306.633,58.176,264.434,4.423,264.013,3.889C262.08,1.433,259.125,0,256,0 c-3.126,0-6.079,1.433-8.013,3.889c-0.422,0.535-42.621,54.287-84.229,121.083c-56.485,90.679-85.127,161.219-85.127,209.66 C78.632,432.433,158.199,512,256,512c97.802,0,177.368-79.567,177.368-177.369C433.368,286.19,404.728,215.65,348.242,124.971z M256,491.602c-86.554,0-156.97-70.416-156.97-156.97c0-93.472,123.907-263.861,156.971-307.658 C289.065,70.762,412.97,241.122,412.97,334.632C412.97,421.185,342.554,491.602,256,491.602z" /></g></g><g><g><path d="M275.451,86.98c-1.961-2.815-3.884-5.555-5.758-8.21c-3.249-4.601-9.612-5.698-14.215-2.45 c-4.601,3.249-5.698,9.613-2.45,14.215c1.852,2.623,3.75,5.328,5.688,8.108c1.982,2.846,5.154,4.369,8.377,4.369 c2.012,0,4.046-0.595,5.822-1.833C277.536,97.959,278.672,91.602,275.451,86.98z" /></g></g><g><g><path d="M362.724,231.075c-16.546-33.415-38.187-70.496-64.319-110.213c-3.095-4.704-9.42-6.01-14.126-2.914 c-4.706,3.096-6.01,9.421-2.914,14.127c25.677,39.025,46.9,75.379,63.081,108.052c1.779,3.592,5.391,5.675,9.148,5.675 c1.521,0,3.064-0.342,4.517-1.062C363.159,242.241,365.224,236.123,362.724,231.075z" /></g></g></svg>
                                                {{ hour.precipIntensity.toFixed(2) | rainDisplay($store.getters.units) }}</blue1>
                                                <unit> {{ rainLabel }}</unit>
                                            </darkskyrainhome1>
                                    </div>
                                </div>
                            </div>
                        </article>
                        <div class="weather34browser-footer">
                            <span class="copyw">&nbsp;
                                <svg id="i-external" viewBox="0 0 32 32" width="10" height="10" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M14 9 L3 9 3 29 23 29 23 18 M18 4 L28 4 28 14 M28 4 L14 18"></path></svg>
                                <a href="https://darksky.net/poweredby/" title="https://darksky.net/poweredby/" target="_blank">Powered by Dark Sky </a>
                            </span>
                        </div>
                    </main>
                </div>
                <button class="lity-close" type="button" aria-label="Close (Press escape to close)" @click="close">×</button>
            </section>
        </div>
    </div>
</template>


<script>
    import moment from 'moment';
    export default {
        name: 'forecasthourlymodel',
        props: {
            forecast: Object
        },
        methods: {
            close: function() {
                this.$parent.closeModal('forecasthourly');
            }
        },
        mounted() {

        },
        computed: {
            title: function() {
                let from = moment.unix(this.forecast.hourly.data[0].time).format("ddd MMM Do HH:mm");
                let to = moment.unix(this.forecast.hourly.data[15].time).format("ddd MMM Do HH:mm");
                return `Hourly Forecast - ${from } - ${to}`;
            }

        },
        filters: {
            hourDisplay: function (date) {
                return moment.unix(date).format('HH:mm');
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
        }
    };
</script>

<style scoped>
    article {
        margin: auto 5px;
    }

    .weather34browser-footer {
        margin-top:10px;
        margin-left:45px;
        font-size: .6rem;
        color: silver;
    }

    .darkskyforecasthome {
         float: left;
         display: block;
         margin: 5px 0 2px 0;
         width: 95%;
         border-radius: 4px;
         font-family: Arial, monospace;
         height: 480px;
         padding: 5px;
         color: #777
     }


    darkskyweekday {
        position: absolute;
        margin: 3px;
        background-color: #494a4b;
        text-align: center;
        padding: 5px;
        color: #aaa;
        font-family: Arial;
        font-size: 12px;
        margin-bottom: 15px;
        border-radius: 4px
    }

    darkskytemphi {
        margin-top: 5px;
        font-size: 13px;
        color: #ff8841;
        font-family: Arial;
        margin-left: 10%
    }

    darkskytemphi span {
        font-size: 13px;
        color: #aaa
    }

    darkskytemplo {
        margin-top: 5px;
        font-size: 13px;
        color: #00a4b4;
        font-family: Arial
    }

    darkskytemplo span {
        font-size: 13px;
        color: #aaa;
        font-family: Arial
    }

    darkskysummary {
        font-size: 12px;
        color: #aaa;
        font-family: Arial
    }

    darkskywindspeed {
        font-size: 12px;
        color: #ff8841;
        font-family: Arial;
        line-height: 10px
    }

    .darkskydiv {
        position: relative;
        width: 800px;
        overflow: hidden!important;
        height: 450px;
        float: none;
        margin-left: 7%;
        margin-top: 1%
    }

    .darkskyforecasthome darkskytemphihome {
        margin-top: 5px;
        font-size: 12px;
        color: #ff8841;
        font-family: Arial;
        margin-left: 15%;
        width: 200px
    }

    .darkskyforecasthome darkskytemphihome span {
        font-size: 12px;
        font-family: Arial;
        color: #ff8841;
        width: 300px
    }

    .darkskyforecasthome darkskytemplohome {
        font-size: 12px;
        color: #01a4b5;
        font-family: Arial
    }

    .darkskyforecasthome darkskytemplohome span {
        font-size: 12px;
        color: #01a4b5;
        font-family: Arial
    }

    .darkskyforecasthome darkskytempwindhome {
        position: absolute;
        font-size: 12px;
        color: #aaa;
        font-family: Arial;
        margin-top: 35px;
        margin-left: 5px;
    }

    .darkskyforecasthome darkskyrainhome {
        position: absolute;
        font-size: 12px;
        color: #aaa;
        font-family: Arial;
        margin-top: 50px;
        margin-left: 15px;
        line-height: 11px
    }

    .darkskyforecasthome darkskyrainhome1 {
        position: absolute;
        font-size: 12px;
        color: #aaa;
        font-family: Arial, monospace;
        margin-top: 50px;
        margin-left: 15px;
        line-height: 11px;
        width: 200px
    }

    .darkskyforecasthome darkskytempwindhome span {
        font-size: 12px;
        color: #01a4b5;
        font-family: Arial;
        margin-top: 30px;
        text-transform: capitalize
    }

    .darkskyforecasthome darkskytempwindhome span2 {
        font-size: 12px;
        color: #ff8841;
        font-family: Arial;
        margin-top: 30px;
        margin-left: 5px
    }

    darkskyiconcurrent {
        margin-left: 30px;
        position: absolute;
        margin-top: 5px;
        margin-bottom: 30px
    }

    .darkskyiconcurrent span1 {
        font-size: 12px;
        color: #ff8841;
        margin-left: 10px;
        display: block
    }

    .darkskyiconcurrent span2 {
        font-size: 12px;
        color: #01a4b5;
        margin-left: 10px
    }

    .darkskyiconcurrent img {
        position: relative;
        width: 110px;
        margin: -40px 40% -10px 0;
        padding-right: 0;
        float: left
    }

    .darkskynexthours {
        position: absolute;
        margin-top: 30px;
        font-family: arial, sans-serif;
        text-rendering: optimizeLegibility;
        -webkit-font-smoothing: antialiased;
        -moz-osx-font-smoothing: grayscale;
        -moz-font-smoothing: antialiased;
        -o-font-smoothing: antialiased;
        -ms-font-smoothing: antialiased;
        font-size: 12px;
        line-height: 10px;
        margin-left: 0
    }

    .darkskynexthours span2 {
        font-size: 12px;
        color: #00a4b4;
        margin-left: 0;
        margin-top: -5px;
        padding: 0
    }

    body {
        line-height: 11px
    }

    grey {
        color: #aaa
    }

    blue1 {
        color: rgba(0, 154, 170, 1.000);
        line-height: 11px
    }

    orange {
        color: #ff8841
    }

    green {
        color: rgba(144, 177, 42, 1.000)
    }

    gust {
        color: #ff8841;
        line-height: 11px
    }

    unit {
        color: #aaa;
        font-size: 10px
    }

    windunit {
        color: #aaa;
        font-size: 10px
    }

    a {
        font-size: 10px;
        color: #aaa;
        text-decoration: none!important;
        font-family: arial
    }

    .forecastupdated {
        position: absolute;
        font-size: 10px;
        color: #aaa;
        font-family: arial;
        bottom: 5px;
        float: right;
        margin-left: 575px
    }

    .darkskyforecastinghome {
        float: left;
        display: inline;
        width: 15%;
        border-radius: 4px;
        margin: 2px 2px;
        font-family: Arial,monospace;
        height: 140px;
        padding: 5px;
        background: rgba(29, 32, 34, 1.000);
        background: linear-gradient(to bottom, rgba(255, 124, 57, 1.000) 12%, rgba(29, 32, 34, 1.000) 11%, rgba(29, 32, 34, 1.000) 100%, rgba(229, 77, 11, 0) 0%);
        background: -webkit-linear-gradient(to bottom, rgba(255, 124, 57, 1.000) 12%, rgba(29, 32, 34, 1.000) 11%, rgba(29, 32, 34, 1.000) 100%, rgba(229, 77, 11, 0) 0%);
        background: -moz-linear-gradient(to bottom, rgba(255, 124, 57, 1.000) 12%, rgba(29, 32, 34, 1.000) 11%, rgba(29, 32, 34, 1.000) 100%, rgba(229, 77, 11, 0) 0%);
        background: -o-linear-gradient(to bottom, rgba(255, 124, 57, 1.000) 12%, rgba(29, 32, 34, 1.000) 11%, rgba(29, 32, 34, 1.000) 100%, rgba(229, 77, 11, 0) 0%);
        color: #aaa;
        overflow: hidden!important;
        border: solid 1px #444;
        border-bottom: solid 1px #444;
        border-top: 1px solid rgba(255, 124, 57, 1.000);
    }

    .darkskyweekdayhome {
        text-align: center;
        padding: 0;
        color: #fff;
        font-family: Arial,monospace;
        font-size: 11px;
        margin: -4px 0 12px;
        background: 0;
    }
</style>