<template>
<div class="content-box">
    <HoTab @change="tabChange" :tabs="tabs"></HoTab>

    <p class="hot-ts" v-if="t != '' ">更新时间: {{ t }}</p>
    <div class="columns">
        <div class="column hot-list">
            <div class="hot" v-for="(hot, idx) in list" :key="idx">
                <div class="hot-item">
                    <a :href="hot.origin_url" :title="hot.title" target="_blank">{{ hot.title }}</a>
                </div>
                <div class="divider"></div>
                <div class="hot-opt"  @click="toggleFavor(idx)">
                    <template v-if="hot.mark">
                        <fa-icon :icon="['fas', 'heart']" class="has-text-danger"></fa-icon>
                    </template>
                    <template v-else>
                        <fa-icon :icon="['far', 'heart']" class="has-text-grey"></fa-icon>
                    </template>
                </div>
            </div>
        </div>
    </div>

    <Footer></Footer>
</div>
</template>

<script>
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'
import 'bulma/css/bulma.css'
import Get, {Post} from "../tools/http"
import HoTab from "./HoTab"
import Footer from "./Footer"

const API = {
    config: "/config",
    list: "/list",
};

export default {
    name: "Content",
    created() {
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
            t: ""
        }
    },
    methods: {
        fetchConfig(callback) {
            Get(API.config).then(function (resp) {
                if (resp.data.code === 10000) {
                    this.tabs = resp.data.data;
                } else {
                    alert(resp.data.msg);
                }
                if (typeof callback == "function") {
                    callback();
                }
            }.bind(this))
        },
        fetchList() {
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
            Get(API.list, {
                params: {
                    key: this.tabs[this.selected.tab]["key"],
                    hkey: this.tabs[this.selected.tab]["tags"][this.selected.tag]["key"]
                }
            }).then(function (resp) {
                if (resp.data.code === 10000) {
                    this.list = resp.data.data.list;
                    this.t = resp.data.data.t;
                } else {
                    this.list = [];
                }
                NProgress.done();
            }.bind(this))
        },
        tabChange(data) {
            this.selected = data;
            this.fetchList();
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
        Footer
    }
}
</script>

<style lang="scss" scoped>
.hot-ts {
    color: #939393;
    font-size: 0.8rem;
}

.hot-list {
    flex-basis: unset;
    width: 100%;
}
.hot {
    width: 100%;
    min-height: 2rem;
    height: auto;
    margin: 0.5rem 0;
    display: flex;
    flex-direction: row;
    .hot-opt {
        padding-left: 3px;
        display: flex;
        flex-direction: row;
        align-items: center;

        &:hover {
             cursor: pointer;
        }

        .tag {
            cursor: pointer;
        }
    }
    .hot-item {
        width: 98%;
        margin-right: 2px;
        display: flex;
        align-items: center;
        word-break: break-word;
    }
    .divider {
        width: 2px;
        margin: 10px 4px;
        background: hsl(0, 0%, 71%)
    }
    &:hover {
        .divider {
            background: hsl(0, 0%, 29%);
        }
    }
}
</style>

