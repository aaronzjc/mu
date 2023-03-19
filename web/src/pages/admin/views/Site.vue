<template>
    <BoxMain :title="['网站管理']">
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
                                <button class="button is-info">新增</button>
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
                                    <th>名称</th>
                                    <th>Cron</th>
                                    <th>类型</th>
                                    <th>节点</th>
                                    <th>状态</th>
                                    <th>操作</th>
                                </tr>
                            </thead>
                            <tbody>
                                <tr v-for="item in state.list" :key="item.id">
                                    <td></td>
                                    <td>{{ item.name }}</td>
                                    <td><span class="tag is-warning" v-if="item.cron !== ''">{{ item.cron }}</span></td>
                                    <td><span class="tag is-light" v-if="item.type != 0">{{ crawType[item.type] }}</span>
                                    </td>
                                    <td>
                                        <template v-if="item.node_option === 1">
                                            <span class="tag is-light">{{ nodeType[item.node_type] }}</span>
                                        </template>
                                        <template v-if="item.node_option === 2">
                                            <template v-if="item.node_hosts.length > 1">
                                                <div class="tags">
                                                    <span class="tag is-info" v-for="nodeId in item.node_hosts"
                                                        :key="nodeId">{{ state.nodes[nodeId]["name"] }}</span>
                                                </div>
                                            </template>
                                            <template v-if="item.node_hosts.length == 1">
                                                <span class="tag is-info" v-for="nodeId in item.node_hosts" :key="nodeId">{{
                                                    state.nodes[nodeId]["name"] }}</span>
                                            </template>
                                        </template>
                                    </td>
                                    <td>
                                        <span class="tag is-light" v-if="!item.enable">未启用</span>
                                        <span class="tag is-success" v-if="item.enable">启用</span>
                                    </td>
                                    <td class="is-actions-cell">
                                        <div class="buttons">
                                            <button class="button is-small is-primary" @click="router.push({name: 'siteEdit', query: {id: item.id}})">
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
import { nodeType, crawType } from '@adm/config';

const API = {
    list: "/admin/sites",
    craw: function (id) {
        return "/admin/sites/" + id + "/craw"
    },
    
}

const state = reactive({
    list: [],
    nodes: []
})

async function fetchList() {
    try {
        let resp = await Get(API.list)
        if (resp.data.code === 10000) {
            state.nodes = resp.data.data.nodeList;
            state.list = resp.data.data.siteList;
        } else {
            console.log(resp.data.msg);
            state.nodes = [];
            state.list = [];
        }
    } catch (err) {
        console.log(err)
    }
}

const router = useRouter()

onMounted(fetchList)
</script>
