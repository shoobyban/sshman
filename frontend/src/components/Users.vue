<script>
import CRUD from './CRUD.vue'

export default {
    name: 'Users',
    components: {
        CRUD
    },
}
</script>

<template>
    <div>
        <CRUD 
            resourceName="Users" 
            endpoint="/api/users"
            orderBy="email" 
            :fields="[
                {label: 'Email', index: 'email', placeholder: 'sam@host.com', type:'email'},
                {label: 'Key', hidefromlist:true, index: 'name', placeholder: 'sam', type:'text'},
                {label: 'Name in key', index: 'name', placeholder: 'sam', type:'text'},
                {label: 'Groups', index: 'groups', placeholder: 'group1,group2', type:'multiselect'},
            ]" />
    </div>
</template>