<template>
<div class="container section">
    <div class="columns">
        <div class="column is-one-fifth">
            <div class="columns">
                <div class="column">
                    <div class="logo has-text-centered"><img :src="state.login.avatar" alt=""></div>
                    <div class="has-text-centered">
                        <strong>欢迎 <span class="has-text-info">{{ state.login.username }}</span></strong>
                    </div>
                </div>
            </div>
            <Menu></Menu>
        </div>

        <div class="column is-four-fifths content">
            <router-view></router-view>
        </div>
    </div>
</div>
</template>

<script>
import 'bulma/css/bulma.css'
import Menu from "./Menu"
import {Get} from "@/tools/http";
import { onBeforeMount, reactive } from 'vue';

const API = {
    loginInfo: "/auth/info/admin"
}

export default {
    name: 'Dashboard',
    setup() {
        const state = reactive({
            login: {
                id: "",
                username: "",
                avatar: ""
            }
        })
        async function fetchLogin() {
            try {
                let resp = await Get(API.loginInfo)
                if (resp.data.code === 10000) {
                    state.login = resp.data.data
                } else {
                    console.log(resp.data.msg)
                }
            } catch(err) {
                console.log(err)
            }
        }
        onBeforeMount(fetchLogin)
        return {
            state
        }
    },
    components: {
        Menu
    }
}
</script>

<style lang="scss">
.logo{
    margin: 0 0 1rem 0;
    img {
        display: inline-block;
        width: 5rem;
        border-radius: 10rem 10rem;
        height: auto;
    }
}
.content {
    padding-top: 3rem;
}
</style>