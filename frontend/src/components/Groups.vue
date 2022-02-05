<script>
import VuexCRUD from './VuexCRUD.vue'
import { mapState, mapActions, mapGetters } from 'vuex'

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
    mounted() {
        this.fetchAll()
    },
}
</script>

<template>
    <div>
        <VuexCRUD 
            v-if="groups"
            resourceName="Groups" 
            v-model="groups.groups"
            orderBy="label" 
            @create="createGroups"
            @update="updateGroups"
            @delete="deleteGroups"
            @fetch="fetchAll"
            idField="label"
            :searchFields="['label', 'hosts', 'users']"
            :fields="[
                {label: 'Label', index: 'label', placeholder: 'group1', type:'text'},
                {label: 'Users', index: 'users', placeholder: 'email@host1,email@host2', type:'multiselect', options: users.allEmails},
                {label: 'Groups', index: 'hosts', placeholder: 'host1,host2', type:'multiselect', options: hosts.allLabels},
                ]" />
    </div>
</template>