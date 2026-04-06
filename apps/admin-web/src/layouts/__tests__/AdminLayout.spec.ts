import { mount, RouterLinkStub } from "@vue/test-utils"

import AdminLayout from "../AdminLayout.vue"

test("renders admin shell navigation", () => {
  const wrapper = mount(AdminLayout, {
    global: {
      stubs: {
        RouterLink: RouterLinkStub,
        RouterView: true,
      },
    },
  })

  expect(wrapper.text()).toContain("召唤之王")
  expect(wrapper.text()).toContain("运营后台")
  expect(wrapper.text()).toContain("工作台")
  expect(wrapper.text()).toContain("配置中心")
  expect(wrapper.text()).toContain("GM 操作台")
})
