<template>
  <nav>
    <v-app-bar flat app dark fixed>
      <v-icon class="grey--text"
      >🔒</v-icon>
      <v-toolbar-title class="text-uppercase grey--text">
        <span>SSH</span>
        <span class="font-weight-light">Man</span>
      </v-toolbar-title>
      <v-spacer></v-spacer>

      <v-avatar
  color="grey"
  @click="logout"
  size="44">Exit</v-avatar>
    </v-app-bar>

  </nav>
</template>

<script>
export default {
  data() {
    return {
      links: [
        { icon: 'dashboard', text: 'Dashboard', route: '/' },
        { icon: 'folder', text: 'My Projects', route: '/projects' },
        { icon: 'person', text: 'Team', route: '/team' },
      ]
    }
  },
  methods: {
    exit() {

    },
    logout() {
        
    },
  },
};
</script>

<style></style>
