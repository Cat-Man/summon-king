import { createRouter, createWebHistory } from 'vue-router'

import ArenaPage from '@/pages/arena/ArenaPage.vue'
import GameLayout from '@/layouts/GameLayout.vue'
import InventoryPage from '@/pages/assets/InventoryPage.vue'
import CultivationPage from '@/pages/cultivation/CultivationPage.vue'
import DungeonRunPage from '@/pages/dungeons/DungeonRunPage.vue'
import BonePage from '@/pages/growth/BonePage.vue'
import ManorPage from '@/pages/growth/ManorPage.vue'
import HomePage from '@/pages/home/HomePage.vue'
import WorldMapPage from '@/pages/maps/WorldMapPage.vue'
import PetDetailPage from '@/pages/pets/PetDetailPage.vue'
import PetListPage from '@/pages/pets/PetListPage.vue'
import PetTeamPage from '@/pages/pets/PetTeamPage.vue'
import RankingPage from '@/pages/ranking/RankingPage.vue'
import SoulPage from '@/pages/growth/SoulPage.vue'
import SpiritPage from '@/pages/growth/SpiritPage.vue'
import PagodaPage from '@/pages/tower/PagodaPage.vue'
import SpiritTowerPage from '@/pages/tower/SpiritTowerPage.vue'

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
        },
        {
          path: 'maps/world',
          name: 'world-map',
          component: WorldMapPage
        },
        {
          path: 'dungeons/run',
          name: 'dungeon-run',
          component: DungeonRunPage
        },
        {
          path: 'cultivation',
          name: 'cultivation',
          component: CultivationPage
        },
        {
          path: 'growth/bone',
          name: 'growth-bone',
          component: BonePage
        },
        {
          path: 'growth/spirit',
          name: 'growth-spirit',
          component: SpiritPage
        },
        {
          path: 'growth/soul',
          name: 'growth-soul',
          component: SoulPage
        },
        {
          path: 'growth/manor',
          name: 'growth-manor',
          component: ManorPage
        },
        {
          path: 'tower/pagoda',
          name: 'tower-pagoda',
          component: PagodaPage
        },
        {
          path: 'tower/spirit',
          name: 'tower-spirit',
          component: SpiritTowerPage
        },
        {
          path: 'arena',
          name: 'arena',
          component: ArenaPage
        },
        {
          path: 'ranking',
          name: 'ranking',
          component: RankingPage
        }
      ]
    }
  ]
})

export default router
