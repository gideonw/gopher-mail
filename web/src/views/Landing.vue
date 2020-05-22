<template>
  <div>
    <a to="/i" v-on:click="handleLogin()">Login</a>|
  </div>
</template>

<script>
import { mapActions } from "vuex";

import { LOGIN } from "@/store/actions";

export default {
  methods: {
    handleLogin() {
      this.$store.dispatch(LOGIN);
      this.$router.push("/i");
    },
    ...mapActions([LOGIN])
  }
};
</script>