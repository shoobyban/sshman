export default {
    state: () => ({ 
        name: 'Sam'
    }),
    getters: {},
    mutations: {
        SET_NAME(state, payload) {
            state.name = payload
        },
    },
    actions: {
        saveName({ commit }, payload) {
            commit('SET_NAME', payload)
        },
    },
}

