<template>
    <div class="container">
        <div class="logo"><img src="../assets/logo.png" alt=""></div>
        <div class="columns">
            <div class="column">
                <div class="tabs">
                    <ul>
                        <li v-for="(tab, idx) in tabs" :class="{ 'is-active' : idx == selected.tab }" @click="switchTab(idx)" :key="idx"><a>{{ tab.name }}</a></li>
                    </ul>
                </div>
            </div>
        </div>
        <div class="columns">
            <div class="column tab-tag tags">
                <span @click="switchTag(idx)" :class="[ 'tag', { 'is-primary' : idx == selected.tag } ]" v-for="(tag, idx) in tabs[selected.tab]['tags']" :key="idx">{{ tag.name }}</span>
            </div>
        </div>
        <p class="hot-ts" v-if="t != '' ">更新时间: {{ t }}</p>
        <div class="columns">
            <div class="column hot-list">
                <div class="hot" v-for="(hot, idx) in list" :key="idx">
<!--                    <div class="hot-opt">-->
<!--                        <span class="tag is-warning"><i class="fas fa-heart"></i></span>-->
<!--                    </div>-->
                    <div class="hot-item">
                        <a :href="hot.origin_url" :title="hot.title" target="_blank">{{ hot.title }}</a>
                    </div>
                </div>
            </div>
        </div>
        <div class="columns">
            <div class="column copyright has-text-centered">
                <p>
                    <a href="https://github.com/aaronzjc">@aaronzjc</a>开发, 源码<a href="https://github.com/aaronzjc/crawler">在此</a>，欢迎Star.
                </p>
                <p class="backtop">
                    <a href="javascript:scrollTo(0,0);">回到顶部</a>
                </p>
            </div>
        </div>
    </div>
</template>

<script>
import NProgress from 'nprogress'
import 'bulma/css/bulma.css'
import 'nprogress/nprogress.css'
// import '@fortawesome/fontawesome-free/js/all.min'
import Get from "../tools/http"

const API = {
    config: "/config",
    list: "/aj"
};

export default {
    name: "Main",
    created: function () {
        this.fetchConfig(this.fetchList);
    },
    data: function () {
        return {
            tabs: [
                {
                    "name": "",
                    "key": undefined,
                    "tags": [
                        {
                            "name": "",
                            "key": undefined
                        }
                    ]
                }
            ],
            selected: {
                tab: 0,
                tag: 0
            },
            list: [],
            t: ""
        }
    },
    methods: {
        fetchConfig: function (callback) {
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
        fetchList: function () {
            if (this.tabs.length === 0) {
                return false;
            }
            NProgress.start();
            var key = this.tabs[this.selected.tab]["key"];
            var hkey = this.tabs[this.selected.tab]["tags"][this.selected.tag]["key"];
            if (hkey === undefined || key === undefined) {
                return false;
            }
            Get(API.list, {
                params: {
                    key: this.tabs[this.selected.tab]["key"],
                    hkey: this.tabs[this.selected.tab]["tags"][this.selected.tag]["key"]
                }
            }).then(function (resp) {
                if (resp.data.code === 10000) {
                    var list = resp.data.data.list;
                    this.list = list;
                    this.t = resp.data.data.t;
                } else {
                    this.list = [];
                }
                NProgress.done();
            }.bind(this))
        },
        switchTab: function (idx) {
            this.selected.tab = idx;
            this.selected.tag = 0;
            this.fetchList()
        },
        switchTag: function (idx) {
            this.selected.tag = idx;
            this.fetchList()
        }
    }
}
</script>

<style lang="scss" scoped>
.logo {
    text-align: center;
    font-size: 1.2rem;
    font-family: Monaco;
}
.logo img {
    width: 36px;
    height: 36px;
}
.container {
    max-width: 960px;
}
.tab-tag {
    padding-top: 0px;
    padding-bottom: 0px;
    cursor: pointer;
}
.hot-ts {
    color: #939393;
    font-size: 0.8rem;
}
.copyright {
    font-size: 0.85rem;
}
.backtop {
    padding-top:1rem;
    a {
        cursor: pointer;
    }
}
.hot-list {
    flex-basis: unset;
    width: 100%;
}
.hot {
    width: 100%;
    height: 2.5rem;
    line-height: 2.5rem;

    display: flex;
    flex-direction: row;
    .hot-opt {
        margin-right: 4px;
        display: flex;
        flex-direction: row;
        align-items: center;
        .tag {
            cursor: pointer;
        }
    }
    .hot-item {
        max-width: 98%;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }
}
</style>