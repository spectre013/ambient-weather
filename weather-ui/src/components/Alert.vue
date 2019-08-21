<template>
    <div class="weather34box alert">
        <div class="title"><svg id="i-info" viewBox="0 0 32 32" width="9" height="9" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M16 14 L16 23 M16 8 L16 10"></path><circle cx="16" cy="16" r="14"></circle></svg>
            Weather <ored>Advisory </ored>
        </div>
        <div class="value" v-if="forecast">
            <div id="position4" v-if="hasName">
                <div class="eqcirclehomeregional">
                    <div class="eqtexthomeregional">
                        <spanelightning>
                            <alertvalue><a href="#" v-on:click="openModal('alert')">
                                <redu> {{ forecast.alerts[0].title }} <svg id="firealert" viewBox="0 0 32 32" width="11px" height="11px" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path d="M16 3 L30 29 2 29 Z M16 11 L16 19 M16 23 L16 25"></path></svg>
                                    {{ forecast.alerts[0].severity }} </redu>
                            </a><br>Expires {{ forecast.alerts[0].expires | expire }}</alertvalue>
                        </spanelightning>
                    </div>
                </div>
            </div>
        </div>
        <div id="position4" v-else>
            <div class="eqcirclehomeregional">
                <div class="eqtexthomeregional">
                    <spanelightning>
                        <alertvalue> No Alerts</alertvalue>
                    </spanelightning>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import moment from 'moment';

export default {
  name: 'alerts',
  props: {
    forecast: Object
  },
  mounted() {
   
  },
  methods: {
    containsKey(obj, key ) {
        return Object.keys(obj).includes(key);
    },
    openModal: function(type) {
        this.$parent.showModal(type);
    },
  },
  filters: {
   expire: function(date) {
       return moment.unix(date).format("HH:mm DD MMM")
   }
  },
  computed: {
    hasName() {
        if (!this.forecast) return false;
        return this.containsKey(this.forecast, 'alerts');
    }
  },
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
