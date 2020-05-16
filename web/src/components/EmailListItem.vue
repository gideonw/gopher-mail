<template>
  <div>{{ email.subject }}</div>
</template>

<script>
export default {
  name: "EmailListItem",
  props: {
    email: Object
  }
};
</script>