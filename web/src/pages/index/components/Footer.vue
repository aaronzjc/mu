<template>
<div class="columns footer">
    <div class="column copyright has-text-centered">
        <p class="user-declare">站点仅供学习交流使用，如有侵权，请联系下线</p>
        <p><a href="https://github.com/aaronzjc">@aaronzjc</a>开发，源码<a href="https://github.com/aaronzjc/mu">在此</a>，欢迎Star v{{ version }}</p>
        <p class="backtop"><a href="javascript:scrollTo(0,0);">回到顶部</a></p>
    </div>
</div>
</template>

<script>
import { readonly } from "vue"
export default {
    name: "footer",
    setup() {
        const version = readonly(process.env.VUE_APP_VERSION)
        return {
            version
        }
    }
}
</script>