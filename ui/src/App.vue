<template>
  <div class="app-container">
    <!-- Toolbar Area -->
    <div class="toolbar">
      <div class="connection-indicator">
        <div class="status-circle" :class="{ connected: isConnected }"></div>
      </div>
      <div class="app-title">firestarter</div>
    </div>

    <!-- Main Content Area -->
    <div class="main-content">
      <TabsComponent
          :tabs="tabs"
          initialTab="tab1"
          @tab-changed="handleTabChange">

        <template #tab1>
          <ListenersTable :socket="sharedSocket" />
        </template>

        <template #tab2>
          <ConnectionsTable :socket="sharedSocket" />
        </template>
      </TabsComponent>
    </div>

    <!-- Hidden WebSocketConnection component for functionality -->
    <WebSocketConnection @socket-ready="handleSocketReady" />
  </div>
</template>
<script setup>
import { ref, computed } from 'vue';
import ListenersTable from './components/ListenersTable.vue';
import WebSocketConnection from './components/WebSocketConnection.vue';
import TabsComponent from './components/TabsComponent.vue';
import ConnectionsTable from './components/ConnectionsTable.vue';

// Define reactive data directly at the top level
const tabs = [
  { id: 'tab1', name: 'Listeners' },
  { id: 'tab2', name: 'Connections' },
];

const sharedSocket = ref(null);
const isConnected = ref(false);

const handleSocketReady = (socket) => {
  sharedSocket.value = socket;
  isConnected.value = true;

  // Listen for connection close events
  socket.addEventListener('close', () => {
    isConnected.value = false;
  });

  socket.addEventListener('error', () => {
    isConnected.value = false;
  });
};

const handleTabChange = (tabId) => {
  // Handle tab changes if needed
};
</script>

<style>
.app-container {
  width: 100%;
  min-height: 100vh;
  position: relative; /* For absolute positioning children */
  padding-top: 60px; /* Add space for the fixed toolbar */
}

.toolbar {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 50px;
  background-color: #232325;
  z-index: 1000; /* Ensures toolbar stays on top */
  display: flex;
  justify-content: flex-end;
  align-items: center;
}

.app-title {
  position: absolute;
  right: 20px; /* Fixed distance from right edge */
}

.connection-indicator {
  position: absolute;
  right: 125px; /* Adjust based on your preference */
}

/* Connection status indicator */
.connection-indicator {
  display: flex;
  align-items: center;
  margin-right: 10px;
}

.status-circle {
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background-color: #e74c3c; /* Red when disconnected */
  transition: background-color 0.3s ease;
}

.status-circle.connected {
  background-color: #32b253; /* Green when connected */
}

/* App title styling */
.app-title {
  font-size: 16px;
  font-weight: bold;
  color: #f1fa8c; /* Green text color */
}

/* Main content area */
.main-content {
  position: relative;
  width: 1000px; /* Fixed width as mentioned before */
  margin: 0 auto;
  padding: 20px;
}

</style>