<script setup>
import axios from 'axios'

var data = {
    users: {},
    hosts: {},
}
axios.get("/api/").then((resp) => {
    data.users = resp.data.users
    data.hosts = resp.data.hosts
})
</script>

<template>
<div class="grid grid-cols-2 bg-gray-50 font-sans">
  <div class="p-6 max-w-sm mx-auto bg-white rounded-xl shadow-md flex space-x-4 align-top">
    <div class="grid grid-cols-2 align-top">
      <div class="h-12 border align-middle text-center bg-gray-50 shadow ml-3 mr-3 mb-2" v-for="user in data.users">
        <div class="m-4">{{user.email}}</div>
      </div>
    </div>
  </div>
  <div class="p-6 max-w-sm mx-auto bg-white rounded-xl shadow-md flex space-x-4 align-top">
    <div class="grid grid-cols-2 align-middle">
      <div class="h-12 border text-center bg-gray-50 shadow ml-3 mr-3 mb-2" v-for="host, alias in data.hosts">
        <div class="m-4">{{alias}}</div>
      </div>
    </div>
  </div>
 </div>
</template>
