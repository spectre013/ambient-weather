import React from 'react';

/**
 * Interface for the CircleGauge component props.
 * @property {number} size - The overall size of the SVG container (e.g., 100 for a 100x100px SVG).
 * @property {number} value - The number to display in the center of the circle.
 * @property {string} color - The stroke color for the outer ring.
 */
interface CircleGaugeProps {
    size: number;
    value: number;
    color: string;
}

/**
 * A React component that renders an SVG of a circle with a number inside.
 * The outer ring's color and the central number are customizable via props.
 */
const CircleGauge: React.FC<CircleGaugeProps> = ({ size, value, color }) => {
    const strokeWidth = size / 2 * 0.1;
    const radius = (size / 2) - (strokeWidth / 2);
    const center = size / 2;

    return (
        <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`}>
            {/* Outer ring */}
            <circle
                cx={center}
                cy={center}
                r={radius}
                fill="none"
                stroke={color}
                strokeWidth={strokeWidth}
            />

            {/* Background circle for the number */}
            <circle
                cx={center}
                cy={center}
                r={radius - (strokeWidth / 2)}
                fill="#121212"
            />

            {/* Number in the center */}
            <text
                x={center}
                y={center}
                textAnchor="middle"
                dominantBaseline="middle"
                fontSize={`${size * 0.4}px`}
                fill={color}
                fontWeight="bold"
            >
                {value}
            </text>
        </svg>
    );
};

export default CircleGauge;
