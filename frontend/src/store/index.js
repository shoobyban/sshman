import { createStore } from 'vuex'

import user from '../store/modules/user'

const store = createStore({
    modules: {
        user,
    }
})

export default store