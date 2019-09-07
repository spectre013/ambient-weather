<template>
<div class="weather34box indoor">
    <div class="title">
        <svg id="i-info" viewBox="0 0 32 32" width="9" height="9" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M16 14 L16 23 M16 8 L16 10"></path><circle cx="16" cy="16" r="14"></circle></svg>
        &nbsp;<oblue>Min</oblue> | <ored>Max</ored>  Temperatures</div>
    <div class="value" v-if="temp">
        <div id="position2">
            <div class="topmin">
                <topblue1>{{ temp.min.year.value | tempDisplay($store.getters.units) }}<smallwindunit>&deg;{{ tempLabel }}</smallwindunit></topblue1>
            </div>
            <div class="minword">{{temp.min.year.date | mFormat }}</div>
            <div class="mintimedate">{{temp.min.year.date | dayFormat }}</div>
            <div class="yearwordbig">{{temp.min.year.date | year }}</div>
            <div class="topmax">
                <toporange1>{{ temp.max.year.value | tempDisplay($store.getters.units) }}<smallwindunit>&deg;{{ tempLabel }}</smallwindunit></toporange1>
            </div>
            <div class="maxword">{{temp.max.year.date | year }}</div>
            <div class="maxtimedate">{{temp.max.year.date | dayFormat }}</div>
        </div>
    </div>
</div>
</template>

<script>

import moment from 'moment';

export default {
  name: 'minmax',
    props: {
        temp: Object,
    },
    methods: {

    },
    mounted: function() {

    },
    filters: {
        at: function (date) {
            return moment(date).format('HH:mm');
        },
        mFormat: function (date) {
            return moment(date).format('MMM');
        },
        dayFormat: function (date) {
            return moment(date).format('MMM Do');
        },
        Year: function (date) {
            return moment(date).format('YYYY');
        },
    },
    computed: {
    },
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
