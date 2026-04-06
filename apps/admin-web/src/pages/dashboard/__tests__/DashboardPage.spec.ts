import { mount, RouterLinkStub } from "@vue/test-utils"

import DashboardPage from "../DashboardPage.vue"

test("renders admin dashboard entries", () => {
  const wrapper = mount(DashboardPage, {
    global: {
      stubs: {
        RouterLink: RouterLinkStub,
      },
    },
  })

  expect(wrapper.text()).toContain("运营后台")
  expect(wrapper.text()).toContain("配置中心")
  expect(wrapper.text()).toContain("GM 操作台")
  expect(wrapper.text()).toContain("玩家查询")
  expect(wrapper.text()).toContain("运营报表")
  expect(wrapper.text()).toContain("调度与审计")
  expect(wrapper.text()).toContain("幻兽、技能、地图副本、修行与成长配置")
  expect(wrapper.text()).toContain("发奖、补偿、封禁、解卡与运营处置")
  expect(wrapper.text()).toContain("玩家信息、资源、背包与战斗记录查询")
  expect(wrapper.text()).toContain("新增、留存、付费、VIP 与资源产消")
  expect(wrapper.text()).toContain("任务补跑、执行状态、GM 日志与资源流水")
  expect(wrapper.text()).toContain("当前只提供后台壳入口")
  expect(wrapper.text()).toContain("真实 GM 执行与审计闭环待后续接入")
  expect(wrapper.text()).toContain("本轮不开放玩家检索与详情查询")
  expect(wrapper.text()).toContain("报表与图表数据接入不在本轮范围")
  expect(wrapper.text()).toContain("任务与审计后台待后续模块落地后接入")

  const routes = wrapper
    .findAllComponents(RouterLinkStub)
    .map((component) => component.props("to"))

  expect(routes).toContainEqual({ name: "config" })
  expect(routes).toContainEqual({ name: "gm" })
})
