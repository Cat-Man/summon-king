import { flushPromises, mount } from "@vue/test-utils"

import { getWorldMap } from "@/api/modules/dungeon"

import WorldMapPage from "../WorldMapPage.vue"

const push = vi.fn()

vi.mock("@/api/modules/dungeon", () => ({
  getWorldMap: vi.fn(),
}))

vi.mock("vue-router", () => ({
  useRouter: () => ({
    push,
  }),
}))

test("renders world map from api", async () => {
  vi.mocked(getWorldMap).mockResolvedValue({
    name: "玄境",
    cities: [
      {
        city_id: 1,
        name: "晨曦城",
        region: "东境",
        loc_x: 110.5,
        loc_y: 220.4,
        dungeons: [
          {
            dungeon_id: 1,
            dungeon_name: "妖窟试炼",
          },
          {
            dungeon_id: 2,
            dungeon_name: "寒渊裂隙",
          },
        ],
      },
      {
        city_id: 2,
        name: "霞光堡",
        region: "南境",
        loc_x: 190.8,
        loc_y: 180.1,
        dungeons: [
          {
            dungeon_id: 2,
            dungeon_name: "寒渊裂隙",
          },
          {
            dungeon_id: 1,
            dungeon_name: "妖窟试炼",
          },
        ],
      },
    ],
  })

  const wrapper = mount(WorldMapPage)
  await flushPromises()

  expect(wrapper.text()).toContain("玄境")
  expect(wrapper.text()).toContain("晨曦城")
  expect(wrapper.text()).toContain("妖窟试炼")
  expect(wrapper.text()).toContain("寒渊裂隙")

  await wrapper.get('[data-city-id="1"][data-dungeon-id="2"]').trigger("click")

  expect(push).toHaveBeenCalledWith({
    name: "dungeon",
    query: {
      dungeon_id: "2",
    },
  })
})
