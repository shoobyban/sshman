import { createRouter, createWebHistory } from 'vue-router'
import MainMenu from '/src/components/MainMenu.vue'
const routes = [
    {
        path: '/',
        name: 'SSH Man',
        component: MainMenu,
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes,
})

export default router