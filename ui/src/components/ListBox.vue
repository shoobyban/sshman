<template>
    <div>
      <v-list-item>
        <v-list-item-content>
          <v-list-item-title class="text-h6">
            {{title}} <v-btn
      class="mx-2"
      fab
      small
      dark
      color="indigo"
      @click="newItem"
    >
      <v-icon dark>
        mdi-plus
      </v-icon>
    </v-btn> <v-select
      v-model="groupfilter"
      :items="groups"
      label="Groups"
      clearable
    ></v-select>

          </v-list-item-title>
          <v-list-item-subtitle>
            <v-text-field v-model="search" placeholder="Search..."></v-text-field>
          </v-list-item-subtitle>
        </v-list-item-content>
      </v-list-item>

      <v-divider></v-divider>

      <v-list class="overflow-y-auto list">
        <v-list-item v-for="(item,k) in filteredlist" :key="k" link @click="selectHandler(item)">
          <v-list-item-content>
            <v-list-item-title>{{ item[label] }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>

    </div>
</template>

<script>
import _ from 'lodash'

export default {
    props: {
        type: { type: String },
        list: {  },
        groups: {  },
        title: { type: String },
        label: { type: String },
        ts: { },
    },
    data() { 
        return {
            search: '',
            groupfilter: null,
        }
    },
    watch: { 
      ts: function(){
        console.log('list changed')
        this.$forceUpdate()
      },
    },
    methods: {
        newItem() {
          this.$root.$emit('select-item', {
                item: null,
                type: this.type,
          })
        },
        selectHandler(item) {
            this.$root.$emit('select-item', {
                item: item,
                type: this.type,
            })
        },
    },
    computed: {
      filteredlist: function() {
        if (this.search == '' && this.groupfilter == null) {
          return this.list
        }
        if (this.groupfilter == null) {
            return _.filter(this.list, currentItem => currentItem[this.label].includes(this.search) )
        } else if (this.search == '') {
            return _.filter(this.list, currentItem => {
                if (currentItem.groups != null) {
                    console.log('groups', currentItem.groups, this.groupfilter)
                    return currentItem.groups.includes(this.groupfilter)
                } else {
                    return false
                }
            })
        }
        return _.filter(this.list, currentItem => currentItem.groups.includes(this.groupfilter) && currentItem[this.label].includes(this.search) )
      },
    },
}
</script>

<style scoped lang="css">

.list {
    height: calc(100vh - 200px);
}

</style>