<template>
    <div class="content-box">
        <div class="columns">
            <div class="column">
                <h4 class="title is-4 has-text-centered">俺的收藏夹</h4>
                <div class="field is-grouped search-form">
                    <p class="control is-expanded">
                        <input class="input" type="text" placeholder="搜一搜" v-model="keyword">
                    </p>
                    <p class="control">
                        <a class="button is-info" @click="search">
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
                        <p class="hot-ts has-text-grey">{{ hot.create_at }}</p>
                        <a :href="hot.origin_url" :title="hot.title" target="_blank">{{ hot.title }}</a>
                    </div>
                    <div class="divider"></div>
                    <div class="hot-opt" @click="remove(idx)">
                        <svg style="width:20px;height:20px" viewBox="0 0 24 24">
                            <path fill="#ff7474" d="M9,3V4H4V6H5V19A2,2 0 0,0 7,21H17A2,2 0 0,0 19,19V6H20V4H15V3H9M7,6H17V19H7V6M9,8V17H11V8H9M13,8V17H15V8H13Z" />
                        </svg>
                    </div>
                </div>
            </div>
        </div>

        <Footer></Footer>
    </div>
</template>

<script>
export default {
    name: "Favor",
    data() {
        return {
            keyword: "",
            selected: {
                tab: 0,
                tag: 0
            },
            list: [],
            tabs: []
        }
    },
    created() {
    },
    methods: {},
    components: {}
}
</script>