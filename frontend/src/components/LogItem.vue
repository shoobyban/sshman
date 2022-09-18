<template>
  <div class="resize-x mt-1 pa-1" :class="source.type">
    {{ source.message }}
  </div>
</template>

<script>
export default {
  name: "Item",
  props: {
    source: {
      type: Object,
      default() {
        return {}
      },
    },
  },
}
</script>

<style scoped>
.error {
    border-left: 3px solid rgb(220,38,38) !important;
    padding-left: 3px;
    margin-left : 2px;
    width: 100%;
    font-size: 10px;
}
.info {
    border-left: 3px solid #BADA55 !important;
    padding-left: 3px;
    margin-left : 2px;
    width: 100%;
    font-size: 10px;
}
</style>
