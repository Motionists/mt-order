<template>
  <div>
    <h2>登录</h2>
    <input v-model="mobile" placeholder="手机号" />
    <input v-model="password" placeholder="密码" type="password" />
    <button @click="doLogin">登录</button>
  </div>
</template>
<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'

const auth = useAuthStore()
const router = useRouter()
const mobile = ref('13800000000')
const password = ref('123456')

const doLogin = async () => {
  await auth.login(mobile.value, password.value)
  router.push('/')
}
</script>