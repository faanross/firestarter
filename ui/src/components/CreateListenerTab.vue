<template>
  <div class="create-listener-container">
    <h2>Create New Listener</h2>

    <form @submit.prevent="createListener" class="create-form">
      <!-- ID Field (Optional) -->
      <div class="form-group">
        <label for="listener-id">ID (Optional):</label>
        <input
            type="text"
            id="listener-id"
            v-model="formData.id"
            placeholder="Leave empty for auto-generated ID"
            class="form-input"
        >
      </div>

      <!-- Port Field -->
      <div class="form-group">
        <label for="listener-port">Port:</label>
        <div class="port-input-group">
          <input
              type="number"
              id="listener-port"
              v-model="formData.port"
              placeholder="e.g. 8000"
              class="form-input port-input"
              required
              min="1024"
              max="65535"
          >
          <button
              type="button"
              @click="checkPortAvailability"
              class="check-button"
              :disabled="!formData.port || checkingPort"
          >
            Check
          </button>
          <span class="port-status" v-if="portStatus">
  <span v-if="portStatus === 'available'" class="status-available">✓ Available</span>
  <span v-else-if="portStatus === 'unavailable'" class="status-unavailable">✗ In Use</span>
  <span v-else-if="portStatus === 'error'" class="status-error">! Error</span>
  <span v-else-if="portStatus === 'invalid'" class="status-error">! Invalid</span>
</span>
          <span v-if="checkingPort" class="checking-status">Checking...</span>
        </div>
        <small class="port-hint">Port must be between 1024 and 65535</small>
        <span class="validation-error" v-if="formErrors.port">{{ formErrors.port }}</span>
      </div>

      <!-- Protocol Field -->
      <div class="form-group">
        <label for="listener-protocol">Protocol:</label>
        <select
            id="listener-protocol"
            v-model="formData.protocol"
            class="form-select"
            required
        >
          <option value="" disabled selected>Select a protocol</option>
          <option value="1">HTTP/1.1 Clear (H1C)</option>
          <option value="2">HTTP/1.1 TLS (H1TLS)</option>
          <option value="3">HTTP/2 Clear (H2C)</option>
          <option value="4">HTTP/2 TLS (H2TLS)</option>
          <option value="5">HTTP/3 (H3)</option>
        </select>
        <span class="validation-error" v-if="formErrors.protocol">{{ formErrors.protocol }}</span>
      </div>

      <!-- Submit Button -->
      <div class="form-actions">
        <button
            type="submit"
            class="create-button"
            :disabled="!isFormValid"
        >
          Create Listener
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, defineProps } from 'vue';


const props = defineProps({
  socket: Object
});

// Form data
const formData = ref({
  id: '',
  port: '',
  protocol: ''
});

// Port check status
const portStatus = ref(null);
const checkingPort = ref(false);

// Validate form
const isFormValid = computed(() => {
  const errors = validateForm();
  return !errors.port && !errors.protocol;
});

// Check port availability
const checkPortAvailability = () => {
  if (!formData.value.port) return;

  // Validate port number format before sending
  const portNum = parseInt(formData.value.port);
  if (isNaN(portNum) || portNum < 1024 || portNum > 65535) {
    portStatus.value = 'invalid';
    return;
  }

  portStatus.value = null;
  checkingPort.value = true;

  // Create check_port command
  const checkCommand = {
    action: 'check_port',
    payload: { port: String(formData.value.port) }
  };

  // Ensure socket is connected
  if (!props.socket || props.socket.readyState !== WebSocket.OPEN) {
    console.error('WebSocket not connected');
    checkingPort.value = false;
    portStatus.value = 'error';
    return;
  }

  // Send command to server
  props.socket.send(JSON.stringify(checkCommand));

  // We'll handle the response in a separate function
};

