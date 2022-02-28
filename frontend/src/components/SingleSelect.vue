<script>
import Multiselect from '@vueform/multiselect'
import { nextTick } from 'vue'

export default {
    name: 'SingleSelect',
    components: {
        Multiselect,
    },
    props: {
        modelValue: { // v-model, unique identifier for each row
            type: String,
            required: true,
        },
        field: { // field data
            type: Object,
            required: true,
        },
    },
    emits: ['update','fetch','input','update:modelValue'],
    data() {
        return {
            current: [],
            old: '',
        }
    },
    mounted: function () {
        if (this !== undefined ) {
            console.log('mounted', this.modelValue)
            this.current = [ this.modelValue ]
        }
    },
    methods: {
        async changeField() {
            if (this.current.length > 0) {
                this.old = this.current[0]
            }
        },
        async replaceSingleItem(val) {
            const el = this.$refs.vms
            el.iv.value=[]
            await nextTick()
            if (this.old == '') {
                this.$emit('input', val[0])
                this.$emit('update', val[0])
                this.$emit('update:modelValue', val[0])
                return
            }
            const newVal = val.filter(x => x != this.old)
            this.current = [newVal[0]]
            await nextTick()
            this.$emit('input', newVal[0])
            this.$emit('update', newVal[0])
            this.$emit('update:modelValue', newVal[0])
        },
    },
}
</script>

<template>
    <Multiselect 
        ref="vms"
        v-model="current"
        class="singletag"
        mode="tags" 
        :create-tag="true" 
        :searchable="true" 
        :native-support="true"
        :options="field.options"
        @open="changeField"
        @input="replaceSingleItem" 
        @update="replaceSingleItem" 
    />
</template>

<style src="@vueform/multiselect/themes/default.css"></style>

<style>

</style>