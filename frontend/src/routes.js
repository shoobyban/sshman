import Home from "./components/Home.vue"
import Login from "./components/Login.vue"
import Users from "./components/Users.vue"
import Hosts from "./components/Hosts.vue"
import Register from "./components/Register.vue"
import NotFound from "./components/NotFound.vue"

export default [
    { path: '/', component: Home },
    { path: '/users', component: Users },
    { path: '/hosts', component: Hosts },
    { path: '/login', component: Login },
    { path: '/register', component: Register },
    { path: '/404', component: NotFound },
]