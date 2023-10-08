 import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue({
    template: {
      compilerOptions: {
        isCustomElement: (tag) => [
            /^darksky/,
            /^temp/,
            /^weather34/,
            /^span/,
            'oorange',
            'orange',
            'ogreen',
            'spanelightning',
            'ored',
            'heatindex',
            'spanmaxwind',
            'rainf',
            'spanwindtitle',
            'spanmaxwind',
            'oblue',
            'yesterdaytimemax',
            'spaneindoortemp',
            'suptemp',
            'trendmovementrising',
            'trendmovementfalling',
            'rainblue',
            'grey',
            'spancalm',
            'supmb',
            'max',
            'supmb',
            'windrun',
            'darkgrey',
            'convtext',
            'steady',
            'strongnumbers',
            'period',
            'min',
            'minutes',
            'hrs',
            'spanindoortempfalling',
            'indoororange1',
            'spanindoortempfalling',
            'indoororange1',
            'spanindoortempfalling',
            'luminance1',
            'uppercase',
            'span1',
            'uviforecasthourgreen',
            'value',
            'valuetext',
            'smalluvunit',
            'alertvalue',
            'noalert',
            'spanyellow',
            'valuetitleunit',
            'topblue1',
            'smallwindunit',
            'smallrainunita',
            'valuetextheading1',
            'raiblue',
            'smallrainunit2',
            'rainratetextheading',
            'maxred',
            'hours',
            'blueu',
            'moonrisecolor',
            'moonm',
            'moonsetcolor',
            'tgreen',
            'topgreen1',
            'smallwindunit',
            'toporange1',
            'minblue',
            'smalltempunit',
            'smalltempunit2',
            'smallrainunit',
            'valuetext1',
            'uviforecasthouryellow',
            'trendmovementfallingx',
            'redu',
            'blue',
            'articlegraph',
            'uviforecasthourred',
            'uviforecasthourorange',
            'stationid',
            'yellow',
            'green',
            'darkskytempwindhome',
            'darkskyiconcurrent',
            'darkskyrainhome1',
            'unit',
            'blue1',
            'windunit',
            'gust',
            'purpleu',
            'orange1',
            'uv',
            'topbarmetricc',
            'topbarimperialf',
            'topbarimperial',
            'topbarmetric',
            'strikeicon',
            'laststrike',
            'weather34windrunspan',
            'weather34bftspan',
            'weather34-barometerlimitmax',
            'weather34-barometerlimitmin',
            'weather34darkdaycircle',
            'weather34daylightdaycircle',
            'weather34',
        ].includes(tag)
      }
    }
  })],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:3000',
        changeOrigin: true,
      },
      '/api/ws': {
          target: 'ws://localhost:3000',
          ws: true
      },
    },
  },
})