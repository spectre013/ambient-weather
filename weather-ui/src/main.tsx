import React from 'react'
import ReactDOM from 'react-dom/client'
import {BrowserRouter, Route, Routes} from "react-router-dom";
import './index.css'
import Home from "./components/Home.tsx";
import NoPage from "./404.tsx";
import TempAlmanac from "./components/almanac/TempAlmanac.tsx";

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
      <BrowserRouter>
          <Routes>
              <Route>
                  <Route index element={<Home />} />
                  <Route path="/almanac/temperature" element={<TempAlmanac />} />
                  <Route path="*" element={<NoPage />} />
              </Route>
          </Routes>
      </BrowserRouter>
  </React.StrictMode>,
)
