import { createRouter, createWebHistory, type RouteRecordRaw } from "vue-router"

import AdminLayout from "@/layouts/AdminLayout.vue"
import DashboardPage from "@/pages/dashboard/DashboardPage.vue"
import SectionPlaceholderPage from "@/pages/sections/SectionPlaceholderPage.vue"

const routes: RouteRecordRaw[] = [
  {
    path: "/",
    component: AdminLayout,
    children: [
      {
        path: "",
        redirect: { name: "dashboard" },
      },
      {
        path: "dashboard",
        name: "dashboard",
        component: DashboardPage,
      },
      {
        path: "config",
        name: "config",
        component: SectionPlaceholderPage,
        props: {
          title: "配置中心",
          summary: "聚合幻兽、技能、地图副本、修行、塔与成长等配置域，后续承接草稿编辑、校验、发布与回滚。",
          status: "当前阶段只完成后台壳接入，真实配置读写链路待后续模块补齐。",
        },
      },
      {
        path: "gm",
        name: "gm",
        component: SectionPlaceholderPage,
        props: {
          title: "GM 操作台",
          summary: "聚合人工运营处置入口，后续承接发奖、补偿、封禁、解卡、发邮件与异常修复。",
          status: "当前阶段只保留入口位置，GM 执行链路、二次确认与审计闭环待后续模块接入。",
        },
      },
    ],
  },
]

export const router = createRouter({
  history: createWebHistory(),
  routes,
})
