import { ref } from "vue"
import { defineStore } from "pinia"

export const useResourceSyncStore = defineStore("resource-sync", () => {
  const version = ref(0)

  function touch() {
    version.value += 1
  }

  return {
    version,
    touch,
  }
})
