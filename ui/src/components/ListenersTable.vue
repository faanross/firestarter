<template>
  <div class="counter active-listeners">
    Active Listeners: {{ listeners.length }}
  </div>

  <table>
    <thead>
    <tr>
      <th>CreatedAt</th>
      <th>ID</th>
      <th>Port</th>
      <th>Protocol</th>
      <th>Stop</th>
    </tr>
    </thead>

    <tbody>
    <tr v-if="listeners.length === 0">
      <td colspan="5">No active listeners</td>
    </tr>
    <tr v-for="listener in listeners" :key="listener.id">
      <td>
        <span class="timestamp">{{ formatTimestamp(listener.createdAt) }}</span>
      </td>
      <td>{{ listener.id }}</td>
      <td>{{ listener.port }}</td>
      <td>{{ listener.protocol }}</td>

      <td>
        <button class="btn-stop" @click="stopListener(listener.id)">
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

const listeners = ref([]);

const formatTimestamp = (timestamp) => {
  if (!timestamp) return 'N/A';

  const date = new Date(timestamp);
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });
};

// Process incoming WebSocket messages
const processMessage = (event) => {
  try {
    const message = JSON.parse(event.data);
    console.log('Received message:', message);

    switch (message.type) {
      case 'listener_created':
        // Add new listener to the list
        addListener(message.payload);
        break;

      case 'listener_stopped':
        // Remove listener from the list
        removeListener(message.payload.id);
        break;

      case 'listeners_snapshot':
        // Replace entire list with snapshot (will implement in Phase 2)
        break;

      default:
        console.log('Unknown message type:', message.type);
    }
  } catch (error) {
    console.error('Error processing WebSocket message:', error);
  }
};

// Add a listener to the list
const addListener = (listener) => {
  // Check if listener already exists (by ID)
  const exists = listeners.value.some(l => l.id === listener.id);
  if (!exists) {
    listeners.value.push(listener);
  }
};

// Remove a listener from the list
const removeListener = (id) => {
  listeners.value = listeners.value.filter(listener => listener.id !== id);
};

// Placeholder for the stop function (we'll implement this in Phase 2)
const stopListener = (id) => {
  console.log('Stop listener requested for:', id);
  // TODO: Send stop request to server
};

// Add message listener when socket becomes available
watch(() => props.socket, (newSocket) => {
  if (newSocket) {
    console.log('Socket connected in ListenersTable');
    newSocket.addEventListener('message', processMessage);
  }
}, { immediate: true });

// Clean up on component unmount
onUnmounted(() => {
  if (props.socket) {
    props.socket.removeEventListener('message', processMessage);
  }
});
</script>

<style>

/* Button styling */
.btn-stop {
  background-color: #e74c3c;
  color: white;
  border: none;
  padding: 5px 10px;
  border-radius: 3px;
  cursor: pointer;
}

.btn-stop:hover {
  background-color: #c0392b;
}

.btn-start {
  background-color: #32b253;
  color: white;
  border: none;
  padding: 5px 10px;
  border-radius: 3px;
  cursor: pointer;
}

.btn-start:hover {
  background-color: #32b253;
}

/* Table styling */
table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
  font-size: 0.9rem; /* Smaller font size for the table */
}

th, td {
  border: 1px solid #ddd;
  padding: 6px; /* Slightly reduced padding for more compact display */
  text-align: left;
}

th {
  background-color: #5e5e5e;
  color: white;
}

tr:nth-child(even) {
  background-color: #2a2a2a; /* Alternating row colors for better readability */
}

tr:hover {
  background-color: #3a3a3a; /* Highlight row on hover */
}

/* Counter styling */
.counter {
  margin: 20px 0;
  padding: 10px;
  background-color: #5e5e5e;
  color: white;
  font-weight: bold;
  border-radius: 4px;
}

/* Message when no data is available */
td[colspan="4"] {
  padding: 15px;
  text-align: center;
  color: #aaa;
}

.counter.active-listeners {
  background-color: #8a70d6; /* Purple color similar to the tab we were using */
  color: white;
}

.timestamp {
  font-size: 0.8rem;
  color: #aaa;
  margin-right: 10px;
}


</style>