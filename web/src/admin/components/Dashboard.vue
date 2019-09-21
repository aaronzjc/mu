<template>
    <div class="container">
        <div class="columns">
            <div class="column is-one-fifth">
                <Menu></Menu>
            </div>

            <div class="column is-four-fifths">
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

export default {
    name: "Dashboard",

    data: () => {
        return {}
    },
    components: {
        Menu
    }
}
</script>