// Add a function to handle WebSocket messages
const processMessage = (event) => {
  try {
    const message = JSON.parse(event.data);
    console.log("CreateListenerTab received message:", message.type);

    // Handle port check result
    if (message.type === 'port_check_result') {
      console.log("CreateListenerTab processing port check result:", message.payload);
      checkingPort.value = false;

      if (message.payload && message.payload.port && message.payload.port === String(formData.value.port)) {
        portStatus.value = message.payload.isAvailable ? 'available' : 'unavailable';
        console.log(`Port ${message.payload.port} availability set to: ${portStatus.value}`);
      } else {
        console.log("Port mismatch or missing payload properties:", {
          messagePort: message.payload?.port,
          formPort: formData.value.port
        });
      }
    }
  } catch (error) {
    console.error('Error processing WebSocket message in CreateListenerTab:', error);
  }
};

// Add socket event listener when component is mounted
onMounted(() => {
  console.log("CreateListenerTab mounted, socket exists:", !!props.socket);
  if (props.socket) {
    props.socket.addEventListener('message', processMessage);
    console.log("WebSocket message listener added in CreateListenerTab");
  }
});

// Clean up when component is unmounted
onUnmounted(() => {
  if (props.socket) {
    props.socket.removeEventListener('message', processMessage);
  }
});

// Create listener
const createListener = () => {
  if (!isFormValid.value) return;

  console.log('Creating listener with:', formData.value);
  // We'll implement the actual creation in a later step

  // Reset form after submission
  formData.value = {
    id: '',
    port: '',
    protocol: ''
  };
  portStatus.value = null;
};

const validateForm = () => {
  const errors = {
    port: null,
    protocol: null
  };

  // Validate port
  if (!formData.value.port) {
    errors.port = 'Port is required';
  } else if (formData.value.port < 1024 || formData.value.port > 65535) {
    errors.port = 'Port must be between 1024 and 65535';
  } else if (portStatus.value === 'unavailable') {
    errors.port = 'This port is already in use';
  }

  // Validate protocol
  if (!formData.value.protocol) {
    errors.protocol = 'Protocol is required';
  }

  return errors;
};

// Checks for any validation errors
const formErrors = computed(() => {
  return validateForm();
});


</script>

<style scoped>
.create-listener-container {
  max-width: 600px;
  margin: 0 auto;
  padding: 20px;
}

h2 {
  margin-bottom: 20px;
  color: #f8f8f2;
}

.create-form {
  background-color: #2c2d30;
  padding: 20px;
  border-radius: 5px;
}

.form-group {
  margin-bottom: 15px;
}

label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
}

.form-input, .form-select {
  width: 100%;
  padding: 8px;
  border: 1px solid #444;
  border-radius: 4px;
  background-color: #383838;
  color: #f8f8f2;
}

.port-input-group {
  display: flex;
  align-items: center;
}

.port-input {
  flex: 1;
}

.port-hint {
  display: block;
  font-size: 0.8rem;
  color: #aaa;
  margin-top: 3px;
}

.check-button {
  padding: 8px 12px;
  margin-left: 8px;
  background-color: #6272a4;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
}

.check-button:hover {
  background-color: #7382b4;
}

.check-button:disabled {
  background-color: #4a4a4a;
  cursor: not-allowed;
}

.port-status {
  margin-left: 10px;
  font-size: 1.2rem;
  font-weight: bold;
}

.form-actions {
  margin-top: 20px;
}

.create-button {
  padding: 10px 20px;
  background-color: #50fa7b;
  color: #282a36;
  border: none;
  border-radius: 4px;
  font-weight: bold;
  cursor: pointer;
}

.create-button:hover {
  background-color: #60fa8b;
}

.create-button:disabled {
  background-color: #4a4a4a;
  color: #999;
  cursor: not-allowed;
}

.validation-error {
  color: #ff5555;
  font-size: 0.8rem;
  display: block;
  margin-top: 4px;
}

.status-available {
  color: #50fa7b;
}

.status-unavailable {
  color: #ff5555;
}

.status-error {
  color: #ff79c6;
}

.checking-status {
  color: #f1fa8c;
  margin-left: 10px;
  font-style: italic;
}

</style>