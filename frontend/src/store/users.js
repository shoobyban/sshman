import axios from "axios"

export default {
    state: {
        users: [],
        allEmails: [],
    },
    mutations: {
        setUsers(state, users) {
            state.users = users
            state.allEmails = []
            for (let i = 0; i < users.length; i++) {
                state.allEmails.push(users[i].email)
            }
        },
        addUser(state, user) {
            state.users.push(user)
        },
        updateUser(state, user) {
            let index = state.users.findIndex((c) => c.id == user.id)
            if (index > -1) {
                state.users[index] = user
            }
        },
        deleteUser(state, userID) {
            let index = state.users.findIndex((c) => c.id == userID)
            if (index > -1) {
                state.users.splice(index, 1)
            }
        }
    },
    actions: {
        fetchUsers(context) {
            axios.get("api/users")
                .then((response) => {
                    context.commit("setUsers", response.data)
                })
        },
        async addUser(context, user) {
            return axios.post("api/users", JSON.stringify(user))
                .then((response) => {
                    context.commit("addUser", {
                        id: response.data.insert_id,
                        ...user
                    })
                })
        },
        async updateUser(context, user) {
            return axios.put("api/users/" + user.email, JSON.stringify(user))
                .then((response) => {
                    context.commit("updateUser", user)
                })
        },
        async deleteUser(context, user) {
            return axios.delete("api/user/" + user.email)
                .then((response) => {
                    context.commit("deleteUser", user.Email)
                })
        }
    },
}
