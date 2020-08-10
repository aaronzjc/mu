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
        </div>
        <div class="columns">
            <div class="table-container column">
                <table class="table is-bordered">
                    <thead>
                    <tr>
                        <th width="10%">网站</th>
                        <th width="15%">地址</th>
                        <th width="10%">cron</th>
                        <th width="5%">类型</th>
                        <th width="10%">节点</th>
                        <th width="5%">状态</th>
                        <th width="25%">操作</th>
                    </tr>
                    </thead>
                    <tbody>
                    <tr v-for="(item, idx) in list" :key="idx">
                        <td>{{ item.name }}</td>
                        <td><span class="tag is-white">{{ item.root }}</span></td>
                        <td><span class="tag is-warning" v-if="item.cron !== '' ">{{ item.cron }}</span></td>
                        <td><span class="tag is-light" v-if="item.type != 0">{{ typeMap[item.type] }}</span></td>
                        <td>
                            <template v-if="item.node_option === 1">
                                <span class="tag is-light">{{ nodeTypeMap[item.node_type] }}</span>
                            </template>
                            <template v-if="item.node_option === 2">
                                <template v-if="item.node_hosts.length > 1">
                                    <div class="tags">
                                        <span class="tag is-info" v-for="nodeId in item.node_hosts" :key="nodeId">{{ nodes[nodeId]["name"] }}</span>
                                    </div>
                                </template>
                                <template v-if="item.node_hosts.length == 1">
                                    <span class="tag is-info" v-for="nodeId in item.node_hosts" :key="nodeId">{{ nodes[nodeId]["name"] }}</span>
                                </template>
                            </template>
                        </td>
                        <td>
                            <span class="tag is-light" v-if="!item.enable">未启用</span>
                            <span class="tag is-success" v-if="item.enable">启用</span>
                        </td>
                        <td>
                            <div class="buttons are-small">
                                <a class="button is-primary" @click="view(idx)">查看</a>
                                <a class="button is-warning" @click="edit(idx)">编辑</a>
                                <a class="button is-info" @click="craw(idx)">立刻拉取</a>
                            </div>
                        </td>
                    </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <div :class='[ "modal", { "is-active" : editModal } ]'>
            <div class="modal-background"></div>
            <div class="modal-card">
                <header class="modal-card-head">
                    <p class="modal-card-title">{{ editForm.id > 0 ? "编辑网站" : "添加网站" }}</p>
                    <button class="delete" aria-label="close" @click="closeEdit"></button>
                </header>
                <section class="modal-card-body">
                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">站点</label>
                        </div>
                        <div class="field-body">
                            <div class="control">
                                <input class="input" type="text" placeholder="名字" v-model="editForm.name">
                            </div>
                        </div>
                    </div>

                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">地址</label>
                        </div>
                        <div class="field-body fix-align">
                            <strong>{{ editForm.root }}</strong>
                        </div>
                    </div>

                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">Key</label>
                        </div>
                        <div class="field-body">
                            <div class="control">
                                <input class="input" type="text" placeholder="key" v-model="editForm.key">
                            </div>
                        </div>
                    </div>

                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">描述</label>
                        </div>
                        <div class="field-body">
                            <div class="field">
                                <div class="control">
                                    <input class="input" type="text" placeholder="描述" v-model="editForm.desc">
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">抓取类型</label>
                        </div>
                        <div class="field-body fix-align">
                            <strong>{{ typeMap[editForm.type] }}</strong>
                        </div>
                    </div>

                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">请求头</label>
                        </div>
                        <div class="field-body tag-field">
                            <template v-if="editForm.req_headers.length == 0">
                                <div class="tag-control">
                                    <button class="button is-small is-primary" @click="addHeader">新建</button>
                                </div>
                            </template>
                            <template v-if="editForm.req_headers.length > 0">
                                <div class="tag-item" v-for="(item, idx) in editForm.req_headers" :key="idx">
                                    <div class="tag-control">
                                        <input class="input is-small" type="text" placeholder="名称" v-model="editForm.req_headers[idx].key">
                                    </div>
                                    <div class="tag-control">
                                        <input class="input is-small" type="text" placeholder="值" v-model="editForm.req_headers[idx].val">
                                    </div>
                                    <div class="tag-control">
                                        <button class="button is-small is-primary" @click="addHeader">新建</button>
                                        <button class="button is-small is-danger" @click="delHeader(idx)">删除</button>
                                    </div>
                                </div>
                            </template>
                        </div>
                    </div>

                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">标签</label>
                        </div>
                        <div class="field-body tag-field">
                            <div class="tag-item" v-for="(tag, tagIdx) in editForm.tags" :key="tagIdx">
                                <div class="tag-control">
                                    <input class="input is-small" type="text" placeholder="标识" readonly v-model="editForm.tags[tagIdx].key">
                                </div>
                                <div class="tag-control">
                                    <input class="input is-small" type="text" placeholder="名字" v-model="editForm.tags[tagIdx].name">
                                </div>
                                <div class="tag-control">
                                    <button :class="[ 'button', 'is-small',  { 'is-danger' : tag.enable === 1 }, { 'is-success' : tag.enable === 0 } ]" @click="toggle(tagIdx)">{{ tag.enable === 1 ? "关闭" : "开启" }}</button>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">Cron</label>
                        </div>
                        <div class="field-body">
                            <div class="control">
                                <input class="input" type="text" placeholder="标准Cron表达式" v-model="editForm.cron">
                            </div>
                        </div>
                    </div>

                    <div class="field is-horizontal">
                        <div class="field-label">
                            <label class="label">节点配置</label>
                        </div>
                        <div class="field-body">
                            <div class="field is-narrow">
                                <div class="control">
                                    <label class="radio">
                                        <input type="radio" name="node_option" :value="parseInt(1)" v-model="editForm.node_option">
                                        按节点类型
                                    </label>
                                    <label class="radio">
                                        <input type="radio" name="node_option" :value="parseInt(2)" v-model="editForm.node_option">
                                        按服务器IP
                                    </label>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="field is-horizontal" v-if="editForm.node_option === 1">
                        <div class="field-label">
                            <label class="label">节点类型</label>
                        </div>
                        <div class="field-body">
                            <div class="field is-narrow">
                                <div class="control">
                                    <label class="radio">
                                        <input type="radio" name="node_type" :value="parseInt(1)" v-model="editForm.node_type">
                                        大陆
                                    </label>
                                    <label class="radio">
                                        <input type="radio" name="node_type" :value="parseInt(2)" v-model="editForm.node_type">
                                        海外
                                    </label>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="field is-horizontal" v-if="editForm.node_option === 2">
                        <div class="field-label">
                            <label class="label">节点列表</label>
                        </div>
                        <div class="field-body">
                            <div class="field is-narrow">
                                <div class="control">
                                    <div class="node_item" v-for="(node, nodeIdx) in nodes" :key="nodeIdx">
                                        <label class="checkbox">
                                            <input type="checkbox" :value="node.id" v-model="editForm.node_hosts">
                                            {{ node.name }}
                                        </label>
                                    </div>
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
                                        <input type="radio" name="enable" :value="parseInt(0)" v-model="editForm.enable">
                                        关闭
                                    </label>
                                    <label class="radio">
                                        <input type="radio" name="enable" :value="parseInt(1)" v-model="editForm.enable">
                                        开启
                                    </label>
                                </div>
                            </div>
                        </div>
                    </div>

                </section>
                <footer class="modal-card-foot">
                    <button class="button is-primary" @click="save">保存</button>
                    <button class="button" @click="closeEdit">取消</button>
                </footer>
            </div>
        </div>

        <div :class='[ "modal", { "is-active" : viewModal } ]'>
            <div class="modal-background"></div>
            <div class="modal-card">
                <header class="modal-card-head">
                    <p class="modal-card-title">查看站点</p>
                    <button class="delete" aria-label="close" @click="closeView"></button>
                </header>
                <section class="modal-card-body">
                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">站点</label>
                        </div>
                        <div class="field-body fix-align">
                            <strong>{{ viewForm.name }}</strong>
                        </div>
                    </div>

                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">地址</label>
                        </div>
                        <div class="field-body fix-align">
                            <strong>{{ viewForm.root }}</strong>
                        </div>
                    </div>

                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">Key</label>
                        </div>
                        <div class="field-body">
                            <div class="field-body fix-align">
                                <strong>{{ viewForm.key }}</strong>
                            </div>
                        </div>
                    </div>

                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">描述</label>
                        </div>
                        <div class="field-body">
                            <div class="field">
                                <div class="field-body fix-align">
                                    <strong>{{ viewForm.desc }}</strong>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">抓取类型</label>
                        </div>
                        <div class="field-body">
                            <div class="field">
                                <div class="field-body fix-align">
                                    <strong>{{ typeMap[viewForm.type] }}</strong>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">标签</label>
                        </div>
                        <div class="field-body">
                            <div class="field">
                                <div class="field-body fix-align tags">
                                    <template v-for="tag in viewForm.tags">
                                        <span :key="tag.key" v-if="tag.enable === 1" class="tag is-link">{{ tag.key }}-{{ tag.name }}</span>
                                    </template>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">Cron</label>
                        </div>
                        <div class="field-body">
                            <div class="field">
                                <div class="field-body fix-align">
                                    <span class="tag is-warning" v-if="viewForm.cron != ''">{{ viewForm.cron }}</span>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="field is-horizontal" v-if="viewForm.node_option === 1">
                        <div class="field-label is-normal">
                            <label class="label">节点配置</label>
                        </div>
                        <div class="field-body">
                            <div class="field">
                                <div class="field-body fix-align">
                                    <strong>{{ nodeTypeMap[viewForm.node_type] }}</strong>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="field is-horizontal" v-if="viewForm.node_option === 2">
                        <div class="field-label is-normal">
                            <label class="label">节点列表</label>
                        </div>
                        <div class="field-body">
                            <div class="field">
                                <div class="field-body fix-align tags">
                                    <span :key="nodeId" v-for="nodeId in viewForm.node_hosts" class="tag is-info">{{ nodes[nodeId].name }}</span>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">是否开启</label>
                        </div>
                        <div class="field-body">
                            <div class="field">
                                <div class="field-body fix-align">
                                    <strong>{{ viewForm.enable === 1 ? "开启" : "未开启" }}</strong>
                                </div>
                            </div>
                        </div>
                    </div>

                </section>
                <footer class="modal-card-foot">
                    <button class="button is-primary" @click="closeView">确认</button>
                </footer>
            </div>
        </div>
    </div>
