<template>
    <div class="weatherfooter-item">
        <div class="hardwarelogo1">
            <a href="http://ambientweather.com" alt="http://ambientweather.com" title="http://ambientweather.com">
                <img src="img/ambient.png" alt="Ambient Weather" title="Ambient" width="25px" height="25px">
                <div class="hardwarelogo1text">
                    Ambient WS-2000
                </div>
            </a>
        </div>
        <div class="hardwarelogo2">
            <a href="https://weather34.com/homeweatherstation/" title="https://weather34.com/homeweatherstation/" target="_blank"><br>
                <img src="img/weather34logo.svg" width="30px" alt="https://weather34.com/homeweatherstation/" class="homeweatherstationlogo">
                <weather34>designed by weather34 2015-{{ year }}</weather34>
            </a>
        </div>
        <div class="footertext">
           <svg id="i-info" viewBox="0 0 32 32" width="9" height="9" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M16 14 L16 23 M16 8 L16 10"></path><circle cx="16" cy="16" r="14"></circle></svg>&nbsp;(<value>WS-2000 <maxred>Osprey</maxred></value>) Ambient Weather(<value><maxred>1.0</maxred></value>)&nbsp;
           <svg id="i-info" viewBox="0 0 32 32" width="9" height="9" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M16 14 L16 23 M16 8 L16 10"></path><circle cx="16" cy="16" r="14"></circle></svg>&nbsp;WH32B x3  WH31B<br>
        </div>
        <div class="footertext">
           <svg id="i-info" viewBox="0 0 32 32" width="9" height="9" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M16 14 L16 23 M16 8 L16 10"></path><circle cx="16" cy="16" r="14"></circle></svg>&nbsp;Charts&nbsp;<img src="img/Copyright.svg" height="10px" width="10px"> <a target="_blank" rel="noopener" href="https://www.canvasjs.com">CanvasJS</a> -  Forecasts &nbsp;<img src="img/Copyright.svg" height="10px" width="10px"> <a target="_blank" rel="noopener" href="https://darksky.net/poweredby/">Powered by Dark Sky</a>
        </div>
        <div class="footertext">
            <svg id="weather34" github="" icon="" width="9pt" height="9pt" fill-rule="evenodd" clip-rule="evenodd" image-rendering="optimizeQuality" shape-rendering="geometricPrecision" text-rendering="geometricPrecision" viewBox="0 0 640 640"><path fill="#3b9cac" d="M320 8a320 320 0 0 0-101 624c16 3 22-7 22-16l-1-59c-89 19-108-38-108-38-14-37-35-47-35-47-29-20 2-19 2-19 32 2 49 33 49 33 29 48 75 34 93 26 3-21 11-35 21-43-71-8-146-35-146-158 0-35 12-63 33-86-4-8-14-40 3-84 0 0 27-9 88 32a307 307 0 0 1 160 0c61-41 88-32 88-32 17 44 7 76 3 84 21 23 33 51 33 86 0 123-75 150-146 158 11 10 22 29 22 59l-1 88c0 9 6 19 22 15A320 320 0 0 0 320 8z"></path></svg>&nbsp;Ambient Version Repository at https://github.com/spectre013/ambient-weather&nbsp;<img src="img/flags/us.svg" width="20px">
        </div>
    </div>
</template>

<script>
    import moment from 'moment';

    export default {
        name: 'minmax',
        mounted() {

        },
        methods: {

        },
        filters: {

        },
        computed: {
            year: function() {
                return moment().format("Y");
            }
        },
    }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
