<template>
  <div class="tabs-container">
    <!-- Tab Headers -->
    <div class="tabs">
      <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="setActiveTab(tab.id)"
          :class="{ active: activeTab === tab.id }"
          class="tab-button">
        {{ tab.name }}
      </button>
    </div>

    <!-- Tab Content -->
    <div class="tab-content">
      <!-- We loop through tabs and conditionally render each content panel -->
      <div
          v-for="tab in tabs"
          :key="`content-${tab.id}`"
          v-show="activeTab === tab.id"
          class="tab-pane">
        <!-- Use named slots to render dynamic content for each tab -->
        <slot :name="tab.id">
          <!-- Fallback content if no slot is provided -->
          <h2>{{ tab.name }}</h2>
          <p>Content for {{ tab.name }}</p>
        </slot>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'TabsComponent',
  props: {
    // Allow parent to pass in tab configuration
    tabs: {
      type: Array,
      required: true,
      // Each tab object should have at least id and name
      validator: (tabs) => tabs.every(tab => tab.id && tab.name)
    },
    // Allow parent to set which tab is initially active
    initialTab: {
      type: String,
      default: null
    }
  },
  data() {
    return {
      // Set active tab, defaulting to first tab if initialTab not provided
      activeTab: this.initialTab || (this.tabs.length > 0 ? this.tabs[0].id : null)
    }
  },
  methods: {
    setActiveTab(tabId) {
      this.activeTab = tabId
      // Emit an event so parent component can react to tab changes if needed
      this.$emit('tab-changed', tabId)
    }
  }
}
</script>

<style scoped>
.tabs {
  display: flex;
  border-bottom: 1px solid #ccc;
}

.tab-button {
  padding: 10px 20px;
  background: none;
  border: none;
  cursor: pointer;
  margin-right: 5px;
  border-radius: 5px 5px 0 0;
  font-size: 14px;
}

.tab-button.active {
  background: #50fa7b;
  border: 1px solid #ccc;
  border-bottom: none;
  color: #2C2D30;
}

.tab-content {
  padding: 20px;
  border: 1px solid #ccc;
  border-top: none;
}

.tab-pane {
  animation: fadeIn 0.5s;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
</style>