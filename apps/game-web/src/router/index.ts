import { createRouter, createWebHistory } from "vue-router"

import GameLayout from "@/layouts/GameLayout.vue"
import CultivationPage from "@/pages/cultivation/CultivationPage.vue"
import DungeonRunPage from "@/pages/dungeons/DungeonRunPage.vue"
import BonePage from "@/pages/growth/BonePage.vue"
import ManorPage from "@/pages/growth/ManorPage.vue"
import SoulPage from "@/pages/growth/SoulPage.vue"
import SpiritPage from "@/pages/growth/SpiritPage.vue"
import HomePage from "@/pages/home/HomePage.vue"
import WorldMapPage from "@/pages/maps/WorldMapPage.vue"
import RankingPage from "@/pages/ranking/RankingPage.vue"
import PagodaPage from "@/pages/tower/PagodaPage.vue"
import SpiritTowerPage from "@/pages/tower/SpiritTowerPage.vue"

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/",
      component: GameLayout,
      children: [
        {
          path: "",
          redirect: { name: "home" },
        },
        {
          path: "home",
          name: "home",
          component: HomePage,
        },
        {
          path: "map",
          name: "map",
          component: WorldMapPage,
        },
        {
          path: "dungeon",
          name: "dungeon",
          component: DungeonRunPage,
        },
        {
          path: "cultivation",
          name: "cultivation",
          component: CultivationPage,
        },
        {
          path: "growth/spirit",
          name: "growth-spirit",
          component: SpiritPage,
        },
        {
          path: "growth/bone",
          name: "growth-bone",
          component: BonePage,
        },
        {
          path: "growth/soul",
          name: "growth-soul",
          component: SoulPage,
        },
        {
          path: "growth/manor",
          name: "growth-manor",
          component: ManorPage,
        },
        {
          path: "tower/pagoda",
          name: "tower-pagoda",
          component: PagodaPage,
        },
        {
          path: "tower/spirit",
          name: "tower-spirit",
          component: SpiritTowerPage,
        },
        {
          path: "ranking",
          name: "ranking",
          component: RankingPage,
        },
      ],
    },
  ],
})
