<template>
    <div class="weather34box alert">
        <div class="title"><svg id="i-info" viewBox="0 0 32 32" width="9" height="9" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M16 14 L16 23 M16 8 L16 10"></path><circle cx="16" cy="16" r="14"></circle></svg>
            Weather <ored>Advisory </ored>
        </div>
        <div class="value" v-if="hasName">
            <div id="position4" v-if="forecast && alert">
                <div class="eqcirclehomeregional">
                    <div class="eqtexthomeregional">
                        <div class="uparrow" v-if="multipleAlerts && minAlerts" v-on:click="switchAlert('up')">
                            <svg xmlns="http://www.w3.org/2000/svg" version="1.1" width="16" height="16"><path d="m 13,6 -5,5 -5,-5 z" fill="#797979"/></svg>
                        </div>
                        <spanelightning>
                            <alertvalue><a href="#" v-on:click="openModal('alert',{alert:currentAlert})">
                                <span v-bind:class="alertColor"><svg id="firealert" viewBox="0 0 32 32" width="11px" height="11px" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"><path d="M16 3 L30 29 2 29 Z M16 11 L16 19 M16 23 L16 25"></path></svg>
                                    {{ alert.title }} <!-- : {{ alert.severity }} --> </span>
                            </a><br><span v-bind:class="alertColor">Expires {{ alert.expires | expire }}</span></alertvalue>
                        </spanelightning>
                        <div class="downarrow" v-if="multipleAlerts &&  maxAlerts" v-on:click="switchAlert('down')">
                            <svg xmlns="http://www.w3.org/2000/svg" version="1.1" width="16" height="16"><path d="m 13,6 -5,5 -5,-5 z" fill="#797979"/></svg>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div id="position4" v-else>
            <div class="eqcirclehomeregional">
                <div class="eqtexthomeregional">
                    <spanelightning>
                        <alertvalue> <span>No Alerts</span></alertvalue>
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
    data() {
        return {
            alert: null,
            currentAlert: 0,
        }
    },
    watch: {
        forecast: function() {
            if(Object.prototype.hasOwnProperty.call(this.forecast,'alerts')) {
                this.alert = this.forecast.alerts[this.currentAlert];
            }
        },
        currentAlert: function() {
            this.alert = this.forecast.alerts[this.currentAlert];
        }

    },
    mounted() {
   
    },
    methods: {
        containsKey(obj, key ) {
            return Object.keys(obj).includes(key);
        },
        openModal: function(type,options) {
            this.$parent.showModal(type,options);
        },
        switchAlert(dir) {
            if(dir === 'up' && this.currentAlert > 0) {
                this.currentAlert--;
            } else if(this.currentAlert <= this.forecast.alerts.length) {
                this.currentAlert++;
            }
        }
    },
    filters: {
        expire: function(date) {
           return moment.unix(date).format("HH:mm DD MMM")
        }
    },
  computed: {
    multipleAlerts: function() {
     return this.forecast.alerts.length > 1;
    },
    maxAlerts: function() {
        let ret = true;
        if(this.currentAlert + 1 >= this.forecast.alerts.length) {
            ret = false;
        }
      return ret;
    },
    minAlerts: function() {
      let ret = true;
      if(this.currentAlert - 1 < 0) {
          ret = false;
      }
      return ret;
    },
    hasName() {
        if (!this.forecast) return false;
        return this.containsKey(this.forecast, 'alerts');
    },
    alertColor: function() {
        if(this.alert.title.startsWith("911")) {
            return "Telephone Outage 911".replace(/\s+/g, "-");
        }
        return this.alert.title.toLowerCase().replace(/\s+/g, "-")
    }
  },
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
    .uparrow {
        position: absolute;
        top: 10px;
        left: 210px;
        transform: rotate(180deg);
        cursor: pointer;
    }
    .downarrow {
        position: absolute;
        top: 50px;
        left: 210px;
        cursor: pointer;
    }
</style>
