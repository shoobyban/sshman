import axios from "axios"

export default {
    state: {
        keys: [],
    },
    mutations: {
        setKeys(state, keys) {
            state.keys = keys
        },
    },
    actions: {
        fetchKeys(context, keyType) {
            axios.get("api/keys?type="+keyType)
                .then((response) => {
                    context.commit("setKeys", response.data)
                })
        },
    },
}
