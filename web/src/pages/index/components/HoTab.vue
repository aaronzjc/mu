<template>
<div class="columns switch">
    <div class="column" v-if="tabs.length > 0">
        <div class="tabs">
            <ul>
                <li v-for="(tab, idx) in tabs" :class="{ 'is-active' : idx == selected.tab }" @click="switchTab(idx)" :key="idx"><a>{{ tab.name }}</a></li>
            </ul>
        </div>
        <div class="tags" v-if="tabs[selected.tab].tags.length > 0">
            <span @click="switchTag(idx)" :class="[ 'tag', { 'is-light-dark' : idx == selected.tag } ]" v-for="(tag, idx) in tabs[selected.tab]['tags']" :key="idx">{{ tag.name }}</span>
        </div>
    </div>
</div>
</template>

<script>
export default {
    name: "HoTab",
    mounted() {
        // 监听吸顶事件
        document.addEventListener("scroll", this.sticky);
    },
    data() {
        return {
            selected: {
                tab: 0,
                tag: 0
            }
        }
    },
    props: {
        tabs: {
            type: Array,
            default: () => []
        },
    },
    methods: {
        sticky() {
            const ele = document.getElementsByClassName("switch")[0]
            ele.classList.toggle("sticky", ele.getBoundingClientRect().top === 0)
        },
        switchTab(idx) {
            this.selected = {
                tab: idx,
                tag: 0
            };
            this.$emit("change", this.selected);
        },
        switchTag(idx) {
            this.selected.tag = idx;
            this.$emit("change", this.selected);
        }
    }
}
</script>