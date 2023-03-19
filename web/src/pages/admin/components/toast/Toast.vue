<template>
    <Transition enter-active-class="fadeInDown" leave-active-class="fadeOut">
        <div :class="['toast', level]" v-show="active" ref="toastEl">
            <p>{{ text }}</p>
        </div>
    </Transition>
</template>

<script setup>
import { ref, onMounted } from 'vue'
const props = defineProps({
    text: {
        type: String,
        required: true
    },
    level: {
        type: String,
        required: false
    }
})
const active = ref(false)
onMounted(() => {
    active.value = true
    setTimeout(() => (active.value = false), 3000)
})
</script>
