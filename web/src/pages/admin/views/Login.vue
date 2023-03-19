<template>
    <section class="section hero login">
        <div class="columns is-centered">
            <div
                class="column is-two-thirds-tablet is-half-desktop is-4-widescreen"
            >
                <div class="card">
                    <header class="card-header">
                        <p class="card-header-title is-centered">
                            <span>后台管理登录</span>
                        </p>
                    </header>
                    <div class="card-content">
                        <form action="javascript:void(0);">
                            <FormField group center>
                                <div class="buttons">
                                    <a class="button is-dark" v-for="(item, idx) in state.auth" :href="item.url" :title="item.name">
                                        <BasicIcon
                                            :name="mdiGithub"
                                            :size="48"
                                            v-if="item.type === 'github'"
                                        ></BasicIcon>
                                        <BasicIcon
                                            :name="mdiSinaWeibo"
                                            :size="48"
                                            v-if="item.type === 'weibo'"
                                        ></BasicIcon>
                                    </a>
                                </div>
                            </FormField>
                        </form>
                    </div>
                </div>
            </div>
        </div>
    </section>
</template>

<script setup>
import FormField from '@adm/components/FormField.vue'
import FormControl from '@adm/components/FormControl.vue'
import BasicIcon from '@adm/components/BasicIcon.vue'
import { mdiAccount, mdiGithub, mdiGoogle, mdiLock, mdiSinaWeibo } from '@mdi/js'

import {Get} from "@/lib/http";
import { onMounted, reactive } from 'vue';

const API = {
    config: "/auth/config"
}

const state = reactive({
    auth: []
})

async function fetchConfig() {
    try {
        let resp = await Get(API.config, {from: "admin"})
        if (resp.data.code === 10000) {
            state.auth = resp.data.data;
        } else {
            console.log(resp.data.msg);
        }
    } catch(err) {
        console.log(err)
    }
}

onMounted(fetchConfig)

</script>
