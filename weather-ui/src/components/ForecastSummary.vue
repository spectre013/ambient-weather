<template>
    <div class="modal-backdrop">
        <div class="modal">
            <section class="modal-body">
                <div id="body">
                    <div class="weather34darkbrowser" url="Forecast Summary"></div>
                    <main class="grid">
                            <div style="position:absolute;width:800px;background:transparent; margin: 10px auto 0 50px;">
                                    <div class="fsdarkskydiv" v-for="day in forecast.daily.data" v-bind:key="day.time">
                                        <value>
                                            <div class="darkskyforecastinghome">
                                                <div class="darkskyweekdayhome">{{ day.time | dayFormat }} </div>
                                                <img :src="'css/darkskyicons/'+day.icon+'.svg'" width="45"><br>
                                                <darkskytemphihome>
                                                    <span>
                                                        <temp :class="tempColor(day.temperatureHigh)">{{ day.temperatureHigh }}°</temp>
                                                        <grey> | </grey>
                                                    </span>
                                                </darkskytemphihome>
                                                <darkskytemplohome>
                                                    <span>
                                                        <temp :class="tempColor(day.temperatureLow)">{{ day.temperatureLow }}°</temp>
                                                    </span>
                                                </darkskytemplohome>
                                                <darkskytemplohome>
                                                    <grey> UVI
                                                        <orange1>
                                                            <uv :class="uvColor(day.uvIndex)">{{ day.uvIndex }}</uv>
                                                        </orange1>
                                                    </grey>
                                                </darkskytemplohome><br>
                                                <br>
                                                <div class='darkskywindspeedicon'>
                                                    <img src="css/windicons/avgw.svg" width="20" v-bind:style="{transform:'rotate('+day.windBearing+'deg)'}">
                                                </div>
                                                <div class='darkskywindgust'>{{ day.windSpeed }} | {{ day.windGust }} mph</div>
                                                <darkskytempwindhome>
                                                    <span>{{ day.summary }}</span>
                                                </darkskytempwindhome><br>&nbsp;<br>
                                                <darkskypreciphome v-if="day.precipType == 'snow'">
                                                        <span>Snow <blue1>{{ day.precipAccumulation }}</blue1> in <blue1>{{ day.precipProbability }}</blue1>%</span>
                                                </darkskypreciphome>
                                                <darkskypreciphome v-if="day.precipType == 'rain'">
                                                        <span>Rain <blue1>&nbsp; {{ day.precipIntensity }}</blue1> in <blue1>{{ day.precipProbability }}</blue1>%</span>
                                                </darkskypreciphome>

                                            </div>
                                        </value>
                                    </div>
                                </div>
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
        name: 'forecastsummaryymodel',
        props: {
            forecast: Object
        },
        methods: {
            close: function() {
                this.$parent.closeModal('forecastsummary');
            },
            tempColor: function(temp) {
                if (temp < 44.6) {
                    return 'bluetds';
                } else if (temp > 104) {
                    return 'purplet';
                } else if (temp > 80.6) {
                    return 'redtds';
                } else if (temp > 64) {
                    return 'orangetds';
                } else if (temp > 55) {
                    return 'yellowtds';
                } else if (temp >= 44.6) {
                    return 'greentds';
                }
                return 'greentds';
            },
            uvColor: function(uv) {
                if (uv>=10) {
                    return 'purpleu';
                } else if (uv>7) {
                    return 'redu';
                } else if (uv>5) {
                    return 'orangeu';
                } else if (uv>2) {
                    return 'yellowu';
                } else if (uv>0) {
                    return 'greenu';
                } else if (uv==0) {
                    return 'zerou';
                }
                return 'zerou';
            }
        },
        mounted() {

        },
        computed: {

        },
        filters: {
            dayFormat: function(date) {
                return moment.unix(date).format('dddd');
            }
        }
    };
</script>

