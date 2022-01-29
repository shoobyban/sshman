<script>
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
    computed: {
        buttonLabel: function () {
            return (this.loading ? 'Loading…' : 'Go');
        }
    },
    methods: {
        run: function () {
            this.reset();
            evtSource = new EventSource(this.url);
            this.loading = true;

            var that = this;

            evtSource.addEventListener('header', function (e) {
                var header = JSON.parse(e.data);
                that.total_items = header.total_items;
                that.actual_msg = header.msg;
            }, false);

            evtSource.addEventListener('item', function (e) {
                var item = JSON.parse(e.data);
                that.items.push(item);
            }, false);

            evtSource.addEventListener('close', function (e) {
                evtSource.close();
                that.loading = false;
            }, false);
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
