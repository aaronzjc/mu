<template>
  <div class="content-box">
    <HoTab @change="switchTab" :tabs="state.tabs"></HoTab>

    <p class="hot-ts" v-if="state.t !== ''">更新时间: {{ state.t }}</p>
    <div class="columns hot-container">
      <div class="column hot-list">
        <component v-for="(hot, idx) in state.list" :is="CardMap[hot['card_type']]" :item="hot" :idx="idx" :key="idx">
        </component>
      </div>
    </div>

    <Footer></Footer>
  </div>
</template>

<script setup>
import { computed, provide, onMounted, reactive } from "vue";
import HoTab from "./HoTab.vue";
import Footer from "./Footer.vue";

import { Get } from "@/lib/http";
import NProgress from "nprogress";

import MText from "./cards/MText.vue";
import MRichText from "./cards/MRichText.vue";
import MVideo from "./cards/MVideo.vue";
const CardMap = {
  0: MText,
  1: MRichText,
  2: MVideo,
};

const API = {
  config: "/api/sites",
  list: "/api/news",
};

const state = reactive({
  tabs: [
    {
      name: "Moo-Yuu",
      key: "新闻",
      tags: [
        { key: "昨天", name: "昨天", enable: 1 },
        { key: "今天", name: "今天", enable: 1 },
        { key: "明天", name: "明天", enable: 1 },
      ],
    },
  ],
  selected: {
    tab: 0,
    tag: 0,
  },
  list: [],
  t: "还没更新呢",
});

async function fetchConfig(callback) {
  let resp = await Get(API.config);
  if (resp.data.code === 10000) {
    state.tabs = Object.freeze(resp.data.data);
  } else {
    alert(resp.data.msg);
  }
  if (typeof callback == "function") {
    callback();
  }
}

// 获取卡片列表
async function fetchList() {
  if (state.tabs.length === 0) {
    return false;
  }
  let key = state.tabs[state.selected.tab]["key"];
  let hkey = undefined;
  if (state.tabs[state.selected.tab]["tags"].length > 0) {
    hkey = state.tabs[state.selected.tab]["tags"][state.selected.tag]["key"];
  }
  if (hkey === undefined || key === undefined) {
    return false;
  }
  NProgress.start();
  let resp = await Get(API.list, {
    key: state.tabs[state.selected.tab]["key"],
    hkey: state.tabs[state.selected.tab]["tags"][state.selected.tag]["key"],
  });
  if (resp.data.code === 10000) {
    state.list = Object.freeze(resp.data.data.list);
    state.t = resp.data.data.t;
  } else {
    state.list = [];
  }
  NProgress.done();
}

let switchTab = (data) => {
  state.selected = data;
  fetchList();

  window.scrollTo({
    top: 0,
    behavior: "smooth",
  });
};

onMounted(() => {
  fetchConfig(fetchList);
});

provide("updateMark", (idx, res) => {
  state.list[idx]["mark"] = res;
});
provide(
  "currentSite",
  computed(() => state.tabs[state.selected.tab]["key"])
);
</script>
