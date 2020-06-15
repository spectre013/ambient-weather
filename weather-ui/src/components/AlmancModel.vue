<template>
    <div class="modal-backdrop">
        <div class="modal">
            <section class="modal-body">
                <div id="body" v-if="temp && humid && maxtemp && mintemp ">
                    <div class="weather34darkbrowser" url="Temperature Almanac"></div>
                    <main class="grid">
                        <article>
                            <div class="actualt">Temperature Today </div>
                            <div class="temperaturecontainer">
                                <div v-bind:class="temp.max.day.value | temperaturetoday">{{ temp.max.day.value }}<smalluvunit>&deg;F</smalluvunit>
                                </div>
                                <div class="temperaturetrend">{{ temp.max.day.date | hourFormat }}</div>
                                <div class="temperaturecontainer">
                                    <div v-bind:class="temp.min.day.value | temperaturetoday">{{ temp.min.day.value }}<smalluvunit>&deg;F</smalluvunit>
                                    </div>
                                </div>
                                <div class="temperaturetrend">{{ temp.min.day.date | hourFormat }}</div>
                                <div class="hitempypos">
                                    <div class="hitempd" v-if="dew">Dew Max<orange> {{ dew.max.day.value }}</orange>&deg;F {{ dew.max.day.date | hourFormat }}<br></div>
                                    <div class="hitempd" style="margin-top:-5px;" v-if="dew">Dew Min<blue>&nbsp;{{ dew.min.day.value }}</blue>&deg;F {{ dew.min.day.date | hourFormat }}<br></div>
                                </div>
                                <div class="hitempypos">
                                    <div class="hitempd1" style="margin-top:15px;">Hum Max<orange>&nbsp;{{ humid.max.day.value }}</orange>% {{ humid.max.day.date | hourFormat }}<br></div><br>
                                    <div class="hitempd1" style="margin-top:16px;">Hum Min<blue>&nbsp;{{ humid.min.day.value }}</blue>% {{ humid.min.day.date | hourFormat }}<br></div><br>
                                </div>
                            </div>
                        </article>
                        <article>
                            <div class="actualt">Temperature Yesterday </div>
                            <div class="temperaturecontainer">
                                <div v-bind:class="temp.max.yesterday.value | temperaturetoday">{{ temp.max.yesterday.value }}<smalluvunit>&deg;F</smalluvunit>
                                </div>
                                <div class="temperaturetrend">{{ temp.max.yesterday.date | hourFormat }}</div>
                                <div class="temperaturecontainer">
                                    <div v-bind:class="temp.min.yesterday.value | temperaturetoday">{{ temp.min.yesterday.value }}<smalluvunit>&deg;F</smalluvunit>
                                    </div>
                                </div>
                                <div class="temperaturetrend">{{ temp.min.yesterday.date | hourFormat }}</div>
                                <div class="hitempypos">
                                    <div class="hitempd" v-if="dew">Dew Max<orange> {{ dew.max.yesterday.value }}</orange>&deg;F {{ dew.max.yesterday.date | hourFormat }}<br></div>
                                    <div class="hitempd" style="margin-top:-5px;" v-if="dew">Dew Min<blue>&nbsp;{{ dew.min.yesterday.value }}</blue>&deg;F {{ dew.min.yesterday.date | hourFormat }}<br></div>
                                </div>
                                <div class="hitempypos">
                                    <div class="hitempd1" style="margin-top:15px;">Hum Max<orange>&nbsp;{{ humid.max.yesterday.value }}</orange>% {{ humid.max.yesterday.date | hourFormat }}<br></div><br>
                                    <div class="hitempd1" style="margin-top:16px;">Hum Min<blue>&nbsp;{{ humid.min.yesterday.value }}</blue>% {{ humid.min.yesterday.date | hourFormat }}<br></div><br>
                                </div>
                            </div>
                        </article>
                        <article>
                            <div class="actualt">Temperature {{ temp.max.month.date | monthYear }} </div>
                            <div class="temperaturecontainer">
                                <div v-bind:class="temp.max.month.value | temperaturetoday">{{ temp.max.month.value }}<smalluvunit>&deg;F</smalluvunit>
                                </div>
                                <div class="temperaturetrend">{{ temp.max.month.date | dayMonth }}</div>
                                <div class="temperaturecontainer">
                                    <div v-bind:class="temp.min.month.value | temperaturetoday">{{ temp.min.month.value }}<smalluvunit>&deg;F</smalluvunit>
                                    </div>
                                </div>
                                <div class="temperaturetrend">{{ temp.min.month.date | dayMonth }}</div>
                                <div class="hitempypos" >
                                    <div class="hitempd" v-if="dew">Dew Max<orange> {{ dew.max.month.value }}</orange>&deg;F {{ dew.max.month.date | dayMonth }}<br></div>
                                    <div class="hitempd" style="margin-top:-5px;" v-if="dew">Dew Min<blue>&nbsp;{{ dew.min.yesterday.value }}</blue>&deg;F {{ dew.min.month.date | dayMonth }}<br></div>
                                </div>
                                <div class="hitempypos">
                                    <div class="hitempd1" style="margin-top:15px;">Hum Max<orange>&nbsp;{{ humid.max.month.value }}</orange>% {{ humid.max.month.date | dayMonth }}<br></div><br>
                                    <div class="hitempd1" style="margin-top:16px;">Hum Min<blue>&nbsp;{{ humid.min.month.value }}</blue>% {{ humid.min.month.date | dayMonth }}<br></div><br>
                                </div>
                            </div>
                        </article>
                        <article>
                            <div class="actualt">Temperature {{ temp.max.year.date | yearFormat }} </div>
                            <div class="temperaturecontainer">
                                <div v-bind:class="temp.max.year.value | temperaturetoday">{{ temp.max.year.value }}<smalluvunit>&deg;F</smalluvunit>
                                </div>
                                <div class="temperaturetrend">{{ temp.max.year.date | dayMonth }}</div>
                                <div class="temperaturecontainer">
                                    <div v-bind:class="temp.min.year.value | temperaturetoday">{{ temp.min.year.value }}<smalluvunit>&deg;F</smalluvunit>
                                    </div>
                                </div>
                                <div class="temperaturetrend">{{ temp.min.year.date | dayMonth }}</div>
                                <div class="hitempypos">
                                    <div class="hitempd" v-if="dew">Dew Max<orange> {{ dew.max.year.value }}</orange>&deg;F {{ dew.max.year.date | dayMonth }}<br></div>
                                    <div class="hitempd" style="margin-top:-5px;" v-if="dew">Dew Min<blue>&nbsp;{{ dew.min.year.value }}</blue>&deg;F {{ dew.min.year.date | dayMonth }}<br></div>
                                </div>
                                <div class="hitempypos">
                                    <div class="hitempd1" style="margin-top:15px;">Hum Max<orange>&nbsp;{{ humid.max.year.value }}</orange>% {{ humid.max.year.date | dayMonth }}<br></div><br>
                                    <div class="hitempd1" style="margin-top:16px;">Hum Min<blue>&nbsp;{{ humid.min.year.value }}</blue>% {{ humid.min.year.date | dayMonth }}<br></div><br>
                                </div>
                            </div>
                        </article>
                        <article style="height:110px;">
                            <div class="actualt">Temperature All-Time </div>
                            <div class="temperaturecontainer">
                                <div v-bind:class="maxtemp.value | temperaturetoday">{{ maxtemp.value }}<smalluvunit>&deg;F</smalluvunit>
                                </div>
                                <div class="temperaturecontainer">
                                    <div v-bind:class="mintemp.value | temperaturetoday">{{ mintemp.value }}<smalluvunit>&deg;F</smalluvunit>
                                    </div>
                                </div>
                                <div class="hitempypos">
                                    <div class="hitempd2">{{ maxtemp.date | fullFormat }}<br></div>
                                    <div class="hitempd2" style="margin-top:25px;">{{ mintemp.date | fullFormat }}</div>
                                </div>
                            </div>
                        </article>
                        <article style="height:110px;" v-if="maxdew && mindew">
                            <div class="actualt">Dewpoint All-Time </div>
                            <div class="temperaturecontainer">
                                <div v-bind:class="maxdew.value | temperaturetoday">{{ maxdew.value }}<smalluvunit>&deg;F</smalluvunit>
                                </div>
                                <div class="temperaturecontainer">
                                    <div v-bind:class="mindew.value | temperaturetoday">{{ mindew.value }}<smalluvunit>&deg;F</smalluvunit>
                                    </div>
                                </div>
                                <div class="hitempypos">
                                    <div class="hitempd2">{{ maxdew.date | fullFormat }}<br></div>
                                    <div class="hitempd2" style="margin-top:25px;">{{ mindew.date | fullFormat }}</div>
                                </div>
                            </div>
                        </article>
                        <article style="height:60px;width:250px;">
                            <div class="lotemp">
                                <svg id="i-info" viewBox="0 0 32 32" width="9" height="9" fill="none" stroke="currentcolor"
                                    stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%">
                                    <path d="M16 14 L16 23 M16 8 L16 10"></path>
                                    <circle cx="16" cy="16" r="14"></circle>
                                </svg>
                                <a href="https://canvasjs.com" title="https://canvasjs.com" target="_blank" style="font-size:8px;">
                                    Charts rendered and compiled using CanvasJs.com v2.3.1 GA (CC BY-NC 3.0) Non-Commercial-Version </a>
                            </div><br>
                        </article>
                    </main>
                    <main class="grid1">
                        <articlegraph>
                            <div ref="chart" id="chartContainer2"></div>
                        </articlegraph>
                    </main>
                </div>
                <button class="lity-close" type="button" aria-label="Close (Press escape to close)" @click="close">×</button>
            </section>
        </div>
    </div>
