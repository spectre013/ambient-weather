import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../components/HomeView.vue'
//import ForecastVue from "../components/ForecastVue.vue";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    // {
    //   path: '/forecast',
    //   name: 'forecast',
    //   component: ForecastVue
    // }
  ]
})

export default router
