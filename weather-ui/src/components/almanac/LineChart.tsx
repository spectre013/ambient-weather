import { useEffect, useRef } from 'react';
import axios from 'axios';
import * as d3 from 'd3';

type Data = {
  data1: DataPoint[];
}

type DataPoint = {
  label: string;
  y: number;
};

type LineChartProps = {
  apiUrl: string;
};

const LineChart = ({ apiUrl }: LineChartProps) => {
  const svgRef = useRef<SVGSVGElement>(null);
  const dataRef = useRef<Data[]>([]);


  useEffect(() => {
    async function fetchData() {
      try {
        const response = await axios.get(apiUrl);
        dataRef.current = response.data.data1;
      } catch (error) {
        console.error('Error fetching data:', error);
      }
    }
    fetchData().then(() => {

      const parseTime = d3.timeParse('%H:%M:%S');
      const formatTime = d3.timeFormat('%H:%M:%S');

      const xScale = d3
          .scaleTime()
          .domain(d3.extent(dataRef.current, (d) => parseTime(d.label)))
          .range([50, 400]);

      const yScale = d3
          .scaleLinear()
          .domain([Math.min(...dataRef.current.map((d) => d.y)), Math.max(...dataRef.current.map((d) => d.y))])
          .range([200, 50]);

      const lineGenerator = d3
          .line<DataPoint>()
          .x((d) => xScale(parseTime(d.label)))
          .y((d) => yScale(d.y));
      console.log("Drawing chart");
      drawChart(xScale, yScale, lineGenerator);
    });
  }, [apiUrl]);

  function drawChart(
      xScale: d3.ScaleTime<number, number>,
      yScale: d3.ScaleLinear<number, number>,
      lineGenerator: d3.Line<DataPoint>
  ) {
    console.log("Drawing chart");
    const svg = d3.select(svgRef.current);

    // Clear previous chart
    svg.selectAll('*').remove();

    // Add axes and grid lines
    svg
        .append('g')
        .attr('transform', 'translate(50,200)')
        .call(d3.axisBottom(xScale).ticks(5));

    svg
        .append('g')
        .attr('transform', 'translate(50,0)')
        .call(d3.axisLeft(yScale))
        .append('text')
        .attr('text-anchor', 'middle')
        .attr('x', 0)
        .attr('y', -10)
        .text('Value');

    svg
        .append('g')
        .attr('class', 'grid')
        .call(
            d3
                .axisBottom(xScale)
                .ticks(5)
                .tickSizeOuter(0)
                .tickFormat()
                .tickSize(-5)
                .tickPadding(10)
        );

    // Add line path
    svg
        .append('path')
        .datum(dataRef.current)
        .attr('d', lineGenerator)
        .style('fill', 'none')
        .style('stroke', '#3366cc');

    // Add data points
    const dots = svg
        .selectAll('circle')
        .data(dataRef.current)
        .enter()
        .append('circle')
        .attr('r', 4)
        .attr('fill', 'white');

    dots
        .transition()
        .duration(100)
        .ease(d3.easeBounceOut)
        .delay((d, i) => i * 50)
        .attr('cx', (d) => xScale(parseTime(d.label)))
        .attr('cy', (d) => yScale(d.y));

    // Add mouse-over functionality
    dots.on('mouseover', (event, d) => {
      showTooltip(d3.pointer(event)[0], d);
    });

    dots.on('mouseout', hideTooltip);

    function showTooltip(x: number, dataPoint: DataPoint): void {
      svg.append('text')
          .attr('x', x)
          .attr('y', yScale(dataPoint.y) - 20)
          .style('alignment-baseline', 'middle')
          .style('fill', '#fff')
          .style('pointer-events', 'none')
          .text(`${formatTime(dataPoint.label)}: ${dataPoint.y.toFixed(1)}`);
    }

    function hideTooltip(): void {
      svg.selectAll('.tooltip').remove();
    }
  }

  return <svg ref={svgRef} width={800} height={600}></svg>;
};

export default LineChart;