</template>
<script>
    /* global CanvasJS */
    import moment from 'moment';
    import axios from 'axios';

    export default {
        name: 'almanacmodel',
        props: {
            minmax: Object,
            appName: {
                type: String,
                default: String
            },
            options: Object,
        },
        data () {
            return {
                data: null,
                maxtemp: null,
                mintemp: null,
                maxdew: null,
                mindew: null,
                temp: null,
                dew: null,
                humid: null
            }
        },
        mounted: function() {
            axios.get('/api/chart/'+this.options.tempfield+'/year').then(response => (this.data = response.data));
            axios.get('/api/alltime/max/'+this.options.tempfield).then(response => (this.maxtemp = response.data));
            axios.get('/api/alltime/min/'+this.options.tempfield).then(response => (this.mintemp = response.data));
            axios.get('/api/minmax/'+this.options.tempfield).then(response => (this.temp = response.data));
            axios.get('/api/minmax/'+this.options.humidfield).then(response => (this.humid = response.data));

            if(this.options.tempfield === 'tempf') {
                axios.get('/api/alltime/max/dewPoint').then(response => (this.maxdew = response.data));
                axios.get('/api/alltime/min/dewPoint').then(response => (this.mindew = response.data));
                axios.get('/api/minmax/dewpoint').then(response => (this.dew = response.data));
            }

        },
        filters: {
            hourFormat: function (date) {
                return moment(date).format('HH:mm');
            },
            monthYear: function (date) {
                return moment(date).format('MMMM Y');
            },
            dayMonth: function (date) {
                return moment(date).format('Do MMM');
            },
            yearFormat: function (date) {
                return moment(date).format('Y');
            },
            fullFormat: function (date) {
                return moment(date).format('Do MMM Y');
            },
        },
        methods: {

            close: function() {
               this.$parent.closeModal('almanac');
            },
            drawChart: function( dataPoints1 , dataPoints2 ) {

                let chart = new CanvasJS.Chart("chartContainer2", {
                    backgroundColor: 'rgba(40, 45, 52,.4)',
                    animationEnabled: true,
                    animationDuration: 500,
                    margin: 0,
                    height: 230,
                    width: 870,
                    title: {
                        text: " ",
                        fontSize: 11,
                        fontColor: '#ccc',
                        fontFamily: "arial",
                    },
                    toolTip:{
                        fontStyle: "normal",
                        cornerRadius: 4,
                        backgroundColor: 'rgba(37, 41, 45, 0.95)',
                        fontSize: 11,
                        contentFormatter: function(e) {
                            var str = '<span style="color: #ccc;">' + e.entries[0].dataPoint.label + '</span><br/>';
                            for (var i = 0; i < e.entries.length; i++) {
                                var temp = '<span style="color: ' + e.entries[i].dataSeries.color + ';">' + e.entries[i].dataSeries.name + '</span> <span style="color: #ccc;">' + e.entries[i].dataPoint.y.toFixed(1) + " &deg;F" + '</span> <br/>';
                                str = str.concat(temp);
                            }
                            return (str);
                        },
                        shared: true,
                    },
                    axisX: {
                        gridColor: 'RGBA(64, 65, 66, 0.8)',
                        labelFontSize: 10,
                        labelFontColor: '#ccc',
                        lineThickness: 1,
                        gridDashType: "dot",
                        gridThickness: 1,
                        titleFontFamily: "arial",
                        labelFontFamily: "arial",
                        minimum:-0.5,
                        interval:'auto',
                        intervalType:"month",
                        xValueType: "dateTime",
                        crosshair: {
                            enabled: true,
                            snapToDataPoint: true,
                            color: '#009bab',
                            labelFontColor: "#F8F8F8",
                            labelFontSize:11,
                            labelBackgroundColor: '#009bab',
                        }

                    },
                    axisY:{
                        title: "",
                        titleFontColor: '#ccc',
                        titleFontSize: 10,
                        titleWrap: false,
                        margin: 0,
                        interval: 'auto',
                        lineThickness: 1,
                        gridThickness: 1,
                        includeZero: false,
                        minimum:-10,
                        gridColor: 'RGBA(64, 65, 66, 0.8)',
                        gridDashType: "dot",
                        labelFontSize: 10,
                        labelFontColor: '#ccc',
                        titleFontFamily: "arial",
                        labelFontFamily: "arial",
                        labelFormatter: function ( e ) {
                            return e.value .toFixed(0) + "&deg;F" ;
                        },
                        crosshair: {
                            enabled: true,
                            snapToDataPoint: true,
                            color: '#ff832f',
                            labelFontColor: "#F8F8F8",
                            labelFontSize:10,
                            labelMaxWidth: 70,
                            labelBackgroundColor: '#ff832f',
                            valueFormatString: "#0.0&deg;F",
                        }

                    },
                    legend:{
                        fontFamily: "arial",
                        fontColor: '#ccc',

                    },

                    data: [
                        {
                            type: "splineArea",
                            color: 'rgba(255, 148, 82, 0.95)',
                            lineColor: 'rgba(255, 131, 47, 1)',
                            markerSize:0,
                            showInLegend:false,
                            legendMarkerType: "circle",
                            lineThickness: 2,
                            markerType: "circle",
                            name:" Hi Temp",
                            dataPoints: dataPoints1,
                            yValueFormatString: "#0.# &deg;F",

                        },
                        {

                            type: "splineArea",
                            color: 'rgba(0, 164, 180, 1)',
                            markerSize:0,
                            markerColor: '#007181',
                            showInLegend:false,
                            legendMarkerType: "circle",
                            lineThickness: 2,
                            lineColor:  '#007181',
                            markerType: "circle",
                            name:" Lo Temp",
                            dataPoints: dataPoints2,
                            yValueFormatString: "#0.# &deg;F",

                        }
                    ]
                });

                chart.render();
            }
        },
        computed: {

        },
        watch: {
           data: function() {
               function checkDom(that){
                   if(document.querySelector("#chartContainer2")) {
                       that.drawChart(that.data.data1,that.data.data2);
                   } else {
                       setTimeout(function() { checkDom(that); }, 500);
                   }
               }
               checkDom(this);
           }
        },
    };
