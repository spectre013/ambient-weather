import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import App from './App.tsx'
import './index.css'
import TempAlmanac from "./components/almanac/TempAlmanac.tsx";

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
      <BrowserRouter>
          <Routes>
              <Route path="/" element={<App />} />
              <Route path="/almanac/temp" element={<TempAlmanac />} />
          </Routes>
      </BrowserRouter>
  </React.StrictMode>,
)
