<template>
    <div>
  <v-container class="grey lighten-5" fluid>
    <v-row>
    <v-col class="col-md-2">
    <v-navigation-drawer permanent>
      <ListBox type="user" :list="users" label="email" title="Users" :groups="groups" :ts="ts" />
    </v-navigation-drawer>
    </v-col>
    <v-col class="col-md-8" >
      <MainForm :type="selectedType" :selected="selected" :groups="groups" :ts="ts" />
      <v-card color="#f2fde4" height="30%">
        <v-container fluid class="text-xs-center justify-center ma-2 pa-0 fill-height" color="white"
        @drop.prevent='onDrop($event)' @dragover.prevent @dragenter.prevent>
          <v-icon x-large color="green">mdi-tray-arrow-down</v-icon>
          <div>Drop ssh public keys here (with extension .pub) for mass-upload</div>
      </v-container>
      </v-card>
    </v-col>
    <v-col class="col-md-2">
    <v-navigation-drawer permanent>
      <ListBox type="host" :list="hosts" label="alias" title="Hosts" :groups="groups" :ts="ts" />
    </v-navigation-drawer>
    </v-col>
    </v-row>
  </v-container>
    </div>
</template>

<script>
import axios from 'axios'
import ListBox from '@/components/ListBox.vue'
import MainForm from '@/components/MainForm.vue'

export default {
    components: {
      ListBox,
      MainForm,
    },
    data() { return {
        users: {},
        hosts: {},
        groups: [],
        search: '',
        selectedType: '',
        selected: {},
        files: [],
        ts: null,
    }},
    mounted() {
        this.$root.$on('select-item', this.selectHandler)
        this.$root.$on('reload-config', this.reloadConfig)
        this.$root.$on('reload-form', this.reloadForm)
        this.reloadConfig()
    },
    methods: {
      selectHandler (selected) {
        if  (selected.item == null) {
          if (selected.type == 'user') {
            selected.item = { email: '', key: '', name: '', type: '' }
          } else {
            selected.item = {}
          }
        }
        this.selectedType = selected.type
        this.selected = selected.item
      },
      reloadForm(v) {
        console.log('reload',v.type)
        this.ts = Date.now()
      },
      reloadConfig () {
        axios.get('/api/config').then(res => {
            this.users = res.data.users
            this.hosts = res.data.hosts
        })
        axios.get('/api/groups').then(res => {
            this.groups = Object.keys(res.data)
        })
      },
      onDrop(e) {
        var self = this
        const crypto = require('crypto')
        let droppedFiles = e.dataTransfer.files
        if(!droppedFiles) return
        ([...droppedFiles]).forEach(f => {
          if (f.name.endsWith('.pub')) {
            console.log('processing', f.name)
            var reader = new FileReader();
            reader.onload = function() {
              var parts = reader.result.trim('\n').split(' ')
              const sha = crypto.createHash('sha1')
              sha.update(parts[1])
              const hash = sha.digest('base64')
              if (self.users[hash] == undefined) {
                console.log('new', f.name, hash, parts[2])
                self.users[hash] = {
                  email: f.name,
                  type: parts[0],
                  key: parts[1],
                  name: parts[2],
                }
                self.ts = Date.now()
              } else {
                console.log('existing', f.name)
              }
            }
            reader.onerror=function () {
              self.files.push('error')
            }
            reader.readAsText(f)
          } else {
            console.log('ignoring', f.name)
          }
        })
        e.stopPropagation()
        return false
      },
    },
}
</script>
