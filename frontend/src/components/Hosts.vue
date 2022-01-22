<script>
import CRUD from './CRUD.vue'

export default {
    name: 'Hosts',
    components: {
        CRUD
    },
}
</script>

<template>
    <div>
        <CRUD 
            resourceName="Host" 
            endpoint="/api/hosts" 
            :fields="[
                {label: 'Alias', index: 'alias', placeholder: 'home.host', type:'text'},
                {label: 'Hostname', index: 'host', placeholder: '127.0.0.1:22', type:'text'},
                {label: 'Username', index: 'user', placeholder: 'root', type:'text'},
                {label: 'Keyfile', index: 'key', placeholder: '~/.ssh/keys.key', type:'file'},
                {label: 'Groups', index: 'groups', placeholder: 'group1,group2', type:'multiselect'},
                ]" 
        />
    </div>
</template>