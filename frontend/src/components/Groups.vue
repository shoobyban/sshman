<script>
import VuexCRUD from './VuexCRUD.vue'
import { mapState, mapActions } from 'vuex'

export default {
    name: 'Groups',
    components: {
        VuexCRUD
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
            'createGroups',
            'updateGroups',
            'deleteGroups',
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
        <VuexCRUD 
            v-if="groups"
            v-model="groups.groups"
            resource-name="Groups" 
            order-by="label" 
            id-field="label"
            :search-fields="['label', 'hosts', 'users']"
            :fields="[
                {label: 'Label', index: 'label', placeholder: 'group1', type:'text'},
                {label: 'Users', index: 'users', placeholder: 'email@host1,email@host2', type:'multiselect', options: users.allEmails},
                {label: 'Groups', index: 'hosts', placeholder: 'host1,host2', type:'multiselect', options: hosts.allLabels},
                ]" 
            @create="createGroups"
            @update="updateGroups"
            @delete="deleteGroups"
            @fetch="fetchAll"
            />
    </div>
</template>