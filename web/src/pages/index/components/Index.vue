<template>
<div class="content-box">
    <p class="hot-ts" v-if="t !== '' ">更新时间: {{ t }}</p>
    <div class="columns hot-container">
        <div class="column hot-list">

        </div>
    </div>

    <Footer></Footer>
</div>
</template>

<script>
export default {
    name: "Content",
    mounted() {
    },
    data() {
        return {
            tabs: [
                {"name":"Moo-Yuu","key":"新闻","tags":[{"key":"昨天","name":"昨天","enable":1},{"key":"今天","name":"今天","enable":1},{"key":"明天","name":"明天","enable":1}]},
            ],
            selected: {
                tab: 0,
                tag: 0,
            },
            list: [],
            t: "还没更新呢",
        }
    },
    methods: {},
    components: {}
}
</script>