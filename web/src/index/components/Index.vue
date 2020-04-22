<template>
<div class="content-box">
    <HoTab @change="tabChange" :tabs="tabs"></HoTab>

    <p class="hot-ts" v-if="t != '' ">更新时间: {{ t }}</p>
    <div class="columns hot-container">
        <div class="column hot-list">
            <component v-for="(hot, idx) in list" :is="cardMap[hot['card_type']]" :item="hot" :idx="idx" :key="idx" @toggle-favor="toggleFavor"></component>
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
import MText from "./cards/MText"
import MRichText from "./cards/MRichText"

const API = {
    config: "/config",
    list: "/list",
};

const CARD_MAP = {
    0: MText.name,
    1: MRichText.name,
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
            t: "",

            cardMap: CARD_MAP
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
        Footer,

        /* eslint-disable vue/no-unused-components */
        MText,
        MRichText,
    }
}
</script>