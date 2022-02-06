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
        }
    },
}
</script>

<template>
    <div>       
        <h3 class="text-xl pl-3 pt-3">Home</h3>
        <div class="flex">
            <div class="w-1/3">
                <h4 class="font-bold">Users</h4>
                <div v-for="(user, idx) in users.users" :key="idx">{{user.email}}</div>
            </div>
            <div class="w-1/3">
                <h4 class="font-bold">Hosts</h4>
                <div v-for="(host, idx) in hosts.hosts" :key="idx">{{host.alias}}</div>
            </div>
            <div class="w-1/3">
                <h4 class="font-bold">Groups</h4>
                <div v-for="(group, idx) in groups.groups" :key="idx">{{group.label}}</div>
            </div>
        </div>
    </div>      
</template>