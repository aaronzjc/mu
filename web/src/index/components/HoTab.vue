<template>
<div class="columns">
    <div class="column" v-if="tabs.length > 0">
        <div class="columns">
            <div class="column">
                <div class="tabs">
                    <ul>
                        <li v-for="(tab, idx) in tabs" :class="{ 'is-active' : idx == selected.tab }" @click="switchTab(idx)" :key="idx"><a>{{ tab.name }}</a></li>
                    </ul>
                </div>
            </div>
        </div>
        <div class="columns" v-if="tabs[selected.tab].tags.length > 0">
            <div class="column tab-tag tags">
                <span @click="switchTag(idx)" :class="[ 'tag', { 'is-primary' : idx == selected.tag } ]" v-for="(tag, idx) in tabs[selected.tab]['tags']" :key="idx">{{ tag.name }}</span>
            </div>
        </div>
    </div>
</div>
</template>

<script>
export default {
    name: "HoTab",
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

<style scoped>
.tab-tag {
    padding-top: 0px;
    padding-bottom: 0px;
    cursor: pointer;
}
</style>