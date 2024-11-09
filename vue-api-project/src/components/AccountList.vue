<template>
  <div>
    <h2>Account Summary</h2>
    <table>
      <thead>
        <tr>
          <th>Account</th>
          <th>Balance</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="account in accounts" :key="account.id">
          <td>{{ account.username }}</td>
          <td>{{ account.balance }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import axios from 'axios';
import { ref, onMounted } from 'vue';

export default {
  setup() {
    const accounts = ref([]);

    // Fetch data from backend on component mount
    onMounted(async () => {
      try {
        const response = await axios.get('http://localhost:8080/account/list');
        accounts.value = response.data;
      } catch (error) {
        console.error("Error fetching accounts:", error);
      }
    });

    return { accounts };
  }
}
</script>
