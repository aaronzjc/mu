<template>
    <div class="columns footer">
        <div class="column copyright has-text-centered">
            <p class="user-declare">站点仅供学习交流使用，如有侵权，请联系下线</p>
            <p>
                <a href="https://github.com/aaronzjc">@aaronzjc</a>
                开发，源码<a href="https://github.com/aaronzjc/mu">在此</a>，欢迎Star v{{ version }}
                <span v-if="state.count != ''">当前在线 <strong class="online">{{ state.count }}</strong> 人</span>
            </p>
            <p class="backtop"><a href="javascript:scrollTo(0,0);">回到顶部</a></p>
        </div>
    </div>
</template>

<script setup>
import { readonly, reactive, onMounted } from "vue"
import { Get } from "@/lib/http"

const state = reactive({
    count: "",
})

async function GetOnline() {
    let resp = await Get("/stat/online")
    state.count = resp.data.data.count
}
setInterval(GetOnline, 30 * 1000)

const version = readonly(import.meta.env.VITE_APP_VERSION)
onMounted(GetOnline)
</script>