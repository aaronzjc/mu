<template>
<div class="hot card0">
    <div class="hot-item">
        <a :href="item.origin_url" :title="item.title" target="_blank" rel="noopener">{{ item.title }}</a>
    </div>
    <div class="divider"></div>
    <Opt :item="item" :idx="idx"></Opt>
</div>
</template>

<script setup>
import Opt from "./Opt.vue"
const props = defineProps(['idx', 'item'])
</script>