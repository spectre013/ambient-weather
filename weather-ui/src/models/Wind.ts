export interface Wind {
    avg: Avg
    dir: Dir
    gust: Gust
    wind: Wind
}

export interface Avg {
    date: string
    value: number
}

export interface Dir {
    date: string
    value: number
}

export interface Gust {
    date: string
    value: number
}

export interface Wind {
    date: string
    value: number
}