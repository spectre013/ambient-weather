<template>
  <div class="weather-item">
    <div class="chartforecast"></div>
    <span class="moduletitle">Daylight | Darkness</span><br />
    <div id="moonphase" v-if="astro">
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
          {{ now }}</span
        >
      </div>
      <div class="daylightmoduleposition">
        <div class="weather34sunlightday">
          {{ astro.day_length
          }}<weather34daylightdaycircle></weather34daylightdaycircle> hrs<br />Total Daylight
        </div>
        <div class="weather34sundarkday">
          {{ darkness }} <weather34darkdaycircle></weather34darkdaycircle> hrs <br />Total Darkness
        </div>

        <div class="weather34sunriseday">
          <svg
            id="weather34 sunrise"
            width="8pt"
            height="8pt"
            viewBox="0 0 200 120"
            version="1.1"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              fill="#ff8841"
              opacity="1.00"
              d=" M 73.12 50.01 C 90.23 43.20 110.01 43.20 127.10 50.08 C 134.54 53.53 141.90 57.50 148.03 63.03 C
              154.51 69.43 160.41 76.59 164.39 84.84 C 168.81 94.26 171.40 104.57 171.60 114.99 C 123.86 115.00 76.13
              115.01 28.40 114.99 C 28.79 97.79 35.36 80.71 47.17 68.13 C 53.98 59.83 63.50 54.42 73.12 50.01 Z"
            ></path></svg
          >Sun Rise<br />{{ todayTomorrow('sunrise') }} {{ astro.sunrise
          }}<!--<br>First Light (<blueu>04:55</blueu>)-->
        </div>
        <div class="weather34sunsetday">
          <svg
            id="weather34 sunset"
            width="8pt"
            height="8pt"
            viewBox="0 0 200 120"
            version="1.1"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              fill="#d65b4a"
              opacity="1.00"
              d=" M 28.39 48.01 C 76.13 47.99 123.87 48.00 171.60 48.01 C 171.21 64.95 164.90 81.83 153.35 94.34 C
              146.33 103.04 136.44 108.75 126.37 113.22 C 109.44 119.78 89.98 119.71 73.11 112.99 C 63.49 108.57 53.96
              103.15 47.15 94.85 C 35.35 82.27 28.80 65.20 28.39 48.01 Z"
            ></path></svg
          >Sun Set<br />{{ todayTomorrow('sunset') }} {{ astro.sunset
          }}<!--<br>Last Light (<blueu>21:24</blueu>)-->
        </div>
        <div class="daylightword">
          <value>Daylight</value>
        </div>

        <div class="elevationword">
          <value
            >Sun Elevation<span>
              <value>
                <maxred>
                  {{ astro.sun_altitude.toFixed(2) }}&nbsp;
                  <div :class="sunBelow">&nbsp;</div>
                </maxred>
              </value>
            </span></value
          >
        </div>
        <div class="circleborder"></div>
        <div class="sundialcontainerdiv2">
          <div id="sundialcontainer" class="sundialcontainer">
            <canvas id="sundial" class="suncanvasstyle" v-sundial="astro"></canvas>
            <div class="weather34sunclock" :style="{ transform: 'rotate(' + sunClock() + 'deg)' }">
              <div id="poscircircle" :style="sunColor"></div>
            </div>
          </div>
          <div class="daylightvalue1">
            <hrs>hrs</hrs>
            <hours>&nbsp;&nbsp;{{ hoursTilSunSet() }}</hours>
            <minutes>{{ minTilSunSet() }}</minutes> <br />&nbsp;<period>
              <value>&nbsp;{{ isSunSet }}</value>
            </period>
            <min>min</min>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import moment from 'moment';

