/**
 * Interface for a single data point with a date and value.
 */
export interface DataPoint {
    date: string;
    value: number;
}

/**
 * Interface for a data series, containing a key, color, and an array of data points.
 */
export interface DataSeries {
    values: DataPoint[];
    key: string;
    color: string;
}

/**
 * Type representing the full dataset, which is an array of data series.
 */
export type ChartData = DataSeries[];