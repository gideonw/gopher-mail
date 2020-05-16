<template>
  <div>
    <div id="mailbox-list" v-if="!loading">
      <div
        is="MailboxListItem"
        v-for="mailbox in mailboxes"
        v-bind:key="mailbox.ID"
        v-bind:mailbox="mailbox"
      ></div>
    </div>
  </div>
</template>

<script>
import MailboxListItem from "./MailboxListItem.vue";

export default {
  name: "MailboxList",
  components: { MailboxListItem },
  data() {
    return {
      loading: true,
      mailboxes: [
        {
          ID: "00001",
          name: "Hello"
        },
        {
          ID: "00002",
          name: "Hello"
        },
        {
          ID: "00003",
          name: "Hello"
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

<style scoped>
#mailbox-list {
  min-width: 80px;
  margin-left: 8px;
  margin-right: 8px;
}
</style>