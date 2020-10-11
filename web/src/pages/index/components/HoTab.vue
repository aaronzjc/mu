<template>
<div class="columns switch">
    <div class="column" v-if="state.tabs.length > 0">
        <div class="tabs">
            <ul>
                <li v-for="(tab, idx) in state.tabs" :class="{ 'is-active' : idx == state.selected.tab }" @click="switchTab(idx)" :key="idx"><a>{{ tab.name }}</a></li>
            </ul>
        </div>
        <div class="tags" v-if="state.tabs[state.selected.tab].tags.length > 0">
            <span @click="switchTag(idx)" :class="[ 'tag', { 'is-light-dark' : idx == state.selected.tag } ]" v-for="(tag, idx) in state.tabs[state.selected.tab]['tags']" :key="idx">{{ tag.name }}</span>
        </div>
    </div>
</div>
</template>

<script>
import { onMounted, reactive } from 'vue';
export default {
    name: "HoTab",
    props: ["tabs"],
    setup(props, ctx) {
        const state = reactive({
            selected: {
                tab: 0,
                tag: 0
            },
            tabs: props.tabs
        })

        let sticky = () => {
            const ele = document.getElementsByClassName("switch")[0]
            ele.classList.toggle("sticky", ele.getBoundingClientRect().top === 0)
        }
        let switchTab = (idx) => {
            state.selected = {
                tab: idx,
                tag: 0
            };
            ctx.emit("change", state.selected);
        }
        let switchTag = (idx) => {
            state.selected.tag = idx;
            ctx.emit("change", state.selected);
        }

        onMounted(() => {
            document.addEventListener("scroll", sticky);
        })

        return {
            state,
            switchTab,
            switchTag
        }
    }
}
</script>