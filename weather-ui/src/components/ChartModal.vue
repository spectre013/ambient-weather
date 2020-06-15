<template>
    <div class="modal-backdrop">
        <div class="modal">
            <section class="modal-body">
                <div id="body" v-if="data && minmax">
                    <div class="weather34darkbrowser" :url="url"></div>
                    <div class="charDiv">
                        <div id="chartContainer" class="chartContainer"></div>
                    </div>
                    <div class="weather34browser-footer">
                        <span class="copyw">&nbsp;
                            <svg id="i-external" viewBox="0 0 32 32" width="10" height="10" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M14 9 L3 9 3 29 23 29 23 18 M18 4 L28 4 28 14 M28 4 L14 18"></path></svg>
                            <a href="https://canvasjs.com" title="https://canvasjs.com" target="_blank">CanvasJs.com v2.3.1 GA (CC BY-NC 3.0) Non-Commercial-Version </a>
                        </span>
                    </div>
                </div>
                <button class="lity-close" type="button" aria-label="Close (Press escape to close)" @click="close" data-lity-close="">×</button>
            </section>
        </div>
    </div>
</template>

<script>
    /* global CanvasJS */
    import axios from 'axios';
    import moment from 'moment';

    export default {
        name: 'chartmodal',
        props: {
            options: Object,
        },
        data () {
            return {
                data: null,
                minmax: null,
                interval: null
            }
        },
        mounted: function() {
            if (this.options.type === 'temp') {
                axios.get('/api/chart/' + this.options.field + '/' + this.options.time).then(response => (this.data = response.data));
                axios.get('/api/minmax/tempf').then(response => (this.minmax = response.data));
            }

            if(this.options.time === 'year') {
                this.interval = 'month';
            } else if(this.options.time === 'month') {
                this.interval = 'day';
            } else {
                this.interval = 'hour';
            }
        },
        filters: {
          y: function(date) {
              moment(date).format("Y");
          }
        },
        computed: {
            url: function() {
                return `Temperature | High: ${this.minmax.max[this.options.time].value} &deg;F Low: ${this.minmax.min[this.options.time].value} &deg;F`
            }
        },
        methods: {
            close: function() {
               this.$parent.closeModal('chart');
            },
            drawChart: function(dataPoints1, dataPoints2,interval) {
                let chart = new CanvasJS.Chart("chartContainer", {
                    backgroundColor: 'rgba(40, 45, 52,.4)',
                    animationEnabled: true,
                    animationDuration: 500,


                    title: {
                        text: " ",
                        fontSize: 11,
                        fontColor: '#ccc',
                        fontFamily: "arial",
                    },
                    toolTip: {
                        fontStyle: "normal",
                        cornerRadius: 4,
                        backgroundColor: 'rgba(37, 41, 45, 0.95)',
                        contentFormatter: function (e) {
                            let str = '<span style="color: #ccc;">' + e.entries[0].dataPoint.label + '</span><br/>';
                            for (let i = 0; i < e.entries.length; i++) {
                                let temp = '<span style="color: ' + e.entries[i].dataSeries.color + ';">' + e.entries[i].dataSeries.name + '</span> <span style="color: #ccc;">' + e.entries[i].dataPoint.y.toFixed(1) + " &deg;F" + '</span> <br/>';
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
                        minimum: -0.5,
                        interval: 'auto',
                        intervalType: interval,
                        xValueType: "dateTime",
                        crosshair: {
                            enabled: true,
                            snapToDataPoint: true,
                            color: '#009bab',
                            labelFontColor: "#F8F8F8",
                            labelFontSize: 11,
                            labelBackgroundColor: '#009bab',
                        }

                    },

                    axisY: {
                        title: "Temperature (&deg;F) Recorded",
                        titleFontColor: '#ccc',
                        titleFontSize: 10,
                        titleWrap: false,
                        margin: 10,
                        interval: 'auto',
                        lineThickness: 1,
                        gridThickness: 1,
                        includeZero: false,
                        gridColor: 'RGBA(64, 65, 66, 0.8)',
                        gridDashType: "dot",
                        labelFontSize: 11,
                        labelFontColor: '#ccc',
                        titleFontFamily: "arial",
                        labelFontFamily: "arial",
                        labelFormatter: function (e) {
                            return e.value.toFixed(0) + " &deg;F ";
                        },
                        crosshair: {
                            enabled: true,
                            snapToDataPoint: true,
                            color: '#ff832f',
                            labelFontColor: "#F8F8F8",
                            labelFontSize: 11,
                            labelBackgroundColor: '#ff832f',
                            valueFormatString: "#0.# &deg;F",
                        }

                    },

                    legend: {
                        fontFamily: "arial",
                        fontColor: '#ccc',

                    },


                    data: [
                        {
                            type: "splineArea",
                            color: 'rgba(255, 148, 82, 0.95)',
                            lineColor: 'rgba(255, 131, 47, 1)',
                            markerSize: 0,
                            showInLegend: true,
                            legendMarkerType: "circle",
                            lineThickness: 2,
                            markerType: "circle",
                            name: " Hi Temp",
                            dataPoints: dataPoints1,
                            yValueFormatString: "#0.# &deg;F",

                        },
                        {

                            type: "splineArea",
                            color: 'rgba(0, 164, 180, 1)',
                            markerSize: 0,
                            markerColor: '#007181',
                            showInLegend: true,
                            legendMarkerType: "circle",
                            lineThickness: 2,
                            lineColor: '#007181',
                            markerType: "circle",
                            name: " Lo Temp",
                            dataPoints: dataPoints2,
                            yValueFormatString: "#0.0 &deg;F",

                        }

                    ]
                });

                chart.render();
            }
        },
        watch: {
            data: function() {
                function checkDom(that){
                    if(document.querySelector("#chartContainer")) {
                        that.drawChart(that.data.data1,that.data.data2,that.interval);
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
    .copyw {
        position:absolute;
        color:#aaa;font-family:arial;
        padding-top:5px;margin-left:25px;
        display:block;
        margin-top:12px;
    }

    .charDiv {
        width:880px;
        background:0;
        padding:0px;
        margin-left:5px;
        font-size: 12px;
        border-radius:3px;
    }

    a {
        color: #aaa;
        text-decoration: none
    }

    body {
        font-size: 12px
    }


    .weather34browser-footer {
        font-size: 0.6rem
    }


    .chartContainer {
        width: 99%;
        height: 450px;
        border: 3px solid rgba(86, 95, 103, 0.7);
        -webkit-border-radius: 3px;
        border-radius: 3px;
        margin-left: 10px;
    }

    @media screen and (max-width:640px) {
        .chartContainer {
            width: 155vh;
            height: 60vh;
            border: 3px solid rgba(86, 95, 103, 0.7);
            -webkit-border-radius: 3px;
            border-radius: 3px;
        }
        blue {
            color: #009bab
        }
        .weather34darkbrowser {
            width: 170vh;
            height: 1vh;
        }
        .weather34darkbrowser[url]:after {
            font-size: .4rem;
            padding: 5px 10px;
        }
        .weather34browser-footer {
            font-size: .4rem
        }
    }
</style>