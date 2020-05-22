<template>
  <div>
    <ul>
      <li>{{$route.params.id}}</li>
      <li>{{email.Date}}</li>
      <li>{{email.To}}</li>
      <li>{{email.From}}</li>
      <li>{{email.Subject}}</li>
    </ul>
    <div>
      <pre>{{email.TextBody}}</pre>
    </div>
  </div>
</template>

<script>
// import { mapGetters } from "vuex";
export default {
  name: "Email",
  data() {
    return {
      email: {}
    };
  },
  mounted() {
    console.log("EMAIL");
    console.log(this.$store);
  },
  // computed: {

  // },
  watch: {
    route() {
      this.updateEmail(this.$store.getters.getEmailByID(this.$route.params.id));
    }
  },
  methods: {
    updateEmail(newEmail) {
      this.email = newEmail;
    }
  }
};
</script>