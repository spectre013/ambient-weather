import React, { useEffect, useRef } from 'react';
import * as d3 from 'd3';
import nv from 'nvd3';
import 'nvd3/build/nv.d3.css';

const TemperatureChart = () => {
    const chartRef = useRef(null);

    useEffect(() => {
        fetch("api/chart/temperature/1h")
            .then(response => response.json())
            .then(data => {
                data.forEach(series => {
                    series.values.forEach(d => {
                        d.x = new Date(d.date); // Convert date string to Date object
                        d.y = d.value;
                    });
                });

                const chart = nv.models.lineChart()
                    .useInteractiveGuideline(true) // Enable tooltips
                    .showLegend(true)
                    .showYAxis(true)
                    .showXAxis(true);

                chart.xAxis
                    .axisLabel('Time')
                    .tickFormat(d => d3.timeFormat('%H:%M')(new Date(d)));

                chart.yAxis
                    .axisLabel('Temperature (°F)')
                    .tickFormat(d3.format('.1f'));

                d3.select(chartRef.current)
                    .datum(data)
                    .call(chart);

                nv.utils.windowResize(chart.update);
            });

        // Apply dark mode styles
        document.body.style.backgroundColor = '#121212';
        document.body.style.color = '#ffffff';
        d3.selectAll('text').style('fill', '#ffffff');
    }, []);

    return (
        <div id="chart">
            <svg ref={chartRef} width="100%" height="500px"></svg>
        </div>
    );
};

export default TemperatureChart;
