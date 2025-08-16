version: '3.8'
services:
  mysql:
    image: mysql:8.0
    restart: always
    environment:
      MYSQL_ROOT_PASSWORD: 'root_password'
      MYSQL_DATABASE: 'mt_order'
    ports:
      - "3306:3306"
    volumes:
      - ./server/migrations:/docker-entrypoint-initdb.d
  redis:
    image: redis:7-alpine
    restart: always
    ports:
      - "6379:6379"
<template>
  <div>
    <h2>购物车</h2>
    <div v-if="cart.items.length === 0">
      购物车为空
    </div>
    <div v-else>
      <div v-for="item in cart.items" :key="item.id" style="display:flex;gap:8px;align-items:center;margin:8px 0">
        <span>商家ID: {{ item.merchant_id }}</span>
        <span>菜品ID: {{ item.dish_id }}</span>
        <span>数量: {{ item.quantity }}</span>
        <button @click="updateQuantity(item.id, item.quantity + 1)">+</button>
        <button @click="updateQuantity(item.id, item.quantity - 1)" :disabled="item.quantity <= 1">-</button>
        <button @click="removeItem(item.id)">删除</button>
      </div>
      <router-link to="/checkout">
        <button style="margin-top:16px">去结算</button>
      </router-link>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useCartStore } from '../stores/cart'
import { onMounted } from 'vue'

const cart = useCartStore()

onMounted(() => {
  cart.fetch()
})

const updateQuantity = async (id: number, quantity: number) => {
  if (quantity < 1) return
  await cart.update(id, quantity)
}

const removeItem = async (id: number) => {
  await cart.remove(id)
}
</script>