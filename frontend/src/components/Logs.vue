<template>
    <div class="overflow-y-auto h-full">
        <div v-for="item in items" class="mt-1 whitespace-nowrap" :class="item.type">
            {{item.message}}
        </div>
    </div>
</template>
<script>
var evtSource = false;

export default {
    name: 'Logs',
    props: {
        url: String,
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
    mounted() {
        this.run();
    },
    updated() {
        var elem = this.$el
        elem.scrollTop = elem.clientHeight;
    },
    computed: {
        buttonLabel: function () {
            return (this.loading ? 'Loading…' : 'Go');
        }
    },
    methods: {
        run: function () {
            this.reset();
            evtSource = new EventSource(this.url);
            evtSource.onmessage = (event) => {
                console.log('message', event);
                this.items.push(JSON.parse(event.data));
            }
            evtSource.onerror = (event) => {
                console.log('error', event);
                this.reset();
            }
            evtSource.onclose = (event) => {
                console.log('close', event);
                this.reset();
            }
        },
        reset: function () {
            if (evtSource !== false) {
                evtSource.close();
            }

            this.loading = false;
            this.items = [];
            this.total_items = -1;
        }
    }
}
</script>

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
