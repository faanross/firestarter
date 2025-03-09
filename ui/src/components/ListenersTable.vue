<!-- ListenersTable.vue -->
<template>

  <!-- Listeners Table -->
  <table>
    <thead>
    <tr>
      <th>Start</th>
      <th>ID</th>
      <th>Port</th>
      <th>Protocol</th>
    </tr>
    </thead>
    <tbody>
    <tr>
      <td></td>
      <td></td>
      <td></td>
      <td></td>
    </tr>
    </tbody>
  </table>

</template>

<script setup>

</script>

<style>
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
</style>