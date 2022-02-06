<script>
import _ from 'lodash'
import Multiselect from '@vueform/multiselect'

export default {
    name: 'VuexCRUD',
    components: {
        Multiselect,
    },
    props: {
        modelValue: { // v-model, unique identifier for each row
            type: Function,
            default: () => {},
        },
        resourceName: { // e.g. 'Users'
            type: String,
            default: 'Items'
        },
        fields: { // format: [{label: 'Email', index: 'email', placeholder: 'sam@test1.com', type:'email'}] 
            type: Array,
            default: () => [],
        },
        idField: { // identifier field index for update, delete
            type: String,
            default: 'id',
        },
        orderBy: { // default order by
            type: String,
            default: 'id',
        },
        orderDir: { // default order direction
            type: String,
            default: 'asc',
        },
        searchFields: { // fields to search
            type: Array,
            default: () => [],
        },
    },
    emits: ['create', 'update', 'delete', 'fetch'],
    data() {
        return {
            searchInput: '',
            deleteModal: false,
            editModal: false,
            addModal: false,
            currentID: '',
            selection: [],
            sortBy: '',
            sortDir: '',
            listItems: {},
        }
    },
    computed: {
        listFields: function () {
            return this.fields.filter(item => (this.visible(item, 'list')))
        },
        addFields: function () {
            return this.fields.filter(item => (this.visible(item, 'add')))
        },
        editFields: function () {
            return this.fields.filter(item => (this.visible(item, 'edit')))
        },
        current: function () {
            if (this.listItems == undefined) {
                return {}
            }
            return this.listItems[this.currentID]
        },
    },
    watch: {
        modelValue: function () {
            this.recalcSearch()
        },
        searchInput: function () {
            this.recalcSearch()
        },
        sortBy: function () {
            this.recalcSearch()
        },
        sortDir: function () {
            this.recalcSearch()
        },
    },
    mounted: function () {
        this.sortBy = this.orderBy
        this.sortDir = this.orderDir
        if (this.orderDir == undefined || this.orderDir == '') {
            this.sortDir = 'asc'
        }
    },

    methods: {
        visible(item, place) {
            return item.hide == false || item.hide == undefined || item.hide.indexOf(place) == -1
        },
        findIn(item, searchFields) {
            if (searchFields == undefined) {
                searchFields = this.listFields
            }
            let value = ''
            for (var i = 0; i < searchFields.length; i++) {
                var field = searchFields[i]
                if (item[field] == undefined) {
                    continue
                }
                value += JSON.stringify(item[field]).toLowerCase()
            }
            if (value.indexOf(this.searchInput.toLowerCase()) != -1) {
                return true
            }
            return false
        },
        recalcSearch() {
            this.listItems = {}
            for (let key in this.modelValue) {
                const value = this.modelValue[key]
                this.listItems[key] = {
                    ...value,
                    __key: key,
                }
            }
            if (this.searchInput != '') {
                let a = _(this.listItems).filter(item => {
                    return this.findIn(item, this.searchFields)
                }).value()
                this.listItems = {}
                for (let k in a) {
                    const value = a[k]
                    const key = value.__key
                    this.listItems[key] = value
                }
            }
            if (this.sortBy != '') {
                let sorted = {}
                let keys = _.keys(this.listItems)
                keys.sort((x, y) => {
                    if (this.sortDir == "asc") {
                        return this.listItems[x][this.sortBy] > this.listItems[y][this.sortBy] ? 1 : -1
                    } else {
                        return this.listItems[x][this.sortBy] < this.listItems[y][this.sortBy] ? 1 : -1
                    }
                })
                _.forEach(keys, (key) => {
                    return sorted[key] = this.listItems[key]
                })
                this.listItems = sorted
            }
        },
        toggleSelected(idx, e) {
            e.stopPropagation()
            if (this.selection.includes(idx)) {
                this.selection.splice(this.selection.indexOf(idx), 1)
            } else {
                this.selection.push(idx)
            }
        },
        isSelected(idx) {
            return this.selection.includes(idx)
        },
        toggleAll(e) {
            e.stopPropagation()
            for (const key in this.value) {
                if (this.selection.includes(key)) {
                    this.selection.splice(this.selection.indexOf(key), 1)
                } else {
                    this.selection.push(key)
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
        async prepareItem(item, prefix) {
            for (const field of this.addFields) {
                if (field.type == 'multiselect') {
                    const value = this.$refs[prefix + field.index][0].plainValue
                    item[field.index] = value
                } else if (field.type == 'file') {
                    const eTarget = this.$refs[prefix + field.index][0]
                    if (eTarget.files.length == 0) {
                        continue
                    }
                    let result = await this.readFile(eTarget.files[0])
                    item[field.index] = result
                } else {
                    console.log('field', field.index, item[field.index])
                }
            }
            return item
        },
        readFile(file){
            return new Promise((resolve, reject) => {
                let reader = new FileReader()
                reader.onload = function (evt) {
                    if (evt.target.readyState !== 2) return
                    if (evt.target.error) {
                        reject(evt.target.error)
                        return
                    }
                    resolve(evt.target.result)
                }
                reader.readAsText(file)
            })
        },
        async createItem(e) {
            this.addModal = false
            this.currentID = ''
            e.stopPropagation()
            let data = new FormData(e.target.form)
            let item = Object.fromEntries(data.entries())
            let id
            if (this.idField != '.') {
                id = item[this.idField]
            }
            item = await this.prepareItem(item, 'add:')
            this.$emit('create', {
                id,
                item
            })
            setTimeout(() => {
                this.$emit('fetch')
            }, 500)
        },
        async updateItem() {
            this.editModal = false
            let id
            if (this.idField != '.') {
                id = this.current[this.idField]
            } else {
                id = this.currentID
            }
            let item = await this.prepareItem(this.current, 'edit:')
            this.$emit('update', {
                id: id,
                item: item
            })
            this.currentID = ''
            setTimeout(() => {
                this.$emit('fetch')
            }, 500)
        },
        deleteItem() {
            console.log('delete', this.currentID)
            let id
            if (this.idField != '.') {
                id = this.current[this.idField]
            } else {
                id = this.currentID
            }
            this.deleteModal = false
            this.currentID = ''
            this.$emit('delete', id)
            setTimeout(() => {
                this.$emit('fetch')
            }, 500)
        },
    },
}
</script>

<template>
<div>
    <h1 class="text-xl pb-5">
        {{ resourceName }}
    </h1>

    <div class="p-4 bg-white block sm:flex items-center justify-between border-b border-gray-200 lg:mt-1.5">
        <div class="mb-1 w-full">
            <div class="sm:flex">
                <div class="hidden sm:flex items-center sm:divide-x sm:divide-gray-100 mb-3 sm:mb-0">
                    <form class="lg:pr-3" action="#" method="GET">
                        <label for="search-input" class="sr-only">Search</label>
                        <div class="mt-1 relative lg:w-52 xl:w-96">
                            <input id="search-input" v-model="searchInput" type="text" class="bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5" :placeholder="'Search for '+resourceName.toLowerCase()">
                        </div>
                    </form>
                </div>
                <div class="flex items-center space-x-2 sm:space-x-3 ml-auto">
                    <button id="add-items" class="w-1/2 text-white bg-blue-600 hover:bg-blue-700 focus:ring-4 focus:ring-blue-200 font-medium inline-flex items-center justify-center rounded-lg text-sm px-3 py-2 text-center sm:w-auto" @click="addModal = true">
                        <i class="fas fa-plus mr-2" />
                        Add {{ resourceName }}
                    </button>
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
                                        <input id="checkbox-all" aria-describedby="checkbox-1" type="checkbox" class="bg-gray-50 border-gray-300 focus:ring-3 focus:ring-blue-200 h-4 w-4 rounded" @click="toggleAll($event)">
                                        <label for="checkbox-all" class="sr-only">checkbox</label>
                                    </div>
                                </th>
                                <th v-for="(field, index) in listFields" :key="index" scope="col" class="p-4 text-left text-xs font-medium text-gray-500 uppercase select-none">
                                    <div @click="toggleSort(field.index)">
                                        {{ field.label }}
                                        <i v-if="sortBy == field.index && sortDir == 'asc'" class="fas fa-sort-up text-blue-600 ml-2 align-bottom" />
                                        <i v-if="sortBy == field.index && sortDir == 'desc'" class="fas fa-sort-down text-blue-600 ml-2 align-top" />
                                    </div>
                                </th>
                                <th scope="col" class="p-4 text-left text-xs font-medium text-gray-500 uppercase">
                                    Actions
                                </th>
                            </tr>
                        </thead>
                        <tbody id="list-items" class="bg-white divide-y divide-gray-200">
                            <tr v-for="(item,idx) in listItems" :key="idx" :data-rowid="idx" class="hover:bg-gray-100" @click="toggleSelected(idx, $event)">
                                <td class="p-4 w-4">
                                    <div class="flex items-center">
                                        <input :id="'checkbox-'+idx" aria-describedby="checkbox-1" type="checkbox" :checked="isSelected(idx)" class="bg-gray-50 border-gray-300 focus:ring-3 focus:ring-blue-200 h-4 w-4 rounded" @click="$event.stopPropagation()">
                                        <label :for="'checkbox-'+idx" class="sr-only">checkbox</label>
                                    </div>
                                </td>
                                <td v-for="field in listFields" :key="field.index" class="p-4 items-center space-x-6 mr-12 lg:mr-0 max-w-lg">
                                    <div v-if="field.type == 'multiselect'">
                                        <div v-for="(grp, index) in item[field.index]" :key="index" class="px-2 bg-green-600 inline hover:bg-red-700 text-white text-sm font-small rounded-full mb-1 mr-1">
                                            {{ grp }}
                                        </div>
                                    </div>
                                    <div v-else class="text-sm font-normal text-gray-500">
                                        <div class="text-sm font-normal text-gray-500">
                                            {{ item[field.index] }}
                                        </div>
                                    </div>
                                </td>
                                <td class="p-4 whitespace-nowrap space-x-2">
                                    <button class="text-white bg-blue-600 hover:bg-blue-700 focus:ring-4 focus:ring-blue-200 font-medium rounded-lg text-sm inline-flex items-center px-3 py-2 text-center" @click="currentID = idx; editModal = true">
                                        <i class="fas fa-pen" />
                                    </button>
                                    <button class="text-white bg-red-600 hover:bg-red-800 focus:ring-4 focus:ring-red-300 font-medium rounded-lg text-sm inline-flex items-center px-3 py-2 text-center" @click="currentID = idx; deleteModal = true">
                                        <i class="fas fa-trash" />
                                    </button>
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </div>
            </div>
        </div>
    </div>

    <!-- Edit Modal -->
    <div v-show="editModal" id="edit-modal" class="overflow-x-hidden overflow-y-auto fixed top-4 left-0 right-0 md:inset-0 z-50 justify-center items-center h-modal sm:h-full flex">
        <div class="modal-overlay absolute w-full h-full bg-gray-900 opacity-50" />
        <div class="relative w-full max-w-2xl px-4 h-full md:h-auto">
            <!-- Modal content -->
            <div class="bg-white rounded-lg shadow relative">
                <!-- Modal header -->
                <div class="flex items-start justify-between p-5 border-b rounded-t">
                    <h3 class="text-xl font-semibold">
                        Edit {{ resourceName }}
                    </h3>
                    <div class="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm p-1.5 ml-auto inline-flex items-center" data-modal-toggle="resource-modal" @click="editModal=false">
                        <i class="fas fa-times" />
                    </div>
                </div>
                <!-- Modal body -->
                <div class="p-6 space-y-6">
                    <form @submit.prevent="true">
                        <div v-if="current" class="grid grid-cols-6 gap-6">
                            <div v-for="(field,index) in editFields" :key="index" class="col-span-6 sm:col-span-3">
                                <label :for="field.index" class="text-sm font-medium text-gray-900 block mb-2">{{ field.label }}</label>
                                <input v-if="field.type == 'text'" :id="'edit-'+field.index" :ref="'edit:'+field.index" v-model="current[field.index]" type="text" :name="field.index" class="shadow-sm bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5" :placeholder="field.placeholder" :required="field.required?true:false">
                                <input v-else-if="field.type == 'email'" :id="'edit-'+field.index" :ref="'edit:'+field.index" v-model="current[field.index]" :name="field.index" type="email" class="shadow-sm bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5" :placeholder="field.placeholder" :required="field.required?true:false">
                                <input v-else-if="field.type == 'file'" :id="'edit-'+field.index" :ref="'edit:'+field.index" type="file" :name="field.index" class="shadow-sm bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5" :placeholder="field.placeholder" :required="field.required?true:false">
                                <Multiselect v-else-if="field.type == 'multiselect'" :id="'edit-'+field.index" :ref="'edit:'+field.index" v-model="current[field.index]" mode="tags" :create-tag="true" :append-new-tag="true" :searchable="true" :options="field.options" />
                                <Multiselect v-else-if="field.type == 'select'" :id="'edit-'+field.index" :ref="'edit:'+field.index" v-model="current[field.index]" mode="single" :searchable="true" :options="field.options" />
                                <div v-else>
                                    Unhandled {{ field.type }}
                                </div>
                            </div>
                        </div>
                    </form>
                </div>
                <!-- Modal footer -->
                <div class="items-center p-6 border-gray-200 rounded-b">
                    <button id="edit-save" class="text-white bg-blue-600 hover:bg-blue-700 focus:ring-4 focus:ring-blue-200 font-medium rounded-lg text-sm px-5 py-2.5 text-center" @click="updateItem">
                        Save
                    </button>
                </div>
            </div>
        </div>
    </div>

    <!-- Add Modal -->
    <div v-show="addModal" id="add-modal" class="overflow-x-hidden overflow-y-auto fixed top-4 left-0 right-0 md:inset-0 z-50 justify-center items-center h-modal sm:h-full flex">
        <div class="modal-overlay absolute w-full h-full bg-gray-900 opacity-50" />
        <div class="relative w-full max-w-2xl px-4 h-full md:h-auto">
            <!-- Modal content -->
            <div class="bg-white rounded-lg shadow relative">
                <!-- Modal header -->
                <div class="flex items-start justify-between p-5 border-b rounded-t">
                    <h3 class="text-xl font-semibold">
                        Add new {{ resourceName }}
                    </h3>
                    <div class="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm p-1.5 ml-auto inline-flex items-center" data-modal-toggle="add-resource-modal" @click="addModal=false">
                        <i class="fas fa-times" />
                    </div>
                </div>
                <!-- Modal body -->
                <div class="p-6 space-y-6">
                    <form @submit.prevent="true">
                        <div class="grid grid-cols-6 gap-6">
                            <div v-for="(field,index) in addFields" :key="index" class="col-span-6 sm:col-span-3">
                                <label :for="field.index" class="text-sm font-medium text-gray-900 block mb-2">{{ field.label }}</label>
                                <input v-if="field.type == 'text'" :id="'add-'+field.index" type="text" :name="field.index" class="shadow-sm bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5" :placeholder="field.placeholder" :required="field.required?true:false">
                                <input v-else-if="field.type == 'email'" :id="'add-'+field.index" type="email" :name="field.index" class="shadow-sm bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5" :placeholder="field.placeholder" :required="field.required?true:false">
                                <input v-else-if="field.type == 'file'" :id="'add-'+field.index" :ref="'add:'+field.index" type="file" :name="'add:'+field.index" class="shadow-sm bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5" :placeholder="field.placeholder" :required="field.required?true:false">
                                <Multiselect v-else-if="field.type == 'multiselect'" :id="'add-'+field.index" :ref="'add:'+field.index" mode="tags" :create-tag="true" :append-new-tag="true" :searchable="true" :options="field.options" />
                                <Multiselect v-else-if="field.type == 'select'" :id="'add-'+field.index" :ref="'add:'+field.index" mode="single" :searchable="true" :options="field.options" />
                                <div v-else>
                                    Unhandled {{ field.type }}
                                </div>
                            </div>
                        </div>
                        <!-- Modal footer -->
                        <div class="items-center mt-6 border-gray-200 rounded-b">
                            <button id="add-save" class="text-white bg-blue-600 hover:bg-blue-700 focus:ring-4 focus:ring-blue-200 font-medium rounded-lg text-sm px-5 py-2.5 text-center" @click="createItem">
                                Save
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </div>

    <!-- Delete Modal -->
    <div v-show="deleteModal" id="delete-modal" class="overflow-x-hidden overflow-y-auto fixed top-4 left-0 right-0 md:inset-0 z-50 justify-center items-center h-modal sm:h-full flex">
        <div class="modal-overlay absolute w-full h-full bg-gray-900 opacity-50" />
        <div class="relative w-full max-w-md px-4 h-full md:h-auto">
            <!-- Modal content -->
            <div class="bg-white rounded-lg shadow relative">
                <!-- Modal header -->
                <div class="flex justify-end p-2">
                    <div class="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm p-1.5 ml-auto inline-flex items-center" data-modal-toggle="delete-resource-modal" @click="deleteModal=false">
                        <i class="fas fa-times" />
                    </div>
                </div>
                <!-- Modal body -->
                <div class="p-6 pt-0 text-center">
                    <i class="fas fa-trash-alt text-5xl text-red-600" />
                    <h3 class="text-xl font-normal text-gray-500 mt-5 mb-6">
                        Are you sure you want to delete this {{ resourceName.toLowerCase() }}?
                    </h3>
                    <button class="text-white bg-red-600 hover:bg-red-800 focus:ring-4 focus:ring-red-300 font-medium rounded-lg text-base inline-flex items-center px-3 py-2.5 text-center mr-2" @click="deleteItem()">
                        Yes, I'm sure
                    </button>
                    <button class="text-gray-900 bg-white hover:bg-gray-100 focus:ring-4 focus:ring-blue-200 border border-gray-200 font-medium inline-flex items-center rounded-lg text-base px-3 py-2.5 text-center" data-modal-toggle="delete-resource-modal" @click="deleteModal=false">
                        No, cancel
                    </button>
                </div>
            </div>
        </div>
    </div>
</div>
</template>

<style src="@vueform/multiselect/themes/default.css"></style>
