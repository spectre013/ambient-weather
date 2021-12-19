<template>
  <div class="weather-item">
    <div class="chartforecast">
      <span class="yearpopup">
        <a
          alt="nearby metar station"
          title="nearby metar station"
          href="#"
          v-on:click="openModal('metar')"
        >
          <svg
            viewBox="0 0 32 32"
            width="8"
            height="8"
            fill="none"
            stroke="#777"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="6.25%"
          >
            <path d="M14 9 L3 9 3 29 23 29 23 18 M18 4 L28 4 28 14 M28 4 L14 18"></path>
          </svg>
          Nearby Metar
        </a>
      </span>
      <span class="monthpopup">
        <a
          href="#"
          v-on:click="openModal('radar')"
          title="Windy.com Radar"
          alt="Windy.com Radar"
          data-lity=""
        >
          <svg
            viewBox="0 0 32 32"
            width="8"
            height="8"
            fill="none"
            stroke="#777"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="6.25%"
          >
            <path d="M14 9 L3 9 3 29 23 29 23 18 M18 4 L28 4 28 14 M28 4 L14 18"></path>
          </svg>
          Radar</a
        >
      </span>
    </div>
    <span class="moduletitle">Lightining</span>
    <div id="lightning-container" v-if="current && lightning">
      <div class="updatedtime">
        <span
          ><svg
            id="i-info"
            viewBox="0 0 32 32"
            width="6"
            height="6"
            fill="#9aba2f"
            stroke="#9aba2f"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="6.25%"
          >
            <path d="M16 14 L16 23 M16 8 L16 10"></path>
            <circle cx="16" cy="16" r="14"></circle>
          </svg>
          {{ current.date | now }}
        </span>
      </div>

      <!--            <div class="lightningsvg">-->
      <!--                <span><svg version="1.1" id="weather34_lightning_alert_notification"
      xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" x="0px" y="0px" viewBox="0 0 30 30"
       fill="current" xml:space="preserve">-->
      <!--                <path class="st0" d="M8.6,16.4c3.9-3.8,7.7-7.7,11.6-11.5c-1,2.9-2,5.8-3,8.7c1.4,0-17.1,0,4.3,
      0c-3.9,3.8-7.9,7.7-11.8,11.5-->
      <!--                    c1.1-2.9,2.2-5.8,3.3-8.7C11.4,16.4,10,16.4,8.6,16.4z"/>-->
      <!--                </svg>-->

      <!--                </span>-->
      <!--            </div>-->
      <!--            -->
      <div id="lightning">
        <div class="sunblock">
          <div class="button button-dial-small">
            <strikeicon>Hour</strikeicon>
            <div class="button-dial-top-small"></div>
            <div class="button-dial-label-small">{{ current.lightninghour }}</div>
          </div>
          <laststrike>
            Last Strike
            <br />{{ current.lightningtime | lightningFormat }} <br />
            {{ current.lightningdistance }} Miles away
          </laststrike>
        </div>
        <div class="sunblock">
          <div class="button button-dial-small">
            <strikeicon>Today</strikeicon>
            <div class="button-dial-top-small"></div>
            <div class="button-dial-label-small">{{ current.lightningday }}</div>
          </div>
        </div>
        <div class="sunblock">
          <div class="button button-dial-small">
            <strikeicon class="two">Yesterday</strikeicon>
            <div class="button-dial-top-small"></div>
            <div class="button-dial-label-small">{{ lightning.max.yesterday.value }}</div>
          </div>
        </div>
      </div>
      <div class="lightningyear-mod2">
        Strikes {{ current.date | month }} <orange> {{ lightning.max.month.value }}</orange>
      </div>
      <div class="lightningmonth-mod2">
        Strikes {{ current.date | year }} <orange> {{ lightning.max.year.value }}</orange>
      </div>
    </div>
  </div>
</template>

