<template>
<div class="hot-item">
    <div class="hot-title">
        <a :href="item.origin_url" :title="item.title" target="_blank">{{ item.title }}</a>
    </div>
    <div class="hot-desc">
        <p>{{ item.desc }}</p>
    </div>
</div>
</template>

<script>
export default {
    name: "MRichText",
    props: ["item"]
}
</script>

<style lang="scss" scoped>
.hot-item {
    width: 98%;
    margin-right: 2px;
    display: flex;
    align-items: flex-start;
    word-break: break-word;
    flex-direction: column;

    .hot-desc {
        padding: 2px 0;
        p {
            color: #586069;
            font-size: 0.8rem;
        }
    }
}
</style>