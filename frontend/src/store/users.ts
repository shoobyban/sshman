import axios from "axios"

export default {
    state: {
        users: {},
        allEmails: [],
    },
    mutations: {
        setUsers(state, users) {
            console.log('setUsers', users)
            state.users = users
            state.allEmails = []
            if (users != null) {
                let obkeys = Object.keys(users)
                for (let i = 0; i < obkeys.length; i++) {
                    state.allEmails.push(users[obkeys[i]].email)
                }
            }
        },
        createUser(state, payload) {
            state.users[payload.id] = payload.item
        },
        updateUser(state, payload) {
            state.users[payload.id] = payload.item
        },
        deleteUser(state, userID) {
            delete(state.users[userID])
        }
    },
    actions: {
        fetchUsers(context) {
            axios.get("api/users")
                .then((response) => {
                    context.commit("setUsers", response.data)
                })
        },
        async createUser(context, payload) {
            return axios.post("api/users", JSON.stringify(payload.item))
                .then((response) => {
                    context.commit("createUser", {
                        id: response.data.insert_id,
                        item: payload.item
                    })
                })
        },
        async updateUser(context, payload) {
            return axios.put("api/users/" + payload.id, JSON.stringify(payload.item))
                .then(() => {
                    context.commit("updateUser", {
                        id: payload.id,
                        item: payload.item
                    })
                })
        },
        async deleteUser(context, id) {
            return axios.delete("api/users/" + id)
                .then(() => {
                    context.commit("deleteUser", id)
                })
        }
    },
}
