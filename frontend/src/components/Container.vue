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
    <div>
        <Header></Header>
        <div class="flex overflow-hidden bg-white pt-16">
            <Navbar></Navbar>
            <div class="h-full w-full bg-gray-50 relative overflow-y-auto lg:ml-52">
            <router-view></router-view>
            </div>
            <div class="bg-gray-100 text-blue-100 md:w-80 pa-2 md:relative">
                <h2 class="text-center text-black mb-2 mt-2 font-bold"><span>Server Logs</span></h2>
                <div class="bg-white w-full h-full text-black">
                    <Logs url="http://localhost/api/logs" />
                </div>
            </div>
        </div>
    </div>
</template>
