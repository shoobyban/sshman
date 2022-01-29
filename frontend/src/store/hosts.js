import axios from "axios"
import _ from "lodash"

export default {
    state: {
        hosts: [],
        allLabels: [],
    },
    mutations: {
        setHosts(state, hosts) {
            state.hosts = hosts
            state.allLabels = []
            _.forEach(hosts, (host, label) => { state.allLabels.push(label) })
        },
        addHost(state, host) {
            state.hosts.push(host)
        },
        updateHost(state, host) {
            let index = state.hosts.findIndex((c) => c.id == host.id)
            if (index > -1) {
                state.hosts[index] = host
            }
        },
        deleteHost(state, hostID) {
            let index = state.hosts.findIndex((c) => c.id == hostID)
            if (index > -1) {
                state.hosts.splice(index, 1)
            }
        }
    },
    actions: {
        fetchHosts(context) {
            axios.get("api/hosts")
                .then((response) => {
                    context.commit("setHosts", response.data)
                })
        },
        async addHost(context, host) {
            return axios.post("api/hosts", JSON.stringify(host))
                .then((response) => {
                    context.commit("addHost", {
                        id: response.data.insert_id,
                        ...host
                    })
                })
        },
        async updateHost(context, host) {
            return axios.put("api/hosts/" + host.email, JSON.stringify(host))
                .then((response) => {
                    context.commit("updateHost", host)
                })
        },
        async deleteHost(context, host) {
            return axios.delete("api/host/" + host.email)
                .then((response) => {
                    context.commit("deleteHost", host.Email)
                })
        }
    },
}
