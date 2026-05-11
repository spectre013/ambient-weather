import React from 'react';
import {WeatherContextModel } from './models/App-Context.ts'
export const weatherContext = {
    theme:"light",
    longname:"Lorson Ranch - Colorado Springs",
    shortname:"Lorson Ranch",
    state:"CO",
    country:"US",
    light:{
        "bg": "#F9F6EE",
        "card-bg": "#6c6c6c",
        "main-card": "#f3f4f6",
        "text-color": "#1a202c",
        "muted-text": "#9E9E9E",
        "border-color": "rgba(0, 0, 0, 0.5)",
        "gradient-1": "#38bdf8",
        "gradient-2": "#818cf8",
        "gradient-3": "#fcd34d",
        "dial-bg": "rgba(255, 255, 255, 0.5)"
    },
    dark:{
        "bg": "#0d1219",
        "card-bg": "#141b21",
        "main-card": "#141b21",
        "text-color": "#e3e9f4",
        "muted-text": "#6b7280",
        "border-color": "rgba(255, 255, 255, 0.05)",
        "gradient-1": "#38bdf8",
        "gradient-2": "#818cf8",
        "gradient-3": "#fcd34d",
        "dial-bg": "rgba(0, 0, 0, 0.5)"
    }
} as WeatherContextModel;
//rgba(255, 255, 255, 0.8)
export const WeatherContext = React.createContext(weatherContext)