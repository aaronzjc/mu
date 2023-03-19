<template>
    <BoxMain :title="title">
        <div class="card">
            <header class="card-header">
                <p class="card-header-title">
                    <BasicIcon :name="mdiBallot"></BasicIcon>
                    网站信息
                </p>
            </header>
            <div class="card-content">
                <form action="jav" method="get">
                    <FormField label="地址" horizontal>
                        <FormControl>
                            <p>{{ state.info.root }}</p>
                        </FormControl>
                    </FormField>
                    <FormField label="标识" horizontal>
                        <FormControl>
                            <p>{{ state.info.key }}</p>
                        </FormControl>
                    </FormField>
                    <FormField label="名称" group horizontal>
                        <FormControl>
                            <input class="input" type="text" v-model="state.info.name" placeholder="名称" />
                        </FormControl>
                    </FormField>
                    <FormField label="描述" horizontal>
                        <FormControl expanded>
                            <input class="input" type="text" v-model="state.info.desc" placeholder="描述" />
                        </FormControl>
                    </FormField>
                    <FormField label="Cron" horizontal>
                        <FormControl expanded>
                            <input class="input" type="text" v-model="state.info.cron" placeholder="Cron" />
                        </FormControl>
                    </FormField>

                    <FormField label="标签" horizontal>
                        <FormControl>
                            <div class="tag-field">
                                <div class="tag-item" v-for="(tag,idx) in state.info.tags">
                                    <div class="tag-control">
                                        <input class="input is-small" disabled type="text" v-model="state.info.tags[idx].key" placeholder="标识">
                                    </div>
                                    <div class="tag-control">
                                        <input class="input is-small" type="text" v-model="state.info.tags[idx].name" placeholder="名称">
                                    </div>
                                    <div class="tag-control">
                                        <button type="button" :class="['button is-small', tag.enable == 1 ? 'is-danger' : 'is-primary' ]">
                                        {{ tag.enable == 1 ? '关闭' : '启用' }}
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </FormControl>
                    </FormField>

                    <FormField label="请求头" horizontal>
                        <FormControl>
                            <div class="tag-field">
                                <div class="tag-item">
                                    <div class="tag-control">
                                        <button type="button" class="button is-small is-primary" @click="addHeader">新增</button>
                                    </div>
                                </div>
                                <div class="tag-item" v-for="(header, idx) in state.info.req_headers">
                                    <div class="tag-control">
                                        <input class="input is-small" type="text" v-model="state.info.req_headers[idx].key" placeholder="名称">
                                    </div>
                                    <div class="tag-control">
                                        <input class="input is-small" type="text" v-model="state.info.req_headers[idx].val" placeholder="值">
                                    </div>
                                    <div class="tag-control">
                                        <button type="button" class="button is-small is-danger" @click="delHeader(idx)">删除</button>
                                    </div>
                                </div>
                            </div>
                        </FormControl>
                    </FormField>

                    <FormField label="节点配置" horizontal>
                        <FormControl>
                            <div class="radio-group">
                                <label class="b-radio radio"><input type="radio" :value="1" v-model="state.info.node_option" />
                                    <span class="check"></span>
                                    <span class="control-label">按类型</span>
                                </label>
                                <label class="b-radio radio"><input type="radio" :value="2" v-model="state.info.node_option" />
                                    <span class="check"></span>
                                    <span class="control-label">按IP</span>
                                </label>
                            </div>
                        </FormControl>
                    </FormField>

                    <FormField label="节点类型" horizontal v-if="state.info.node_option == 1">
                        <FormControl>
                            <div class="radio-group">
                                <label class="b-radio radio"><input type="radio" :value="1" v-model="state.info.node_type" />
                                    <span class="check"></span>
                                    <span class="control-label">国内</span>
                                </label>
                                <label class="b-radio radio"><input type="radio" :value="2" v-model="state.info.node_type" />
                                    <span class="check"></span>
                                    <span class="control-label">海外</span>
                                </label>
                            </div>
                        </FormControl>
                    </FormField>

                    <FormField label="可用节点" horizontal v-if="state.info.node_option == 2">
                        <FormControl>
                            <div class="check-group">
                                <label class="b-checkbox checkbox" v-for="(node, idx) in state.nodes">
                                    <input type="checkbox" :value="node.id" />
                                    <span class="check is-primary"></span>
                                    <span class="control-label">{{ node.name }}</span>
                                </label>
                            </div>
                        </FormControl>
                    </FormField>

                    <FormField label="是否开启" horizontal>
                        <FormControl>
                            <div class="radio-group">
                                <label class="b-radio radio">
                                    <input type="radio" :value="0" v-model="state.info.enable" />
                                    <span class="check"></span>
                                    <span class="control-label">关闭</span>
                                </label>
                                <label class="b-radio radio">
                                    <input type="radio" :value="1" v-model="state.info.enable" />
                                    <span class="check"></span>
                                    <span class="control-label">开启</span>
                                </label>
                            </div>
                        </FormControl>
                    </FormField>
                    <hr />
                    <FormField label=" " group horizontal>
                        <FormControl>
                            <button type="button" class="button is-primary" @click="saveInfo">
                                <span>保存数据</span>
                            </button>
                        </FormControl>
                    </FormField>
                </form>
            </div>
        </div>
    </BoxMain>
</template>

<script setup>
import BoxMain from '@adm/components/BoxMain.vue'
import FormField from '@adm/components/FormField.vue'
import FormControl from '@adm/components/FormControl.vue'
import { mdiBallot } from '@mdi/js'
import BasicIcon from '@adm/components/BasicIcon.vue'

import { Get, Post } from "@/lib/http";
import { reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import {useToast} from '@adm/components/toast'

const API = {
    list: "/admin/sites",
    upsert: function (id) {
        return "/admin/sites/" + id + "/upsert"
    }
}

const route = useRoute()
const id = route.query.id ?? 0

const title = ['网站管理', id > 0 ? '编辑' : '新增']

const state = reactive({
    nodes: [],
    info: {
        id: 0,
        name: "",
        root: "",
        key: "",
        desc: "",
        type: 1,
        tags: [],
        cron: "",
        enable: 0,
        node_option: 1,
        node_type: 1,
        node_hosts: [],
        req_headers: [],
    }
})

async function fetchInfo() {    
    if (id == 0) {
        return false
    }
    try {
        let resp = await Get(API.list, {id: id})
        if (resp.data.code === 10000) {
            state.info = Object.assign({}, resp.data.data.siteList[0])
            state.nodes = resp.data.data.nodeList
        }
    } catch (err) {
        console.log(err)
    }
}

function addHeader() {
    state.info.req_headers.push({
        key: '',
        val: ''
    })
}

function delHeader(idx) {
    state.info.req_headers.splice(idx, 1)
}

const toast = useToast()
async function saveInfo() {
    try {
        let resp = await Post(API.upsert(state.info.id), state.info)
        if (resp.data.code === 10000) {
            toast.show('保存成功')
        } else {
            toast.error('保存成功')
        }
    } catch(err) {
        console.log(err)
    }
}

onMounted(fetchInfo)
</script>