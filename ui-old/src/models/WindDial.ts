export interface WindDialOptions {
    speed: number;
    direction: number; // In degrees, where 0 is North, 90 is East, etc.
    gusts?: number;
    color: string;
    radiusColor: string;
    tickColor: string;
}