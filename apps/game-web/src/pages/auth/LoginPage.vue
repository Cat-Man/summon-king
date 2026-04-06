<template>
  <section class="login-page">
    <div class="login-panel">
      <p class="eyebrow">访客登录</p>
      <h1>进入召唤之王</h1>
      <p class="description">
        当前阶段先用最小访客登录打通会话链路。输入一个昵称，前端将拿到 `player_id` 和 `token`，随后进入首页主循环。
      </p>

      <form class="login-form" @submit.prevent="submit">
        <label for="nickname">昵称</label>
        <input id="nickname" v-model.trim="nickname" type="text" maxlength="12" placeholder="例如：北境召唤师" />
        <p v-if="errorMessage" class="error-text">{{ errorMessage }}</p>
        <button class="submit-btn" type="submit" :disabled="submitting">
          {{ submitting ? "登录中..." : "进入世界" }}
        </button>
      </form>
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed, ref } from "vue"
import { useRoute, useRouter } from "vue-router"

import { APIError, apiRequest } from "@/api/http"
import { useSessionStore } from "@/stores/session"

type GuestLoginResponse = {
  player_id: number
  token: string
  nickname: string
}

const router = useRouter()
const route = useRoute()
const sessionStore = useSessionStore()

const nickname = ref("")
const submitting = ref(false)
const errorMessage = ref("")

const normalizedNickname = computed(() => nickname.value.trim())
const redirectTarget = computed(() => {
  const redirect = route.query.redirect
  if (typeof redirect === "string" && redirect.startsWith("/")) {
    return redirect
  }
  return "/home"
})

async function submit() {
  if (!normalizedNickname.value) {
    errorMessage.value = "请输入昵称"
    return
  }

  submitting.value = true
  errorMessage.value = ""

  try {
    const payload = await apiRequest<GuestLoginResponse>("auth/guest-login", {
      method: "POST",
      body: JSON.stringify({
        nickname: normalizedNickname.value,
      }),
    })

    sessionStore.setSession({
      token: payload.token,
      playerId: payload.player_id,
      nickname: payload.nickname,
    })

    await router.push(redirectTarget.value)
  } catch (error) {
    if (error instanceof APIError) {
      errorMessage.value = error.message
    } else {
      errorMessage.value = "登录失败，请稍后重试"
    }
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: calc(100vh - 160px);
  display: grid;
  place-items: center;
}

.login-panel {
  width: min(100%, 520px);
  padding: 34px;
  border-radius: 28px;
  border: 1px solid rgba(247, 239, 225, 0.14);
  background:
    radial-gradient(circle at top right, rgba(104, 171, 255, 0.12), transparent 30%),
    linear-gradient(180deg, rgba(17, 20, 31, 0.94), rgba(10, 12, 18, 0.96));
  box-shadow: 0 36px 80px rgba(0, 0, 0, 0.36);
}

.eyebrow {
  margin: 0 0 10px;
  color: #8fb5ff;
  letter-spacing: 0.24em;
  font-size: 12px;
  text-transform: uppercase;
}

.login-panel h1,
.description {
  margin: 0;
}

.login-panel h1 {
  font-size: clamp(30px, 5vw, 44px);
}

.description {
  margin-top: 14px;
  color: rgba(247, 239, 225, 0.72);
  line-height: 1.7;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-top: 28px;
}

.login-form label {
  font-size: 14px;
  color: rgba(247, 239, 225, 0.78);
}

.login-form input {
  width: 100%;
  padding: 14px 16px;
  border: 1px solid rgba(247, 239, 225, 0.16);
  border-radius: 14px;
  background: rgba(247, 239, 225, 0.04);
  color: #fff8f0;
  font-size: 16px;
}

.login-form input:focus {
  outline: none;
  border-color: rgba(143, 181, 255, 0.6);
  box-shadow: 0 0 0 3px rgba(143, 181, 255, 0.12);
}

.error-text {
  margin: 0;
  color: #ff9c9c;
  font-size: 14px;
}

.submit-btn {
  margin-top: 8px;
  padding: 14px 18px;
  border: none;
  border-radius: 16px;
  background: linear-gradient(135deg, #f0a348, #db5d1f);
  color: #190f04;
  font-size: 16px;
  font-weight: 700;
  cursor: pointer;
}

.submit-btn:disabled {
  opacity: 0.65;
  cursor: not-allowed;
}
</style>
