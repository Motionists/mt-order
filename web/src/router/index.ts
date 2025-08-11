import { createRouter, createWebHistory } from 'vue-router'

const Menu = () => import('../pages/Menu.vue')
const Cart = () => import('../pages/Cart.vue')
const Checkout = () => import('../pages/Checkout.vue')
const Login = () => import('../pages/Login.vue')

export default createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: Menu },
    { path: '/cart', component: Cart },
    { path: '/checkout', component: Checkout },
    { path: '/login', component: Login }
  ]
})