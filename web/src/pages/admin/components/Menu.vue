<template>
    <nav id="navbar-main" class="navbar is-fixed-top">
        <div class="navbar-brand">
            <a
                class="navbar-item is-hidden-desktop jb-aside-mobile-toggle"
                @click="styleStore.toggleAside()"
            >
                <BasicIcon
                    :name="
                        styleStore.isAsideMobileOpen
                            ? mdiBackburger
                            : mdiForwardburger
                    "
                    :size="24"
                ></BasicIcon>
            </a>
            <div class="navbar-item has-control">
                <div class="control"></div>
            </div>
        </div>
        <div class="navbar-brand is-right">
            <a
                class="navbar-item is-hidden-desktop jb-navbar-menu-toggle"
                data-target="navbar-menu"
                @click="styleStore.toggleNav()"
            >
                <BasicIcon :name="mdiDotsVertical"></BasicIcon>
            </a>
        </div>
        <div
            :class="[
                'navbar-menu fadeIn animated faster',
                { 'is-active': styleStore.isNavMobileOpen }
            ]"
            id="navbar-menu"
        >
            <div class="navbar-end">
                <div
                    class="navbar-item has-dropdown has-dropdown-with-icons has-divider has-user-avatar is-hoverable"
                >
                    <a class="navbar-link is-arrowless">
                        <div class="is-user-avatar">
                            <img
                                :src="mainStore.userInfo.avatar"
                                :alt="mainStore.userInfo.username"
                            />
                        </div>
                        <div class="is-user-name">
                            <span>{{ mainStore.userInfo.username }}</span>
                        </div>
                    </a>
                </div>
                <a
                    href="https://github.com/aaronzjc"
                    title="Github"
                    class="navbar-item has-divider is-desktop-icon-only"
                >
                    <BasicIcon :name="mdiGithub"></BasicIcon>
                    <span>Github</span>
                </a>
                <a title="退出登录" class="navbar-item is-desktop-icon-only" @click="logout">
                    <BasicIcon :name="mdiLogout"></BasicIcon>
                    <span>退出登录</span>
                </a>
            </div>
        </div>
    </nav>

    <aside class="aside is-placed-left is-expanded">
        <div class="aside-tools">
            <div class="aside-tools-label">
                <span>后台管理</span>
            </div>
        </div>
        <div class="menu is-menu-main">
            <p class="menu-label">全部</p>
            <ul class="menu-list">
                <MenuItem :menu="menu" v-for="menu in menus"></MenuItem>
            </ul>
        </div>
    </aside>
</template>

<script setup>
import { menus } from '@adm/config.js'
import {
    mdiForwardburger,
    mdiBackburger,
    mdiLogout,
    mdiDotsVertical,
    mdiChevronDown,
    mdiGithub,
mdiAccount
} from '@mdi/js'
import BasicIcon from '@adm/components/BasicIcon.vue'
import { useStyleStore } from '@adm/stores/style.js'
import MenuItem from '@adm/components/MenuItem.vue'
import { useRouter } from 'vue-router'
import { useMainStore } from '@adm/stores/main'

const styleStore = useStyleStore()

const router = useRouter()
router.beforeEach(() => {
    styleStore.isAsideMobileOpen = false
    styleStore.isNavMobileOpen = false
})

const mainStore = useMainStore()

function logout() {
    localStorage.removeItem(import.meta.env.VITE_TOKEN_KEY)
    router.push({name: 'login'})
}
</script>
