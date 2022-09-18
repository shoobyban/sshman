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
            title: 'Container',
            console: true,
        }
    },
    mounted() {
        this.$store.dispatch("initTheme")
        const noConsole = window.localStorage.getItem('console')=='false'
        this.console = !noConsole
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
        },
        toggleConsole() {
            console.log('toggle console', this.console)
            this.console = !this.console
            window.localStorage.setItem('console', this.console)
        },
    }
}
</script>

<template>
    <div class="flex flex-col h-screen">
        <Header :console="console" @toggle-console="toggleConsole" />
        <div class="flex flex-grow overflow-hidden dark:bg-gray-900 dark:text-white bg-white pt-16">
            <div class="w-full bg-gray-50 dark:bg-gray-900 dark:text-white relative overflow-y-auto p-5">
                <router-view />
            </div>
            <div v-show="console" class="bg-gray-100 dark:bg-gray-900 dark:text-white text-blue-100 md:w-80 pa-2 md:relative">
                <div class="text-black dark:text-white w-full overflow-y-auto">
                    <Logs url="/api/logs" />
                </div>
            </div>
        </div>
    </div>
</template>
