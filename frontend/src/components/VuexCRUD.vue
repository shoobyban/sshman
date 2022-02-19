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
            type: Object,
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
            sortBy: '',
            sortDir: '',
            listItems: {},
            expandAll: false,
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
    created() {
        const onEscape = (e) => {
            if (e.keyCode === 27) {
                this.editModal = false
                this.addModal = false
                this.deleteModal = false
            }
        }
        document.addEventListener('keyup', onEscape)
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
                } else if (field.type == 'select') {
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

    <div class="p-4 appbg block sm:flex items-center justify-between border-b border-gray-200 lg:mt-1.5">
        <div class="mb-1 w-full">
            <div class="sm:flex">
                <div class="hidden sm:flex items-center sm:divide-x sm:divide-gray-100 mb-3 sm:mb-0">
                    <form class="lg:pr-3" action="#" method="GET">
                        <label for="search-input" class="sr-only">Search</label>
                        <div class="mt-1 relative lg:w-52 xl:w-96">
                            <input id="search-input" v-model="searchInput" type="text" class="appbg border border-gray-300 sm:text-sm rounded-lg block w-full p-2.5" :placeholder="'Search for '+resourceName.toLowerCase()">
                        </div>
                    </form>
                </div>
                <div class="flex items-center space-x-2 sm:space-x-3 ml-auto">
                    <slot name="extra-buttons" />
                    <label for="expandall">Expand All</label>
                    <input id="expandall" v-model="expandAll" type="checkbox">
                    <button id="add-items" class="btn" @click="addModal = true">
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
                        <thead class="bg-gray-100 dark:bg-gray-800 dark:text-white cursor-pointer">
                            <tr>
                                <th v-for="(field, index) in listFields" :key="index" scope="col" class="headerlabel select-none">
                                    <div @click="toggleSort(field.index)">
                                        {{ field.label }}
                                        <i v-if="sortBy == field.index && sortDir == 'asc'" class="fas fa-sort-up text-blue-600 ml-2 align-bottom" />
                                        <i v-if="sortBy == field.index && sortDir == 'desc'" class="fas fa-sort-down text-blue-600 ml-2 align-top" />
                                    </div>
                                </th>
                                <th scope="col" class="headerlabel">
                                    Actions
                                </th>
                            </tr>
                        </thead>
                        <tbody id="list-items" class="appbg divide-y divide-gray-200">
                            <tr v-for="(item,idx) in listItems" :key="idx" :data-rowid="idx" class="hover:bg-gray-100 dark:hover:bg-gray-800">
                                <td v-for="field in listFields" :key="field.index" class="p-2">
                                    <div v-if="field.type == 'multiselect'" :class="{'max-h-14': !expandAll}" class="overflow-y-scroll flex flex-wrap">
                                        <div v-for="(grp, index) in item[field.index]" :key="index" :class="{'max-w-[150px]': !expandAll}" class="overflow-hidden pl-2 pr-2 multiselect-tag">
                                            <span>{{ grp }}</span>
                                        </div>
                                    </div>
                                    <div v-else class="text-sm font-normal graytext">
                                        <div class="text-sm font-normal graytext">
                                            {{ item[field.index] }}
                                        </div>
                                    </div>
                                </td>
                                <td class="p-4 whitespace-nowrap space-x-2">
                                    <button aria-label="Edit item" class="edit-item btn" @click="currentID = idx; editModal = true">
                                        <i class="fas fa-pen" />
                                    </button>
                                    <button aria-label="Delete item" class="delete-item delbtn" @click="currentID = idx; deleteModal = true">
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
    <div v-show="editModal" id="edit-modal" class="modal-backdrop">
        <div class="modal-overlay absolute w-full h-full bg-gray-900 opacity-50" />
        <div class="appbg relative w-full max-w-2xl px-4 h-full md:h-auto">
            <!-- Modal content -->
            <div class="appbg rounded-lg light:shadow dark:border-white dark:border-2 relative">
                <!-- Modal header -->
                <div class="flex items-start justify-between p-5 border-b rounded-t">
                    <h3 class="text-xl font-semibold">
                        Edit {{ resourceName }}
                    </h3>
                    <div class="closebtn" data-modal-toggle="resource-modal" @click="editModal=false">
                        <i class="fas fa-times" />
                    </div>
                </div>
                <!-- Modal body -->
                <div class="p-6 space-y-6">
                    <form @submit.prevent="true">
                        <div v-if="current" class="grid grid-cols-6 gap-6">
                            <div v-for="(field,index) in editFields" :key="index" :class="field.double?'col-span-6':'col-span-3'">
                                <label :for="field.index" class="appbg text-sm font-medium block mb-2">{{ field.label }}</label>
                                <input v-if="field.type == 'text'" :id="'edit-'+field.index" :ref="'edit:'+field.index" v-model="current[field.index]" type="text" :name="field.index" class="appbg shadow-sm border border-gray-300 sm:text-sm rounded-lg block w-full p-2.5" :placeholder="field.placeholder" :required="field.required?true:false">
                                <input v-else-if="field.type == 'email'" :id="'edit-'+field.index" :ref="'edit:'+field.index" v-model="current[field.index]" :name="field.index" type="email" class="shadow-sm border border-gray-300 sm:text-sm rounded-lg block w-full p-2.5" :placeholder="field.placeholder" :required="field.required?true:false">
                                <input v-else-if="field.type == 'file'" :id="'edit-'+field.index" :ref="'edit:'+field.index" type="file" :name="field.index" :placeholder="field.placeholder" :required="field.required?true:false">
                                <Multiselect v-else-if="field.type == 'multiselect'" :id="'edit-'+field.index" :ref="'edit:'+field.index" v-model="current[field.index]" mode="tags" :create-tag="true" :append-new-tag="true" :searchable="true" :options="field.options" />
                                <Multiselect v-else-if="field.type == 'select'" :id="'edit-'+field.index" :ref="'edit:'+field.index" v-model="current[field.index]" mode="single" :append-new-option="true" :searchable="true" :options="field.options" />
                                <div v-else>
                                    Unhandled {{ field.type }}
                                </div>
                            </div>
                        </div>
                    </form>
                </div>
                <!-- Modal footer -->
                <div class="items-center p-6 border-gray-200 rounded-b">
                    <button id="edit-save" class="btn" @click="updateItem">
                        Save
                    </button>
                </div>
            </div>
        </div>
    </div>

    <!-- Add Modal -->
    <div v-show="addModal" id="add-modal" class="modal-backdrop">
        <div class="modal-overlay absolute w-full h-full bg-gray-900 opacity-50" />
        <div class="relative w-full max-w-2xl px-4 h-full md:h-auto">
            <!-- Modal content -->
            <div class="appbg rounded-lg light:shadow dark:border-white dark:border-2 relative">
                <!-- Modal header -->
                <div class="flex items-start justify-between p-5 border-b rounded-t">
                    <h3 class="text-xl font-semibold">
                        Add new {{ resourceName }}
                    </h3>
                    <div class="closebtn" data-modal-toggle="add-resource-modal" @click="addModal=false">
                        <i class="fas fa-times" />
                    </div>
                </div>
                <!-- Modal body -->
                <div class="p-6 space-y-6">
                    <form @submit.prevent="true">
                        <div class="grid grid-cols-6 gap-6">
                            <div v-for="(field,index) in addFields" :key="index" class="col-span-6 sm:col-span-3">
                                <label :for="field.index" class="appbg text-sm font-medium block mb-2">{{ field.label }}</label>
                                <input v-if="field.type == 'text'" :id="'add-'+field.index" type="text" :name="field.index" class="shadow-sm appbg border border-gray-300 sm:text-sm rounded-lg block w-full p-2.5" :placeholder="field.placeholder" :required="field.required?true:false">
                                <input v-else-if="field.type == 'email'" :id="'add-'+field.index" type="email" :name="field.index" class="shadow-sm bg-gray-50 border border-gray-300 sm:text-sm rounded-lg block w-full p-2.5" :placeholder="field.placeholder" :required="field.required?true:false">
                                <input v-else-if="field.type == 'file'" :id="'add-'+field.index" :ref="'add:'+field.index" type="file" :name="'add:'+field.index" class="appbg" :placeholder="field.placeholder" :required="field.required?true:false">
                                <Multiselect v-else-if="field.type == 'multiselect'" :id="'add-'+field.index" :ref="'add:'+field.index" class="appbg" mode="tags" :create-tag="true" :append-new-tag="true" :searchable="true" :options="field.options" />
                                <Multiselect v-else-if="field.type == 'select'" :id="'add-'+field.index" :ref="'add:'+field.index" class="appbg" mode="single" :searchable="true" :options="field.options" />
                                <div v-else>
                                    Unhandled {{ field.type }}
                                </div>
                            </div>
                        </div>
                        <!-- Modal footer -->
                        <div class="items-center mt-6 border-gray-200 rounded-b">
                            <button id="add-save" class="btn" @click="createItem">
                                Save
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    </div>

    <!-- Delete Modal -->
    <div v-show="deleteModal" id="delete-modal" class="modal-backdrop">
        <div class="modal-overlay absolute w-full h-full bg-gray-900 opacity-50" />
        <div class="relative w-full max-w-md px-4 h-full md:h-auto">
            <!-- Modal content -->
            <div class="appbg rounded-lg light:shadow dark:border-white dark:border-2 relative">
                <!-- Modal header -->
                <div class="flex justify-end p-2">
                    <div class="closebtn" data-modal-toggle="delete-resource-modal" @click="deleteModal=false">
                        <i class="fas fa-times" />
                    </div>
                </div>
                <!-- Modal body -->
                <div class="p-6 pt-0 text-center">
                    <i class="fas fa-trash-alt text-5xl text-red-600" />
                    <h3 class="text-xl font-normal graytext mt-5 mb-6">
                        Are you sure you want to delete this {{ resourceName.toLowerCase() }}?
                    </h3>
                    <button id="delete-item" class="delbtn" @click="deleteItem()">
                        Yes, I'm sure
                    </button>
                    <button class="cancelbtn" data-modal-toggle="delete-resource-modal" @click="deleteModal=false">
                        No, cancel
                    </button>
                </div>
            </div>
        </div>
    </div>
</div>
</template>

<style src="@vueform/multiselect/themes/default.css"></style>
<style>
.btn {
    @apply text-white bg-blue-600 hover:bg-blue-700  font-medium rounded-lg text-sm inline-flex items-center px-3 py-2 text-center;
}
.closebtn {
    @apply text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm p-1.5 ml-auto inline-flex items-center;
}
.delbtn {
    @apply text-white bg-red-600 hover:bg-red-800  font-medium rounded-lg text-sm inline-flex items-center px-3 py-2 text-center;
}
.cancelbtn {
    @apply text-gray-900 bg-white hover:bg-gray-100  border border-gray-200 font-medium inline-flex items-center rounded-lg text-base px-3 py-2.5 text-center;
}
.appbg {
    @apply bg-white dark:bg-gray-900 dark:text-white;
}
.modal-backdrop {
    @apply overflow-x-hidden overflow-y-auto fixed top-4 left-0 right-0 md:inset-0 z-50 justify-center items-center h-modal sm:h-full flex;
}
.graytext {
    @apply text-gray-500 dark:text-gray-300;
}
.headerlabel {
    @apply p-4 text-left text-xs font-medium graytext uppercase;
}

.multiselect {
    @apply appbg;
}
.multiselect-tags-search {
    @apply appbg;
}
.multiselect-option {
    @apply appbg;
}
.multiselect-option:hover {
    @apply bg-gray-100 dark:bg-gray-800 dark:text-white;
}
::-webkit-scrollbar {
    @apply bg-gray-900 dark:bg-gray-100;
    width: 5px;
    height: 5px;
}
::-webkit-scrollbar-track {
    @apply bg-gray-200 dark:bg-gray-800;
}
::-webkit-scrollbar-thumb {
    background: #888;
}
</style>