<script>
import axios from 'axios'
import _ from 'lodash'

export default {
    name: 'CRUD',
    props: {resourceName: String, endpoint: String, fields: Array, orderBy: String, orderDir: String},
    data() {
        return {
            allData: [],
            searchInput: '',
            deleteModal: false,
            editModal: false,
            addModal: false,
            dataIndex: null,
            selected: [],
            sortBy: '',
            sortDir: '',
        }
    },
    mounted: function() {
        this.sortBy = this.orderBy
        this.sortDir = this.orderDir
        if (this.orderDir == undefined || this.orderDir == '') {
            this.sortDir = 'asc'
        }
        axios.get(this.endpoint).then(response => {
                this.allData = response.data
            })
            .catch(error => {
                console.log(error)
            })
    },
    computed: {
        searchResult: function() {
            var data;
            if (this.searchInput != '') {
                data = this.allData.filter(item => {
                    var itemGroups = JSON.stringify(item.groups)
                    return (item.email.toLowerCase().indexOf(this.searchInput.toLowerCase()) !== -1) || (item.name.toLowerCase().indexOf(this.searchInput.toLowerCase()) !== -1) || (itemGroups.toLowerCase().indexOf(this.searchInput.toLowerCase()) !== -1)
                })
            } else {
                data = this.allData
            }
            if (this.sortBy != '') {
                return _.orderBy(data, this.sortBy, this.sortDir)
            }
            return data
        },
        listFields: function() {
            return this.fields.filter(item => {
                return item.hidefromlist == false || item.hidefromlist == undefined
            })
        },
        selectedItem: function() {
            return this.searchResult[this.dataIndex]
        },
    },
    methods: {
        toggleSelected(str, e) {
            e.stopPropagation()
            if (this.selected.includes(str)) {
                this.selected.splice(this.selected.indexOf(str), 1)
            } else {
                this.selected.push(str)
            }
        },
        isSelected(str) {
            return this.selected.includes(str)
        },
        toggleAll(e) {
            e.stopPropagation()
            for (const key in this.allData) {
                if (this.selected.includes(key)) {
                    this.selected.splice(this.selected.indexOf(key), 1)
                } else {
                    this.selected.push(key)
                }
            }
        },
        toggleSort(fieldname) {
            if (this.sortBy == fieldname) {
                if (this.sortDir == 'asc' || this.sortDir == '' || this.sortDir == undefined) {
                    this.sortDir = 'desc'
                } else {
                    this.sortDir = 'asc'
                }
            } else {
                this.sortBy = fieldname
                this.sortDir = 'asc'
            }
        },
    },
}
</script>

