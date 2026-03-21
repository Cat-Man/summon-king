import { createRouter, createWebHistory } from 'vue-router'

import GameLayout from '@/layouts/GameLayout.vue'
import InventoryPage from '@/pages/assets/InventoryPage.vue'
import HomePage from '@/pages/home/HomePage.vue'
import PetDetailPage from '@/pages/pets/PetDetailPage.vue'
import PetListPage from '@/pages/pets/PetListPage.vue'
import PetTeamPage from '@/pages/pets/PetTeamPage.vue'

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
        },
        {
          path: 'assets',
          name: 'assets',
          component: InventoryPage
        },
        {
          path: 'pets/catalog',
          name: 'pet-catalog',
          component: PetListPage
        },
        {
          path: 'pets/detail',
          name: 'pet-detail',
          component: PetDetailPage
        },
        {
          path: 'pets/team',
          name: 'pet-team',
          component: PetTeamPage
        }
      ]
    }
  ]
})

export default router
