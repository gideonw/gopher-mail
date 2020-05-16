<template>
  <div id="header">
    <div class="greeting">Gopher Mail!</div>
  </div>
</template>
<script>
export default {
  name: "Header"
};
</script>

<style scoped>
#header {
  border: 1px;
  margin: 0;
  width: 100%;
  height: 60px;
  align-content: flex-start;
  display: flex;
  background-color: cadetblue;
  margin-bottom: 10px;
}

.greeting {
  color: #a5eeea;
  font-family: "Inconsolata", monospace;
  font-size: 40px;
  padding: 8px;
}
</style>