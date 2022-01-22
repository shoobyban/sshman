import Home from "./components/Home.vue"
import Login from "./components/Login.vue"
import Users from "./components/Users.vue"
import Hosts from "./components/Hosts.vue"
import Groups from "./components/Groups.vue"
import Register from "./components/Register.vue"
import NotFound from "./components/NotFound.vue"

export default [
    { path: '/', component: Home, label: 'Home' },
    { path: '/users', component: Users, label: 'Users' },
    { path: '/hosts', component: Hosts, label: 'Hosts' },
    { path: '/groups', component: Groups, label: 'Groups' },
    { path: '/login', component: Login },
    { path: '/register', component: Register },
    { path: '/404', component: NotFound },
]