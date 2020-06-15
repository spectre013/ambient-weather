<template>
    <div class="weather-item">
        <div class="chartforecast">
        <span class="yearpopup">
            <span class="yearpopup">
                <a alt="UV Guide" title="UV Guide" href="uvindex.php" data-lity="">
                    <svg viewBox="0 0 32 32" width="8" height="8" fill="none" stroke="#777" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"> <path d="M14 9 L3 9 3 29 23 29 23 18 M18 4 L28 4 28 14 M28 4 L14 18"></path>
                    </svg> UV Guide </a></span> <span class="yearpopup"><a alt="UV Almanac" title="UV Almanac" href="uvalmanac.php" data-lity=""> <svg viewBox="0 0 32 32" width="8" height="8" fill="none" stroke="#777" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%">
                        <path d="M14 9 L3 9 3 29 23 29 23 18 M18 4 L28 4 28 14 M28 4 L14 18"></path>
                    </svg> UV Almanac </a></span><span class="yearpopup"> <a alt="Solar Almanac" title="Solar Almanac" href="solaralmanac.php" data-lity=""><svg viewBox="0 0 32 32" width="8" height="8" fill="none" stroke="#777" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%">
                        <path d="M14 9 L3 9 3 29 23 29 23 18 M18 4 L28 4 28 14 M28 4 L14 18"></path>
                    </svg> Solar Almanac </a></span></span>
        </div>
        <span class="moduletitle">UV | Solar Radiation</span>
        <div id="uvi" v-if="current && minmax && astro">
            <div class="updatedtime"><span><svg id="i-info" viewBox="0 0 32 32" width="6" height="6" fill="#9aba2f" stroke="#9aba2f" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"> <path d="M16 14 L16 23 M16 8 L16 10"></path><circle cx="16" cy="16" r="14"></circle>
                </svg> {{current.date | now }} </span></div>
            <div class="weather34solarword">
                <valuetext>W/m<sup>2</sup></valuetext>
            </div>
            <div class="weather34solarvalue">
                <div class="solartodaycontainer1">
                    <div class="solarluxtoday">{{current.solarradiation }}</div>
                </div>
            </div>
            <div class="solarluxtodayword">
                <valuetext>Solar Radiation</valuetext>
            </div>
            <div class="solarwrap"></div>
            <div class="uvsolarwordbig">Currently</div>

            <div class="uvcontainer1">
                <div v-bind:class="uvToday">{{current.uv }}<smalluvunit> &nbsp;UVI</smalluvunit>
                </div>
            </div>
            <div class="uvtrend">UV INDEX</div>
            <div class="uvcaution">&nbsp;&nbsp;&nbsp;&nbsp;<valuetext>Max {{ minmax.max.day.value }} <time>
                <value>({{ minmax.max.day.date | timeFormat }}</value>)
            </time></valuetext>
            </div>
            <div class="uvcautionbig"><svg id="weather34 uvi icon" width="10pt" height="10pt" viewBox="0 0 302 255">
                <path fill="#ff8841" stroke="#ff8841" stroke-width=".1"
                      d="M147.5 5h8v29.2h-8V5zM96.6 34.5l6-6 16.8 17c-2 1.8-3.8 4.2-6.2 5.7-5.4-5.7-11-11-16.6-16.7zM184.6 45.4l17-16.8 5.8 6-17 16.8-5.8-6zM143.6 46.8c8.3-1.5 17-.7 24.6 3 11 5.2 19 15.8 21.2 27.7 1.3 8 .4 16.5-3.2 23.8-5 10.6-15.3 18.6-27 21-8.3 1.5-17.4.6-25-3.5-10.8-5.3-18.7-16-20.6-27.8-1.2-7.8-.3-16 3.2-23.3 5-10.5 15.3-18.6 26.8-21zM72 80.5h29.2v8H72v-8zM201.8 80.5H231v8h-29.2v-8zM96.6 133.5l17-17 5.8 6-17 17-5.8-6zM184.6 122.5l6-6 16.8 17-6 6c-5.5-5.8-11.2-11.4-16.8-17zM147.5 134.8h8V164h-8v-29.2z">
                </path>
            </svg><span>UVI</span> {{ this.current.uv | uvCaution }}</div>
        </div>
    </div>
</template>

<script>
    import moment from 'moment';
    import axios from 'axios';

    export default {
        name: 'temp',
        props: {
            current: Object,
            astro: Object
        },
        data() {
            return {
                hasSunset: false,
                minmax: null,
            }
        },
        mounted: function() {
            function updateData(self){
                axios.get('/api/minmax/uv').then(response => (self.minmax = response.data));
                setTimeout(function() { updateData(self); }, 60000);
            }
            updateData(this);
            axios.get('/api/minmax/uv').then(response => (this.minmax = response.data));
        },
        watch: {
            astro: function() {
                this.sunHasSet();
            }
        },
        methods: {
            openModal: function() {
                this.$root.showModal();
            },
            sunHasSet: function() {
                let h = this.astro.sunset.split(":");
                let sunset = moment().startOf('day').hour(h[0]).minute(h[1]);
                let m = sunset.diff(moment(),'minutes');
                if(m <= 0) {
                    this.hasSunset = true;
                } else {
                    this.hasSunset = false;
                }
            }
        },
        filters: {
            now: function (date) {
                return moment(date).format('HH:mm:ss');
            },
            timeFormat: function (date) {
                return moment(date).format('HH:mm');
            },
            uvCaution: function(uv) {
                if(uv >= 10) {
                    return "Extreme";
                } else if(uv >= 8) {
                    return "Very High";
                } else if(uv >= 5) {
                    return "High";
                } else if(uv >= 3) {
                    return "Moderate";
                } else if(uv >= 0) {
                    return "Low";
                } else if(this.hasSunset && uv <= 0) {
                    return "Below Horizon";
                }
                return "";
            },
        },
        computed: {
            luxToday: function() {
                if(this.current.solarradiation === 0) {
                    return "solarluxtodaydark";
                }
                return "solarluxtoday"
            },
            uvToday: function() {
                if(this.current.uv >= 10) {
                    return "uvtoday11";
                } else if(this.current.uv >= 8) {
                    return "uvtoday9-10";
                } else if(this.current.uv >= 5) {
                    return "uvtoday6-8";
                } else if(this.current.uv >= 3) {
                    return "uvtoday4-5";
                } else if(this.current.uv >= 0) {
                    return "uvtoday1-3";
                }
                return "";
            },
            solarToday: function() {
                if(this.current.solarradiation >= 1000) {
                    return "solartoday1000";
                } else if(this.current.solarradiation >= 500) {
                    return "solartoday500";
                } else if(this.current.solarradiation >= 10) {
                    return "solartoday200";
                } else if(this.current.solarradiation >= 0) {
                    return "solartoday1";
                }
                return "";
            }
        }
    }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
