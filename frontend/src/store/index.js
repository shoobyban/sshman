import { createStore } from "vuex"
// import { auth } from "./auth.js"
import users from "./users.js"
import groups from "./groups.js"
import hosts from "./hosts.js"

const store = createStore({
    namespaced: true,
    modules: {
        //      auth,
        users, 
        groups,
        hosts,
    },
  })
  
  export default store
  