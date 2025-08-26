import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import App from './App.tsx'
import './index.css'
import { ThemeProvider, createTheme } from '@mui/material/styles';
import CssBaseline from '@mui/material/CssBaseline';
import Temperature from "./components/details/Temperature.tsx";
import Wind from "./components/details/Wind.tsx";
import Lightning from "./components/details/Lightning.tsx";
import Rain from "./components/details/Rain.tsx";
import Forecast from "./components/details/Forecast.tsx";

const darkTheme = createTheme({
    palette: {
        mode: 'dark',
    },
});
ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
      <ThemeProvider theme={darkTheme}>
          <CssBaseline />
              <BrowserRouter>
                  <Routes>
                      <Route path="/" element={<App />} />
                      <Route path="/details/temperature" element={<Temperature />} />
                      <Route path="/details/wind" element={<Wind />} />
                      <Route path="/details/lightning" element={<Lightning />} />
                      <Route path="/details/rain" element={<Rain />} />
                      <Route path="/details/forecast/:day" element={<Forecast />} />

                  </Routes>
              </BrowserRouter>
      </ThemeProvider>
  </React.StrictMode>,
)
