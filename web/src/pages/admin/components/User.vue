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
            <div class="column table-container">
                <table class="table is-bordered">
                    <thead>
                    <tr>
                        <th width="10%">ID</th>
                        <th width="20%">头像</th>
                        <th width="20%">用户名</th>
                        <th width="10%">认证类型</th>
                        <th width="40%">认证时间</th>
                    </tr>
                    </thead>
                    <tbody>
                    <tr v-for="(item, idx) in state.list" :key="idx">
                        <td>{{ item.id }}</td>
                        <td>
                            <figure>
                                <img class="avatar" :src="item.avatar">
                            </figure>
                        </td>
                        <td>{{ item.username }}</td>
                        <td><span class="tag is-light">{{ item.auth_type }}</span></td>
                        <td><span class="tag is-light">{{ item.auth_time }}</span></td>
                    </tr>
                    </tbody>
                </table>
            </div>
        </div>
    </div>
</template>

<script>
import {Get} from "@/tools/http";
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
import { onMounted, reactive } from 'vue';

const API = {
    list: "/admin/users"
}

export default {
    name: "user",
    setup() {
        const state = reactive({
            list: []
        })

        async function fetchList() {
            NProgress.start();
            try {
                let resp = await Get(API.list)
                if (resp.data.code === 10001) {
                    console.log(resp.data.msg);
                } else {
                    state.list = resp.data.data;
                }
            } catch(err) {
                console.log(err)
            }
            NProgress.done();
        }

        onMounted(fetchList)

        return {
            state
        }
    }
}
</script>

<style scoped>
table td, table th {
    text-align: center!important;
    vertical-align: middle;
}
td figure {
    text-align: center;
    padding:0 0;
    margin: 0 0;
    display: flex;
    flex-direction: row;
    justify-content: center;
}
td img.avatar {
    width: 32px;
    height: 32px;
    border-radius: 32px 32px;
}
</style>