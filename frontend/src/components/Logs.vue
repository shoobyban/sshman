<script>
import VirtualList from 'vue3-virtual-scroll-list'
import Item from './LogItem.vue'
var evtSource = false

export default {
    name: 'Logs',
    components: {
        VirtualList,
    },
    props: {
        url: {
            type: String,
            default: '/logs',
        },
    },
    data() {
        return {
            msg: 'Hello world',
            actual_msg: '',
            total_items: -1,
            items: [],
            loading: false,
            itemComponent: Item,
        }
    },
    computed: {
        start: () => 0,
    },
    mounted() {
        this.run()
    },
    updated() {
        var elem = this.$el
        elem.scrollTop = elem.clientHeight
    },
    methods: {
        run: function () {
            this.reset()
            evtSource = new EventSource(this.url)
            evtSource.onmessage = (event) => {
                console.log('log message', event)
                this.items.push(JSON.parse(event.data))
                this.$refs.virturalList.scrollToBottom()
            }
            evtSource.onerror = (event) => {
                console.log('log error', event)
                this.reset()
            }
            evtSource.onclose = (event) => {
                console.log('log close', event)
                this.reset()
            }
        },
        reset: function () {
            if (evtSource !== false) {
                evtSource.close()
            }
            this.loading = false
        }
    }
}
</script>

<template>
    <virtual-list
    ref="virturalList"
    class="overflow-scroll h-screen"
    v-bind="virtualAttrs"
    :bench="0"
    :start="start"
    :data-key="'id'"
    :data-sources="items"
    :data-component="itemComponent"
    >
        <template #header>
            <h2 class="text-black dark:text-white text-center mb-2 mt-2 font-bold">
                    <span>Console</span>
            </h2>
        </template>
        <template #footer>
            <div style="height: 70px;" />
        </template>
    </virtual-list>
</template>
