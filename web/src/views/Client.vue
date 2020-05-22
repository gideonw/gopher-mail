<template>
  <div>
    <MailboxList />
    <router-view />
  </div>
</template>

<script>
// @ is an alias to /src
import MailboxList from "@/components/MailboxList.vue";

export default {
  name: "Client",
  components: {
    MailboxList
  }
};
</script>
