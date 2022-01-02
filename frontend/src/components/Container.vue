<script>
import Header from './Header.vue'
import Navbar from './Navbar.vue'

export default {
    components: {
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
        </div>
    </div>
</template>
