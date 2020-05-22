<template>
  <div>
    <!-- <div v-if="!loading"> -->
    <div is="EmailListItem" v-for="email in emails" v-bind:key="email" v-bind:email="email"></div>
    <!-- </div> -->
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
          messageID: ""
        }
      ]
    };
  }
};
</script>