<template>
    <div class="modal-backdrop">
        <div class="modal">
            <section class="modal-body">
                <div id="body" v-if="forecast">
                    <div class="weather34darkbrowser" :url="title"></div>
                    <main class="grid">
                        <article>
                            <div class="darkskyforecastinghome" v-for="hour in forecast.hourly.data.slice(0,15)" v-bind:key="hour.time">
                                <div class="darkskyweekdayhome">{{ hour.time | hourDisplay }}</div>
                                <darkskytemphihome><img src="css/icons/temp34.svg" width="10px"><span>{{hour.temperature | toNumber }}°<grey></grey> &nbsp;</span>
                                    <grey>&nbsp;UVI:</grey> <orange>{{ hour.uvIndex }} </orange>
                                </darkskytemphihome><br>
                                <darkskyiconcurrent><img :src="'css/darkskyicons/'+hour.icon+'.svg'" width="50"></darkskyiconcurrent>
                                <darkskytempwindhome>
                                    <span2 style="color:#ff8841;"><img src="css/windicons/avgw.svg" width="20" v-bind:style="{transform:'rotate('+hour.windBearing+'deg)'}"></span2>
                                    <span4> {{ hour.windSpeed | toNumber }} | <gust>{{ hour.windGust | toNumber }}</gust></span4>
                                    <windunit>mph</windunit><br>&nbsp;<darkskyrainhome><span>{{ hour.summary}} </span></darkskyrainhome><br>
                                    <darkskyrainhome1>{{hour.precipProbability.toFixed(2) }}% <blue1><svg id="weather34 raindrop" x="0px" y="0px" viewBox="0 0 512 512" width="8px" fill="#01a4b5" stroke="#01a4b5" stroke-width="3%"><g<g><path d="M348.242,124.971C306.633,58.176,264.434,4.423,264.013,3.889C262.08,1.433,259.125,0,256,0	c-3.126,0-6.079,1.433-8.013,3.889c-0.422,0.535-42.621,54.287-84.229,121.083c-56.485,90.679-85.127,161.219-85.127,209.66 C78.632,432.433,158.199,512,256,512c97.802,0,177.368-79.567,177.368-177.369C433.368,286.19,404.728,215.65,348.242,124.971z M256,491.602c-86.554,0-156.97-70.416-156.97-156.97c0-93.472,123.907-263.861,156.971-307.658 C289.065,70.762,412.97,241.122,412.97,334.632C412.97,421.185,342.554,491.602,256,491.602z"></path></g></g><g><g><path d="M275.451,86.98c-1.961-2.815-3.884-5.555-5.758-8.21c-3.249-4.601-9.612-5.698-14.215-2.45 c-4.601,3.249-5.698,9.613-2.45,14.215c1.852,2.623,3.75,5.328,5.688,8.108c1.982,2.846,5.154,4.369,8.377,4.369 c2.012,0,4.046-0.595,5.822-1.833C277.536,97.959,278.672,91.602,275.451,86.98z"></path></g></g><g><g><path d="M362.724,231.075c-16.546-33.415-38.187-70.496-64.319-110.213c-3.095-4.704-9.42-6.01-14.126-2.914 c-4.706,3.096-6.01,9.421-2.914,14.127c25.677,39.025,46.9,75.379,63.081,108.052c1.779,3.592,5.391,5.675,9.148,5.675 c1.521,0,3.064-0.342,4.517-1.062C363.159,242.241,365.224,236.123,362.724,231.075z"></path></g></g></svg>
                                        {{ hour.precipIntensity.toFixed(2) }}</blue1><unit> in</unit></darkskyrainhome1>
                                </darkskytempwindhome>
                            </div>
                        </article>
                    </main>
                </div>
            </section>
            <button class="lity-close" type="button" aria-label="Close (Press escape to close)" @click="close">×</button>
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
    .darkskyforecastinghome {
        float: left;
        display: inline;
        width: 18%;
        border-radius: 3px;
        margin: 2px;
        height: 140px;
        padding: 1px;
        background: 0;
        border: 1px solid rgba(153, 155, 156, .1);
        color: silver
    }
</style>