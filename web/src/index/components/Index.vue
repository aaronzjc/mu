<template>
<div class="content-box">
    <HoTab @change="tabChange" :tabs="tabs"></HoTab>

    <p class="hot-ts" v-if="t != '' ">更新时间: {{ t }}</p>
    <div class="columns hot-container">
        <div class="column hot-list">
            <div class="hot" v-for="(hot, idx) in list" :key="idx">
                <div class="hot-item">
                    <a :href="hot.origin_url" :title="hot.title" target="_blank">{{ hot.title }}</a>
                </div>
                <div class="divider"></div>
                <div class="hot-opt"  @click="toggleFavor(idx)">
                    <template v-if="!hot.mark">
                        <svg style="width:20px;height:20px" viewBox="0 0 24 24">
                            <path fill="#b5b5b5" d="M12.1,18.55L12,18.65L11.89,18.55C7.14,14.24 4,11.39 4,8.5C4,6.5 5.5,5 7.5,5C9.04,5 10.54,6 11.07,7.36H12.93C13.46,6 14.96,5 16.5,5C18.5,5 20,6.5 20,8.5C20,11.39 16.86,14.24 12.1,18.55M16.5,3C14.76,3 13.09,3.81 12,5.08C10.91,3.81 9.24,3 7.5,3C4.42,3 2,5.41 2,8.5C2,12.27 5.4,15.36 10.55,20.03L12,21.35L13.45,20.03C18.6,15.36 22,12.27 22,8.5C22,5.41 19.58,3 16.5,3Z" />
                        </svg>
                    </template>
                    <template v-else>
                        <svg style="width:20px;height:20px" viewBox="0 0 24 24">
                            <path fill="#ff3860" d="M12,21.35L10.55,20.03C5.4,15.36 2,12.27 2,8.5C2,5.41 4.42,3 7.5,3C9.24,3 10.91,3.81 12,5.08C13.09,3.81 14.76,3 16.5,3C19.58,3 22,5.41 22,8.5C22,12.27 18.6,15.36 13.45,20.03L12,21.35Z" />
                        </svg>
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
.hot-container {
    margin-top: 0;
}

.hot-list {
    padding-top:0;
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