<template>
    <div>
        <h2 class="font-bold pl-3 pt-3">{{resourceName}}</h2>

        <div class="p-4 bg-white block sm:flex items-center justify-between border-b border-gray-200 lg:mt-1.5">
            <div class="mb-1 w-full">
                <div class="sm:flex">
                    <div class="hidden sm:flex items-center sm:divide-x sm:divide-gray-100 mb-3 sm:mb-0">
                        <form class="lg:pr-3" action="#" method="GET">
                        <label for="allData-search" class="sr-only">Search</label>
                        <div class="mt-1 relative lg:w-52 xl:w-96">
                            <input type="text" v-model="searchInput" id="allData-search" class="bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5" :placeholder="'Search for '+resourceName.toLowerCase()">
                        </div>
                        </form>
                    </div>
                    <div class="flex items-center space-x-2 sm:space-x-3 ml-auto">
                        <div @click="this.addModal = true" class="w-1/2 text-white bg-blue-600 hover:bg-blue-700 focus:ring-4 focus:ring-blue-200 font-medium inline-flex items-center justify-center rounded-lg text-sm px-3 py-2 text-center sm:w-auto">
                            <i class="fas fa-plus mr-2"></i>
                            Add {{resourceName}}
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="flex flex-col">
            <div class="overflow-x-auto">
                <div class="align-middle inline-block min-w-full">
                    <div class="shadow overflow-hidden">
                        <table class="table-fixed min-w-full divide-y divide-gray-200">
                            <thead class="bg-gray-100">
                                <tr>
                                    <th scope="col" class="p-4">
                                        <div class="flex items-center">
                                            <input id="checkbox-all" aria-describedby="checkbox-1" type="checkbox" @click="toggleAll($event)"
                                                class="bg-gray-50 border-gray-300 focus:ring-3 focus:ring-blue-200 h-4 w-4 rounded">
                                            <label for="checkbox-all" class="sr-only">checkbox</label>
                                        </div>
                                    </th>
                                    <th v-for="field in listFields" scope="col" class="p-4 text-left text-xs font-medium text-gray-500 uppercase select-none">
                                        <div @click="toggleSort(field.index)">
                                            {{field.label}}
                                            <i v-if="sortBy == field.index && sortDir == 'asc'" class="fas fa-sort-up text-blue-600 ml-2 align-bottom"></i>
                                            <i v-if="sortBy == field.index && sortDir == 'desc'" class="fas fa-sort-down text-blue-600 ml-2 align-top"></i>
                                        </div>
                                    </th>
                                    <th scope="col" class="p-4 text-left text-xs font-medium text-gray-500 uppercase">
                                        Actions
                                    </th>
                                </tr>
                            </thead>
                            <tbody class="bg-white divide-y divide-gray-200">
                                <tr v-for="(item,idx) in searchResult" :key="idx" class="hover:bg-gray-100" @click="toggleSelected(idx, $event)">
                                    <td class="p-4 w-4">
                                        <div class="flex items-center">
                                            <input :id="'checkbox-'+idx" aria-describedby="checkbox-1" type="checkbox" :checked="isSelected(idx)" @click="$event.stopPropagation()"
                                                class="bg-gray-50 border-gray-300 focus:ring-3 focus:ring-blue-200 h-4 w-4 rounded">
                                            <label :for="'checkbox-'+idx" class="sr-only">checkbox</label>
                                        </div>
                                    </td>
                                    <td v-for="field in listFields" class="p-4 items-center space-x-6 mr-12 lg:mr-0 max-w-lg">
                                        <div v-if="field.type == 'multiselect'">
                                            <button v-for="(grp, index) in item[field.index]" class="px-2 bg-green-600 hover:bg-red-700 text-white text-sm font-small rounded-full mb-1 mr-1">
                                            {{ grp }}
                                            </button>
                                        </div>
                                        <div v-else class="text-sm font-normal text-gray-500">
                                            <div class="text-sm font-normal text-gray-500">{{ item[field.index] }}</div>
                                        </div>
                                    </td>
                                    <td class="p-4 whitespace-nowrap space-x-2">
                                        <div @click="this.dataIndex = idx; this.editModal = true;" class="text-white bg-blue-600 hover:bg-blue-700 focus:ring-4 focus:ring-blue-200 font-medium rounded-lg text-sm inline-flex items-center px-3 py-2 text-center">
                                            <i class="fas fa-pen"></i>
                                        </div>
                                        <div @click="this.dataIndex = idx; this.deleteModal = true;" class="text-white bg-red-600 hover:bg-red-800 focus:ring-4 focus:ring-red-300 font-medium rounded-lg text-sm inline-flex items-center px-3 py-2 text-center">
                                            <i class="fas fa-trash"></i>
                                        </div>
                                    </td>
                                </tr>                        
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>

        <!-- Edit Modal -->
        <div v-show="editModal" class="overflow-x-hidden overflow-y-auto fixed top-4 left-0 right-0 md:inset-0 z-50 justify-center items-center h-modal sm:h-full flex">
            <div class="modal-overlay absolute w-full h-full bg-gray-900 opacity-50"></div>
            <div class="relative w-full max-w-2xl px-4 h-full md:h-auto">
                <!-- Modal content -->
                <div class="bg-white rounded-lg shadow relative">
                    <!-- Modal header -->
                    <div class="flex items-start justify-between p-5 border-b rounded-t">
                        <h3 class="text-xl font-semibold">
                            Edit {{resourceName}}
                        </h3>
                        <div @click="this.editModal=false" class="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm p-1.5 ml-auto inline-flex items-center" data-modal-toggle="resource-modal">
                            <i class="fas fa-times"></i>  
                        </div>
                    </div>
                    <!-- Modal body -->
                    <div class="p-6 space-y-6">
                        <form action="#">
                            <div class="grid grid-cols-6 gap-6" v-if="selectedItem">
                                <div v-for="field in fields" class="col-span-6 sm:col-span-3">
                                    <label for="first-name" class="text-sm font-medium text-gray-900 block mb-2">{{field.label}}</label>
                                    <input type="text" name="first-name" id="first-name" v-model="selectedItem[field.index]" class="shadow-sm bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5" :placeholder="field.placeholder" required>
                                </div>
                            </div>
                        </form>
                    </div>
                    <!-- Modal footer -->
                    <div class="items-center p-6 border-t border-gray-200 rounded-b">
                        <button class="text-white bg-blue-600 hover:bg-blue-700 focus:ring-4 focus:ring-blue-200 font-medium rounded-lg text-sm px-5 py-2.5 text-center" type="submit">Save all</button>
                    </div>
                </div>
            </div>
        </div>

        <!-- Add Modal -->
        <div v-show="addModal" class="overflow-x-hidden overflow-y-auto fixed top-4 left-0 right-0 md:inset-0 z-50 justify-center items-center h-modal sm:h-full flex" id="add-resource-modal">
            <div class="modal-overlay absolute w-full h-full bg-gray-900 opacity-50"></div>
            <div class="relative w-full max-w-2xl px-4 h-full md:h-auto">
                <!-- Modal content -->
                <div class="bg-white rounded-lg shadow relative">
                    <!-- Modal header -->
                    <div class="flex items-start justify-between p-5 border-b rounded-t">
                        <h3 class="text-xl font-semibold">
                            Add new {{resourceName}}
                        </h3>
                        <div @click="this.addModal=false" class="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm p-1.5 ml-auto inline-flex items-center" data-modal-toggle="add-resource-modal">
                            <i class="fas fa-times"></i>
                        </div>
                    </div>
                    <!-- Modal body -->
                    <div class="p-6 space-y-6">
                        <form action="#">
                            <div class="grid grid-cols-6 gap-6">
                                <div v-for="field in fields" class="col-span-6 sm:col-span-3">
                                    <label for="first-name" class="text-sm font-medium text-gray-900 block mb-2">{{field.label}}</label>
                                    <input type="text" :name="field.index" :id="field.index" class="shadow-sm bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5" :placeholder="field.placeholder" required>
                                </div>
                            </div> 
                    </form>
                        </div>
                        <!-- Modal footer -->
                        <div class="items-center p-6 border-t border-gray-200 rounded-b">
                            <button class="text-white bg-blue-600 hover:bg-blue-700 focus:ring-4 focus:ring-blue-200 font-medium rounded-lg text-sm px-5 py-2.5 text-center" type="submit">Add {{resourceName}}</button>
                        </div>
                </div>
            </div>
        </div>

        <!-- Delete Modal -->
        <div v-show="deleteModal" class="overflow-x-hidden overflow-y-auto fixed top-4 left-0 right-0 md:inset-0 z-50 justify-center items-center h-modal sm:h-full flex" id="delete-resource-modal">
            <div class="modal-overlay absolute w-full h-full bg-gray-900 opacity-50"></div>
            <div class="relative w-full max-w-md px-4 h-full md:h-auto">
                <!-- Modal content -->
                <div class="bg-white rounded-lg shadow relative">
                    <!-- Modal header -->
                    <div class="flex justify-end p-2">
                        <div @click="this.deleteModal=false" class="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm p-1.5 ml-auto inline-flex items-center" data-modal-toggle="delete-resource-modal">
                            <i class="fas fa-times"></i>
                        </div>
                    </div>
                    <!-- Modal body -->
                    <div class="p-6 pt-0 text-center">
                        <i class="fas fa-trash-alt text-5xl text-red-600"></i>
                        <h3 class="text-xl font-normal text-gray-500 mt-5 mb-6">Are you sure you want to delete this {{resourceName.toLowerCase()}}?</h3>
                        <a href="#" class="text-white bg-red-600 hover:bg-red-800 focus:ring-4 focus:ring-red-300 font-medium rounded-lg text-base inline-flex items-center px-3 py-2.5 text-center mr-2">
                            Yes, I'm sure
                        </a>
                        <a href="#" class="text-gray-900 bg-white hover:bg-gray-100 focus:ring-4 focus:ring-blue-200 border border-gray-200 font-medium inline-flex items-center rounded-lg text-base px-3 py-2.5 text-center" data-modal-toggle="delete-resource-modal">
                            No, cancel
                        </a>
                    </div>
                </div>
            </div>
        </div>

    </div>
</template>