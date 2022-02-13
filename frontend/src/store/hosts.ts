import axios from "axios"

export default {
    state: {
        hosts: {},
        allLabels: [],
    },
    mutations: {
        setHosts(state, hosts) {
            state.hosts = hosts
            state.allLabels = []
            if (hosts != null) {
                let obkeys = Object.keys(hosts)
                for (let i = 0; i < obkeys.length; i++) {
                    state.allLabels.push(obkeys[i])
                }
            }
        },
        createHost(state, payload) {
            state.hosts[payload.id] = payload.item
        },
        updateHost(state, payload) {
            state.hosts[payload.id] = payload.item
        },
        deleteHost(state, hostID) {
            delete(state.hosts[hostID])
        }
    },
    actions: {
        fetchHosts(context) {
            axios.get("api/hosts")
                .then((response) => {
                    context.commit("setHosts", response.data)
                })
        },
        async createHost(context, payload) {
            return axios.post("api/hosts", JSON.stringify(payload.item))
                .then((response) => {
                    context.commit("createHost", {
                        id: response.data.insert_id,
                        item: payload.item
                    })
                })
        },
        async updateHost(context, payload) {
            return axios.put("api/hosts/" + payload.id, JSON.stringify(payload.item))
                .then(() => {
                    context.commit("updateHost", payload)
                })
        },
        async deleteHost(context, id) {
            return axios.delete("api/hosts/" + id)
                .then(() => {
                    context.commit("deleteHost", id)
                })
        }
    },
}
