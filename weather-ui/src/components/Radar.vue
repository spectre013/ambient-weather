<template>
    <div class="modal-backdrop">
        <div class="modal">
            <section class="modal-body">
                <div id="body">
                    <div class="weather34darkbrowser" url="Radar"></div>
                    <main class="grid">
                        <article>
                            <div id="radarmap" style="width: 880px; height: 480px;"></div>
                        </article>
                    </main>
                </div>
                <button class="lity-close" type="button" aria-label="Close (Press escape to close)" @click="close">×</button>
            </section>
        </div>
    </div>
</template>


<script>
    import L from 'leaflet'
    export default {
        name: 'radarmodel',
        props: {
        },
        data() {
            return {
                map: null,
            }
        },
        methods: {
            close: function() {
                this.$parent.closeModal('radar');
            },
            initMap: function(mapEl) {
                let map = L.map(mapEl).setView([38.7360,-104.6284], 7);
                L.tileLayer('https://api.tiles.mapbox.com/v4/{id}/{z}/{x}/{y}.png?access_token=pk.eyJ1IjoibWFwYm94IiwiYSI6ImNpejY4NXVycTA2emYycXBndHRqcmZ3N3gifQ.rJcFIG214AriISLbB6B5aw', {
                    maxZoom: 18,
                    attribution: 'Map data &copy; <a href="https://www.openstreetmap.org/">OpenStreetMap</a> contributors, ' +
                        '<a href="https://creativecommons.org/licenses/by-sa/2.0/">CC-BY-SA</a>, ' +
                        'Imagery © <a href="https://www.mapbox.com/">Mapbox</a>',
                    id: 'mapbox.streets'
                }).addTo(map);
                let greenIcon = L.icon({
                    iconUrl: 'css/circle.svg',
                    iconSize:     [7, 7], // size of the icon
                    iconAnchor:   [3.5, 3.5], // point of the icon which will correspond to marker's location

                });

                let imageUrl = 'https://radar.weather.gov/RadarImg/N0R/PUX_N0R_0.gif',
                    imageBounds = [[35.928214846524, -101.42056117665],[40.9805982693759, -106.932252183394]];
                L.marker([38.7360,-104.6285],{icon:greenIcon}).addTo(map);
                L.imageOverlay(imageUrl, imageBounds,{opacity:.6}).addTo(map);
            },
        },
        mounted() {
            this.$nextTick(() => {
               this.initMap(document.querySelector('#radarmap'));
            })
        },
        computed: {
        },
    };
</script>

<style scoped>

</style>