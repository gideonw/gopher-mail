<template>
  <div>
    <div v-if="!loading">
      <div
        is="EmailListItem"
        v-for="email in emails"
        v-bind:key="email.messageID"
        v-bind:email="email"
      ></div>
    </div>
  </div>
</template>

<script>
import EmailListItem from "./EmailListItem.vue";

export default {
  name: "EmailList",
  components: { EmailListItem },
  data() {
    return {
      loading: true,
      emails: [
        {
          messageID: "0001",
          subject: "Hello"
        },
        {
          messageID: "0002",
          subject: "Hello"
        },
        {
          messageID: "0003",
          subject: "Hello"
        }
      ]
    };
  },
  created() {
    window.setTimeout(
      function() {
        this.loading = false;
      }.bind(this),
      2000
    );
  }
};
</script>