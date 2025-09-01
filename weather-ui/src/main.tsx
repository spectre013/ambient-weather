import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import App from './App.tsx'
import './index.css'
import Temperature from "./components/details/Temperature.tsx";
import Wind from "./components/details/Wind.tsx";
import Lightning from "./components/details/Lightning.tsx";
import Rain from "./components/details/Rain.tsx";
import Forecast from "./components/details/Forecast.tsx";
import About from "./components/About.tsx";
import {WeatherContext, weatherContext} from "./Context.ts";

const theme = localStorage.getItem('theme') || 'dark';
const body: HTMLElement = document.body;
body.dataset['theme'] = theme

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
      <WeatherContext.Provider value={weatherContext}>
              <BrowserRouter>
                  <Routes>
                      <Route path="/" element={<App />} />
                      <Route path="/details/temperature" element={<Temperature />} />
                      <Route path="/details/wind" element={<Wind />} />
                      <Route path="/details/lightning" element={<Lightning />} />
                      <Route path="/details/rain" element={<Rain />} />
                      <Route path="/details/forecast/:day" element={<Forecast />} />
                      <Route path="/about" element={<About />} />
                  </Routes>
              </BrowserRouter>
      </WeatherContext.Provider>
  </React.StrictMode>,
)
