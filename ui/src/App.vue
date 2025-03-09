<template>
  <div class="container">
    <h1>firestarterC2</h1>

    <WebSocketConnection/>
    <ListenersTable/>
  </div>
</template>

<script setup>

import ListenersTable from './components/ListenersTable.vue';
import WebSocketConnection from './components/WebSocketConnection.vue';


</script>

<style>
.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

h1 {
  text-align: center;
  margin-bottom: 30px;
}

</style>