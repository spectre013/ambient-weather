<template>
  <div class="modal-backdrop">
    <div class="modal">
      <section class="modal-body">
        <div id="body" v-if="alert">
          <div class="weather34darkbrowser" :url="url"></div>
          <main class="grid">
            <article>
              <h1 v-bind:class="alertColor">
                <svg
                  id="firealert"
                  viewBox="0 0 32 32"
                  width="22px"
                  height="22px"
                  fill="none"
                  stroke="currentcolor"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                >
                  <path d="M16 3 L30 29 2 29 Z M16 11 L16 19 M16 23 L16 25"></path>
                </svg>
                {{ alert.event }}
              </h1>
              <h3 v-bind:class="alertColor">{{ alert.headline }}</h3>
              <div class="alertinfo">
                <div v-bind:class="alertColor">Effective: {{ dates(alert.effective) }}</div>
                <div v-bind:class="alertColor">Ends: {{ dates(alert.end) }}</div>
                <div v-bind:class="alertColor">Certainty: {{ alert.certainty }}</div>
                <div v-bind:class="alertColor">Response: {{ alert.response }}</div>
              </div>
              <div class="alertText">
                <span v-bind:class="alertColor">Alert Description: </span>
                <span v-html="alert.description"></span>
              </div>
              <div class="targetarea">
                <span v-bind:class="alertColor">Regions Affected: </span>
                <span>{{ alert.areadesc }}</span>
              </div>
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

<script>
// import moment from 'moment';
//
// export default {
//   name: 'alertmodal',
//   props: {
//     alerts: Array,
//     options: Object,
//   },
//   data() {
//     return {
//       alert: null,
//     };
//   },
//   methods: {
//     close: function () {
//       this.$parent.closeModal('alert');
//     },
//   },
//   filters: {
//     renderDesc: function (desc) {
//       console.log(desc);
//       desc.replace('\n', '<br />');
//       return desc;
//     },
//     dates: function (date) {
//       console.log(date);
//       date = moment(date).utcOffset(7);
//       return date.format('HH:mm DD MMM');
//     },
//   },
//   computed: {
//     url: function () {
//       return `Weather Alerts -  ${this.alert.event}`;
//     },
//     alertColor: function () {
//       if (this.alert.event.startsWith('911')) {
//         return 'Telephone Outage 911'.toLowerCase().replace(/\s+/g, '-');
//       }
//       return this.alert.event.toLowerCase().replace(/\s+/g, '-');
//     },
//   },
//   mounted: function () {
//     this.alert = this.alerts[this.options.alert];
//     this.alert.description = this.alert.description.replaceAll('\n', '<br />');
//     console.log(this.alert);
//   },
// };
</script>

<style scoped>
h1 {
  font-size: 22px;
  text-align: center;
}
h3 {
  font-size: 16px;
  text-align: center;
}
.alertText {
  color: silver;
}
.targetarea {
  margin-top: 10px;
  color: silver;
}
table {
  border: none;
  padding: 0;
  margin: 0;
}
td {
  margin: 2px;
  padding: 5px;
}
.alertinfo div {
  display: table;
  width: 23%;
}
</style>
