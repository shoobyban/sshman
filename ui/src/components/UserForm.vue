<template>
  <v-form
    ref="form"
    lazy-validation
  >
    <h2>User {{selected.email}}</h2>
    <v-text-field
      v-model="email"
      :counter="50"
      label="Email"
      required
    ></v-text-field>

    <v-textarea
      v-model="keyfile"
      :label="this.selected.email == ''?'Paste ssh public key (.pub) file here':'Key'"
      class="key"
      :readonly="selected.email != ''"
    ></v-textarea>

    <v-combobox
      clearable hide-selected
      multiple small-chips
      v-model="selected.groups"
      :items="groups"
      color="#388E3C"
      label="Groups"
    ></v-combobox>
    
    <v-btn
      color="error"
      class="mr-4"
      @click="reset"
    >
      Cancel
    </v-btn>

    <v-btn
      color="success"
      :disabled="notChanged"
      @click="submit"
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
        email: '',
        usergroups: [],
      }
    },
    computed: {
      notChanged: function() {
        return  this.selected.email == this.email ||
          this.email == '' || 
          this.selected.key == '' || 
          (this.selected.groups != undefined && this.selected.groups.length !== this.usergroups.length) ||
          !this.usergroups.every((v, i) => v === this.selected.groups[i])
      },
      keyfile: {
        get: function() {
        return (this.selected.type + ' ' +
        this.selected.key + ' ' +
        this.selected.name).trim(' ')
        },
        set: function(val) {
          const parts = val.split(' ')
          if (parts.length == 3) {
            this.selected.type = parts[0]
            this.selected.key = parts[1]
            this.selected.name = parts[2]
            if (this.selected.email == '') {
              this.email = parts[2]
            }
          }
        }
      },
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
        this.email = this.selected.email
        this.usergroups = this.selected.groups
      },
      submit() {
        if (this.selected.email == '' && this.email != '') {
          // new user
          axios.post('/api/user',{ email: this.email, key: this.key, groups: this.usergroups }).then(res => {
            console.log(res)
            this.$root.$emit('reload-config', {})
            this.selected.email = this.email
          })
          return
        }
        if (this.selected.email != this.email && this.email != '') {
          // renaming user
          axios.post('/api/rename/user',{ new: this.email, old: this.selected.email }).then(res => {
            console.log(res)
            this.$root.$emit('reload-config', {})
            this.selected.email = this.email
          })
        }
        if (this.selected.groups != this.usergroups && this.email != '') {
          // setting user groups
          axios.post('/api/groups/user',{ groups: this.usergroups, emails: [this.selected.email, this.email] }).then(res => {
            console.log(res)
            this.$root.$emit('reload-config', {})
            this.selected.groups = this.usergroups
          })
        }
      },
      reset() { this.$root.$emit('reload-form', { type: this.type }) },
    },
}
</script>

<style scoped lang="css">
.key {
  font-size: 12px;
}
</style>