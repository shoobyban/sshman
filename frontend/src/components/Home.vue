<script>
import { mapState, mapActions } from 'vuex'

export default {
    name: 'Home',
    computed: {
        ...mapState({
            hosts: state => state.hosts,
            groups: state => state.groups,
            users: state => state.users,
        }),
    },
    mounted() {
        this.fetchAll()
    },
    methods: {
        ...mapActions([
            'fetchHosts',
            'fetchGroups',
            'fetchUsers',
        ]),
        fetchAll() {
            this.fetchHosts()
            this.fetchGroups()
            this.fetchUsers()
        },
        len(obj) {
            return Object.keys(obj).length
        },
    },
}
</script>

<template>
    <div>       
        <h1 class="text-xl pb-5">
            Home
        </h1>
        <div v-if="len(users.users) == 0 && len(hosts.hosts) == 0" class="w-full grid place-content-center">
        <div class="bg-blue-50 rounded-lg border-2 border-sky-500 p-3 shadow-xl max-w-xl">
            <h2 class="text-large font-bold mb-4">
                Starting with a blank slate?
            </h2>
            <div class="text-large">
                <p class="mb-2">
                    You don't have a configuration file in ~/.ssh/.sshman yet. 
                </p>
                <p class="mb-2">
                    Please create 
                    <router-link to="/users" class="text-blue-500">
                        users
                    </router-link> 
                    and 
                    <router-link to="/hosts" class="text-blue-500">
                        hosts
                    </router-link>. You will see them here and in the CRUD area.
                </p>
                <p class="mb-2">
                    Don't forget to add <b>groups</b> to both your users and hosts, this will connect them together.
                </p>
                <p class="mb-2">
                    Users will be uploaded to the hosts automatically by matching groups.
                </p>
            </div>
            </div>
        </div>
        <div v-else>
            <div class="w-full grid place-content-center">
                {{ len(users.users) }} users
                {{ len(hosts.hosts) }} hosts
                {{ len(groups.groups) }} groups
            </div>
        <div class="flex mt-10">
            <div class="w-1/3 overflow-hidden mr-3">
                <h4 class="font-bold">
                    Users
                </h4>
                <div v-for="(user, idx) in users.users" :key="idx">
                    {{ user.email }}
                </div>
            </div>
            <div class="w-1/3 overflow-hidden mr-3">
                <h4 class="font-bold">
                    Hosts
                </h4>
                <div v-for="(host, idx) in hosts.hosts" :key="idx">
                    {{ host.alias }}
                </div>
            </div>
            <div class="w-1/3 overflow-hidden">
                <h4 class="font-bold">
                    Groups
                </h4>
                <div v-for="(group, idx) in groups.groups" :key="idx">
                    {{ group.label }}
                </div>
            </div>
        </div>
        </div>
    </div>      
</template>