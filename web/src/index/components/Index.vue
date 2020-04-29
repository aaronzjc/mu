<template>
<div class="content-box">
    <HoTab @change="tabChange" :tabs="tabs"></HoTab>

    <p class="hot-ts" v-if="t != '' ">更新时间: {{ t }}</p>
    <div class="columns hot-container">
        <div class="column hot-list">
            <component v-for="(hot, idx) in list" :is="CardMap[hot['card_type']]" :item="hot" :idx="idx" :key="idx" @toggle-favor="toggleFavor"></component>
        </div>
    </div>

    <Footer></Footer>
</div>
</template>

<script>
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

import Get, {Post} from "../tools/http"
import * as ls from "../tools/ls"
import HoTab from "./HoTab"
import Footer from "./Footer"

import {CardMap, Cards} from "../tools/card";

const API = {
    config: "/config",
    list: "/list",
};

export default {
    name: "Content",
    mounted() {
        this.fetchConfig(this.fetchList);
    },
    data() {
        return {
            tabs: [],
            selected: {
                tab: 0,
                tag: 0,
            },
            list: [],
            t: "",

            CardMap: CardMap
        }
    },
    methods: {
        fetchConfig(callback) {
            // 使用本地缓存，减少一个请求
            let cacheKey = "tabs";
            let tabStr = ls.Get(cacheKey)
            if (tabStr !== false) {
                this.tabs = JSON.parse(tabStr)
                if (typeof callback == "function") {
                    callback(true);
                }
                return true
            }
            Get(API.config).then(function (resp) {
                if (resp.data.code === 10000) {
                    this.tabs = Object.freeze(resp.data.data);
                    ls.Set(cacheKey, JSON.stringify(resp.data.data), 5*60)
                } else {
                    alert(resp.data.msg);
                }
                if (typeof callback == "function") {
                    callback(true);
                }
            }.bind(this))
        },
        fetchList(landing) {
            if (this.tabs.length === 0) {
                return false;
            }
            var key = this.tabs[this.selected.tab]["key"];
            var hkey = undefined;
            if (this.tabs[this.selected.tab]["tags"].length > 0) {
                hkey = this.tabs[this.selected.tab]["tags"][this.selected.tag]["key"];
            }
            if (hkey === undefined || key === undefined) {
                return false;
            }
            NProgress.start();
            let cacheKey = "init_list";
            var needCache = parseInt(this.selected.tab + this.selected.tag) === 0;
            if ( needCache && landing) {
                let listStr = ls.Get(cacheKey)
                if (listStr !== false) {
                    let data = JSON.parse(listStr);
                    this.list = data.list;
                    this.t = data.t;
                    NProgress.done();
                    return true;
                }
            }
            Get(API.list, {
                params: {
                    key: this.tabs[this.selected.tab]["key"],
                    hkey: this.tabs[this.selected.tab]["tags"][this.selected.tag]["key"]
                }
            }).then(function (resp) {
                if (resp.data.code === 10000) {
                    this.list = Object.freeze(resp.data.data.list);
                    this.t = resp.data.data.t;

                    if (needCache) {
                        ls.Set(cacheKey, JSON.stringify(resp.data.data), 60);
                    }
                } else {
                    this.list = [];
                }
                NProgress.done();
            }.bind(this))
        },
        tabChange(data) {
            this.selected = data;
            this.fetchList(false);
        },
        toggleFavor(idx) {
            if (this.list[idx].mark) {
                this.remove(idx);
            } else {
                this.add(idx);
            }
        },
        add(idx) {
            var item = this.list[idx];
            Post("/api/favor/add", {
                key: item.key,
                url: item.origin_url,
                title: item.title,
                site: this.tabs[this.selected.tab]["key"]
            }).then(resp => {
                if (resp.data.code != 10000) {
                    alert("操作失败");
                    return false;
                }

                this.list[idx].mark = true;
            })
        },
        remove(idx) {
            var item = this.list[idx];
            Post("/api/favor/remove", {
                key: item.key,
                site: this.tabs[this.selected.tab]["key"]
            }).then(resp => {
                if (resp.data.code != 10000) {
                    alert("操作失败");
                    return false;
                }

                this.list[idx].mark = false;
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