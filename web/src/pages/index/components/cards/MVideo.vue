<template>
    <div class="hot card2">
        <div class="hot-item" @click="redirect(item.origin_url)">
            <div class="hot-cover">
                <span class="rank">{{ idx+1 }}</span>
                <img src="../../assets/video.jpg" alt="">
            </div>
            <div class="hot-content">
                <div class="hot-title">
                    <p>{{ item.title }}</p>
                </div>
                <div class="hot-ext">
                    <p class="dt">发布于<span>{{ item.ext.date }}</span></p>
                    <p class="pc"><span>{{ item.ext.score }}</span>次观看</p>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    name: "MVideo",
    props: ["idx", "item"],
    methods: {
        redirect(url) {
            window.open(url)
        },
        toggleFavor(idx) {
            // Do nothing
            return idx
        }
    }
}
</script>