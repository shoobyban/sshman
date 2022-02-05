import { createStore } from "vuex"
// import { auth } from "./auth.js"
import users from "./users"
import groups from "./groups"
import hosts from "./hosts"
import keys from "./keys"

const store = createStore({
    modules: {
        //      auth,
        users, 
        groups,
        hosts,
        keys,
    },
  })
  
  export default store
  