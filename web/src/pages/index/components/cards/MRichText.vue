<template>
<div class="hot card1">
    <div class="hot-item">
        <div class="hot-title">
            <a :href="item.origin_url" :title="item.title" target="_blank" rel="noopener">{{ item.title }}</a>
        </div>
        <div class="hot-desc">
            <p class="has-text-grey">{{ item.desc }}</p>
        </div>
    </div>
    <div class="divider"></div>
    <Opt :item="item" :idx="idx"></Opt>
</div>
</template>

<script setup>
import Opt from "./Opt.vue"
const props = defineProps(['idx', 'item'])
</script>