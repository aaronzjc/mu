<template>
    <BoxMain :title="['总览']">
        <Tile :items="items"></Tile>
    </BoxMain>
</template>

<script setup>
import BoxMain from '@adm/components/BoxMain.vue'
import { mdiAccount } from '@mdi/js'
import { reactive } from 'vue'
import Tile from '@adm/components/Tile.vue'
const items = reactive([
    {
        title: '当前在线人数',
        sub: 800,
        icon: mdiAccount,
        color: 'primary'
    },
    {
        title: '当前在线人数',
        sub: 800,
        icon: mdiAccount,
        color: 'info'
    },
    {
        title: '当前在线人数',
        sub: 800,
        icon: mdiAccount,
        color: 'warning'
    }
])
</script>
