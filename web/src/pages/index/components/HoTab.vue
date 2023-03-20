<template>
    <div class="columns switch">
        <div class="column" v-if="tabs.length > 0">
            <div class="tabs">
                <ul>
                    <li v-for="(tab, idx) in tabs" :class="{ 'is-active': idx == state.selected.tab }"
                        @click="switchTab(idx)" :key="idx">
                        <a>{{ tab.name }}</a>
                    </li>
                </ul>
            </div>
            <div class="tags" v-if="tabs[state.selected.tab].tags.length > 0"> 
                <span @click="switchTag(idx)" :class="['tag', { 'is-light-dark': idx == state.selected.tag }]"
                    v-for="(tag, idx) in tabs[state.selected.tab]['tags']" :key="idx">{{ tag.name }}</span>
            </div>
        </div>
    </div>
</template>

<script setup>
import { onMounted, reactive } from "vue";
const emit = defineEmits(["change"]);
const props = defineProps({
    tabs: Object,
});
const state = reactive({
    selected: {
        tab: 0,
        tag: 0,
    },
});

const switchTab = (idx) => {
    state.selected.tab = idx
    state.selected.tag = 0
    emit("change", state.selected);
};
const switchTag = (idx) => {
    state.selected.tag = idx;
    emit("change", state.selected);
};

onMounted(() => {
    document.addEventListener("scroll", () => {
        const ele = document.getElementsByClassName("switch")[0];
        ele.classList.toggle("sticky", ele.getBoundingClientRect().top === 0);
    });
});
</script>
