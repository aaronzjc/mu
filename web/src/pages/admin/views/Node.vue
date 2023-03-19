<template>
    <BoxMain :title="['节点管理']">
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
                            <p class="control">
                                <button class="button is-info" @click="router.push({name: 'nodeEdit'}, {id: 0})">新增</button>
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
                                    <th>名称</th>
                                    <th>Addr</th>
                                    <th>类型</th>
                                    <th>Ping</th>
                                    <th>状态</th>
                                    <th>操作</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="item in state.list" :key="item.id">
                                    <td></td>
                                    <td>{{ item.id }}</td>
                                    <td>{{ item.name }}</td>
                                    <td>{{ item.addr }}</td>
                                    <td><span class="tag is-warning">{{ nodeType[item.type] }}</span></td>
                                    <td>
                                        <span style="width: 15px;height: 15px;padding:0;" :class="['ping tag', item.ping == 1 ? 'is-success' : 'is-danger' ]"></span>
                                    </td>
                                    <td>
                                        <span class="tag is-light" v-if="!item.enable">未开启</span>
                                        <span class="tag is-success" v-if="item.enable">开启</span>
                                    </td>
                                    <td class="is-actions-cell">
                                        <div class="buttons">
                                            <button class="button is-small is-primary"
                                                @click="router.push({ name: 'nodeEdit', query: { id: item.id } })">
                                                <BasicIcon :name="mdiPencil"></BasicIcon>
                                            </button>
                                            <button class="button is-small is-warning">
                                                <BasicIcon :name="mdiDownload"></BasicIcon>
                                            </button>
                                        </div>
                                    </td>
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
import { mdiDownload, mdiPencil, mdiTable } from '@mdi/js'

import { Get } from "@/lib/http";
import { reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router';
import { nodeType } from '@adm/config';

const API = {
    list: "/admin/nodes",
    upsert: function (id) {
        return "/admin/nodes/" + id + "/upsert"
    },
    del: function(id) {
        return "/admin/nodes/" + id + "/del"
    }
}

const state = reactive({
    list: []
})

async function fetchList() {
    try {
        let resp = await Get(API.list)
        if (resp.data.code === 10001) {
            console.log(resp.data.msg)
        } else {
            state.list = resp.data.data;
        }
    } catch(err) {
        // eslint-disable-next-line
        console.log(err)
    }
}

const router = useRouter()

onMounted(fetchList)
</script>