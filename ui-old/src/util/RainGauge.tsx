import React from 'react';

/**
 * Interface for the RainGauge component props.
 * @property {number} rainAmount - The amount of rain in inches (e.g., 0.5 for half an inch).
 * @property {number} size - The overall size of the SVG container.
 * @property {string} gaugeColor - The stroke color for the gauge outline.
 * @property {string} waterColor - The fill color for the water.
 * @property {string} tickColor - The stroke color for the tick marks.
 */
interface RainGaugeProps {
    rainAmount: number;
    size: number;
    gaugeColor: string;
    waterColor: string;
    tickColor: string;
}

/**
 * A React component that renders a dynamic SVG of a rain gauge filling with water.
 * The gauge has an open top and includes tick marks for every tenth of an inch.
 */
const RainGauge: React.FC<RainGaugeProps> = ({ rainAmount, size, gaugeColor, waterColor, tickColor }) => {
    // Define dimensions relative to the size prop for scalability
    const gaugeWidth = size * 0.4;
    const gaugeHeight = size * 0.8;
    const gaugeX = (size - gaugeWidth) / 2;
    const gaugeY = 10;
    const maxRainfallInches = 2.0; // The maximum rainfall the gauge can hold
    const tickMarksPerInch = 10;
    // Calculate the height of the water based on the rainfall amount
    const waterHeight = Math.min(
        (rainAmount / maxRainfallInches) * gaugeHeight,
        gaugeHeight
    );

    // Calculate the vertical position of the water
    const waterY = gaugeY + gaugeHeight - waterHeight;

    // Array to hold the tick mark elements
    const tickMarks = [];
    const totalTicks = maxRainfallInches * tickMarksPerInch;

    // Loop to create tick marks and labels
    for (let i = 0; i <= totalTicks; i++) {
        // Calculate the Y position for the current tick
        const tickY = gaugeY + gaugeHeight - (i / totalTicks) * gaugeHeight;
        const tickLength = i % tickMarksPerInch === 0 ? size * 0.05 : size * 0.02; // longer for inch marks

        // Create a line for the tick mark
        tickMarks.push(
            <line
                key={`tick-${i}`}
                x1={gaugeX + gaugeWidth}
                y1={tickY}
                x2={gaugeX + gaugeWidth + tickLength}
                y2={tickY}
                stroke={tickColor}
                strokeWidth="1"
            />
        );

        // Create a label for every full inch
        if (i % tickMarksPerInch === 0) {
            const labelText = (i / tickMarksPerInch).toFixed(1);
            tickMarks.push(
                <text
                    key={`label-${i}`}
                    x={gaugeX + gaugeWidth + tickLength + size * 0.01}
                    y={tickY + size * 0.015} // Adjust text position
                    fill={tickColor}
                    fontSize={size * 0.04}
                    textAnchor="start"
                >
                    {labelText}
                </text>
            );
        }
    }

    return (
        <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`}>
            {/* Water element (dynamic rectangle) */}
            <rect
                x={gaugeX}
                y={waterY}
                width={gaugeWidth}
                height={waterHeight}
                fill={waterColor}
            />

            {/* Rain gauge outline */}
            <rect
                x={gaugeX}
                y={gaugeY}
                width={gaugeWidth}
                height={gaugeHeight}
                fill="transparent"
                stroke={gaugeColor}
                strokeWidth="2"
                rx="5"
                ry="5"
            />

            {/* Render the tick marks */}
            {tickMarks}
        </svg>
    );
};

export default RainGauge;