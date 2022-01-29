<script>
import VuexCRUD from './VuexCRUD.vue'
import { mapState, mapActions, mapGetters } from 'vuex'

export default {
    name: 'Hosts',
    components: {
        VuexCRUD
    },
    computed: {
        ...mapState({
            hosts: state => state.hosts,
            groups: state => state.groups,
        }),
    },
    methods: {
        ...mapActions([
            'fetchHosts',
            'createHost',
            'updateHost',
            'deleteHost',
            'fetchGroups',
        ]),
        fetchAll() {
            this.fetchHosts()
            this.fetchGroups()
        }
    },
    mounted() {
        this.fetchAll()
    },
}
</script>

<template>
    <div>
        <VuexCRUD
            v-if="hosts"
            resourceName="Hosts" 
            v-model="hosts.hosts"
            orderBy="alias"
            @create="createHost"
            @update="updateHost"
            @delete="deleteHost"
            @fetch="fetchAll"
            idField="email"
            :fields="[
                {label: 'Alias', index: 'alias', placeholder: 'home.host', type:'text'},
                {label: 'Hostname', index: 'host', placeholder: '127.0.0.1:22', type:'text'},
                {label: 'Username', index: 'user', placeholder: 'root', type:'text'},
                {label: 'Keyfile', index: 'key', placeholder: '~/.ssh/keys.key', type:'text'},
                {label: 'Groups', index: 'groups', placeholder: 'group1,group2', type:'multiselect', options: groups.allLabels},
                ]" /> 
    </div>
</template>