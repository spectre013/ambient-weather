/* eslint-disable @typescript-eslint/no-explicit-any */
import * as d3 from "d3";
import { useEffect, useRef } from "react";

const Barchart = () => {
    const svgRef = useRef();

    useEffect(() => {
        // set the dimensions and margins of the graph
        const margin = {top: 10, right: 30, bottom: 30, left: 60},
            width = 460 - margin.left - margin.right,
            height = 400 - margin.top - margin.bottom;

// append the svg object to the body of the page
        const svg = d3.select(svgRef.current)
            .append("svg")
            .attr("width", width + margin.left + margin.right)
            .attr("height", height + margin.top + margin.bottom)
            .append("g")
            .attr("transform",
                "translate(" + margin.left + "," + margin.top + ")");

        // Parse the Data
        d3.json("http://localhost:5173/api/chart/tempf/day"
        ).then(function(result:any ) {
            const rows = result.data1;


            // Add X axis --> it is a date format
            const x = d3.scaleTime()
                .domain(d3.extent(rows, function (d) {
                    return d.date;
                }))
                .range([0, width]);

            svg.append("g")
                .attr("transform", "translate(0," + height + ")")
                .call(d3.axisBottom(x));

            // Add Y axis
            const y = d3.scaleLinear()
                .domain([0, d3.max(rows, function (d) {
                    return +d.value;
                })])
                .range([height, 0]);
            svg.append("g")
                .call(d3.axisLeft(y));

            // Add the line
            svg.append("path")
                .datum(result)
                .attr("fill", "none")
                .attr("stroke", "steelblue")
                .attr("stroke-width", 1.5)
                .attr("d", d3.line()
                    .x(function (d) {
                        return x(d.label)
                    })
                    .y(function (d) {
                        return y(d.label)
                    })
                )
        });
    }, []);

    return <svg width={460} height={400} id="barchart" ref={svgRef} />;
};

export default Barchart;