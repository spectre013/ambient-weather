<template>
<div class="weather34box indoor">
    <div class="title"><svg id="i-info" viewBox="0 0 32 32" width="9" height="9" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M16 14 L16 23 M16 8 L16 10"></path><circle cx="16" cy="16" r="14"></circle></svg> Min | Max  <orange> Temps</orange></div>
    <div class="value" v-if="temp">
        <div id="position2">
            <div class="topmin">
                <topgreen1>{{ temp.min.month.value }}<smallwindunit>&deg;</smallwindunit></topgreen1>
            </div>
            <div class="minword">{{temp.min.month.date | mFormat }}</div>
            <div class="mintimedate">{{temp.min.month.date | dayFormat }}</div>
            <div class="yearwordbig">{{temp.min.year.date | year }}</div>
            <div class="topmax">
                <toporange1>{{ temp.max.year.value }}<smallwindunit>&deg;</smallwindunit></toporange1>
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
