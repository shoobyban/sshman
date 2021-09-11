<script setup>
import { computed, ref } from 'vue'
import { useStore } from 'vuex'

const store = useStore()

const name = computed(() => {
    return store.state.user.name
})

const newName = ref('')

function saveName() {
    store.dispatch('saveName', newName.value)
    newName.value = ''
}
</script>

<template>
    <h1 class="mb-6 text-3xl font-extrabold">Home</h1>
    <p class="mb-6">{{name}}</p>
    <input type="text" class="border mr-4 border-gray rounded" v-model="newName" />
    <button class="p2 text-white bg-indigo-600 rounded" @click="saveName">Save</button>
</template>
