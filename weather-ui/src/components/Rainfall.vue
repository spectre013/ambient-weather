<template>
    <div class="weather34box earthquake"><div class="title"><svg id="i-info" viewBox="0 0 32 32" width="9" height="9" fill="none" stroke="currentcolor" stroke-linecap="round" stroke-linejoin="round" stroke-width="6.25%"><path d="M16 14 L16 23 M16 8 L16 10"></path><circle cx="16" cy="16" r="14"></circle></svg> Rainfall<oblue> Data </oblue>
        </div>
        <div class="value" v-if="current">
            <div id="position3">
                <div class="topmin">
                    <topblue1>{{ current.monthlyrainin }}<smallwindunit>in</smallwindunit></topblue1>
                </div>
                <div class="minword">{{ current.date | month }}</div>
                <div class="mintimedate">Total
                </div>
                <div class="yearwordbig">{{ current.date | year}}</div>
                <div class="topmax">
                    <topblue1>{{ current.yearlyrainin }}<smallwindunit>in</smallwindunit></topblue1>
                </div>
                <div class="maxword">{{ current.date | year }}</div>
                <div class="maxtimedate">Total</div>
            </div>
        </div>
    </div>
</template>

<script>
import moment from 'moment';

export default {
  name: 'rainfall',
  props: {
    current: Object
  },
  mounted() {
   
  },
  methods: {
   
  },
  filters: {
    year: function (date) {
      return moment(date).format('YYYY');
    },
    month: function (date) {
      return moment(date).format('MMM');
    },
  },
  computed: {
    
  },
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
