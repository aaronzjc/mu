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
    <Opt :mark="item.mark" @toggle="toggleFavor(idx)"></Opt>
</div>
</template>

<script>
import Opt from "./Opt"

export default {
    name: "MRichText",
    props: ["idx", "item"],
    methods: {
        toggleFavor(idx) {
            this.$emit("toggle-favor", idx);
        }
    },
    components: {
        Opt
    }
}
</script>