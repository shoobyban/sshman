<script>
import Header from './Header.vue'
import Navbar from './Navbar.vue'
import Logs from './Logs.vue'

export default {
    components: {
        Logs,
        Header,
        Navbar,
    },
    data() {
        return {
            title: 'Container'
        }
    }
}
</script>

<template>
    <div class="flex flex-col h-screen">
        <Header />
        <div class="flex flex-grow overflow-hidden bg-white pt-16">
            <Navbar />
            <div class="w-full bg-gray-50 relative overflow-y-auto lg:ml-40 p-5">
                <router-view />
            </div>
            <div class="bg-gray-100 max-h-full text-blue-100 md:w-80 pa-2 md:relative">
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
