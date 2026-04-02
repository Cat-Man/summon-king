import { createRouter, createWebHistory } from "vue-router"

import GameLayout from "@/layouts/GameLayout.vue"
import HomePage from "@/pages/home/HomePage.vue"

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
      ],
    },
  ],
})
