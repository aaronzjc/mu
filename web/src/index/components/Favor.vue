<template>
    <div class="content-box">
        <div class="columns">
            <div class="column">
                <h4 class="title is-4 has-text-centered">俺的收藏夹</h4>
                <div class="field is-grouped search-form">
                    <p class="control is-expanded">
                        <input class="input" type="text" placeholder="搜一搜">
                    </p>
                    <p class="control">
                        <a class="button is-info">
                            搜一搜
                        </a>
                    </p>
                </div>
            </div>
        </div>

        <HoTab @change="tabChange" :tabs="tabs"></HoTab>

        <div class="columns">
            <div class="column hot-list">
                <div class="hot" v-for="(hot, idx) in list" :key="idx">
                    <div class="hot-item">
                        <p class="hot-ts has-text-grey">2019.12.01</p>
                        <a :href="hot.origin_url" :title="hot.title" target="_blank">{{ hot.title }}</a>
                    </div>
                    <div class="divider"></div>
                    <div class="hot-opt" @click="remove(idx)">
                        <i class="far fa-trash-alt has-text-danger"></i>
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
import HoTab from "./HoTab"
import Footer from "./Footer"
import {Get, Post} from "../tools/http";

const API = {
    config: "/api/favor/config",
    list: "/api/favor/list",
    remove: "/api/favor/remove"
};

export default {
    name: "Favor",
    data() {
        return {
            selected: {
                tab: 0,
                tag: 0
            },
            list: [],
            tabs: []
        }
    },
    created() {
        this.fetchConfig(this.fetchFavor);
    },
    methods: {
        tabChange(playload) {
            this.selected = playload;
            this.fetchFavor();
        },
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
        fetchFavor() {
            if (this.tabs.length === 0) {
                return false;
            }
            var key = this.tabs[this.selected.tab]["key"];
            if (key === undefined) {
                return false;
            }
            NProgress.start();
            Get(API.list, {
                s: this.tabs[this.selected.tab]["key"],
            }).then(function (resp) {
                if (resp.data.code === 10000) {
                    this.list = resp.data.data;
                } else {
                    this.list = [];
                }
                NProgress.done();
            }.bind(this))
        },
        remove(idx) {
            if (!confirm("确定移除该条吗？")) {
                return false;
            }

            var item = this.list[idx];

            Post(API.remove, {
                key: item.key,
                site: this.tabs[this.selected.tab]["key"]
            }).then(resp => {
                if (resp.data.code != 10000) {
                    alert("操作失败");
                    return false;
                }
                this.fetchFavor();
            });

        }
    },
    components: {
        HoTab,
        Footer
    }
}
</script>

<style lang="scss" scoped>
@import 'bulma/sass/utilities/mixins';

.search-form {
    margin: 0 auto;
    padding: 0 1rem;
}

@include desktop {
    .search-form {
        width: 50%;
    }
}

.hot-list {
    padding-top: 0;
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
    }

    .hot-item {
        width: 98%;
        margin-right: 2px;
        display: flex;
        word-break: break-word;
        flex-direction: column;

        .hot-ts {
            font-size: 0.6rem;
        }
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