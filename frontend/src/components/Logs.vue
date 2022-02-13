<script>
var evtSource = false

export default {
    name: 'Logs',
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
            loading: false
        }
    },
    computed: {
        buttonLabel: function () {
            return (this.loading ? 'Loading…' : 'Go')
        }
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
            this.items = []
            this.total_items = -1
        }
    }
}
</script>

<template>
    <div class="overflow-y-auto h-full">
        <div v-for="(item, i) in items" :key="i" class="mt-1 whitespace-nowrap" :class="item.type">
            {{ item.message }}
        </div>
    </div>
</template>

<style scoped>
.error {
    background-color: #f8d7da !important;
    border-left: 3px solid rgb(220,38,38) !important;
    padding-left: 5px;
    width: 100%;
}
.info {
    background-color: #d4edfa !important;
    border-left: 3px solid #BADA55 !important;
    padding-left: 5px;
    width: 100%;
}
</style>
