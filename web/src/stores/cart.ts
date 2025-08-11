import { defineStore } from 'pinia'
import http from '../api/http'

export interface CartItem {
  id: number
  merchant_id: number
  dish_id: number
  quantity: number
}

export const useCartStore = defineStore('cart', {
  state: () => ({ items: [] as CartItem[] }),
  actions: {
    async fetch() {
      const { data } = await http.get('/cart')
      this.items = data
    },
    async add(merchantId: number, dishId: number, quantity = 1) {
      await http.post('/cart/items', { merchant_id: merchantId, dish_id: dishId, quantity })
      await this.fetch()
    },
    async update(id: number, quantity: number) {
      await http.patch(`/cart/items/${id}`, { quantity })
      await this.fetch()
    },
    async remove(id: number) {
      await http.delete(`/cart/items/${id}`)
      await this.fetch()
    }
  }
})