import axios from "axios"

export default {
    state: {
        groups: {},
        allLabels: [],
    },
    mutations: {
        setGroups(state, groups) {
            state.groups = groups
            state.allLabels = []
            for (let key in groups) {
                state.allLabels.push(key)
            }
        },
        createGroup(state, payload) {
            state.groups[payload.id] = payload.item
        },
        updateGroup(state, payload) {
            state.groups[payload.id] = payload.item
        },
        deleteGroup(state, groupID) {
            delete(state.groups[groupID])
        }
    },
    actions: {
        fetchGroups(context) {
            axios.get("api/groups")
                .then((response) => {
                    context.commit("setGroups", response.data)
                })
        },
        async createGroup(context, payload) {
            return axios.post("api/groups", JSON.stringify(payload.item))
                .then((response) => {
                    context.commit("createGroup", {
                        id: response.data.insert_id,
                        item: payload.item
                    })
                })
        },
        async updateGroup(context, payload) {
            return axios.put("api/groups/" + payload.id, JSON.stringify(payload.item))
                .then(() => {
                    context.commit("updateGroup", payload)
                })
        },
        async deleteGroup(context, id) {
            return axios.delete("api/groups/" + id)
                .then(() => {
                    context.commit("deleteGroup", id)
                })
        }
    },
}
