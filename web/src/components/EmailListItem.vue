<template>
  <div>{{ emailObj.messageID }} - {{emailObj.email.Subject}}</div>
</template>

<script>
export default {
  name: "EmailListItem",
  props: {
    email: String
  },
  data() {
    return {
      emailObj: {}
    };
  }
};
</script>