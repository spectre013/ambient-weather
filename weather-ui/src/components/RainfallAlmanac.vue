<template>
    <div class="modal-backdrop">
        <div class="modal">
            <section class="modal-body">
                <div id="body" v-if="current">
                    <div class="weather34darkbrowser" url="Rainfall Almanac"></div>

                    <main class="grid">
                        <article>
                            <div class="actualt">Rainfall Today</div>
                            <div class="rainfalltoday1">3.3<smalluvunit>mm</smalluvunit>
                                <div class="w34convertrain">0.13<smalluvunit>in</smalluvunit></div>
                            </div>
                            <div class="hitempy"><svg id="i-info" viewBox="0 0 32 32" width="8" height="8" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M16 14 L16 23 M16 8 L16 10"></path><circle cx="16" cy="16" r="14"></circle></svg>
                                Last Hour<blue> 0.0</blue> mm
                            </div>
                        </article>

                        <article>
                            <div class="actualt">Rainfall Yesterday</div>
                            <div class="rainfalltoday1">0.0<smalluvunit>mm</smalluvunit>
                                <div class="w34convertrain">
                                0.00<smalluvunit>in</smalluvunit><div>
                            </div>

                                <div class="hitempy"><svg id="i-info" viewBox="0 0 32 32" width="8" height="8" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M16 14 L16 23 M16 8 L16 10"></path><circle cx="16" cy="16" r="14"></circle></svg> Last 24 Hours<br><blue>3.3</blue> mm</div>
                            </div></div>
                        </article>


                        <article>
                            <div class="actualt">Rainfall Jul 2020 </div>
                            <div class="rainfalltoday1">3.3<smalluvunit>mm</smalluvunit><div class="w34convertrain">
                                0.13<smalluvunit>in</smalluvunit><div></div>

                                <div class="hitempy">
                                    <svg id="i-info" viewBox="0 0 32 32" width="8" height="8" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M16 14 L16 23 M16 8 L16 10"></path><circle cx="16" cy="16" r="14"></circle></svg> Last Rainfall<br>
                                    <blue>Today</blue>
                                </div>
                            </div></div></article>


                        <article>
                            <div class="actualt">Rainfall 2020 </div>
                            <div class="rainfalltoday1">92.0<smalluvunit>mm</smalluvunit><div class="w34convertrain">
                                3.62<smalluvunit>in</smalluvunit><div></div>

                                <div class="hitempy"><svg id="i-info" viewBox="0 0 32 32" width="8" height="8" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M16 14 L16 23 M16 8 L16 10"></path><circle cx="16" cy="16" r="14"></circle></svg> Since<br>
                                    <blue>Jan 2020</blue>
                                </div>
                            </div></div></article>

                        <article>
                            <div class="actualt">&nbsp;Rainfall All-Time </div>
                            <div class="rainfalltoday1">92.0<smalluvunit>mm</smalluvunit><div class="w34convertrain">
                                3.62 <smalluvunit>in</smalluvunit><div></div>

                                <div class="hitempy"><svg id="i-info" viewBox="0 0 32 32" width="8" height="8" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M16 14 L16 23 M16 8 L16 10"></path><circle cx="16" cy="16" r="14"></circle></svg> Since<br>
                                    <blue>All Time</blue>
                                </div>
                            </div></div></article>
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
        name: 'rainfallalmanac',
        props: {
            current: Object,
            appName: {
                type: String,
                default: String
            },
            options: Object,
        },
        data () {
            return {
                data: null
            }
        },
        mounted: function() {
            axios.get('/api/chart/dailyrainin/year').then(response => (this.data = response.data));
            axios.get('/api/alltime/max/dailyrainin').then(response => (this.maxrain = response.data));
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
                this.$parent.closeModal('rainfallalmanac');
            },
            drawChart: function( dataPoints1 ) {

                let chart = new CanvasJS.Chart("chartContainer2", {
                    backgroundColor: 'rgba(40, 45, 52,.4)',
                    animationEnabled: true,
                    animationDuration: 500,
                    margin: 0,
                    height: 360,
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
                            let str = '<span style="color: #ccc;">' + e.entries[0].dataPoint.label + '</span><br/>';
                            for (let i = 0; i < e.entries.length; i++) {
                                let temp = '<span style="color: ' + e.entries[i].dataSeries.color + ';">' + e.entries[i].dataSeries.name + '</span> <span style="color: #ccc;">' + e.entries[i].dataPoint.y.toFixed(1) + " in" + '</span> <br/>';
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
                        minimum:0,
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
                        minimum:0,
                        gridColor: 'RGBA(64, 65, 66, 0.8)',
                        gridDashType: "dot",
                        labelFontSize: 10,
                        labelFontColor: '#ccc',
                        titleFontFamily: "arial",
                        labelFontFamily: "arial",
                        crosshair: {
                            enabled: true,
                            snapToDataPoint: true,
                            color: '#44a6b5',
                            labelFontColor: "#F8F8F8",
                            labelFontSize:10,
                            labelMaxWidth: 70,
                            labelBackgroundColor: '#44a6b5',
                            valueFormatString: "#0.0 in",
                        }

                    },
                    legend:{
                        fontFamily: "arial",
                        fontColor: '#ccc',

                    },

                    data: [
                        {
                            type: "column",
                            color: 'rgba(66, 166, 181, 0.95)',
                            lineColor: 'rgba(66, 166, 181, 1)',
                            markerSize:0,
                            showInLegend:false,
                            legendMarkerType: "circle",
                            lineThickness: 2,
                            markerType: "circle",
                            name:" in Precip",
                            dataPoints: dataPoints1,
                            yValueFormatString: "#0.# in",

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
                        that.drawChart(that.data.data1);
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

    .grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(140px, 2fr));
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
        height: 85px;
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
        height: 390px
    }



    .actualt {
        position: relative;
        left: 5px;
        -webkit-border-radius: 3px;
        -moz-border-radius: 3px;
        -o-border-radius: 3px;
        border-radius: 3px;
        background: rgba(74, 99, 111, 0.1);
        padding: 5px;
        font-family: Arial, Helvetica, sans-serif;
        width: 120px;
        height: 0.8em;
        font-size: 0.8rem;
        padding-top: 2px;
        color: #aaa;
        align-items: center;
        justify-content: center;
        margin-bottom: 10px;
        top: 0;
        text-align: center;
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

    .rainfallcontainer1 {
        left: 5px;
        top: 0
    }

    .rainfalltoday1 {
        font-family: weathertext2, Arial, Helvetica, system;
        width: 3rem;
        height: 2.5rem;
        -webkit-border-radius: 3px;
        -moz-border-radius: 3px;
        -o-border-radius: 3px;
        font-weight: normal;
        font-size: .9rem;
        padding-top: 5px;
        color: #fff;
        border-bottom: 12px solid rgba(56, 56, 60, 1);
        align-items: center;
        justify-content: center;
        text-align: center;
        border-radius: 3px;
        background: rgba(68, 166, 181, 1.000);
        -webkit-font-smoothing: antialiased;
        -moz-osx-font-smoothing: grayscale;
    }

    .rainfallcaution,
    .rainfalltrend {
        position: absolute;
        font-size: 1rem
    }

    smalluvunit {
        font-size: .6rem;
        font-family: Arial, Helvetica, system;
    }

    metricsblue {
        color: #44a6b5;
        font-family: "weathertext2", Helvetica, Arial, sans-serif;
        background: rgba(86, 95, 103, 0.5);
        -webkit-border-radius: 2px;
        border-radius: 2px;
        align-items: center;
        justify-content: center;
        font-size: .9em;
        left: 10px;
        padding: 0 3px 0 3px;
    }

    .w34convertrain {
        position: relative;
        font-size: .8em;
        top: 1px;
        color: #c0c0c0;
        margin-left: auto;
        margin-right: auto;
    }

    .hitempy {
        position: relative;
        background: rgba(61, 64, 66, 0.5);
        color: #aaa;
        width: 75px;
        padding: 1px;
        -webit-border-radius: 2px;
        border-radius: 2px;
        margin-top: -40px;
        margin-left: 56px;
        padding-left: 3px;
        line-height: 11px;
        font-size: 9px
    }
</style>