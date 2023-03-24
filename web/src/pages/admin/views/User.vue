<template>
    <BoxMain :title="['用户管理']">
        <div class="card has-table">
            <header class="card-header">
                <div class="card-header-action level">
                    <div class="level-left">
                        <div class="field is-grouped">
                            <p class="control">
                                <input class="input" type="text" placeholder="输入关键词..." />
                            </p>
                            <p class="control">
                                <a class="button is-primary">搜索</a>
                            </p>
                        </div>
                    </div>
                </div>
            </header>
            <div class="card-content">
                <div class="b-table has-pagination">
                    <div class="table-wrapper has-mobile-cards">
                        <table class="table is-fullwidth is-striped is-hoverable is-fullwidth">
                            <thead>
                                <tr>
                                    <th></th>
                                    <th>ID</th>
                                    <th>用户名</th>
                                    <th>昵称</th>
                                    <th>认证方式</th>
                                    <th>认证时间</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="item in state.list" :key="item.id">
                                    <td class="is-image-cell">
                                        <div class="image">
                                            <img class="is-rounded" :src="item.avatar">
                                        </div>
                                    </td>
                                    <td>{{ item.id }}</td>
                                    <td>{{ item.username }}</td>
                                    <td>{{ item.nickname }}</td>
                                    <td><span class="tag is-light">{{ item.auth_type }}</span></td>
                                    <td><span class="tag is-light">{{ item.auth_time }}</span></td>
                                </tr>
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>
    </BoxMain>
</template>

<script setup>
import BasicIcon from '@adm/components/BasicIcon.vue'
import BoxMain from '@adm/components/BoxMain.vue'
import { mdiAccount, mdiDownload, mdiPencil, mdiTable } from '@mdi/js'

import { Get } from "@/lib/http";
import { reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router';
import { nodeType, crawType } from '@adm/config';

const API = {
    list: "/admin/users"
}

const state = reactive({
    list: []
})

async function fetchList() {
    try {
        let resp = await Get(API.list)
        if (resp.data.code === 10001) {
            console.log(resp.data.msg);
        } else {
            state.list = resp.data.data;
        }
    } catch (err) {
        console.log(err)
    }
}

onMounted(fetchList)
</script>
