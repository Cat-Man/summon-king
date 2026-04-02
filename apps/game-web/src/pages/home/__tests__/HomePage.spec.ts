import { mount } from "@vue/test-utils"

import HomePage from "../HomePage.vue"

test("renders home title", () => {
  const wrapper = mount(HomePage)

  expect(wrapper.text()).toContain("召唤之王")
})
