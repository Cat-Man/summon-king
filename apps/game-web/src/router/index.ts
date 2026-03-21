import { createRouter, createWebHistory } from 'vue-router'

import GameLayout from '@/layouts/GameLayout.vue'
import HomePage from '@/pages/home/HomePage.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: GameLayout,
      children: [
        {
          path: '',
          name: 'home',
          component: HomePage
        }
      ]
    }
  ]
})

export default router
