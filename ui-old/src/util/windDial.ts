import { WindDialOptions } from '../models/WindDial.ts';

class WindDial {
    private options: WindDialOptions;
    private readonly svgWidth = 250;
    private readonly svgHeight = 250;
    private readonly cx = this.svgWidth / 2;
    private readonly cy = this.svgHeight / 2;
    private readonly outerRadius = 100;
    private readonly innerRadius = 75;

    /**
     * Constructs a new WindDial instance.
     * @param options The wind data to display.
     */
    constructor(options: WindDialOptions) {
        this.options = options;
    }

    /**
     * Generates the SVG string for the wind dial.
     * @returns The complete SVG string.
     */
    public generateSvg(): string {
        const { speed, direction, gusts, color, radiusColor, tickColor} = this.options;

        // The core SVG structure and styling
        const svgContent = `
      <svg width="${this.svgWidth}" height="${this.svgHeight}" viewBox="0 0 ${this.svgWidth} ${this.svgHeight}" xmlns="http://www.w3.org/2000/svg">
        <style>
          .wind-dial-text {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
            fill: #fff;
            text-anchor: middle;
            dominant-baseline: central;
          }
        </style>

        <!-- Outer translucent ring -->
        <circle cx="${this.cx}" cy="${this.cy}" r="${this.outerRadius}" fill="none" stroke=${radiusColor} stroke-width="20" />

        <!-- Compass points (N, E, S, W) -->
        <text class="wind-dial-text" fill="${tickColor}" stroke="${tickColor}" font-size="20" x="${this.cx}" y="${this.cy - this.outerRadius - 15}">N</text>
        <text class="wind-dial-text" fill="${tickColor}" stroke="${tickColor}" font-size="20" x="${this.cx + this.outerRadius + 15}" y="${this.cy}">E</text>
        <text class="wind-dial-text" fill="${tickColor}" stroke="${tickColor}" font-size="20" x="${this.cx}" y="${this.cy + this.outerRadius + 15}">S</text>
        <text class="wind-dial-text" fill="${tickColor}" stroke="${tickColor}" font-size="20" x="${this.cx - this.outerRadius - 15}" y="${this.cy}">W</text>

        <!-- Dynamic tick marks -->
        ${this.generateTickMarks(tickColor)}

        <!-- Central circle -->
        <circle cx="${this.cx}" cy="${this.cy}" r="${this.innerRadius}" fill=${radiusColor} />

        <!-- Wind speed text -->
        <text class="wind-dial-text" fill="${color}" stroke="${color}" font-size="40" font-weight="bold" x="${this.cx}" y="${this.cy - 10}">${speed}</text>
        <text class="wind-dial-text" fill="${color}" stroke="${color}" font-size="20" x="${this.cx}" y="${this.cy + 20}">mph</text>

        <!-- Wind direction arrow -->
        ${this.generateArrow(direction, color)}

        <!-- Gusts text -->
        ${gusts ? `<text class="wind-dial-text" font-size="16" x="${this.cx}" y="${this.cy + this.outerRadius + 40}">Gusts: ${gusts} mph</text>` : ''}
      </svg>
    `;

        return svgContent;
    }

    /**
     * Helper function to generate the tick marks around the dial.
     * @returns A string containing SVG path elements for the ticks.
     */
    private generateTickMarks(tickColor: string): string {
        let ticks = '';
        const numTicks = 36;
        const majorTicks = 4; // N, E, S, W
        const tickLength = 5;
        const majorTickLength = 10;
        const outerRadius = this.outerRadius;
        //const innerRadius = this.outerRadius - tickLength;

        for (let i = 0; i < numTicks; i++) {
            const angle = (i * 360 / numTicks) * (Math.PI / 180);
            const isMajor = i % (numTicks / majorTicks) === 0;
            const currentTickLength = isMajor ? majorTickLength : tickLength;
            const currentInnerRadius = this.outerRadius - currentTickLength;

            const x1 = this.cx + outerRadius * Math.sin(angle);
            const y1 = this.cy - outerRadius * Math.cos(angle);
            const x2 = this.cx + currentInnerRadius * Math.sin(angle);
            const y2 = this.cy - currentInnerRadius * Math.cos(angle);

            ticks += `<line x1="${x1}" y1="${y1}" x2="${x2}" y2="${y2}" stroke=${tickColor} stroke-width="2" stroke-linecap="round" />`;
        }
        return ticks;
    }

    /**
     * Helper function to generate the wind direction arrow.
     * @param direction The direction in degrees.
     * @param color the color of the arrow
     * @returns A string containing the SVG path element for the arrow.
     */
    private generateArrow(direction: number, color: string): string {
        // The arrow is an SVG path. The "transform" attribute is used to rotate it.
        // We create two segments to form a gap matching the inner circle's diameter.
        const transform = `rotate(${direction - 90}, ${this.cx}, ${this.cy})`;
        return `
      <path d="M${this.cx + this.innerRadius},${this.cy} L${this.cx + this.outerRadius},${this.cy}"
            fill="none"
            stroke="${color}"
            stroke-width="5"
            stroke-linecap="round"
            transform="${transform}" />
      <path d="M${this.cx - this.innerRadius},${this.cy} L${this.cx - this.outerRadius},${this.cy}"
            fill="none"
            stroke="${color}"
            stroke-width="5"    
            stroke-linecap="round"
            transform="${transform}" />
      <path d="M${this.cx + this.outerRadius},${this.cy} L${this.cx + this.outerRadius - 10},${this.cy - 5} M${this.cx + this.outerRadius},${this.cy} L${this.cx + this.outerRadius - 10},${this.cy + 5}"
            fill="none"
            stroke="${color}"
            stroke-width="5"
            stroke-linecap="round"
            transform="${transform}" />
    `;
    }
}

export default WindDial;