<template>
    <BoxMain :title="['总览']">
        <div class="columns">
            <div class="column is-one-third">
                <div class="card">
                    <header class="card-header">
                        <p class="card-header-title">在线用户列表</p>
                    </header>
                    <div class="card-content">
                        <ul>
                            <li class="" v-for="item in state.onlineList">{{ item }}</li>
                        </ul>
                    </div>
                </div>
            </div>
        </div>
    </BoxMain>
</template>

<script setup>
import BoxMain from '@adm/components/BoxMain.vue'
import { onMounted, reactive } from 'vue';
import { Get } from "@/lib/http";

const API = {
    online: "/admin/stats/online/list"
}

const state = reactive({
    onlineList: []
})

async function fetchOnlineList() {
    try {
        let resp = await Get(API.online)
        if (resp.data.code === 10000) {
            state.onlineList = resp.data.data.onlineList;
        }
    } catch (err) {
        console.log(err)
    }
}
setInterval(fetchOnlineList, 30 * 1000)
onMounted(fetchOnlineList)
</script>
