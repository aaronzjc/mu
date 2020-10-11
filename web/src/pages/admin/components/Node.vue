<template>
    <div>
        <div class="columns">
            <div class="column is-one-third">
                <div class="field is-grouped">
                    <p class="control is-expanded">
                        <input class="input" type="text" placeholder="关键词">
                    </p>
                    <p class="control">
                        <a class="button is-primary">搜索</a>
                    </p>
                </div>
            </div>
            <div class="column add-btn">
                <a class="button is-info" @click="add">添加节点</a>
            </div>
        </div>
        <div class="columns">
            <div class="column table-container">
                <table class="table is-bordered">
                    <thead>
                    <tr>
                        <th width="5%">ID</th>
                        <th width="10%">名称</th>
                        <th width="20%">Addr</th>
                        <th width="5%">类型</th>
                        <th width="5%">Ping</th>
                        <th width="5%">状态</th>
                        <th width="30%">操作</th>
                    </tr>
                    </thead>
                    <tbody>
                    <tr v-for="(item, idx) in state.list" :key="idx">
                        <td>{{ item.id }}</td>
                        <td>{{ item.name }}</td>
                        <td>{{ item.addr }}</td>
                        <td><span class="tag is-warning">{{ state.typeMap[item.type] }}</span></td>
                        <td>
                            <span :class="[ 'ping',  { 'has-background-success' : item.ping === 1 }, { 'has-background-danger' : item.ping === 0,  } ]"></span>
                        </td>
                        <td>
                            <span class="tag is-light" v-if="!item.enable">未开启</span>
                            <span class="tag is-success" v-if="item.enable">开启</span>
                        </td>
                        <td>
                            <div class="buttons are-small">
                                <a class="button is-primary" @click="view(idx)">查看</a>
                                <a class="button is-warning" @click="edit(idx)">编辑</a>
                                <a class="button is-danger" @click="del(idx)">删除</a>
                            </div>
                        </td>
                    </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <div :class='[ "modal", { "is-active" : state.editModal } ]'>
            <div class="modal-background"></div>
            <div class="modal-card">
                <header class="modal-card-head">
                    <p class="modal-card-title">{{ state.editForm.id > 0 ? "编辑节点" : "添加节点" }}</p>
                    <button class="delete" aria-label="close" @click="cancel"></button>
                </header>
                <section class="modal-card-body">
                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">节点</label>
                        </div>
                        <div class="field-body">
                            <div class="control">
                                <input class="input" type="text" placeholder="节点名字" v-model="state.editForm.name">
                            </div>
                        </div>
                    </div>

                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">Addr</label>
                        </div>
                        <div class="field-body">
                            <div class="control">
                                <input class="input" type="text" placeholder="addr, ip:port" v-model="state.editForm.addr">
                            </div>
                        </div>
                    </div>

                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">类型</label>
                        </div>
                        <div class="field-body">
                            <div class="field">
                                <div class="control select">
                                    <select v-model="state.editForm.type">
                                        <option :value="parseInt(key)" v-for="(val, key) in state.typeMap" :key="key">{{ val }}</option>
                                    </select>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="field is-horizontal">
                        <div class="field-label">
                            <label class="label">是否开启</label>
                        </div>
                        <div class="field-body">
                            <div class="field is-narrow">
                                <div class="control">
                                    <label class="radio">
                                        <input type="radio" v-model="state.editForm.enable" v-bind:value="1">
                                        是
                                    </label>
                                    <label class="radio">
                                        <input type="radio" v-model="state.editForm.enable" v-bind:value="0">
                                        否
                                    </label>
                                </div>
                            </div>
                        </div>
                    </div>

                </section>
                <footer class="modal-card-foot">
                    <button class="button is-primary" @click="save">保存</button>
                    <button class="button" @click="cancel">取消</button>
                </footer>
            </div>
        </div>

        <div :class='[ "modal", { "is-active" : state.viewModal } ]'>
            <div class="modal-background"></div>
            <div class="modal-card">
                <header class="modal-card-head">
                    <p class="modal-card-title">节点详情</p>
                    <button class="delete" aria-label="close" @click="state.viewModal = false"></button>
                </header>
                <section class="modal-card-body">
                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">节点</label>
                        </div>
                        <div class="field-body fix-align">
                            <strong>{{ state.viewForm.name }}</strong>
                        </div>
                    </div>

                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">Addr</label>
                        </div>
                        <div class="field-body fix-align">
                            <strong>{{ state.viewForm.addr }}</strong>
                        </div>
                    </div>

                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">节点类型</label>
                        </div>
                        <div class="field-body fix-align">
                            <span class="tag is-info">{{ state.typeMap[state.viewForm.type] }}</span>
                        </div>
                    </div>

                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">是否启用</label>
                        </div>
                        <div class="field-body fix-align">
                            <span class="tag is-success" v-if="state.viewForm.enable">启用</span>
                            <span class="tag is-warning" v-if="!state.viewForm.enable">未启用</span>
                        </div>
                    </div>

                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">Ping</label>
                        </div>
                        <div class="field-body fix-align">
                            <span class="tag is-light">{{ state.viewForm.ping == 1 ? "在线" : "掉线" }}</span>
                        </div>
                    </div>

                </section>
                <footer class="modal-card-foot">
                    <button class="button is-primary" @click="state.viewModal = false">确定</button>
                </footer>
            </div>
        </div>
    </div>
