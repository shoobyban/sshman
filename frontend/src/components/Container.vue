<script>
import Header from './Header.vue'
import Logs from './Logs.vue'

export default {
    components: {
        Logs,
        Header,
    },
    data() {
        return {
            title: 'Container'
        }
    },
    mounted() {
        this.$store.dispatch("initTheme")
        this.refreshTheme()
    },
    methods: {
        refreshTheme() {
            if (this.$store.state.theme.theme === 'dark') {
                document.getElementById('theme-toggle-dark-icon').classList.add('hidden')
                document.getElementById('theme-toggle-light-icon').classList.remove('hidden')
                document.body.classList.add('dark')
            } else {
                document.getElementById('theme-toggle-dark-icon').classList.remove('hidden')
                document.getElementById('theme-toggle-light-icon').classList.add('hidden')
                document.body.classList.remove('dark')
            }
        }
    }
}
</script>

<template>
    <div class="flex flex-col h-screen">
        <Header />
        <div class="flex flex-grow overflow-hidden dark:bg-gray-900 dark:text-white bg-white pt-16">
            <div class="w-full bg-gray-50 dark:bg-gray-900 dark:text-white relative overflow-y-auto p-5">
                <router-view />
            </div>
            <div class="bg-gray-100 max-h-full dark:bg-gray-900 dark:text-white text-blue-100 md:w-80 pa-2 md:relative">
                <h2 class="text-center text-black mb-2 mt-2 font-bold">
                    <span>Server Logs</span>
                </h2>
                <div class="bg-white w-full overflow-y-auto text-black">
                    <Logs url="/api/logs" />
                </div>
            </div>
        </div>
    </div>
</template>
