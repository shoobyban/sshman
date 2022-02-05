<script>
import VuexCRUD from './VuexCRUD.vue'
import { mapState, mapActions, mapGetters } from 'vuex'

export default {
    name: 'Users',
    components: {
        VuexCRUD
    },
    computed: {
        ...mapState({
            users: state => state.users,
            groups: state => state.groups,
        }),
    },
    methods: {
        ...mapActions([
            'fetchUsers',
            'createUser',
            'updateUser',
            'deleteUser',
            'fetchGroups',
        ]),
        fetchAll() {
            this.fetchUsers()
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
            v-if="users"
            resourceName="Users" 
            v-model="users.users"
            orderBy="email"
            @create="createUser"
            @update="updateUser"
            @delete="deleteUser"
            @fetch="fetchAll"
            idField="."
            :fields="[
                {label: 'Email', index: 'email', placeholder: 'sam@host.com', type:'email'},
                {label: 'Public Key (.pub)', hidefromlist:true, index: 'keyfile', placeholder: '~/.ssh/key.pub', type:'file'},
                {label: 'Name in key', index: 'name', placeholder: 'sam', type:'text'},
                {label: 'Groups', index: 'groups', placeholder: 'group1,group2', type:'multiselect', options: groups.allLabels},
                ]" />
    </div>
</template>