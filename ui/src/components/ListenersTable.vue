<template>
  <div class="table-container">
    <div class="table-wrapper">
  <table>
    <colgroup>
      <col style="width: 15%"> <!-- CreatedAt -->
      <col style="width: 25%"> <!-- ID -->
      <col style="width: 15%"> <!-- Port -->
      <col style="width: 30%"> <!-- Protocol -->
      <col style="width: 80px"> <!-- Stop - fixed width for this column -->
    </colgroup>
    <thead>
    <tr>
      <th>CreatedAt</th>
      <th>ID</th>
      <th>Port</th>
      <th>Protocol</th>
      <th>🛑</th>
    </tr>
    </thead>

    <tbody>
    <tr v-if="listeners.length === 0">
      <td colspan="5">Listeners: 0</td>
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
          ⬢
        </button>
      </td>
    </tr>
    </tbody>
  </table>
  </div>
  </div>
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
        // Replace entire list with snapshot data
        handleSnapshot(message.payload);
        break;

      default:
        console.log('Unknown message type:', message.type);
    }
  } catch (error) {
    console.error('Error processing WebSocket message:', error);
  }
};

// Process "snapshot" - list of all active listeners
const handleSnapshot = (listenersData) => {
  console.log('Received listeners snapshot with', listenersData.length, 'listeners');

  // Replace the entire listeners array with the snapshot data
  listeners.value = listenersData;
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

// Stop a listener by sending a command to the server
const stopListener = (id) => {
  console.log('Requesting to stop listener:', id);

  if (!props.socket || props.socket.readyState !== WebSocket.OPEN) {
    console.error('Cannot stop listener: WebSocket not connected');
    return;
  }

  // Create and send the stop command
  const stopCommand = {
    action: 'stop_listener',
    payload: { id }
  };

  props.socket.send(JSON.stringify(stopCommand));
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

// Request a snapshot of all listeners from the server
const requestSnapshot = () => {
  console.log('Requesting listeners snapshot');

  if (!props.socket || props.socket.readyState !== WebSocket.OPEN) {
    console.error('Cannot request snapshot: WebSocket not connected');
    return;
  }

  // Create and send the get_listeners command
  const getCommand = {
    action: 'get_listeners',
    payload: {}
  };

  props.socket.send(JSON.stringify(getCommand));
};

// And update the watch function to request a snapshot when the socket connects:
watch(() => props.socket, (newSocket) => {
  if (newSocket) {
    console.log('Socket connected in ListenersTable');
    newSocket.addEventListener('message', processMessage);

    // Request a snapshot when the socket connects
    // (This is a backup in case the automatic snapshot on connection fails)
    setTimeout(requestSnapshot, 500);
  }
}, { immediate: true });

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


table {
    width: 900px;
    table-layout: fixed; /* Prevents resizing based on content */
  }

.table-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
}

.table-wrapper {
  display: flex;
  justify-content: center;
  width: 100%;
}

th, td {
  border: 1px solid #ddd;
  padding: 6px; /* Slightly reduced padding for more compact display */
  text-align: center;
  font-size: 14px;
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


/* Message when no data is available */
td[colspan="5"] {
  padding: 15px;
  text-align: center;
  color: #aaa;
}


.timestamp {
  font-size: 0.8rem;
  color: #aaa;
  margin-right: 10px;
}


</style>