</template>

<script>
import {Get, Post} from "@/tools/http";
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
import {nodeType} from "../def";
import { onMounted, reactive } from 'vue';

export default {
    name: "Node",
    setup() {
        const state = reactive({
            editModal: false,
            viewModal: false,

            typeMap: Object.freeze(nodeType),

            list: [],

            editForm: {},

            viewForm: {}
        })

        async function fetchList() {
            NProgress.start();
            try {
                let resp = await Get("/admin/node/list")
                if (resp.data.code === 10001) {
                    console.log(resp.data.msg)
                } else {
                    state.list = resp.data.data;
                }
            } catch(err) {
                // eslint-disable-next-line
                console.log(err)
            }
            NProgress.done();
        }

        async function del(idx) {
            if (!confirm("确认删除吗?")) {
                return false;
            }
            try {
                let resp = await Get("/admin/node/del", {"id": state.list[idx].id})
                if (resp.data.code !== 10000) {
                    alert(resp.data.msg);
                }
                fetchList()
            } catch(err) {
                // eslint-disable-next-line
                console.log(err)
            }
        }

        async function save() {
            try {
                let resp = await Post("/admin/node/upsert", state.editForm)
                if (resp.data.code === 10001) {
                    alert(resp.data.msg);
                } else {
                    cancel();
                    alert("保存成功")
                }
                fetchList()
            } catch(err) {
                // eslint-disable-next-line
                console.log(err)
            }
        }

        let view = (idx) => {
            state.viewForm = state.list[idx]
            state.viewModal = true
        }

        let add = () => {
            state.editForm = {
                "id": 0,
                "name": "",
                "addr": "",
                "type": 1,
                "ping": 0,
                "enable": 0
            }
            state.editModal = true
        }

        let edit = (idx) => {
            state.viewModal = false
            state.editForm = Object.assign({}, state.list[idx])
            state.editModal = true
        }

        let cancel = () => {
            state.viewModal = false
            state.editModal = false
            state.editForm = {
                "id": 0,
                "name": "",
                "addr": "",
                "type": 1,
                "ping": 0,
                "enable": 0
            }
        }

        onMounted(fetchList)

        return {
            state,
            fetchList,
            add,
            view,
            edit,
            save,
            del,
            cancel
        }
    }
}
</script>

<style lang="scss" scoped>
table {
    td, th {
        text-align: center!important;
        vertical-align: middle;
    }
}
.modal {
    justify-content: unset;
}
.modal-card {
    margin-top: 100px;
    width: 600px;
}
.modal-card-head {
    padding: 1rem;

    .modal-card-title {
        font-size: 1.2rem;
    }
}
.modal-card-foot {
    padding: 1rem;
    justify-content: center;
}
.fix-align {
    padding-top: 0.375em;
}
.ping {
    display: inline-block;
    width: 15px;
    height:15px;
    border-radius: 15px 15px;
}
.add-btn {
    padding-left: 0;
}
</style>