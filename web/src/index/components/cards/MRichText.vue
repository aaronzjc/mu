<template>
<div class="hot-item card1">
    <div class="hot-title">
        <a :href="item.origin_url" :title="item.title" target="_blank">{{ item.title }}</a>
    </div>
    <div class="hot-desc">
        <p class="has-text-grey">{{ item.desc }}</p>
    </div>
</div>
</template>

<script>
export default {
    name: "MRichText",
    props: ["item"],
}
</script>