<style scoped>
    .weather34browser-footer {
        margin-top:10px;
        margin-left:45px;
        font-size: .6rem;
        color: silver;
    }
    darkskyweekday {
        position: absolute;
        margin: 3px 3px 20px 3px;
        background-color: rgba(253, 166, 16, 1.000);
        text-align: center;
        padding: 5px;
        color: #aaa;
        font-family: weathertext2, sans-serif;
        border-radius: 4px;
        font-size: 0.6rem;
        line-height: 15px
    }
    darkskytempwindhome {
        height: 50px;
        display: block;
    }
    darkskypreciphome {
        position: relative;
        left:10px;
        bottom: 3px;
    }
    darkskytemphi {
        margin-top: 5px;
        font-size: 0.6rem;
        color: rgba(255, 124, 57, 1);
        font-family: weathertext2, sans-serif;
        margin-left: 10%
    }

    darkskytemphi span {
        font-size: 0.6rem;
        color: #aaa
    }

    darkskytemplo {
        margin-top: 5px;
        font-size: 10px;
        color: #00a4b4;
        font-family: weathertext2, sans-serif;
    }

    darkskytemplo span {
        font-size: 10px;
        color: #aaa;
        font-family: weathertext2, sans-serif;
    }

    darkskysummary {
        font-size: 0.5rem;
        color: #aaa;
        font-family: weathertext2, sans-serif;
        line-height: 11px
    }

    darkskywindspeed {
        font-size: 10px;
        color: #aaa;
        font-family: weathertext2, sans-serif;
        line-height: 11px
    }

    .darkskywindspeedicon {
        position: absolute;
        font-size: 0.6rem;
        color: #aaa;
        font-family: weathertext2, sans-serif;
        line-height: 11px;
        margin-top: -65px;
        margin-left: 50px
    }

    .darkskywindgust {
        position: absolute;
        font-size: 0.6rem;
        color: #aaa;
        font-family: weathertext2, sans-serif;
        line-height: 11px;
        margin-top: -60px;
        margin-left: 72px
    }

    .fsdarkskydiv {

    }

    .darkskyforecastinghome {
        font-size: 10px;
        float: left;
        display: inline;
        width: 21%;
        border-radius: 4px;
        margin: 5px 0;
        font-family: weathertext2, system, sans-serif;
        height: 200px;
        padding: 5px;
        background: rgba(29, 32, 34, 1.000);
        background: linear-gradient(to bottom, rgba(97, 106, 114, 1.000) 12%, rgba(29, 32, 34, 0) 11%, rgba(29, 32, 34, 0) 100%, rgba(229, 77, 11, 0) 0%);
        background: -webkit-linear-gradient(to bottom, rgba(97, 106, 114, 1.000) 12%, rgba(29, 32, 34, 0) 11%, rgba(29, 32, 34, 0) 100%, rgba(229, 77, 11, 0) 0%);
        background: -moz-linear-gradient(to bottom, rgba(97, 106, 114, 1.000) 12%, rgba(29, 32, 34, 0) 11%, rgba(29, 32, 34, 0) 100%, rgba(229, 77, 11, 0) 0%);
        background: -o-linear-gradient(to bottom, rgba(97, 106, 114, 1.000) 12%, rgba(29, 32, 34, 0) 11%, rgba(29, 32, 34, 0) 100%, rgba(229, 77, 11, 0) 0%);
        color: #aaa;
        overflow: hidden!important;
        border: solid 1px #444;
        border-bottom: solid 1px #444;
        border-top: 1px solid rgba(97, 106, 114, 1.000);
    }

    .darkskyweekdayhome {
        text-align: center;
        padding: 0;
        color: #fff;
        font-family: weathertext2, sans-serif;
        font-size: 0.9rem;
        background: 0;
        margin-top: -7px;
        margin-bottom: 12px;
        font-weight: bold;
    }

    .darkskyforecasthome darkskytemphihome {
        font-size: 0.6rem;
        color: #ff7c39;
        font-family: weathertext2, sans-serif;
        line-height: 10px
    }

    .darkskyforecasthome darkskytemphihome span {
        font-size: 0.6rem;
        color: #ff7c39;
        font-family: weathertext2, sans-serif;
        line-height: 10px
    }

    .darkskyforecasthome darkskytemplohome {
        font-size: 0.6rem;
        color: #ff7c39;
        font-family: weathertext2, sans-serif;
        line-height: 10px
    }

    .darkskyforecasthome darkskytemplohome span {
        font-size: 0.6rem;
        color: #01a4b5;
        font-family: weathertext2, sans-serif;
    }

    .darkskyforecasthome darkskytempwindhome {
        font-size: 10px;
        color: #aaa;
        font-family: weathertext2, sans-serif;
        line-height: 10px
    }

    .darkskyforecasthome darkskytempwindhome span {
        font-size: 0.6rem;
        color: #aaa;
        font-family: weathertext2, sans-serif;
        line-height: 10px
    }

    .darkskyforecasthome darkskytempwindhome span2 {
        font-size: 0.6rem;
        color: #aaa;
        font-family: weathertext2, sans-serif;
        line-height: 10px;
        margin-top: 3px
    }

    .darkskyiconcurrent img {
        position: relative;
        width: 80px;
        margin: -50px 0 -10px 0;
        padding-right: 0;
        float: left
    }


    .darkskynexthours span2 {
        line-height: 12px
    }
    grey {
        color: #aaa
    }

    blue1 {
        color: #01a4b5;
        text-transform: capitalize
    }

    orange1 {
        color: #ff7c39
    }

    green {
        color: rgba(144, 177, 42, 1)
    }

    a {
        font-size: 10px;
        color: #aaa;
        text-decoration: none!important;
        font-family: weathertext2, sans-serif;
    }

    precip {
        position: relative;
        top: 5px;
        padding: 2px;
        border-radius: 3px;
        background: rgba(97, 106, 114, 0.2);
    }

    value {
        font-size: .8em;
        font-family: weathertext2, sans-serif;
    }

    value1 {
        font-size: 1em;
        font-family: weathertext2, sans-serif;
    }

    temp {
        color: #fff;
        text-transform: capitalize;
        border-radius: 2px;
        width: 35px;
        padding: 0 3px;
        font-size: 11px;
    }

    .bluetds {
        background: #01a4b5
    }

    .yellowtds {
        background: #e6a141
    }

    .orangetds {
        background: #d05f2d
    }

    .greentds {
        background: #90b12a
    }

    .redtds {
        background: -webkit-linear-gradient(90deg, #d86858, rgba(211, 93, 78, .7));
        background: linear-gradient(90deg, #d86858, rgba(211, 93, 78, .7))
    }

    .purpletds {
        background: -webkit-linear-gradient(90deg, #d86858, rgba(157, 59, 165, .4));
        background: linear-gradient(90deg, #d86858, rgba(157, 59, 165, .4))
    }

    uv {
        color: #fff;
        border-radius: 2px;
        width: 35px;
        font-size: 11px;
        padding: 0 3px
    }

    .blueu {
        background: #01a4b5
    }

    .zerou {
        color: #777
    }

    .yellowu {
        background: #e6a141
    }

    .orangeu {
        background: #d05f2d
    }

    .greenu {
        background: #90b12a
    }

    .redu {
        background: #cd5245
    }

    .purpleu {
        background: #b600b0
    }

    .zerou {
        background: #4a636f
    }
</style>