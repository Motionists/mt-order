<template>
  <div>
    <h2>门店示例</h2>
    <div v-for="m in merchants" :key="m.id" style="margin: 12px 0">
      <h3>{{ m.name }}</h3>
      <button @click="loadDishes(m.id)">查看菜品</button>
      <div v-if="dishes[m.id]">
        <div v-for="d in dishes[m.id]" :key="d.id" style="display:flex;gap:8px;align-items:center">
          <span>{{ d.name }} ¥{{ (d.price/100).toFixed(2) }}</span>
          <button @click="add(m.id, d.id)">加入购物车</button>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import http from '../api/http'
import { useCartStore } from '../stores/cart'
import { onMounted, reactive, ref } from 'vue'

const merchants = ref<any[]>([])
const dishes = reactive<Record<number, any[]>>({})
const cart = useCartStore()

onMounted(async () => {
  const { data } = await http.get('/merchants')
  merchants.value = data
})

const loadDishes = async (mid: number) => {
  const { data } = await http.get(`/merchants/${mid}/dishes`)
  dishes[mid] = data
}
const add = async (mid: number, did: number) => {
  await cart.add(mid, did, 1)
  alert('已加入购物车')
}
</script>