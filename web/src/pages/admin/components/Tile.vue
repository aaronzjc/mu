<template>
    <div class="tile is-ancestor">
        <div class="tile is-parent is-4" v-for="item in items">
            <div class="card tile is-child">
                <div class="card-content">
                    <div class="level is-mobile">
                        <div class="level-item">
                            <div class="is-widget-label">
                                <h3
                                    class="subtitle is-spaced"
                                    v-if="item.title"
                                >
                                    {{ item.title }}
                                </h3>
                                <h1 class="title" v-if="item.sub">
                                    {{ item.sub }}
                                </h1>
                            </div>
                        </div>
                        <div class="level-item has-widget-icon">
                            <div class="is-widget-icon">
                                <BasicIcon
                                    :name="item.icon"
                                    :size="48"
                                    :color="item.color"
                                ></BasicIcon>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
<script setup>
import BasicIcon from '@adm/components/BasicIcon.vue'
defineProps({
    items: Array
})
</script>
