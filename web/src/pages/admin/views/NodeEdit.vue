<template>
    <BoxMain :title="title">
        <div class="card">
            <header class="card-header">
                <p class="card-header-title">
                    <BasicIcon :name="mdiBallot"></BasicIcon>
                    节点信息
                </p>
            </header>
            <div class="card-content">
                <form method="get">
                    <FormField label="名称" group horizontal>
                        <FormControl>
                            <input class="input" type="text" v-model="state.info.name" placeholder="名称" />
                        </FormControl>
                    </FormField>
                    <FormField label="地址" horizontal>
                        <FormControl expanded>
                            <input class="input" type="text" v-model="state.info.addr" placeholder="描述" />
                        </FormControl>
                    </FormField>

                    <FormField label="节点类型" horizontal>
                        <FormControl>
                            <div class="radio-group">
                                <label class="b-radio radio"><input type="radio" :value="1" v-model="state.info.type" />
                                    <span class="check"></span>
                                    <span class="control-label">国内</span>
                                </label>
                                <label class="b-radio radio"><input type="radio" :value="2" v-model="state.info.type" />
                                    <span class="check"></span>
                                    <span class="control-label">海外</span>
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
    list: "/admin/nodes",
    upsert: function (id) {
        return "/admin/nodes/" + id + "/upsert"
    }
}

const route = useRoute()
const id = route.query.id ?? 0

const state = reactive({
    info: {}
})

const title = ['节点管理', id > 0 ? '编辑' : '新增']

async function fetchInfo() {
    
    console.log(id)
    if (id == 0) {
        return false
    }
    try {
        let resp = await Get(API.list, {id: id})
        if (resp.data.code === 10000) {
            state.info = Object.assign({}, resp.data.data[0])
        }
    } catch (err) {
        console.log(err)
    }
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