<template>
  <div>
    <!-- Connection Status -->
    <div class="status" :class="{ connected: isConnected }">
      Server: {{ connectionStatus }}
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, defineEmits } from 'vue';

const emits = defineEmits(['socket-ready']);

// WebSocket connection
const socket = ref(null);
const isConnected = ref(false);
const connectionStatus = ref('Disconnected');

// Connect to WebSocket server
const connectWebSocket = () => {
  // Close existing connection if any
  if (socket.value) {
    socket.value.close();
  }

  // Create new WebSocket connection
  const wsUrl = 'ws://localhost:8080/ws';
  socket.value = new WebSocket(wsUrl);

  // Connection opened
  socket.value.addEventListener('open', (event) => {
    console.log('Connected to WebSocket server');
    isConnected.value = true;
    connectionStatus.value = 'Connected';

    // Emit the socket to parent component
    emits('socket-ready', socket.value);

    // Send a message to the server
    socket.value.send('Hello from Vue client!');
  });

  // Connection closed
  socket.value.addEventListener('close', (event) => {
    console.log('Disconnected from WebSocket server');
    isConnected.value = false;
    connectionStatus.value = 'Disconnected';
  });

  // Connection error
  socket.value.addEventListener('error', (event) => {
    console.error('WebSocket error:', event);
    connectionStatus.value = 'Error';
  });
};

// Connect on component mount
onMounted(() => {
  connectWebSocket();
});

// Clean up on component unmount
onUnmounted(() => {
  if (socket.value) {
    socket.value.close();
  }
});
</script>

<style scoped>
.status {
  margin: 20px 0;
  padding: 10px;
  background-color: #5e5e5e;
}
.status.connected {
  background-color: #32b253;
}
</style>