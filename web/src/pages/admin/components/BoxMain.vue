<template>
    <div class="section-box">
        <section class="section is-title-bar">
            <div class="level">
                <div class="level-left">
                    <ul>
                        <li v-for="t in sectionTitle">{{ t }}</li>
                    </ul>
                </div>
            </div>
        </section>
        <section class="section is-main-section">
            <slot />
        </section>
    </div>
</template>

<script setup>
import { computed } from 'vue';
import { useRoute } from 'vue-router'

const route = useRoute()
const props = defineProps({
    title: {
        type: Array,
        default: []
    }
})
const sectionTitle = computed(() => {
    if (props.title.length > 0) {
        return props.title
    }
    return route.meta.title ?? []
})
</script>
