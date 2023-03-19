<template>
    <div class="hot card2">
        <div class="hot-item" @click="window.open(item.origin_url)">
            <div class="hot-cover">
                <span class="rank">{{ idx+1 }}</span>
                <img src="@idx/assets/video.jpg" alt="">
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

<script setup>
const props = defineProps(['idx', 'item'])
</script>