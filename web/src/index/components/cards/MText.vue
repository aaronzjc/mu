<template>
<div class="hot-item card0">
    <a :href="item.origin_url" :title="item.title" target="_blank">{{ item.title }}</a>
</div>
</template>

<script>
export default {
    name: "MText",
    props: ["item"]
}
</script>