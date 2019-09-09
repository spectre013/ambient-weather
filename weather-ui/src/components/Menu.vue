<template>
    <div id="sidebar">
        <div class="weather34toolbar__right">
            <!-- topbarmetric -->
            <a href="#" v-on:click="setUnits    ">
                <topbarimperial v-if="units==='imperial'">°F</topbarimperial>
                <topbarmetric v-if="units==='metric'">°C</topbarmetric>
            </a>
        </div>
        <input type="checkbox" class="openweather34sidebarMenu" id="openweather34sidebarMenu">
        <label for="openweather34sidebarMenu" class="weather34sidebarIconToggle">
            <div class="weather34spinner weather34cross part-1"></div>
            <div class="weather34spinner weather34horizontal"></div>
            <div class="weather34spinner weather34cross part-2"></div>
        </label>
        <div id="weather34sidebarMenu">
            <ul class="weather34sidebarMenuInner">
                <li class="header">ADMIN</li>
                <li>
                    <a href="#" title="WEATHERSTATION SETTINGS PAGE">
                        <weather34menumarkeryellow></weather34menumarkeryellow> Settings
                    </a>
                </li>
<!--                <li class="header">UI THEME</li>-->
<!--                <li><a v-on:click="setActiveStyleSheet()">-->
<!--                        <weather34menumarkerlight></weather34menumarkerlight> {{ styleText() }}-->
<!--                    </a>-->
<!--                </li>-->
                <li class="header">UNITS</li>
                    <li><a href="#" v-on:click="setUnits">
                        <weather34menumarkerorange></weather34menumarkerorange> Imperial <topbarimperialf>°F</topbarimperialf>
                    </a>
                </li>
                <li>
                    <a href="#" v-on:click="setUnits">
                        <weather34menumarkerblue></weather34menumarkerblue> Metric <topbarmetricc>°C</topbarmetricc>
                    </a>
                </li>
                <li> <a href="#" v-on:click="setSpeed">
                    <weather34menumarkerorange></weather34menumarkerorange> MPH <topbarimperialf>°F</topbarimperialf>
                </a>
                </li>
                <li> <a href="#" v-on:click="setSpeed">
                        <weather34menumarkerred></weather34menumarkerred> M/S <topbarmetricc>°C</topbarmetricc>
                    </a>
                </li>
                <li class="header">EXTRAS</li>
                <li>
                    <a href="bio.php" data-lity="" title="Weather Station Owner Contact Card Info">
                        <weather34menumarkerorange></weather34menumarkerorange>Contact Card
                    </a>
                </li>
                <li>
                    <a href="stationinfo.php" data-lity="" title="Hardware Weather Station Hardware Info">
                        <weather34menumarkerred></weather34menumarkerred> Hardware Info
                    </a>
                </li>
<!--                <li><a href=""><weather34menumarkerbluegrey></weather34menumarkerbluegrey> Select Language</a><br>-->
<!--                    <div class="languages34">-->
<!--                        <a href="index.php?lang=en"><img src="img/flags/en.svg" title="English" width="25px" height="25px"></a>&nbsp;-->
<!--                        <a href="index.php?lang=dk"><img src="img/flags/dk.svg" title="Danish" width="25px" height="25px"></a>&nbsp;-->
<!--                        <a href="index.php?lang=gr"><img src="img/flags/gr.svg" title="Greek" width="25px" height="25px"></a>&nbsp;-->
<!--                        <a href="index.php?lang=it"><img src="img/flags/it.svg" title="Italian" width="25px" height="25px"></a>&nbsp;-->
<!--                        <a href="index.php?lang=fr"><img src="img/flags/fr.svg" title="French" width="25px" height="25px"></a>&nbsp;-->
<!--                    </div>-->
<!--                    <br>-->
<!--                    <div class="languages34">-->
<!--                        <a href="index.php?lang=dl"><img src="img/flags/dl.svg" title="German" width="25px" height="25px"></a>&nbsp;-->
<!--                        <a href="index.php?lang=nl"><img src="img/flags/nl.svg" title="Dutch" width="25px" height="25px"></a>&nbsp;-->
<!--                        <a href="index.php?lang=cat"><img src="img/flags/cat.svg" title="Catalan" width="25px" height="25px"></a>&nbsp;-->
<!--                        <a href="index.php?lang=sp"><img src="img/flags/sp.svg" title="Spanish" width="25px" height="25px"></a>&nbsp;-->
<!--                        <a href="index.php?lang=tr"><img src="img/flags/tr.svg" title="Turkish" width="25px" height="25px"></a>&nbsp;-->
<!--                    </div>-->
<!--                    <br>-->
<!--                    <div class="languages34">-->
<!--                        <a href="index.php?lang=hu"><img src="img/flags/hu.svg" title="Hungary" width="25px" height="25px"></a>&nbsp;-->
<!--                        <a href="index.php?lang=pl"><img src="img/flags/pl.svg" title="Polish" width="25px" height="25px"></a>&nbsp;-->
<!--                        <a href="index.php?lang=sp"><img src="img/flags/ar.svg" title="Argentina" width="25px" height="25px"></a>&nbsp;-->
<!--                        <a href="index.php?lang=dl"><img src="img/flags/ch-switzerland.svg" title="Switzerland" width="25px" height="25px"></a>&nbsp;-->
<!--                    </div>-->
<!--                </li>-->
                <li>
                    <a href="https://github.com/spectre013/ambient-weather/"
                       title="Weather34 Ambient version via Github " target="_blank">
                        <weather34menumarkerbluegrey></weather34menumarkerbluegrey> Download Meteobridge Template
                    </a>
                </li>
                <li>
                    <a href="https://github.com/spectre013/ambient-weather/"
                       title="Weather34 Ambient version via Github " target="_blank">
                        <weather34menumarkerbluegrey></weather34menumarkerbluegrey> Maintained by Brian(Spectre013)
                    </a>
                </li>
                <li>
                    <a href="https://github.com/weather34/Meteobridge-Weather34-Template"
                       title="https://github.com/weather34/Meteobridge-Weather34-Template/" target="_blank">
                        <weather34menumarkerbluegrey></weather34menumarkerbluegrey> Designed by weather34.com
                    </a>
                </li>
            </ul>
        </div>
    </div>
</template>

<script>
    export default {
        name: 'minmax',
        data () {
            return {
                currentTheme: 'dark',
                units: 'imperial',
                speed: 'mph',
            }
        },
        mounted() {

        },
        beforeMount: function() {
            this.currentTheme = this.$store.getters.theme;
            this.units = this.$store.getters.units;
        },
        methods: {
            // setActiveStyleSheet: function() {
            //     let theme = this.currentTheme === 'dark' ? 'light' : 'dark';
            //     this.$store.dispatch('setTheme',theme);
            //     this.setStyle();
            // },
            // styleText: function() {
            //     return this.currentTheme === 'dark' ? 'Light Mode' : 'Dark Mode';
            // },
            setUnits: function() {
                this.units = this.units === 'imperial' ? 'metric' : 'imperial';
                this.$store.dispatch('setUnits',this.units);
                this.$store.dispatch('setSettings',this.$store.getters.settings)
            },
            setSpeed: function() {
                this.speed = this.speed === 'mph' ? 'ms' : 'mph';
                this.$store.dispatch('setSpeed',this.speed);
                this.$store.dispatch('setSettings',this.$store.getters.settings)
            }

        },
        filters: {

        },
        computed: {

        },
    }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
    .weather34sidebarMenuInner {
        padding-top:40px;
    }
</style>
