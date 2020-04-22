<template>
<div class="hot card0">
    <div class="hot-item">
        <a :href="item.origin_url" :title="item.title" target="_blank">{{ item.title }}</a>
    </div>
    <div class="divider"></div>
    <Opt :mark="item.mark" @toggle="toggleFavor(idx)"></Opt>
</div>
</template>

<script>
import Opt from "./Opt"

export default {
    name: "MText",
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