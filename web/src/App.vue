<template>
  <div id="app">
    <Header />
    <div id="client-body">
      <router-view />
    </div>
  </div>
</template>

<script>
import Header from "./components/Header.vue";

export default {
  name: "App",
  components: {
    Header
  }
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  display: flex;
  align-items: flex-start;
  justify-content: flex-start;
  flex-direction: column;
  flex: 1;
}
#client-body {
  justify-content: flex-start;
  flex-direction: row;
  display: flex;
}
</style>
