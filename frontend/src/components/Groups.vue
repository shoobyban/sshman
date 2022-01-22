<script>
import CRUD from './CRUD.vue'

export default {
    name: 'Groups',
    components: {
        CRUD
    },
}
</script>

<template>
    <div>
        <CRUD 
            resourceName="Groups" 
            endpoint="/api/groups"
            orderBy="label" 
            :fields="[
                {label: 'Label', index: 'label', placeholder: 'group1', type:'text'},
                {label: 'Users', index: 'users', placeholder: 'email@host1,email@host2', type:'multiselect'},
                {label: 'Hosts', index: 'hosts', placeholder: 'host1,host2', type:'multiselect'},
            ]" />
    </div>
</template>