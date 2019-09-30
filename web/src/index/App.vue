<template>
  <section id="app" class="section">
      <Main></Main>
  </section>
</template>

<script>
import Main from './components/Main'

export default {
  name: 'app',
  components: {
    Main
  }
}
</script>

<style>
.section {
    padding-top: 1rem;
}
* {
    -webkit-tap-highlight-color: transparent;
}
</style>