</script>

<style scoped>
    .canvasjs-chart-credit {
        display:none;
    }
    .modal-backdrop {
        position: fixed;
        top: 0;
        bottom: 0;
        left: 0;
        right: 0;
        display: flex;
        justify-content: center;
        align-items: center;
        z-index: 10000;
    }

    .modal {
        background: #181818;
        height: 550px;
        width: 900px;
        overflow-x: auto;
        display: flex;
        flex-direction: column;
        font-size: 13px;
        color: black;
    }

    .modal-body {
        position: relative;
        padding: 20px 10px;
        font-size: 13px;
    }


    #chartContainer2 {
        width:100%;
        height:240px;
        padding:0;
        border-radius:3px;
        border: 1px solid rgba(245, 247, 252,.02);
        box-shadow: 2px 2px 6px 0px  rgba(0,0,0,0.6);
        -webkit-font-smoothing: antialiased;
        -moz-osx-font-smoothing: grayscale;
    }

    .grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(185px, 2fr));
        grid-gap: 5px;
        align-items: stretch;
        color: #f5f7fc;
    }

    .grid>article {
        border: 1px solid rgba(245, 247, 252, .02);
        box-shadow: 2px 2px 6px 0px rgba(0, 0, 0, 0.6);
        padding: 5px;
        font-size: 0.8em;
        -webkit-border-radius: 4px;
        border-radius: 4px;
        background: 0;
        -webkit-font-smoothing: antialiased;
        -moz-osx-font-smoothing: grayscale;
        height: 110px;
    }

    .grid1 {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(100%, 1fr));
        grid-gap: 5px;
        align-items: stretch;
        color: #f5f7fc;
        margin-top: 5px
    }

    .grid1>articlegraph {
        border: 1px solid rgba(245, 247, 252, .02);
        box-shadow: 2px 2px 6px 0px rgba(0, 0, 0, 0.6);
        padding: 5px;
        font-size: 0.8em;
        -webkit-border-radius: 4px;
        border-radius: 4px;
        background: 0;
        -webkit-font-smoothing: antialiased;
        -moz-osx-font-smoothing: grayscale;
        height: 250px;
    }

    a {
        color: #aaa;
        text-decoration: none
    }


    blue {
        color: #01a4b4
    }

    orange {
        color: #009bb4
    }

    orange1 {
        position: relative;
        color: #009bb4;
        margin: 0 auto;
        text-align: center;
        margin-left: 5%;
        font-size: 1.1rem
    }

    green {
        color: #aaa
    }

    red {
        color: #f37867
    }

    red6 {
        color: #d65b4a
    }

    value {
        color: #fff
    }

    yellow {
        color: #CC0
    }

    purple {
        color: #916392
    }


    .hitempypos {
        position: absolute;
        margin-top: -60px;
        margin-left: 40px;
        margin-bottom: 20px;
        display: block;
    }


    .hitempd {
        position: absolute;
        font-family: weathertext2, Arial, Helvetica, sans-serif;
        background: rgba(86, 95, 103, 0.3);
        color: #aaa;
        font-size: 0.6rem;
        width: 130px;
        padding: 0;
        margin-left: 32px;
        padding-left: 3px;
        align-items: center;
        justify-content: center;
        display: block;
        margin-top: -20px;
    }

    .hitempd1 {
        position: absolute;
        font-family: weathertext2, Arial, Helvetica, sans-serif;
        background: rgba(86, 95, 103, 0.3);
        color: #aaa;
        font-size: 0.6rem;
        width: 130px;
        padding: 0;
        margin-left: 32px;
        padding-left: 3px;
        align-items: center;
        justify-content: center;
        display: block;
        margin-top: 10px;
    }

    .hitempd2 {
        position: absolute;
        font-family: weathertext2, Arial, Helvetica, sans-serif;
        background: rgba(86, 95, 103, 0.3);
        color: #aaa;
        font-size: 0.7rem;
        width: 90px;
        padding: 0;
        margin-left: 32px;
        padding-left: 3px;
        align-items: center;
        justify-content: center;
        display: block;
        margin-top: -15px;
    }


    .hitemp span {
        color: rgba(255, 124, 57, 1.000)
    }

    blue {
        color: rgba(0, 154, 171, 1.000)
    }

    .temperaturecontainer1 {
        position: absolute;
        left: 20px;
        margin-top: -5px;
        margin-bottom: 20px;
    }

    .temperaturecontainer2 {
        position: absolute;
        left: 20px;
        margin-top: 60px
    }

    .temperaturetrend1, .temperaturecaution, .temperaturetrend {
        position: absolute;
        font-size: 0.85rem
    }

    .temperaturetodayminus10, .temperaturetodayminus5, .temperaturetodayminus, .temperaturetoday0-5, .temperaturetoday6-10, .temperaturetoday11-15, .temperaturetoday16-20, .temperaturetoday21-25, .temperaturetoday26-30, .temperaturetoday31-35, .temperaturetoday36-40, .temperaturetoday41-45 {
        font-family: weathertext2, Arial, Helvetica, system;
        width: 3.75rem;
        height: 1.75rem;
        -webkit-border-radius: 3px;
        -moz-border-radius: 3px;
        -o-border-radius: 3px;
        display: flex;
        font-size: .95rem;
        padding-top: 2px;
        color: #fff;
        border-bottom: 10px solid rgba(56, 56, 60, 1);
        align-items: center;
        justify-content: center;
        border-radius: 3px;
        margin-bottom: 10px;
        -webkit-font-smoothing: antialiased;
        -moz-osx-font-smoothing: grayscale;
    }

    .temperaturetodayminus10 {
        background:#00a4b4;
        background: -webkit-linear-gradient(270deg, #00a4b4, rgba(80, 69, 188, 1));
        background: linear-gradient(270deg, #00a4b4, rgba(80, 69, 188, 1))
    }

    .temperaturetodayminus5 {
        background: #00a4b4;
        background: -webkit-linear-gradient(270deg, #00a4b4, rgba(80, 69, 188, 1));
        background: linear-gradient(270deg, #00a4b4, rgba(80, 69, 188, 1))
    }

    .temperaturetodayminus {
        background: #00a4b4;
        background: -webkit-linear-gradient(270deg, #00a4b4, rgba(80, 69, 188, 1));
        background: linear-gradient(270deg, #00a4b4, rgba(80, 69, 188, 1))
    }

    .temperaturetoday0-5 {
        background: #00a4b4;
        background: -webkit-linear-gradient(270deg, #00a4b4, #087883);
        background: linear-gradient(270deg, #00a4b4, #087883)
    }

    .temperaturetoday6-10 {
        background: #88b04b;
        background: -webkit-linear-gradient(90deg, #00a4b4, #88b04b);
        background: linear-gradient(90deg, #00a4b4, #88b04b)
    }

    .temperaturetoday11-15 {
        background: #e6a141;
        background: -webkit-linear-gradient(90deg, #90b12a, #e6a141);
        background: linear-gradient(90deg, #90b12a, #e6a141)
    }

    .temperaturetoday16-20 {
        background: #ff7c39;
        background: -webkit-linear-gradient(90deg, #90b12a, #ff7c39);
        background: linear-gradient(90deg, #90b12a, #ff7c39)
    }

    .temperaturetoday21-25 {
        background: #ff7c39;
        background: -webkit-linear-gradient(90deg, #e6a141, #ff7c39);
        background: linear-gradient(90deg, #e6a141, #ff7c39)
    }

    .temperaturetoday26-30 {
        background: #d05f2d;
        background: -webkit-linear-gradient(90deg, #d05f2d, rgba(236, 102, 21, 1));
        background: linear-gradient(90deg, #d05f2d, rgba(236, 102, 21, 1))
    }

    .temperaturetoday31-35 {
        background: #d86858;
        background: -webkit-linear-gradient(90deg, #d86858, rgba(211, 93, 78, .7));
        background: linear-gradient(90deg, #d86858, rgba(211, 93, 78, .7))
    }

    .temperaturetoday36-40 {
        background: #fd7641;
        background: -webkit-linear-gradient(90deg, #fd7641, #637ff6);
        background: linear-gradient(90deg, #fd7641, #637ff6)
    }

    .temperaturetoday41-45 {
        background: #de2c52;
        background: -webkit-linear-gradient(90deg, #de2c52, #637ff6);
        background: linear-gradient(90deg, #de2c52, #637ff6)
    }

    .temperaturetrend {
        margin-left: 17px;
        margin-top: -21px;
        z-index: 1;
        color: #fff;
        font-size: .6rem;
    }

    .temperaturetrend1 {
        margin-left: 15px;
        margin-top: -22px;
        z-index: 1;
        color: #fff;
        font-size: .6rem;
    }

    .temperaturetrend2 {
        margin-left: 0px;
        margin-top: -20px;
        z-index: 1;
        color: #fff;
        font-size: .5rem;
        position: absolute
    }

    smalluvunit {
        font-size: .7rem;
        font-family: weathertext2, Arial, Helvetica, system;
    }

    .w34convertrain {
        position: relative;
        font-size: .5em;
        top: 10px;
        color: #c0c0c0;
        margin-left: 5px
    }

    .hitempy {
        position: relative;
        background: rgba(61, 64, 66, 0.5);
        color: #aaa;
        width: 90px;
        padding: 1px;
        -webit-border-radius: 2px;
        border-radius: 2px;
        margin-top: -20px;
        margin-left: 98px;
        padding-left: 3px;
        line-height: 11px;
        font-size: 9px
    }

    .actualt {
        position: relative;
        left: 5px;
        -webkit-border-radius: 3px;
        -moz-border-radius: 3px;
        -o-border-radius: 3px;
        border-radius: 3px;
        background: rgba(74, 99, 111, 0.1);
        padding: 2px 5px 5px 5px;
        font-family: Arial, Helvetica, sans-serif;
        width: 195px;
        height: 1.7em;
        font-size: 0.8rem;
        color: #aaa;
        align-items: center;
        justify-content: center;
        margin-bottom: 10px;
        top: 0
    }

    .actualw {
        position: relative;
        left: 5px;
        -webkit-border-radius: 3px;
        -moz-border-radius: 3px;
        -o-border-radius: 3px;
        border-radius: 3px;
        background: rgba(74, 99, 111, 0.1);
        padding: 5px;
        font-family: Arial, Helvetica, sans-serif;
        width: 100px;
        height: 0.8em;
        font-size: 0.8rem;
        padding-top: 2px;
        color: #aaa;
        border-bottom: 2px solid rgba(56, 56, 60, 1);
        align-items: center;
        justify-content: center;
        margin-bottom: 10px;
        top: 0
    }

    .actual {
        position: relative;
        left: 5px;
        -webkit-border-radius: 3px;
        -moz-border-radius: 3px;
        -o-border-radius: 3px;
        border-radius: 3px;
        padding: 5px;
        font-family: Arial, Helvetica, sans-serif;
        width: 95%;
        height: 0.8em;
        font-size: 0.8rem;
        padding-top: 2px;
        color: #aaa;
        align-items: center;
        justify-content: center;
        margin-bottom: 10px;
        top: 0
    }
</style>