export default {
  name: 'daylight',
  props: {
    astro: Object,
  },
  data() {
    return {
      hasSunset: false,
    };
  },
  watch: {
    astro: function () {
      this.sunHasSet();
    },
  },
  mounted: function () {},
  directives: {
    sundial: function (canvas, data) {
      let d_crcl = (24 * 60) / 2;
      let now = moment().unix();
      let sunRise = moment().startOf('day').add(timeToDecimal(data.value.sunrise), 'hours').unix();
      let sunSet = moment().startOf('day').add(timeToDecimal(data.value.sunset), 'hours').unix();

      function timeToDecimal(t) {
        let arr = t.split(':');
        let dec = parseInt((arr[1] / 6) * 10, 10);
        return parseFloat(parseInt(arr[0], 10) + '.' + (dec < 10 ? '0' : '') + dec);
      }
      function clc_crcl(ts) {
        let h = parseInt(moment.unix(ts).format('H'));
        let m = parseInt(moment.unix(ts).format('m'));
        let calc = m + h * 60;
        calc = 0.5 + calc / d_crcl;
        if (calc > 2.0) {
          calc = calc - 2;
        }
        return calc.toFixed(5);
      }
      let start = parseFloat(clc_crcl(sunRise));
      let end = parseFloat(clc_crcl(sunSet));
      let pos = parseFloat(clc_crcl(moment().unix()));
      let sn_clr = '';
      if (now > sunSet || now < sunRise) {
        sn_clr = 'rgba(86,95,103,0)';
      } else {
        sn_clr = 'rgba(255, 112,50,1)';
      }

      let ctx = canvas.getContext('2d');
      ctx.imageSmoothingEnabled = false;
      ctx.beginPath();
      ctx.arc(63, 65, 55, 0, 2 * Math.PI);
      ctx.lineWidth = 0;
      ctx.strokeStyle = '#565f67';
      ctx.stroke();
      ctx.beginPath();
      ctx.arc(63, 65, 55, start * Math.PI, end * Math.PI);
      ctx.lineWidth = 2;
      ctx.strokeStyle = '#3b9cac';
      ctx.stroke();
      ctx.beginPath();
      ctx.arc(63, 65, 55, pos * Math.PI, pos * Math.PI);
      ctx.lineWidth = 0;
      ctx.strokeStyle = `"${sn_clr}"`;
      ctx.stroke();
    },
  },
  methods: {
    sunClock: function () {
      return moment().hour() * 15 + moment().minute() / 4 - 86;
    },
    timeToDecimal: function (t) {
      let arr = t.split(':');
      let dec = parseInt((arr[1] / 6) * 10, 10);
      return parseFloat(parseInt(arr[0], 10) + '.' + (dec < 10 ? '0' : '') + dec);
    },
    sunRiseTime: function () {
      let sunRiseTime = moment()
        .startOf('day')
        .add(this.timeToDecimal(this.astro.sunrise), 'hours');
      return sunRiseTime;
    },
    sunSetTime: function () {
      let sunSetTime = moment().startOf('day').add(this.timeToDecimal(this.astro.sunset), 'hours');
      return sunSetTime;
    },
    hoursTilSunSet: function () {
      let sunsetTime = moment(this.astro.date + ' ' + this.astro.sunset + ':00');
      let now = moment();
      let h = moment.duration(sunsetTime.diff(now)).hours();
      if (h < 1 && this.hasSunset) {
        sunsetTime = moment(this.astro.tomorrow.date + ' ' + this.astro.tomorrow.sunrise + ':00');
      }
      let duration = sunsetTime.diff(now);
      return moment.duration(duration).hours();
    },
    minTilSunSet: function () {
      let sunsetTime = moment(this.astro.date + ' ' + this.astro.sunset + ':00');
      let now = moment();
      let t = moment.duration(sunsetTime.diff(now)).minutes();

      if (t <= 0 && this.hasSunset) {
        sunsetTime = moment(this.astro.tomorrow.date + ' ' + this.astro.tomorrow.sunrise + ':00');
      }
      let duration = sunsetTime.diff(now);

      return moment.duration(duration).minutes();
    },
    setDateTime(hours) {
      let h = hours.split(':');
      return moment().startOf('day').hour(h[0]).minute(h[1]);
    },
    todayTomorrow: function (type) {
      let event = this.setDateTime(this.astro[type]);
      if (moment() > event) {
        return 'Tomorrow';
      } else {
        return 'Today';
      }
    },
    sunHasSet: function () {
      let ss = this.astro.sunset.split(':');
      let sunset = moment().startOf('day').hour(ss[0]).minute(ss[1]);
      let s = moment.duration(sunset.diff(moment())).minutes();

      let sr = this.astro.sunrise.split(':');
      let sunrise = moment().startOf('day').hour(sr[0]).minute(sr[1]);
      let r = moment.duration(sunrise.diff(moment())).minutes();

      if (s <= 0 || r >= 0) {
        this.hasSunset = true;
      } else {
        this.hasSunset = false;
      }
    },
  },
  filters: {},
  computed: {
    now: function () {
      return moment().format('HH:mm:ss');
    },
    darkness: function () {
      let day = moment('1970-01-01 23:59:59');
      let hours = day.subtract(this.timeToDecimal(this.astro.day_length), 'hours');
      return hours.format('HH:mm');
    },
    isSunSet: function () {
      if (this.hasSunset) {
        return 'Til Sunrise';
      } else {
        return 'Til Sunset';
      }
    },
    sunBelow: function () {
      if (this.hasSunset) {
        return 'sunbelowweather34';
      } else {
        return 'sunaboveweather34';
      }
    },
    sunColor: function () {
      let value = '';
      if (this.astro.sun_altitude <= 0.5 && this.astro.sun_altitude > -4) {
        value = 'rgba(255, 112, 50, 0.5)';
      } else if (this.astro.sun_altitude <= 0) {
        value = 'rgba(86, 95, 103, 0.7)';
      } else {
        value = 'rgba(255,124,57,1)';
      }
      return { background: value };
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.weather34sunclock {
  border: 5px solid rgba(255, 255, 255, 0);
  width: 110px;
  height: 110px;
  top: -9px;
  margin-left: 104px;
}

.weather34sunclock #poscircircle {
  top: 50%;
  left: calc(48% - 52%);
  z-index: 1;
  height: 8px;
  width: 8px;
  border: 0;
  -webkit-border-radius: 50%;
  border-radius: 50%;
}
</style>
