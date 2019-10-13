<template>
    <section class="section" id="app">
        <router-view></router-view>
    </section>
</template>

<script>
export default {
  name: 'app',
  components: {}
}
</script>

<style>
section.section {
    padding:5px 0px 0px 0px;
}
* {
    -webkit-tap-highlight-color: transparent;
}
</style>
