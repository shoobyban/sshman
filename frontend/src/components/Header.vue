<script>
import routes from '../routes.js'
import { mapActions } from 'vuex'

export default {
    data() {
        return {
            syncRunning: false,
            routes: routes,
            currentRoute: this.$route.path,
        }
    },
    computed: {
        routesWithLabel() {
            return this.routes.filter(route => route.label)
        },
    },
    mounted: () => {
        if (this == undefined) {
            return
        }
        this.refreshTheme()
    },
    methods: {
        ...mapActions([
            'syncHosts',
            'stopSync',
        ]),
        async sync() {
            let syncButton = document.getElementById('sync-button')
            if (this.syncRunning) {
                this.stopSync().then(() => {
                    this.syncRunning = false
                })
                return
            }
            syncButton.classList.add('sync-spin')
            this.syncRunning = true
            this.syncHosts().then(() => {
                syncButton.classList.remove('sync-spin')
                this.syncRunning = false
                this.fetchHosts()
            })
        },
        toggleTheme() {
            this.$store.dispatch("toggleTheme")
            this.refreshTheme()
        },
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
        isDarkTheme() {
            return this.$store.state.theme.theme === 'dark'
        },
    },
}
</script>

<template>
    <nav id="header" class="bg-white dark:bg-gray-800 fixed w-full z-10 top-0 shadow">
        <div class="flex justify-between h-11 mt-0 pt-2">
            <div class="ml-2 md:pl-0 flex">
                <a class="text-gray-900 dark:text-white text-base xl:text-xl no-underline hover:no-underline font-bold" href="/">
                    <i class="fas fa-link text-pink-600 pr-3 mb-4" />
                    SSHMan
                </a>
                <div v-for="(route, i) in routesWithLabel" :key="i" class="ml-10 flex-1">
                    <router-link 
                        :to="route.path"
                        class="align-middle text-gray-900 dark:text-white no-underline"
                        >
                        {{ route.label }}
                    </router-link>
                </div>
            </div>
            <div class="mr-2 flex">
                <button id="sync-button" data-tooltip-target="tooltip-sync" class="headerbtn" @click="sync">
                    <i class="fas fa-sync mr-2" />
                    Update
                </button>
                <div id="tooltip-sync" role="tooltip" class="headertooltip">
                    Download users from servers and update the storage
                    <div class="tooltip-arrow" data-popper-arrow />
                </div>
                <button id="theme-toggle" type="button" class="headerbtn" @click="toggleTheme">
                    <svg id="theme-toggle-dark-icon" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path d="M17.293 13.293A8 8 0 016.707 2.707a8.001 8.001 0 1010.586 10.586z" /></svg>
                    <svg id="theme-toggle-light-icon" class="hidden w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg"><path d="M10 2a1 1 0 011 1v1a1 1 0 11-2 0V3a1 1 0 011-1zm4 8a4 4 0 11-8 0 4 4 0 018 0zm-.464 4.95l.707.707a1 1 0 001.414-1.414l-.707-.707a1 1 0 00-1.414 1.414zm2.12-10.607a1 1 0 010 1.414l-.706.707a1 1 0 11-1.414-1.414l.707-.707a1 1 0 011.414 0zM17 11a1 1 0 100-2h-1a1 1 0 100 2h1zm-7 4a1 1 0 011 1v1a1 1 0 11-2 0v-1a1 1 0 011-1zM5.05 6.464A1 1 0 106.465 5.05l-.708-.707a1 1 0 00-1.414 1.414l.707.707zm1.414 8.486l-.707.707a1 1 0 01-1.414-1.414l.707-.707a1 1 0 011.414 1.414zM4 11a1 1 0 100-2H3a1 1 0 000 2h1z" fill-rule="evenodd" clip-rule="evenodd" /></svg>
                </button>
                <a href="https://github.com/shoobyban/sshman" data-tooltip-target="tooltip-github" target="new">
                    <img v-if="isDarkTheme()" src="/github-inverse.png" alt="Github" class="h-8 w-8">
                    <img v-else src="/github.png" alt="Github" class="h-8 w-8">
                </a>
                <div id="tooltip-github" role="tooltip" class="headertooltip">
                    Link to the Github repository
                    <div class="tooltip-arrow" data-popper-arrow />
                </div>
            </div>
        </div>
    </nav>
</template>

<style>
.headertooltip {
    @apply inline-block absolute invisible z-10 py-2 px-3 text-sm font-medium text-white bg-gray-900 rounded-lg shadow-sm opacity-0 transition-opacity duration-300 tooltip dark:bg-gray-700;
}
.headerbtn {
    @apply mr-3 mb-1 w-1/2 text-black dark:text-white hover:bg-black hover:text-white font-medium inline-flex items-center justify-center rounded-lg text-sm px-3 py-2 text-center sm:w-auto;
}
</style>