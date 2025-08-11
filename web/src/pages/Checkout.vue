<template>
  <div>
    <h2>结算</h2>
    <div>
      <label>商家ID：</label>
      <input v-model.number="mid" type="number" />
    </div>
    <button @click="createOrder">创建订单</button>
    <div v-if="orderId">
      <p>订单ID: {{ orderId }}</p>
      <button @click="pay">支付（模拟）</button>
      <p>状态：{{ status }}</p>
    </div>
  </div>
</template>
<script setup lang="ts">
import http from '../api/http'
import { ref } from 'vue'

const mid = ref<number>(1)
const orderId = ref<number>()
const status = ref('')

const createOrder = async () => {
  const { data } = await http.post('/orders', { merchant_id: mid.value })
  orderId.value = data.order_id
  status.value = data.status
}

const pay = async () => {
  if (!orderId.value) return
  await http.post(`/orders/${orderId.value}/pay`)
  const es = new EventSource(`/api/v1/orders/${orderId.value}/stream`)
  es.addEventListener('status', (e: any) => {
    const d = JSON.parse(e.data)
    status.value = d.status
    if (['done', 'canceled'].includes(status.value)) es.close()
  })
}
</script>