</template>

<script>
import {Get, Post} from "@/tools/http";
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
import {nodeType, crawType} from "../def";

const initEdit = {
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
};
export default {
    name: "Site",
    mounted() {
        this.fetchList();
    },
    data: () => {
        return {
            viewModal: false,
            editModal: false,

            typeMap: crawType,

            nodeTypeMap: nodeType,

            nodes: {},

            list: [],

            editForm: {
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
            },

            viewForm: {

            }
        }
    },
    methods: {
        fetchList() {
            NProgress.start();
            Get("/admin/site/list").then(resp => {
                if (resp.data.code === 10000) {
                    this.nodes = resp.data.data.nodeList;
                    this.list = resp.data.data.siteList;
                } else {
                    console.log(resp.data.msg);
                    this.nodes = [];
                    this.list = [];
                }
                NProgress.done();
            }).catch( () => {
                NProgress.done();
            })
        },
        save() {
            Post("/admin/site/update", this.editForm).then(resp => {
                if (resp.data.code === 10000) {
                    alert("保存成功");
                    this.closeEdit();
                    this.fetchList();
                } else {
                    alert(resp.data.msg);
                }
            });
        },
        edit(idx) {
            this.editModal = true;
            this.editForm = Object.assign({}, this.list[idx]);
        },
        view(idx) {
            this.viewModal = true;
            this.viewForm = this.list[idx];
        },
        closeView() {
            this.viewModal = false;
        },
        closeEdit() {
            this.editModal = false;
            this.editForm = JSON.parse(JSON.stringify(initEdit));
        },
        addHeader: function () {
            this.editForm.req_headers.push({
                key: "",
                val: ""
            });
        },
        delHeader: function (idx) {
            this.editForm.req_headers.splice(idx, 1)
        },
        toggle(idx) {
            var c = this.editForm.tags[idx]["enable"];
            var r = c === 0 ? 1 : 0 ;
            this.editForm.tags[idx]["enable"] = r;
        },
        craw(idx) {
            Post("/admin/site/craw", {"id": this.list[idx]["id"]}).then(resp => {
                if (resp.data.code === 10000) {
                    alert("数据抓取成功")
                } else {
                    alert(resp.data.msg);
                }
            });
        }
    }
}
</script>

<style lang="scss" scoped>
    table {
    td, th {
        text-align: center!important;
    }
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
    .tag-field {
        width: 100%;
        display: flex;
        flex-direction:column;

        .tag-item {
            display: flex;
            flex-direction: row;
            * {
                margin-right: 4px;
                margin-bottom: 2px;
            }
        }
    }
</style>