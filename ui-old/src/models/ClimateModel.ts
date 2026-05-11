/**
 * @interface ClimateData
 * Corresponds to the GoLang ClimateData struct.
 * Contains arrays for each climate metric for an entire year.
 * The index of the array corresponds to the month (e.g., index 1 for January).
 */
export interface ClimateData {
    avgrain: number[];
    avgtemp: number[];
    maxtemp: number[];
    mintemp: number[];
}

/**
 * @interface ClimateRecord
 * Corresponds to the GoLang ClimateRecord struct.
 * Holds all the climate data for a single year.
 */
export interface ClimateRecord {
    Year: number;
    Data: ClimateData;
}

export interface FrostDateRecord {
    year: number;
    spring: string;
    fall: string;
}

export type weatherDisplay = (arg1: number, arg2: string) => string;