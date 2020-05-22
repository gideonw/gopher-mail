<template>
  <div>
    <!-- <div v-if="!loading"> -->
    <div is="EmailListItem" v-for="email in emails" v-bind:key="email" v-bind:email="email"></div>
    <!-- </div> -->
  </div>
</template>

<script>
import EmailListItem from "./EmailListItem.vue";
import { mapState } from "vuex";

export default {
  name: "EmailList",
  components: { EmailListItem },
  computed: mapState(["emails"]),
  mounted() {
    console.log("EMAIL LIST");
  }
};
</script>