<script>
import VuexCRUD from './VuexCRUD.vue'
import { mapState, mapActions } from 'vuex'

export default {
    name: 'Hosts',
    components: {
        VuexCRUD
    },
    computed: {
        ...mapState({
            hosts: state => state.hosts,
            groups: state => state.groups,
            keys: state => state.keys,
        }),
    },
    mounted() {
        this.fetchAll()
    },
    methods: {
        ...mapActions([
            'fetchHosts',
            'createHost',
            'updateHost',
            'deleteHost',
            'fetchGroups',
            'fetchKeys',
        ]),
        fetchAll() {
            this.fetchHosts()
            this.fetchGroups()
            this.fetchKeys('private')
        }
    },
}
</script>

<template>
    <div>
        <VuexCRUD
            v-if="hosts"
            v-model="hosts.hosts"
            resource-name="Hosts" 
            order-by="alias"
            id-field="."
            :search-fields="['alias', 'host', 'user', 'groups', 'key']"
            :fields="[
                {label: 'Alias', apikey: true, index: 'alias', placeholder: 'home.host', type:'text'},
                {label: 'Hostname', index: 'host', placeholder: '127.0.0.1:22', type:'text'},
                {label: 'Username', index: 'user', placeholder: 'root', type:'text'},
                {label: 'Keyfile', index: 'key', placeholder: '~/.ssh/keys.key', type:'select', options: keys.keys},
                {label: 'Groups', index: 'groups', placeholder: 'group1,group2', type:'multiselect', options: groups.allLabels},
                ]"
            @create="createHost"
            @update="updateHost"
            @delete="deleteHost"
            @fetch="fetchAll"
            /> 
    </div>
</template>