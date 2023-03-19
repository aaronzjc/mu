<template>
    <a
        :class="[
            { 'is-active': active },
            { 'has-icon has-dropdown-icon': dropdown }
        ]"
    >
        <BasicIcon v-if="icon" :name="icon"></BasicIcon>
        <span class="menu-item-label">{{ title }}</span>
        <div v-if="dropdown" class="dropdown-icon">
            <BasicIcon :name="active ? mdiMinus : mdiPlus"></BasicIcon>
        </div>
    </a>
</template>

<script setup>
import { mdiPlus, mdiMinus } from '@mdi/js'
import BasicIcon from '@adm/components/BasicIcon.vue'

defineProps({
    title: {
        type: String,
        required: true
    },
    active: {
        type: Boolean,
        required: true
    },
    icon: String,
    dropdown: Boolean
})
</script>
