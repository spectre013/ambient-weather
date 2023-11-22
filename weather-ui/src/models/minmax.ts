export interface MinMax {
    avg: Avg
    max: Max
    min: Min
}

export interface Avg {
    day: Day
    month: Month
    year: Year
    yesterday: Yesterday
}

export interface Max {
    day: Day
    month: Month
    year: Year
    yesterday: Yesterday
}

export interface Min {
    day: Day
    month: Month
    year: Year
    yesterday: Yesterday
}

export interface Day {
    value: number
    date: string
}

export interface Month {
    value: number
    date: string
}

export interface Year {
    value: number
    date: string
}

export interface Yesterday {
    value: number
    date: string
}


