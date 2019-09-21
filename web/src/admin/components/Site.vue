<template>
    <div>
        <div class="columns">
            <div class="column is-half">
                <div class="field is-grouped">
                    <p class="control is-expanded">
                        <input class="input" type="text" placeholder="关键词">
                    </p>
                    <p class="control">
                        <a class="button is-primary">
                            搜索
                        </a>
                    </p>
                    <p class="control">
                        <a class="button is-info" @click="modal = true">
                            添加网站
                        </a>
                    </p>
                </div>
            </div>
        </div>
        <div class="columns">
            <div class="column">
                <table class="table is-fullwidth is-bordered">
                    <thead>
                    <tr>
                        <th width="10%">网站</th>
                        <th width="15%">地址</th>
                        <th width="10%">cron</th>
                        <th width="5%">类型</th>
                        <th width="10%">节点</th>
                        <th width="5%">状态</th>
                        <th width="15%">操作</th>
                    </tr>
                    </thead>
                    <tbody>
                    <tr v-for="(item, idx) in list" :key="idx">
                        <td>{{ item.name }}</td>
                        <td>{{ item.root }}</td>
                        <td><span class="tag is-warning">{{ item.cron }}</span></td>
                        <td><span class="tag is-light">{{ item.type }}</span></td>
                        <td>
                            <span v-for="(node, i) in item.nodes" :key="i"><span class="tag is-light">{{ node.name }}</span></span>
                        </td>
                        <td>
                            <span class="tag is-light">{{ item.enable ? "开启" : "关闭" }}</span>
                        </td>
                        <td>
                            <div class="buttons are-small">
                                <a class="button is-primary">编辑</a>
                                <a class="button is-warning">下线</a>
                                <a class="button is-danger">删除</a>
                            </div>
                        </td>
                    </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <div :class='[ "modal", { "is-active" : modal } ]'>
            <div class="modal-background"></div>
            <div class="modal-card">
                <header class="modal-card-head">
                    <p class="modal-card-title">编辑网站</p>
                    <button class="delete" aria-label="close" @click="modal = false"></button>
                </header>
                <section class="modal-card-body">
                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">站点</label>
                        </div>
                        <div class="field-body fix-align">
                            <strong>V2ex</strong>
                        </div>
                    </div>

                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">地址</label>
                        </div>
                        <div class="field-body fix-align">
                            <strong>https://v2ex.com</strong>
                        </div>
                    </div>

                    <div class="field is-horizontal">
                        <div class="field-label is-normal">
                            <label class="label">Cron</label>
                        </div>
                        <div class="field-body">
                            <div class="field">
                                <div class="control">
                                    <input class="input" type="text" placeholder="标准Cron表达式">
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
                                        <input type="radio" name="member">
                                        是
                                    </label>
                                    <label class="radio">
                                        <input type="radio" name="member">
                                        否
                                    </label>
                                </div>
                            </div>
                        </div>
                    </div>

                </section>
                <footer class="modal-card-foot">
                    <button class="button is-primary">保存</button>
                    <button class="button" @click="modal = false">取消</button>
                </footer>
            </div>
        </div>
    </div>
</template>

<script>
    export default {
        name: "Site",
        data: () => {
            return {
                modal: false,

                list: [
                    {
                        "name": "V2ex",
                        "key": "v2ex",
                        "root": "https://v2ex.com",
                        "tags": [
                            {
                                "key": "all",
                                "name": "全部"
                            },
                            {
                                "key": "hot",
                                "name": "最热"
                            }
                        ],
                        "type": "HTML",
                        "cron": "* * * * *",
                        "nodes": [
                            {
                                "id": 1,
                                "name": "ucloud"
                            }
                        ],
                        "enable": true
                    },
                    {
                        "name": "抽屉新热榜",
                        "key": "chouti",
                        "root": "https://chouti.com",
                        "tags": [
                            {
                                "key": "hot",
                                "name": "新热榜"
                            },
                            {
                                "key": "24hot",
                                "name": "24小时最热"
                            },
                            {
                                "key": "72hot",
                                "name": "3天最热"
                            }
                        ],
                        "type": "JSON",
                        "cron": "*/30 * * * *",
                        "nodes": [
                            {
                                "id": 1,
                                "name": "Vultr"
                            }
                        ],
                        "enable": true
                    }
                ],

                editForm: {

                },
                addForm: {

                }
            }
        },
        methods: {
            add() {

            },
            edit() {

            },
            cancel() {

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
</style>