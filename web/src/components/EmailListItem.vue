<template>
  <div>{{ emailObj.messageID }} - {{emailObj.email.Subject}}</div>
</template>

<script>
import axios from "axios";

export default {
  name: "EmailListItem",
  props: {
    email: String
  },
  data() {
    return {
      emailObj: {}
    };
  },
  mounted() {
    axios.get(`/api/gideon/email/${this.email}`).then(resp => {
      console.log(resp);
      this.emailObj = resp.data;
    });
  }
};
</script>