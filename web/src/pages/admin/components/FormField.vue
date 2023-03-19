<template>
    <div :class="['field', { 'is-horizontal': horizontal }]">
        <div v-if="label" :class="['field-label', normal ? 'is-normal' : '']">
            <label class="label">{{ label }}</label>
        </div>
        <div class="field-body">
            <div
                :class="[
                    'field',
                    group ? 'is-grouped' : '',
                    center ? 'is-justify-content-center' : ''
                ]"
            >
                <slot />
            </div>
        </div>
    </div>
</template>

<script setup>
const props = defineProps({
    label: String,
    horizontal: {
        type: Boolean,
        default: false
    },
    group: {
        type: Boolean,
        default: false
    },
    center: {
        type: Boolean,
        default: false
    },
    normal: {
        type: Boolean,
        default: false
    }
})
</script>
