<template>
  <div>{{ mailbox.name }}</div>
</template>

<script>
export default {
  name: "MailboxListItem",
  props: {
    mailbox: Object
  }
};
</script>