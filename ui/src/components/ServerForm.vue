<template>
  <v-form
    ref="form"
    lazy-validation
  >
    <h2>Host {{selected.alias}}</h2>

    <v-text-field
      v-model="alias"
      :counter="10"
      label="Alias"
      required
    ></v-text-field>

    <v-text-field
      v-model="host"
      :counter="10"
      label="Server Address"
      placeholder="host.com:22"
      required
    ></v-text-field>
    <v-text-field
      v-model="user"
      label="User"
      required
    ></v-text-field>

    <v-combobox
      clearable hide-selected
      multiple small-chips
      v-model="servergroups"
      :items="groups"
      color="#388E3C"
      label="Groups"
    ></v-combobox>

    <v-select
      multiple chips small-chips
      :items="selected.users"
      label="Users"
      placeholder="Users are added by matching user and host groups"
      readonly
    ></v-select>

    <v-btn
      color="error"
      class="mr-4"
      @click="reset"
    >
      Cancel
    </v-btn>

    <v-btn
      color="success"
      @click="submit"
      :disabled="notChanged"
    >
      Submit
    </v-btn>
  </v-form>
</template>

<script>
import axios from 'axios'

export default {
    data() {
      return {
        alias: '',
        user: '',
        host: '',
        servergroups: [],
      }
    },
    computed: {
      notChanged: function() {
        return (this.selected.alias == this.alias || this.alias == '') &&
          (this.selected.user == this.user || this.user == '') &&
          (this.selected.host == this.host || this.host == '') &&
          (this.selected.groups && !(this.selected.groups.length !== this.servergroups.length) &&
          this.servergroups.every((v, i) => v === this.selected.groups[i]))
      }
    },
    mounted() {
      this.reloadForm()
    },
    props: {
        type: { type: String },
        selected: { type: Object },
        groups: {  },
        ts: {  },
    },
    watch: { 
      selected: function() {
        this.reloadForm()
      },
      ts: function() {
        this.reloadForm()
      },
    },
    methods: {
      reloadForm() {
        this.alias = this.selected.alias
        this.servergroups = this.selected.groups
        this.host = this.selected.host
        this.user = this.selected.user
      },
      submit() { 
        if (this.selected.alias != this.alias && this.alias != '') {
          axios.post('/api/rename/server', { new: this.alias, old: this.selected.alias }).then(res => {
            console.log(res)
            this.$root.$emit('reload-config', {})
            this.selected.alias = this.alias
          })
        }
        if (this.selected.groups != this.servergroups && this.alias != '') {
          // setting server groups
          axios.post('/api/groups/server',{ groups: this.servergroups, aliases: [this.selected.alias, this.alias] }).then(res => {
            console.log(res)
            this.$root.$emit('reload-config', {})
            this.selected.groups = this.servergroups
          })
        }
        },
      reset() { this.$root.$emit('reload-form', { type: this.type }) },
    },
}
</script>