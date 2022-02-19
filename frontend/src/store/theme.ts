
export default {
    state: {
        theme: {}
    },
    mutations: {
        SET_THEME(state, theme) {
            state.theme = theme
            localStorage.theme = theme
        }
    },
    actions: {
        initTheme({ commit }) {
            const cachedTheme = localStorage.theme ? localStorage.theme : false
            console.log('cached', cachedTheme)
            const userPrefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
            console.log('userPrefersDark', userPrefersDark)
            if (cachedTheme)
                commit('SET_THEME', cachedTheme)
            else if (userPrefersDark)
                commit('SET_THEME', 'dark')
            else
                commit('SET_THEME', 'light')
        },
        toggleTheme({ commit }) {
            switch (localStorage.theme) {
                case 'dark':
                    commit('SET_THEME', 'light')
                    break
                default:
                    commit('SET_THEME', 'dark')
                    break
            }
        }
    },
    getters: {
        getTheme: (state) => {
            return state.theme
        }
    }
}
