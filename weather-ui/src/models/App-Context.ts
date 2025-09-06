export interface ThemeColors {
    bg: string;
    "card-bg": string;
    "main-card": string;
    "text-color": string;
    "muted-text": string;
    "border-color": string;
    "gradient-1": string;
    "gradient-2": string;
    "gradient-3": string;
    "dial-bg": string;
}

export interface WeatherContextModel {
    theme: "light" | "dark";
    longname: string;
    shortname: string;
    state: string;
    country: string;
    light: ThemeColors;
    dark: ThemeColors;
}