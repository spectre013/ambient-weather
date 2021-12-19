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
                {{ alert.title }}
              </h1>
              <div class="alertText">
                <span v-bind:class="alertColor">Alert Description: </span> {{ alert.description }}
              </div>
              <div class="targetarea">
                <span v-bind:class="alertColor">Regions Affected: </span>
                <span>{{ joinRegions }}</span>
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
export default {
  name: 'alertmodal',
  props: {
    forecast: Object,
    options: Object,
  },
  data() {
    return {
      alert: null,
    };
  },
  methods: {
    close: function () {
      this.$parent.closeModal('alert');
    },
  },
  computed: {
    url: function () {
      return `Weather Alerts -  ${this.alert.title}`;
    },
    alertColor: function () {
      if (this.alert.title.startsWith('911')) {
        return 'Telephone Outage 911'.toLowerCase().replace(/\s+/g, '-');
      }
      return this.alert.title.toLowerCase().replace(/\s+/g, '-');
    },
    joinRegions: function () {
      return this.alert.regions.join(', ');
    },
  },
  mounted: function () {
    this.alert = this.forecast.alerts[this.options.alert];
  },
};
</script>

<style scoped>
h1 {
  font-size: 22px;
  text-align: center;
}
.alertText {
  color: silver;
}
.targetarea {
  margin-top: 10px;
  color: silver;
}
</style>
