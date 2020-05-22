<template>
  <div>
    <ul>
      <li>{{emailObj.messageID}}</li>
      <li>{{emailObj.email.Date}}</li>
      <li>{{emailObj.email.To}}</li>
      <li>{{emailObj.email.From}}</li>
      <li>{{emailObj.email.Subject}}</li>
    </ul>
    <div>
      <pre>{{emailObj.email.TextBody}}</pre>
    </div>
  </div>
</template>

<script>
export default {
  name: "Email",
  mounted() {
    console.log("EMAIL");
  }
};
</script>