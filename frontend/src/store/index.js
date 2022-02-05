import { createStore } from "vuex"
// import { auth } from "./auth.js"
import users from "./users.ts"
import groups from "./groups.js"
import hosts from "./hosts.js"
import keys from "./keys.js"

const store = createStore({
    namespaced: true,
    modules: {
        //      auth,
        users, 
        groups,
        hosts,
        keys,
    },
  })
  
  export default store
  