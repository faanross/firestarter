<template>
  <div class="counter active-connections">
    Active Connections: {{ connections.length }}
  </div>

  <table>
    <thead>
    <tr>
      <th>CreatedAt</th>
      <th>ID</th>
      <th>Agent UUID</th>
      <th>Remote Address</th>
      <th>Port</th>
      <th>Protocol</th>
      <th>Stop</th>
    </tr>
    </thead>

    <tbody>
    <tr v-if="connections.length === 0">
      <td colspan="7">No active connections</td>
    </tr>
    <tr v-for="connection in connections" :key="connection.id">
      <td>
        <span class="timestamp">{{ formatTimestamp(connection.createdAt) }}</span>
      </td>
      <td>{{ connection.id }}</td>
      <td>{{ truncateUUID(connection.agentUUID) }}</td>
      <td>{{ connection.remoteAddr }}</td>
      <td>{{ connection.port }}</td>
      <td>{{ connection.protocol }}</td>
      <td>
        <button class="btn-stop" @click="stopConnection(connection.id)">
          Stop
        </button>
      </td>
    </tr>
    </tbody>
  </table>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, defineProps } from 'vue';

const props = defineProps({
  socket: Object
});

const connections = ref([]);

// Helper functions
const formatTimestamp = (timestamp) => {
  if (!timestamp) return 'N/A';
  const date = new Date(timestamp);
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });
};

const truncateUUID = (uuid) => {
  if (!uuid) return 'N/A';
  // Show first 8 characters of UUID for brevity
  return uuid.substring(0, 8) + '...';
};

// WebSocket message handling
const processMessage = (event) => {
  try {
    const message = JSON.parse(event.data);
    console.log('Received message:', message);

    switch (message.type) {
      case 'connection_created':
        // Add new connection to the list
        addConnection(message.payload);
        break;

      case 'connection_stopped':
        // Remove connection from the list
        removeConnection(message.payload.id);
        break;

      case 'connections_snapshot':
        // Replace entire list with snapshot data
        handleSnapshot(message.payload);
        break;
    }
  } catch (error) {
    console.error('Error processing WebSocket message:', error);
  }
};

// Process "snapshot" - list of all active connections
const handleSnapshot = (connectionsData) => {
  console.log('Received connections snapshot with', connectionsData.length, 'connections');
  // Replace the entire connections array with the snapshot data
  connections.value = connectionsData;
};

// Add a connection to the list
const addConnection = (connection) => {
  // Check if connection already exists (by ID)
  const exists = connections.value.some(c => c.id === connection.id);
  if (!exists) {
    connections.value.push(connection);
  }
};

// Remove a connection from the list
const removeConnection = (id) => {
  connections.value = connections.value.filter(connection => connection.id !== id);
};

// Stop a connection by sending a command to the server
const stopConnection = (id) => {
  console.log('Requesting to stop connection:', id);

  if (!props.socket || props.socket.readyState !== WebSocket.OPEN) {
    console.error('Cannot stop connection: WebSocket not connected');
    return;
  }

  // Create and send the stop command
  const stopCommand = {
    action: 'stop_connection',
    payload: { id }
  };

  props.socket.send(JSON.stringify(stopCommand));
};

// Request a snapshot of all connections from the server
const requestSnapshot = () => {
  console.log('Requesting connections snapshot');

  if (!props.socket || props.socket.readyState !== WebSocket.OPEN) {
    console.error('Cannot request snapshot: WebSocket not connected');
    return;
  }

  // Create and send the get_connections command
  const getCommand = {
    action: 'get_connections',
    payload: {}
  };

  props.socket.send(JSON.stringify(getCommand));
};

// Add message listener when socket becomes available
watch(() => props.socket, (newSocket) => {
  if (newSocket) {
    console.log('Socket connected in ConnectionsTable');
    newSocket.addEventListener('message', processMessage);

    // Request a snapshot when the socket connects
    setTimeout(requestSnapshot, 500);
  }
}, { immediate: true });

// Clean up on component unmount
onUnmounted(() => {
  if (props.socket) {
    props.socket.removeEventListener('message', processMessage);
  }
});
</script>

<style scoped>
/* You can reuse the same styles as ListenersTable or customize as needed */
.counter.active-connections {
  background-color: #3498db; /* Blue color to differentiate from listeners */
  color: white;
  margin: 20px 0;
  padding: 10px;
  font-weight: bold;
  border-radius: 4px;
}

/* Reuse other styles from ListenersTable.vue */
</style>