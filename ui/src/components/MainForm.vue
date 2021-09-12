<template>
  <div>
    <UserForm v-if="type=='user'" :selected="selected" :groups="groups" :ts="ts" />
    <ServerForm v-else  :selected="selected" :groups="groups" :ts="ts" />
  </div>
</template>

<script>
import UserForm from './UserForm.vue'
import ServerForm from './ServerForm.vue'

export default {
    components: {
        UserForm,
        ServerForm,
    },
    props: {
        type: { type: String },
        selected: { type: Object },
        groups: { },
        ts: {  },
    },
}
</script>