import React from 'react';

// Function to convert polar coordinates to Cartesian
const polarToCartesian = (centerX: number, centerY: number, radius: number, angleInDegrees: number) => {
    const angleInRadians = (angleInDegrees - 90) * Math.PI / 180.0;
    return {
        x: centerX + (radius * Math.cos(angleInRadians)),
        y: centerY + (radius * Math.sin(angleInRadians))
    };
};

// Function to generate the path for an arc segment
const describeArc = (x: number, y: number, radius: number, startAngle: number, endAngle: number) => {
    const start = polarToCartesian(x, y, radius, endAngle);
    const end = polarToCartesian(x, y, radius, startAngle);
    const largeArcFlag = endAngle - startAngle <= 180 ? '0' : '1';
    const d = [
        'M', start.x, start.y,
        'A', radius, radius, 0, largeArcFlag, 0, end.x, end.y
    ].join(' ');
    return d;
};

interface GaugeProps {
    size: number;
    min: number;
    max: number;
    value: number;
    startColor: string;
    endColor: string;
    currentColor: string;
}

const Gauge: React.FC<GaugeProps> = ({ size, min, max, value, startColor, endColor, currentColor }) => {
    const strokeWidth = size * 0.075; // dynamic stroke based on size
    const radius = (size / 2) - (strokeWidth / 2);
    const center = size / 2;

    // Angle calculations for a 270-degree arc
    const startAngle = -135;
    const endAngle = 135;

    // Calculate the angle for the current value
    const valueAngle = startAngle + ((value - min) / (max - min)) * (endAngle - startAngle);

    // Calculate the coordinates for the small circle indicator
    const indicator = polarToCartesian(center, center, radius, valueAngle);

    return (
        <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`}>
            <defs>
                <linearGradient id="gaugeGradient" x1="0%" y1="0%" x2="100%" y2="0%">
                    <stop offset="0%" style={{ stopColor: startColor, stopOpacity: 1 }} />
                    <stop offset="100%" style={{ stopColor: endColor, stopOpacity: 1 }} />
                </linearGradient>
            </defs>

            {/* Main Gauge Arc */}
            <path
                d={describeArc(center, center, radius, startAngle, endAngle)}
                fill="none"
                stroke="url(#gaugeGradient)"
                strokeWidth={strokeWidth}
                strokeLinecap="round"
            />

            {/* Value Indicator Circle */}
            <circle
                cx={indicator.x}
                cy={indicator.y}
                r={strokeWidth / 2} // size of the indicator circle
                fill="#fff"
            />

            {/* Main Number */}
            <text
                x={center}
                y={center}
                textAnchor="middle"
                dominantBaseline="middle"
                fontSize={`${size * 0.55}px`} // dynamic font size
                fill={currentColor}
                fontWeight="bold"
            >
                {value}
            </text>

            {/* Lower Numbers */}
            <text
                x={center - (size * 0.15)}
                y={center + (size * 0.4)}
                textAnchor="middle"
                fontSize={`${size * 0.15}px`} // dynamic font size
                fill={startColor}
            >
                {min}
            </text>
            <text
                x={center + (size * 0.15)}
                y={center + (size * 0.4)}
                textAnchor="middle"
                fontSize={`${size * 0.15}px`} // dynamic font size
                fill={endColor}
            >
                {max}
            </text>
        </svg>
    );
};

export default Gauge;