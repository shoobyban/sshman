<script>
import Crud from './Crud.vue'
import { mapState, mapActions } from 'vuex'

export default {
    name: 'Groups',
    components: {
        Crud
    },
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
            'fetchGroups',
            'createGroup',
            'updateGroup',
            'deleteGroup',
            'fetchUsers',
            'fetchHosts',
        ]),
        fetchAll() {
            this.fetchGroups()
            this.fetchHosts()
            this.fetchUsers()
        }
    },
}
</script>

<template>
    <div>
        <Crud 
            v-if="groups"
            v-model="groups.groups"
            resource-name="Groups" 
            order-by="label" 
            id-field="."
            :search-fields="['label', 'hosts', 'users']"
            :fields="[
                {label: 'Label', index: 'label', placeholder: 'group1', type:'text', double: true},
                {label: 'Users', index: 'users', placeholder: 'email@host1,email@host2', type:'multiselect', options: users.allEmails},
                {label: 'Hosts', index: 'hosts', placeholder: 'host1,host2', type:'multiselect', options: hosts.allLabels},
                ]" 
            @create="createGroup"
            @update="updateGroup"
            @delete="deleteGroup"
            @fetch="fetchAll"
            />
    </div>
</template>