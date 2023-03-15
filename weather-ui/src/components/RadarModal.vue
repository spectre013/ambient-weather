<template>
  <div class="modal-backdrop">
    <div class="modal">
      <section class="modal-body">
        <div id="body">
          <div class="weather34darkbrowser" url="Radar"></div>
          <main class="grid">
            <article>
              <div id="radarmap" style="width: 880px; height: 480px"></div>
            </article>
          </main>
        </div>
        <button
          class="lity-close"
          type="button"
          aria-label="Close (Press escape to close)"
          @click="close"
        >
          ×
        </button>
      </section>
    </div>
  </div>
</template>

<script setup>
import {onMounted } from 'vue';
import L from '../../public/leaflet.js';

let map = null
//const emit = defineEmits(['openModal'])

onMounted(() => {
  this.$nextTick(() => {
    initMap(document.querySelector('#radarmap'));
  });
});


// function openModal(type, options) {
//   emit('openModal',{type:type,options:options})
// }

function initMap(mapEl) {
  map = L.map(mapEl).setView([38.736, -104.6284], 10);
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    maxZoom: 18,
    attribution:
      'Map data &copy; <a href="https://www.openstreetmap.org/">OpenStreetMap</a> contributors, ' +
      '<a href="https://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, ' +
      'Imagery © <a href="https://www.mapbox.com/">Mapbox</a>',
    id: 'mapbox.streets',
  }).addTo(map);
  let greenIcon = L.icon({
    iconUrl: 'css/circle.svg',
    iconSize: [7, 7], // size of the icon
    iconAnchor: [3.5, 3.5], // point of the icon which will correspond to marker's location
  });

  let imageUrl = 'https://radar.weather.gov/RadarImg/N0R/PUX_N0R_0.gif',
    imageBounds = [
      [35.928214846524, -101.42056117665],
      [40.9805982693759, -106.932252183394],
    ];
  L.marker([38.736, -104.6285], { icon: greenIcon }).addTo(map);
  L.imageOverlay(imageUrl, imageBounds, { opacity: 0.6 }).addTo(map);
}
</script>

<style scoped></style>
