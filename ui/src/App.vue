<template>
  <v-app>
    <Nav />
    <v-content>
      <router-view></router-view>
    </v-content>
  </v-app>
</template>

<script>
import Nav from './components/Nav'

export default {
  components: { Nav },
  name: 'App',
};
</script>
