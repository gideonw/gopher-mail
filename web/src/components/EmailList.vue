<template>
  <div>
    <div v-if="!emailLoading" class="list">
      <div
        is="EmailListItem"
        v-for="(email, messageID) in emails"
        :key="messageID"
        :messageID="messageID"
        :email="email"
      ></div>
    </div>
  </div>
</template>

<script>
import EmailListItem from "./EmailListItem.vue";
import { mapState } from "vuex";

export default {
  name: "EmailList",
  components: { EmailListItem },
  computed: mapState(["emails", "emailLoading"]),
  watch: {
    $route() {
      // react to route changes...
    }
  },
  mounted() {
    console.log("EMAIL LIST");
  }
};
</script>

<style scoped>
.list {
  flex: 1;
  display: flex;
  flex-direction: column;
}
</style>