<template>
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
</template>

<script>
import moment from 'moment';

export default {
  name: 'boxheader',
  props: {
    current: {},
  },
  mounted() {},
  methods: {},
  filters: {
    now: function (date) {
      return moment(date).format('HH:mm:ss');
    },
  },
  computed: {},
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped></style>
