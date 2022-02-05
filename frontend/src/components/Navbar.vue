<script>
import routes from '../routes.js'

export default {
    name: 'Navbar',
    data() {
        return {
            routes: routes,
            currentRoute: this.$route.path,
        }
    },
    methods: {
        navigate(route) {
            console.log('Navigate to ', route)
 //           this.$router.push(route)
        },
    },
    computed: {
        routesWithLabel() {
            return this.routes.filter(route => route.label)
        },
    },
}
</script>
<template>
    <div class="fixed hidden h-full top-0 left-0 flex lg:flex flex-shrink-0 flex-col w-52 transition-width duration-75">
        <div class="bg-gray-800 shadow-xl h-16 fixed bottom-0 mt-12 md:relative md:h-screen z-10 w-full md:w-48">
            <div class="md:mt-16 md:w-48 md:fixed md:left-0 md:top-0 content-center md:content-start text-left justify-between">
                <ul class="list-reset flex flex-row md:flex-col py-0 md:py-3 px-1 md:px-2 text-center md:text-left">
                    <li v-for="route in routesWithLabel" class="mr-3 flex-1">
                        <router-link class="block py-1 md:py-3 pl-1 align-middle text-gray-400 no-underline hover:text-white border-b-2 border-gray-800" :to="route.path">{{route.label}}</router-link>
                    </li>
                </ul>
            </div>
        </div>
    </div>
</template>