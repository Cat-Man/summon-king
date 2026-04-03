import { flushPromises, mount } from "@vue/test-utils"

import { getWorldMap } from "@/api/modules/dungeon"

import WorldMapPage from "../WorldMapPage.vue"

vi.mock("@/api/modules/dungeon", () => ({
  getWorldMap: vi.fn(),
}))

test("renders world map from api", async () => {
  vi.mocked(getWorldMap).mockResolvedValue({
    name: "玄境",
    cities: [
      { city_id: 1, name: "晨曦城", region: "东境", loc_x: 110.5, loc_y: 220.4 },
      { city_id: 2, name: "霞光堡", region: "南境", loc_x: 190.8, loc_y: 180.1 },
    ],
  })

  const wrapper = mount(WorldMapPage)
  await flushPromises()

  expect(wrapper.text()).toContain("玄境")
  expect(wrapper.text()).toContain("晨曦城")
})
