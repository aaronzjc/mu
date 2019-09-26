<template>
    <div class="container section">
        <div class="columns">
            <div class="column is-one-fifth">
                <div class="columns">
                    <div class="column">
                        <div class="logo has-text-centered"><img :src="login.avatar" alt=""></div>
                        <div class="has-text-centered">
                            <strong>欢迎 <span class="has-text-info">{{ login.username }}</span></strong>
                        </div>
                    </div>
                </div>
                <Menu></Menu>
            </div>

            <div class="column is-four-fifths content">
                <transition>
                    <router-view></router-view>
                </transition>
            </div>
        </div>
    </div>
</template>

<script>
import 'bulma/css/bulma.css'
import Menu from "./Menu"
import {Get} from "../tools/http";

export default {
    name: "Dashboard",
    created() {
        this.fetchLogin()
    },
    data: () => {
        return {
            login: {
                id: "",
                username: "",
                avatar: ""
            }
        }
    },
    methods: {
        fetchLogin() {
            Get("/info").then(resp => {
                if (resp.data.code === 10000) {
                    this.login = resp.data.data
                } else {
                    alert(resp.data.msg)
                }
            })
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
            width: 50%;
            border-radius: 10rem 10rem;
            height: auto;
        }
    }
    .content {
        padding-top: 3rem;
    }
</style>