<script>
import moment from 'moment';
import axios from 'axios';
//import axios from 'axios';
export default {
  name: 'Lightining',
  props: {
    current: Object,
  },
  data() {
    return {
      lightning: null,
    };
  },
  mounted: function () {
    function updateData(self) {
      axios.get('/api/minmax/lightning').then((response) => (self.lightning = response.data));
      setTimeout(function () {
        updateData(self);
      }, 60000);
    }
    updateData(this);
  },
  watch: {},
  methods: {
    openModal: function (type, options) {
      this.$parent.showModal(type, options);
    },
  },
  filters: {
    dayOrHour: function (cnt) {
      return cnt > 0 ? 'Hour' : 'Today';
    },
    year: function (date) {
      return moment(date).format('Y');
    },
    month: function (date) {
      return moment(date).format('MMMM');
    },
    now: function (date) {
      return moment(date).format('HH:mm:ss');
    },
    lightningFormat: function (date) {
      return moment(date).fromNow(); //format('HH:mm');
    },
    timeFormat: function (date) {
      return moment(date).format('HH:mm');
    },
    getCount: function (current) {
      return current.lightninghour > 0 ? current.lightninghour : current.lightningday;
    },
  },
  computed: {},
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
:root {
  --white: #ffffff;
  --light: #f5f5f5;
  --dark: #07090a;
  --dark-popup: #07090a;
  --dark-light: hsla(0, 0%, 0%, 0.251);
  --dark-toggle: #35393b;
  --dark-caption: rgba(66, 70, 72, 0.429);
  --black: #000000;
  --deepblue: #00adbd;
  --blue: #00adbd;
  --rainblue: #00adbd;
  --darkred: #703232;
  --deepred: #703232;
  --red: #d35f50;
  --yellow: #e6a241;
  --green: #90b22a;
  --orange: #d87040;
  --purple: #8680bc;
  --silver: #ecf0f3;
  --border-dark: #3d464d;
  --body-text-dark: #afb7c0;
  --body-text-light: #545454;
  --blocks: #e6e8ef;
  --modules: #1e1f26;
  --blocks-background: rgba(82, 92, 97, 0.19);
  --temp-5-10: #7face6;
  --temp-0-5: #00adbd;
  --temp0-5: #00adbd;
  --temp5-10: #9bbc2f;
  --temp10-15: #e6a241;
  --temp15-20: #f78d03;
  --temp20-25: #d87040;
  --temp25-30: #e64b24;
  --temp30-35: #cc504c;
  --temp35-40: hsl(4, 40%, 48%);
  --temp40-45: #be5285;
  --temp45-50: #b95c95;
  --font-color: grey;
  --bg-color: hsla(198, 14%, 14%, 0.19);
  --button-bg-color: hsla(198, 14%, 14%, 0.19);
  --button-shadow: inset 5px 5px 20px #0c0b0b, inset -5px -5px 20px hsla(198, 14%, 14%, 0.19);
}
.lightningsvg {
  fill: #af1909;
  height: 60px;
  width: 60px;
  float: left;
}
.two {
  margin-left: -10px;
}
#lightning-container {
  height: 140px;
}
.sunblock {
  padding: 15px 0 0 0;
  display: flex;
  margin-left: 0;
  text-align: left;
  width: 310px;
  position: relative;
  color: var(--blue);
  font-size: 2rem;
  text-transform: uppercase;
  height: 130px;
  line-height: 2;
  margin-top: 5px;
  font-family: weathertext2, sans-serif;
}
.lightningyear-mod2 {
  display: -webkit-box;
  display: -ms-flexbox;
  display: flex;
  font-size: 12px;
  color: var(--body-text-dark);
  position: relative;
  top: -30px;
  left: 30px;
  width: 8.25rem;
  line-height: 1;
  height: 25px;
  font-family: weathertext2, sans-serif;
  -webkit-box-pack: center;
  -ms-flex-pack: center;
  justify-content: center;
  -webkit-box-align: center;
  -ms-flex-align: center;
  align-items: center;
  background-color: hsla(204, 31%, 21%, 0.229);
  text-transform: none;
}

.lightningmonth-mod2 {
  display: -webkit-box;
  display: -ms-flexbox;
  display: flex;
  font-size: 12px;
  color: var(--body-text-dark);
  position: relative;
  top: -55px;
  left: 171px;
  width: 8.25rem;
  line-height: 1;
  height: 25px;
  font-family: weathertext2;
  -webkit-box-pack: center;
  -ms-flex-pack: center;
  justify-content: center;
  -webkit-box-align: center;
  -ms-flex-align: center;
  align-items: center;
  background-color: hsla(204, 31%, 21%, 0.229);
  text-transform: none;
}
.lightningdesc,
.moondesc {
  display: flex;
  background-color: hsla(204, 31%, 21%, 0.229);
  color: var(--body-text-dark);
  top: 0px;
  position: absolute;
  margin-left: 195px;
  font-family: weathertext2;
  text-align: left;
  width: 100px;
  font-size: 10px;
  padding: 2px 3px 1px 3px;
  line-height: 1.2;
  align-items: center;
  justify-content: center;
  border-radius: 3px;
  -webkit-border-radius: 3px;
  -moz-border-radius: 3px;
  -ms-border-radius: 3px;
  -o-border-radius: 3px;
  text-transform: none;
}
.button-dial-small {
  border-radius: 50%;
  display: flex;
  height: 85px;
  left: 0px;
  align-items: center;
  justify-content: center;
  width: 85px;
  font-family: weathertext2;
  margin: -20px 0 0;
  color: silver;
}
.button {
  color: grey;
  position: relative;
}
.button-dial-top-small {
  background: var(--button-bg-color);
  box-shadow: var(--button-shadow);
  border-radius: 50%;
  width: 73%;
  height: 73%;
  margin: 0 auto;
  position: absolute;
  top: 15%;
  left: 15%;
  text-align: center;
  z-index: 5;
  background-image: linear-gradient(hsla(0, 0%, 33%, 0.1) 1px, transparent 1px),
    linear-gradient(to right, hsla(0, 0%, 33%, 0.1) 1px, transparent 1px);
  background-size: 2px 2px;
}
.button-dial-label-small {
  font-size: 18px;
  position: relative;
  z-index: 10;
  color: silver;
}
laststrike {
  display: flex;
  background: 0;
  color: var(--silver);
  font-size: 9px;
  line-height: 1;
  top: 70px;
  position: absolute;
  left: 20px;
  font-family: weathertext2;
  text-align: left;
  width: 150px;
  word-wrap: break-word;
  text-transform: uppercase;
}
#lightning {
  align-items: center;
  display: flex;
  justify-content: center;
  min-width: 315px;
  max-width: 315px;
  padding: 0;
  height: 145px;
  -webkit-box-sizing: border-box;
  box-sizing: border-box;
  margin: 3px;
  color: var(--body-text-dark);
  background-color: var(--modules);
  border-radius: 5px;
  -webkit-border-radius: 5px;
  -moz-border-radius: 5px;
  -ms-border-radius: 5px;
  -o-border-radius: 5px;
  /*border: 1px thin hsla(217, 15%, 17%, 0.5);*/
  /*border-bottom: 5px solid hsla(217, 15%, 17%, 0.5);*/
}
strikeicon,
strikeicon2 {
  position: absolute;
  left: 31px;
  top: 49px;
  z-index: 100;
  font-size: 8px;
  color: silver;
}
orange {
  padding-left: 5px;
}
</style>
