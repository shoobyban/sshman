<script>
import { mapActions } from 'vuex'
export default {
    data() {
        return {
            syncRunning: false,
        }
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
        }
    }
}
</script>

<template>
    <nav id="header" class="bg-white fixed w-full z-10 top-0 shadow">
        <div class="flex justify-between h-11 mt-0 pt-2">
            <div class="ml-2 md:pl-0">
                <a class="text-gray-900 text-base xl:text-xl no-underline hover:no-underline font-bold" href="/">
                    <i class="fas fa-link text-pink-600 pr-3 mb-4" />
                    SSHMan
                </a>
            </div>
            <div class="mr-2 flex">
                <button id="sync-button" data-tooltip-target="tooltip-sync" class="mr-3 mb-1 w-1/2 text-black hover:bg-black hover:text-white focus:ring-4 focus:ring-blue-200 font-medium inline-flex items-center justify-center rounded-lg text-sm px-3 py-2 text-center sm:w-auto" @click="sync">
                    <i class="fas fa-sync mr-2" />
                    Update
                </button>
                <div id="tooltip-sync" role="tooltip" class="inline-block absolute invisible z-10 py-2 px-3 text-sm font-medium text-white bg-gray-900 rounded-lg shadow-sm opacity-0 transition-opacity duration-300 tooltip dark:bg-gray-700">
                    Download users from servers and update the storage
                    <div class="tooltip-arrow" data-popper-arrow />
                </div>
                <a href="https://github.com/shoobyban/sshman" data-tooltip-target="tooltip-github" target="new">
                    <img src="/github.png" alt="Github" class="h-8 w-8">
                </a>
                <div id="tooltip-github" role="tooltip" class="inline-block absolute invisible z-10 py-2 px-3 text-sm font-medium text-white bg-gray-900 rounded-lg shadow-sm opacity-0 transition-opacity duration-300 tooltip dark:bg-gray-700">
                    Link to the Github repository
                    <div class="tooltip-arrow" data-popper-arrow />
                </div>
            </div>
        </div>
    </nav>
</template>