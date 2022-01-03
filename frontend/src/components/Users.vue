<script>
import axios from 'axios';

export default {
    data() {
        return {
            title: 'Container',
            users: [],
            searchEmail: '',
            deleteModal: false,
            editModal: false,
            addModal: false,
            userIndex: null,
        }
    },
    mounted: function() {
        axios.get('/api/users')
            .then(response => {
                this.users = response.data;
            })
            .catch(error => {
                console.log(error);
            });
    },
    computed: {
        searchByEmail: function() {
            if (this.searchEmail != '') {
                return this.users.filter(user => {
                    var userGroups = JSON.stringify(user.groups);
                    return (user.email.toLowerCase().indexOf(this.searchEmail.toLowerCase()) !== -1) || (user.name.toLowerCase().indexOf(this.searchEmail.toLowerCase()) !== -1) || (userGroups.toLowerCase().indexOf(this.searchEmail.toLowerCase()) !== -1);
                });
            } else {
                return this.users;
            }
        }
    }
}
</script>

<template>
    <div>
        <h2 class="font-bold pl-3 pt-3">Users</h2>

        <div class="p-4 bg-white block sm:flex items-center justify-between border-b border-gray-200 lg:mt-1.5">
            <div class="mb-1 w-full">
                <div class="sm:flex">
                    <div class="hidden sm:flex items-center sm:divide-x sm:divide-gray-100 mb-3 sm:mb-0">
                        <form class="lg:pr-3" action="#" method="GET">
                        <label for="users-search" class="sr-only">Search</label>
                        <div class="mt-1 relative lg:w-52 xl:w-96">
                            <input type="text" v-model="searchEmail" id="users-search" class="bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5" placeholder="Search for users">
                        </div>
                        </form>
                    </div>
                    <div class="flex items-center space-x-2 sm:space-x-3 ml-auto">
                        <div @click="this.addModal = true" class="w-1/2 text-white bg-blue-600 hover:bg-blue-700 focus:ring-4 focus:ring-blue-200 font-medium inline-flex items-center justify-center rounded-lg text-sm px-3 py-2 text-center sm:w-auto">
                            <i class="fas fa-plus mr-2"></i>
                            Add user
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
                                            <input id="checkbox-all" aria-describedby="checkbox-1" type="checkbox"
                                                class="bg-gray-50 border-gray-300 focus:ring-3 focus:ring-blue-200 h-4 w-4 rounded">
                                            <label for="checkbox-all" class="sr-only">checkbox</label>
                                        </div>
                                    </th>
                                    <th scope="col" class="p-4 text-left text-xs font-medium text-gray-500 uppercase">
                                        Email
                                    </th>
                                    <th scope="col" class="p-4 text-left text-xs font-medium text-gray-500 uppercase">
                                        Groups
                                    </th>
                                    <th scope="col" class="p-4">
                                    </th>
                                </tr>
                            </thead>
                            <tbody class="bg-white divide-y divide-gray-200">
                                <tr v-for="(user,idx) in searchByEmail" :key="idx" class="hover:bg-gray-100">
                                    <td class="p-4 w-4">
                                        <div class="flex items-center">
                                            <input id="checkbox-{{ idx }}" aria-describedby="checkbox-1" type="checkbox"
                                                class="bg-gray-50 border-gray-300 focus:ring-3 focus:ring-blue-200 h-4 w-4 rounded">
                                            <label for="checkbox-{{ idx }}" class="sr-only">checkbox</label>
                                        </div>
                                    </td>
                                    <td class="p-4 flex items-center whitespace-nowrap space-x-6 mr-12 lg:mr-0">
                                        <div class="text-sm font-normal text-gray-500">
                                            <div class="text-base font-semibold text-gray-900">{{ user.name }}</div>
                                            <div class="text-sm font-normal text-gray-500">{{ user.email }}</div>
                                        </div>
                                    </td>
                                    <td class="p-4 whitespace-nowrap text-base font-normal text-gray-900 space-x-1">
                                        <button v-for="(grp, index) in user.groups" class="px-2 py-1 bg-green-600 hover:bg-red-700 text-white text-sm font-small rounded-full">
                                        {{ grp }}
                                        </button>
                                    </td>
                                    <td class="p-4 whitespace-nowrap space-x-2">
                                        <div @click="this.userIndex = idx; this.editModal = true;" class="text-white bg-blue-600 hover:bg-blue-700 focus:ring-4 focus:ring-blue-200 font-medium rounded-lg text-sm inline-flex items-center px-3 py-2 text-center">
                                            <i class="fas fa-pen"></i>
                                        </div>
                                        <div @click="this.userIndex = idx; this.deleteModal = true;" class="text-white bg-red-600 hover:bg-red-800 focus:ring-4 focus:ring-red-300 font-medium rounded-lg text-sm inline-flex items-center px-3 py-2 text-center">
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

        <!-- Edit User Modal -->
        <div v-show="editModal" class="overflow-x-hidden overflow-y-auto fixed top-4 left-0 right-0 md:inset-0 z-50 justify-center items-center h-modal sm:h-full flex">
            <div class="modal-overlay absolute w-full h-full bg-gray-900 opacity-50"></div>
            <div class="relative w-full max-w-2xl px-4 h-full md:h-auto">
                <!-- Modal content -->
                <div class="bg-white rounded-lg shadow relative">
                    <!-- Modal header -->
                    <div class="flex items-start justify-between p-5 border-b rounded-t">
                        <h3 class="text-xl font-semibold">
                            Edit user
                        </h3>
                        <div @click="this.editModal=false" class="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm p-1.5 ml-auto inline-flex items-center" data-modal-toggle="user-modal">
                            <i class="fas fa-times"></i>  
                        </div>
                    </div>
                    <!-- Modal body -->
                    <div class="p-6 space-y-6">
                        <form action="#">
                            <div class="grid grid-cols-6 gap-6">
                                <div class="col-span-6 sm:col-span-3">
                                    <label for="first-name" class="text-sm font-medium text-gray-900 block mb-2">First Name</label>
                                    <input type="text" name="first-name" id="first-name" class="shadow-sm bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5" placeholder="Bonnie" required>
                                </div>
                                <div class="col-span-6 sm:col-span-3">
                                    <label for="last-name" class="text-sm font-medium text-gray-900 block mb-2">Last Name</label>
                                    <input type="text" name="last-name" id="last-name" class="shadow-sm bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5" placeholder="Green" required>
                                </div>
                                <div class="col-span-6 sm:col-span-3">
                                    <label for="email" class="text-sm font-medium text-gray-900 block mb-2">Email</label>
                                    <input type="email" name="email" id="email" class="shadow-sm bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5" placeholder="example@company.com" required>
                                </div>
                                <div class="col-span-6 sm:col-span-3">
                                    <label for="phone-number" class="text-sm font-medium text-gray-900 block mb-2">Phone Number</label>
                                    <input type="number" name="phone-number" id="phone-number" class="shadow-sm bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5" placeholder="e.g. +(12)3456 789" required>
                                </div>
                                <div class="col-span-6 sm:col-span-3">
                                    <label for="department" class="text-sm font-medium text-gray-900 block mb-2">Department</label>
                                    <input type="text" name="department" id="department" class="shadow-sm bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5" placeholder="Development" required>
                                </div>
                                <div class="col-span-6 sm:col-span-3">
                                    <label for="company" class="text-sm font-medium text-gray-900 block mb-2">Company</label>
                                    <input type="number" name="company" id="company" class="shadow-sm bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5" placeholder="123456" required>
                                </div>
                                <div class="col-span-6 sm:col-span-3">
                                    <label for="current-password" class="text-sm font-medium text-gray-900 block mb-2">Current Password</label>
                                    <input type="password" name="current-password" id="current-password" class="shadow-sm bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5" placeholder="••••••••" required>
                                </div>
                                <div class="col-span-6 sm:col-span-3">
                                    <label for="new-password" class="text-sm font-medium text-gray-900 block mb-2">New Password</label>
                                    <input type="password" name="new-password" id="new-password" class="shadow-sm bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5" placeholder="••••••••" required>
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

        <!-- Add User Modal -->
        <div v-show="addModal" class="overflow-x-hidden overflow-y-auto fixed top-4 left-0 right-0 md:inset-0 z-50 justify-center items-center h-modal sm:h-full flex" id="add-user-modal">
            <div class="modal-overlay absolute w-full h-full bg-gray-900 opacity-50"></div>
            <div class="relative w-full max-w-2xl px-4 h-full md:h-auto">
                <!-- Modal content -->
                <div class="bg-white rounded-lg shadow relative">
                    <!-- Modal header -->
                    <div class="flex items-start justify-between p-5 border-b rounded-t">
                        <h3 class="text-xl font-semibold">
                            Add new user
                        </h3>
                        <div @click="this.addModal=false" class="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm p-1.5 ml-auto inline-flex items-center" data-modal-toggle="add-user-modal">
                            <i class="fas fa-times"></i>
                        </div>
                    </div>
                    <!-- Modal body -->
                    <div class="p-6 space-y-6">
                        <form action="#">
                            <div class="grid grid-cols-6 gap-6">
                                <div class="col-span-6 sm:col-span-3">
                                    <label for="first-name" class="text-sm font-medium text-gray-900 block mb-2">First Name</label>
                                    <input type="text" name="first-name" id="first-name" class="shadow-sm bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5" placeholder="Bonnie" required>
                                </div>
                                <div class="col-span-6 sm:col-span-3">
                                    <label for="last-name" class="text-sm font-medium text-gray-900 block mb-2">Last Name</label>
                                    <input type="text" name="last-name" id="last-name" class="shadow-sm bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5" placeholder="Green" required>
                                </div>
                                <div class="col-span-6 sm:col-span-3">
                                    <label for="email" class="text-sm font-medium text-gray-900 block mb-2">Email</label>
                                    <input type="email" name="email" id="email" class="shadow-sm bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5" placeholder="example@company.com" required>
                                </div>
                                <div class="col-span-6 sm:col-span-3">
                                    <label for="phone-number" class="text-sm font-medium text-gray-900 block mb-2">Phone Number</label>
                                    <input type="number" name="phone-number" id="phone-number" class="shadow-sm bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5" placeholder="e.g. +(12)3456 789" required>
                                </div>
                                <div class="col-span-6 sm:col-span-3">
                                    <label for="department" class="text-sm font-medium text-gray-900 block mb-2">Department</label>
                                    <input type="text" name="department" id="department" class="shadow-sm bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5" placeholder="Development" required>
                                </div>
                                <div class="col-span-6 sm:col-span-3">
                                    <label for="company" class="text-sm font-medium text-gray-900 block mb-2">Company</label>
                                    <input type="number" name="company" id="company" class="shadow-sm bg-gray-50 border border-gray-300 text-gray-900 sm:text-sm rounded-lg focus:ring-blue-600 focus:border-blue-600 block w-full p-2.5" placeholder="123456" required>
                                </div>
                            </div> 
                    </form>
                        </div>
                        <!-- Modal footer -->
                        <div class="items-center p-6 border-t border-gray-200 rounded-b">
                            <button class="text-white bg-blue-600 hover:bg-blue-700 focus:ring-4 focus:ring-blue-200 font-medium rounded-lg text-sm px-5 py-2.5 text-center" type="submit">Add user</button>
                        </div>
                </div>
            </div>
        </div>

        <!-- Delete User Modal -->
        <div v-show="deleteModal" class="overflow-x-hidden overflow-y-auto fixed top-4 left-0 right-0 md:inset-0 z-50 justify-center items-center h-modal sm:h-full flex" id="delete-user-modal">
            <div class="modal-overlay absolute w-full h-full bg-gray-900 opacity-50"></div>
            <div class="relative w-full max-w-md px-4 h-full md:h-auto">
                <!-- Modal content -->
                <div class="bg-white rounded-lg shadow relative">
                    <!-- Modal header -->
                    <div class="flex justify-end p-2">
                        <div @click="this.deleteModal=false" class="text-gray-400 bg-transparent hover:bg-gray-200 hover:text-gray-900 rounded-lg text-sm p-1.5 ml-auto inline-flex items-center" data-modal-toggle="delete-user-modal">
                            <i class="fas fa-times"></i>
                        </div>
                    </div>
                    <!-- Modal body -->
                    <div class="p-6 pt-0 text-center">
                        <i class="fas fa-trash-alt text-5xl text-red-600"></i>
                        <h3 class="text-xl font-normal text-gray-500 mt-5 mb-6">Are you sure you want to delete this user?</h3>
                        <a href="#" class="text-white bg-red-600 hover:bg-red-800 focus:ring-4 focus:ring-red-300 font-medium rounded-lg text-base inline-flex items-center px-3 py-2.5 text-center mr-2">
                            Yes, I'm sure
                        </a>
                        <a href="#" class="text-gray-900 bg-white hover:bg-gray-100 focus:ring-4 focus:ring-blue-200 border border-gray-200 font-medium inline-flex items-center rounded-lg text-base px-3 py-2.5 text-center" data-modal-toggle="delete-user-modal">
                            No, cancel
                        </a>
                    </div>
                </div>
            </div>
        </div>

    </div>
</template>