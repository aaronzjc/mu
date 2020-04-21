<template>
<div class="hot-item">
    <a :href="item.origin_url" :title="item.title" target="_blank">{{ item.title }}</a>
</div>
</template>

<script>
export default {
    name: "MText",
    props: ["item"]
}
</script>

<style lang="scss" scoped>
.hot-item {
    width: 98%;
    margin-right: 2px;
    display: flex;
    align-items: center;
    word-break: break-word;
}
</style>