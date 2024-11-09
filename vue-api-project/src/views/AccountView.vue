<template>
  <div>
    <h2>Account Data</h2>
    <!-- Summary Section -->
    <div>
      <p>Total Accounts: {{ totalAccounts }}</p>
      <p>Total Balance: {{ totalBalance }}</p>
      <p>Average Balance: {{ averageBalance }}</p>
    </div>

    <!-- Account List Table -->
    <table>
      <thead>
        <tr>
          <th>ID</th>
          <th>Username</th>
          <th>Email</th>
          <th>Balance</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="account in accounts" :key="account.id">
          <td>{{ account.id }}</td>
          <td>{{ account.username }}</td>
          <td>{{ account.email }}</td>
          <td>{{ account.balance }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import { getAccountData } from "..UPLOAD REPOSITORY GIT/task-golang-db";

export default {
  data() {
    return {
      accounts: [],
      totalBalance: 0,
      totalAccounts: 0,
      averageBalance: 0
    };
  },
  async created() {
    const response = await getAccountData();
    this.accounts = response.data;

    // Calculate summary data
    this.totalAccounts = this.accounts.length;
    this.totalBalance = this.accounts.reduce((sum, account) => sum + account.balance, 0);
    this.averageBalance = this.totalAccounts ? (this.totalBalance / this.totalAccounts) : 0;
  }
};
</script>
