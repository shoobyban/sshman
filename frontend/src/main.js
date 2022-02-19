import { createApp } from 'vue'
import routes from './routes'
import Container from './components/Container.vue'
import { createRouter,createWebHashHistory } from 'vue-router'
import store from './store'
import 'flowbite'

import './index.css'

const router = createRouter({
    history: createWebHashHistory(),
    routes,
})

createApp(Container).use(router).use(store).mount('#app')
