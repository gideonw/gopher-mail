<template>
  <div id="header">
    <div class="greeting">Gopher Mail!</div>
    <div v-if="loggedIn" v-on:click="handleLoggout()" class="nav">Sign Out!</div>
  </div>
</template>
<script>
import { mapActions } from "vuex";

import { LOGOUT } from "@/store/actions";

export default {
  name: "Header",
  methods: {
    handleLoggout() {
      this.$store.dispatch(LOGOUT);
      this.$router.push("/");
    },
    ...mapActions([LOGOUT])
  },
  computed: {
    loggedIn() {
      return this.$store.state.auth.valid;
    }
  }
};
</script>

<style scoped>
#header {
  border: 1px;
  margin: 0;
  width: 100%;
  height: 60px;
  align-content: flex-start;
  justify-content: space-between;
  display: flex;
  background-color: cadetblue;
  margin-bottom: 10px;
}

.greeting {
  color: #a5eeea;
  font-family: "Inconsolata", monospace;
  font-size: 40px;
  padding: 8px;
}

.nav {
  padding: 8px;
}
</style>