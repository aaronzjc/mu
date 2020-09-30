<template>
<div class="content-box">
    <HoTab @change="tabChange" :tabs="state.tabs"></HoTab>

    <p class="hot-ts" v-if="state.t !== '' ">更新时间: {{ state.t }}</p>
    <div class="columns hot-container">
        <div class="column hot-list">
            <component v-for="(hot, idx) in state.list" :is="state.CardMap[hot['card_type']]" :item="hot" :idx="idx" :key="idx" @toggle-favor="toggleFavor"></component>
        </div>
    </div>

    <Footer></Footer>
</div>
</template>

<script>
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

import { Get, Post } from "@/tools/http"
import * as ls from "@/tools/ls"
import HoTab from "./HoTab"
import Footer from "./Footer"

import {CardMap, Cards} from "../ext/card";
import { onMounted, reactive } from 'vue'
export default {
    name: "Content",
    setup() {
        const API = {
            config: "/api/config",
            list: "/api/list",
        }
        const state = reactive({
            tab: [
                {"name":"Moo-Yuu","key":"新闻","tags":[{"key":"昨天","name":"昨天","enable":1},{"key":"今天","name":"今天","enable":1},{"key":"明天","name":"明天","enable":1}]},
            ],
            selected: {
                tab: 0,
                tag: 0,
            },
            list: [],
            t: "还没更新呢",
            CardMap: CardMap
        })
        
        // 获取Tab列表
        async function fetchConfig(callback) {
            // 使用本地缓存，减少一个请求
            let cacheKey = "tabs";
            let tabStr = ls.Get(cacheKey)
            if (tabStr !== false) {
                state.tabs = JSON.parse(tabStr)
                if (typeof callback == "function") {
                    callback(true);
                }
                return true
            }
            let resp = await Get(API.config)
            if (resp.data.code === 10000) {
                state.tabs = Object.freeze(resp.data.data);
                ls.Set(cacheKey, JSON.stringify(resp.data.data), 5*60)
            } else {
                alert(resp.data.msg);
            }
            if (typeof callback == "function") {
                callback(true);
            }
        }

        // 获取卡片列表
        async function fetchList(landing) {
            if (state.tabs.length === 0) {
                return false;
            }
            let key = state.tabs[state.selected.tab]["key"]
            let hkey = undefined;
            if (state.tabs[state.selected.tab]["tags"].length > 0) {
                hkey = state.tabs[state.selected.tab]["tags"][state.selected.tag]["key"];
            }
            if (hkey === undefined || key === undefined) {
                return false;
            }
            NProgress.start();
            let cacheKey = "init_list";
            let needCache = parseInt(state.selected.tab + state.selected.tag) === 0;
            if ( needCache && landing) {
                let listStr = ls.Get(cacheKey)
                if (listStr !== false) {
                    let data = JSON.parse(listStr);
                    state.list = data.list;
                    state.t = data.t;
                    NProgress.done();
                    return true;
                }
            }
            let resp = await Get(API.list, {
                key: state.tabs[state.selected.tab]["key"],
                hkey: state.tabs[state.selected.tab]["tags"][state.selected.tag]["key"]
            })
            if (resp.data.code === 10000) {
                state.list = Object.freeze(resp.data.data.list);
                state.t = resp.data.data.t;

                if (needCache) {
                    ls.Set(cacheKey, JSON.stringify(resp.data.data), 60);
                }
            } else {
                state.list = [];
            }
            NProgress.done();
        }

        onMounted(fetchConfig(fetchList))

        return {
            state,
            fetchConfig,
            fetchList
        }
    },
    methods: {
        tabChange(data) {
            this.state.selected = data;
            this.fetchList(false);

            window.scrollTo({
                top: 0,
                behavior: "smooth"
            })
        },
        toggleFavor(idx) {
            if (this.state.list[idx].mark) {
                this.remove(idx);
            } else {
                this.add(idx);
            }
        },
        add(idx) {
            let item = this.state.list[idx];
            Post("/api/favor/add", {
                key: item.key,
                url: item.origin_url,
                title: item.title,
                site: this.state.tabs[this.state.selected.tab]["key"]
            }).then(resp => {
                if (resp.data.code != 10000) {
                    alert("操作失败");
                    return false;
                }

                this.state.list[idx].mark = true;
            })
        },
        remove(idx) {
            let item = this.state.list[idx];
            Post("/api/favor/remove", {
                key: item.key,
                site: this.state.tabs[this.state.selected.tab]["key"]
            }).then(resp => {
                if (resp.data.code != 10000) {
                    alert("操作失败");
                    return false;
                }

                this.state.list[idx].mark = false;
            })
        }
    },
    components: {
        HoTab,
        Footer,

        /* eslint-disable vue/no-unused-components */
        ...Cards
    }
}
</script>