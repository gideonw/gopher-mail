<template>
  <div>{{ key }} - {{ email.Subject }}</div>
</template>

<script>
export default {
  name: "EmailListItem",
  props: {
    key: String,
    email: String
  }
};
</script>