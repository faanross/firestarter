<template>
  <div class="container">
    <div>
      <h1>firestarterC2</h1>
      <WebSocketConnection @socket-ready="handleSocketReady" />
    </div>
    <div>
      <TabsComponent
          :tabs="tabs"
          initialTab="tab1"
          @tab-changed="handleTabChange">

        <template #tab1>
          <ListenersTable :socket="sharedSocket" />
        </template>

        <template #tab2>
          <p>This will contain info on our active connections.</p>
        </template>
      </TabsComponent>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import ListenersTable from './components/ListenersTable.vue';
import WebSocketConnection from './components/WebSocketConnection.vue';
import TabsComponent from './components/TabsComponent.vue';

// Define reactive data directly at the top level
const tabs = [
  { id: 'tab1', name: 'Listeners' },
  { id: 'tab2', name: 'Connections' },
];

const sharedSocket = ref(null);

const handleSocketReady = (socket) => {
  sharedSocket.value = socket;
};
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