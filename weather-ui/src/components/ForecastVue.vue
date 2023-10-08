<!--<script setup>-->
<!--import {ref, onMounted, onBeforeMount} from 'vue';-->
<!--import axios from 'axios';-->
<!--import {useStore} from "vuex";-->


<!--let loaded = ref(false);-->
<!--let live = ref(null);-->
<!--//let theme = ref('dark');-->
<!--const store = useStore()-->
<!--let modalOptions = ref(null);-->
<!--let connection = null;-->

<!--function show(values) {-->
<!--  modalOptions.value = values.options;-->
<!--  models.value[values.type] = true;-->
<!--}-->
<!--function close(event) {-->
<!--  models.value[event] = false;-->
<!--}-->

<!--onBeforeMount(() => {-->
<!--  //setStyle();-->
<!--  store.dispatch('getSettings');-->
<!--  loaded.value = true;-->
<!--});-->

<!--onMounted(() => {-->
<!--  axios.get('/api/luna').then((response) => (astro.value = response.data));-->
<!--});-->
<!--// created() {-->
<!--//   this.$on('unitsChange', (units) => {-->
<!--//     this.units = units;-->
<!--//   });-->
<!--// },-->

<!--</script>-->
<!--<template>-->
<!--    <div class="weather2-container">-->
<!--      <h1>Forecast</h1>-->
<!--    </div>-->
<!--</template>-->

<!--<style></style>-->
