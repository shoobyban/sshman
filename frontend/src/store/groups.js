import axios from "axios"
import _ from "lodash"

export default {
    state: {
        groups: [],
        allLabels: [],
        messages: [],
    },
    mutations: {
        setGroups(state, groups) {
            state.groups = groups
            state.allLabels = []
            _.forEach(groups, (group, label) => { state.allLabels.push(label) })
        },
        addGroup(state, group) {
            state.groups.push(group)
        },
        updateGroup(state, group) {
            let index = state.groups.findIndex((c) => c.id == group.id)
            if (index > -1) {
                state.groups[index] = group
            }
        },
        deleteGroup(state, groupID) {
            let index = state.groups.findIndex((c) => c.id == groupID)
            if (index > -1) {
                state.groups.splice(index, 1)
            }
        }
    },
    actions: {
        fetchGroups(context) {
            axios.get("api/groups")
                .then((response) => {
                    context.commit("setGroups", response.data)
                })
        },
        async addGroup(context, group) {
            return axios.post("api/groups", JSON.stringify(group))
                .then((response) => {
                    context.commit("addGroup", {
                        id: response.data.insert_id,
                        ...group
                    })
                })
        },
        async updateGroup(context, group) {
            return axios.put("api/groups/" + group.email, JSON.stringify(group))
                .then((response) => {
                    context.commit("updateGroup", group)
                })
        },
        async deleteGroup(context, group) {
            return axios.delete("api/group/" + group.email)
                .then((response) => {
                    context.commit("deleteGroup", group.Email)
                })
        }
    },
}
