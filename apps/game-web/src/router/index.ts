import { createRouter, createWebHistory } from "vue-router"

import GameLayout from "@/layouts/GameLayout.vue"
import ArenaPage from "@/pages/arena/ArenaPage.vue"
import CultivationPage from "@/pages/cultivation/CultivationPage.vue"
import LoginPage from "@/pages/auth/LoginPage.vue"
import DungeonRunPage from "@/pages/dungeons/DungeonRunPage.vue"
import BonePage from "@/pages/growth/BonePage.vue"
import ManorPage from "@/pages/growth/ManorPage.vue"
import SoulPage from "@/pages/growth/SoulPage.vue"
import SpiritPage from "@/pages/growth/SpiritPage.vue"
import HomePage from "@/pages/home/HomePage.vue"
import WorldMapPage from "@/pages/maps/WorldMapPage.vue"
import PetPage from "@/pages/pet/PetPage.vue"
import RankingPage from "@/pages/ranking/RankingPage.vue"
import { readSessionSnapshot } from "@/stores/session"
import PagodaPage from "@/pages/tower/PagodaPage.vue"
import SpiritTowerPage from "@/pages/tower/SpiritTowerPage.vue"

export const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/login",
      name: "login",
      component: LoginPage,
      meta: { public: true },
    },
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
          path: "pet",
          name: "pet",
          component: PetPage,
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
          path: "arena",
          name: "arena",
          component: ArenaPage,
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

router.beforeEach((to) => {
  const session = readSessionSnapshot()
  const isLoggedIn = Boolean(session.token && session.playerId)

  if (to.meta.public) {
    if (to.name === "login" && isLoggedIn) {
      return { name: "home" }
    }
    return true
  }

  if (!isLoggedIn) {
    return {
      name: "login",
      query: { redirect: to.fullPath },
    }
  }

  return true
})
