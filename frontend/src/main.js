import { createApp } from 'vue'
import routes from './routes'
import Container from './components/Container.vue'
import { createRouter,createWebHashHistory } from 'vue-router'
import store from './store'

const router = createRouter({
    // 4. Provide the history implementation to use. We are using the hash history for simplicity here.
    history: createWebHashHistory(),
    routes, // short for `routes: routes`
})

createApp(Container).use(router).use(store).mount('#app')
