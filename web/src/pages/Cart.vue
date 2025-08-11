<template>
  <div>
    <h2>购物车</h2>
    <div v-for="it in cart.items" :key="it.id" style="display:flex;gap:8px;align-items:center">
      <span>菜品ID: {{ it.dish_id }}</span>
      <input type="number" v-model.number="m[it.id]" min="1" style="width:60px" />
      <button @click="update(it.id)">更新</button>
      <button @click="remove(it.id)">删除</button>
    </div>
    <router-link to="/checkout">去结算</router-link>
  </div>
</template>
<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import { useCartStore } from '../stores/cart'
const cart = useCartStore()
const m = reactive<Record<number, number>>({})

onMounted(async () => {
  await cart.fetch()
  cart.items.forEach(i => m[i.id] = i.quantity)
})

const update = async (id: number) => {
  await cart.update(id, m[id])
}
const remove = async (id: number) => {
  await cart.remove(id)
}
</script>