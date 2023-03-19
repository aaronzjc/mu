<template>
    <li :class="{ 'is-active': _menu.active }">
        <MenuLink
            :title="_menu.title"
            :active="_menu.active"
            :icon="_menu.icon"
            :dropdown="!!_menu.children"
            @click="menuClick(_menu)"
        >
        </MenuLink>

        <ul v-show="_menu.children && _menu.active">
            <li v-for="subMenu in _menu.children">
                <MenuLink
                    :title="subMenu.title"
                    :active="subMenu.active"
                    @click="menuClick(subMenu)"
                >
                </MenuLink>
            </li>
        </ul>
    </li>
</template>

<script setup>
import { watch, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import MenuLink from '@adm/components/MenuLink.vue'

const props = defineProps({
    menu: {
        type: Object,
        required: true
    }
})

const _menu = reactive(props.menu)

const router = useRouter()
const route = useRoute()

function menuClick(menu) {
    if (menu.route) {
        if (menu.route != route.name) {
            router.push({ name: menu.route })
        }
        return
    }
    menu.active = !menu.active
}

watch(
    () => route.name,
    to => {
        if (_menu.route) {
            _menu.active = _menu.route == to || _menu.route == route.meta.hl
            return
        }
        let got = false
        for (let i = 0; i < _menu.children.length; i++) {
            const sub = _menu.children[i]
            sub.active = sub.route == to || sub.route == route.meta.hl
            if (sub.active) {
                got = true
            }
        }
        _menu.active = got
    },
    { immediate: true }
)
</script>
