import Vue from 'vue'
import Router from 'vue-router'
import Main from './pages/Main.vue'

Vue.use(Router)

export default new Router({
    mode: 'history',
    base: process.env.BASE_URL,
    routes: [
        {
            path: '/',
            name: 'Home',
            component: Main,
        }
    ]
})