<template>
  <!-- <div v-on:click="openEmail()">{{ messageID }} - {{ email.Subject }}</div> -->
  <router-link :to="`/i/${this.messageID}`">{{ messageID }} - {{ email.Subject }}</router-link>
</template>

<script>
export default {
  name: "EmailListItem",
  props: {
    messageID: String,
    email: Object
  }
};
</script>