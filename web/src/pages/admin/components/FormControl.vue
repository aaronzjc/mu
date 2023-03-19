<template>
    <div :class="controlCls">
        <slot></slot>
        <BasicIcon v-if="icon" :name="icon" size="16"></BasicIcon>
        <p :class="['help is-danger']" v-if="help">
            {{ help }}
        </p>
    </div>
</template>

<script setup>
import { computed } from 'vue'
import BasicIcon from '@adm/components/BasicIcon.vue'
const props = defineProps({
    expanded: Boolean,
    help: String, // 提示文字
    icon: String // 输入框图标
})

const controlCls = computed(() => {
    let cls = 'control '
    cls += props.expanded ? 'is-expanded ' : ''
    cls += props.icon ? 'has-icons-left ' : ''
    return cls
})
</script>
