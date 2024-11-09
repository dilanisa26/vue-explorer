<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';

const headers = [
  { text: 'Nama Account', value: 'name', align: 'center', sortable: true },
  { text: 'Balance', value: 'balance', align: 'center', sortable: true },
  { text: 'Action', value: 'actions', align: 'center', sortable: false }
];

const accounts = ref([]);
const totalBalance = ref(0);
const averageBalance = ref(0);
const totalAccounts = ref(0);

async function fetchAccounts() {
  try {
    const response = await axios.get('http://localhost:8080/account/list');
    if (response.data && response.data.data) {
      accounts.value = response.data.data;
      totalAccounts.value = accounts.value.length;
      totalBalance.value = accounts.value.reduce((sum, account) => sum + account.balance, 0);
      averageBalance.value = totalAccounts.value > 0 ? totalBalance.value / totalAccounts.value : 0;
    } else {
      console.error('Invalid response structure:', response.data);
    }
  } catch (error) {
    console.error('Error fetching accounts:', error);
  }
}

onMounted(() => {
  fetchAccounts();
});
</script>

<template>
  <v-container>
    <!-- Header -->
    <v-row class="mb-6">
      <v-col cols="12">
        <v-sheet class="py-4 text-center" color="green" dark>
          <h1 class="font-weight-bold text-h4">Account Dashboard</h1>
        </v-sheet>
      </v-col>
    </v-row>

    <!-- Account Summary -->
    <v-row justify="center" class="mb-6">
      <v-col cols="12" md="8" lg="6">
        <v-card outlined elevation="3" class="pa-3">
          <v-card-title class="text-center" style="font-size: 1.5rem; font-weight: bold;">
            Summary
          </v-card-title>
          <v-divider></v-divider>
          <v-card-text>
            <v-row justify="space-around" align="center">
              <v-col cols="12" sm="4" class="text-center">
                <div class="summary-title">Total Accounts</div>
                <div class="text-h6">{{ totalAccounts }}</div>
              </v-col>
              <v-divider vertical class="hidden-sm-and-down"></v-divider>
              <v-col cols="12" sm="4" class="text-center">
                <div class="summary-title">Total Balance</div>
                <div class="text-h6">Rp{{ totalBalance.toFixed(2) }}</div>
              </v-col>
              <v-divider vertical class="hidden-sm-and-down"></v-divider>
              <v-col cols="12" sm="4" class="text-center">
                <div class="summary-title">Average Balance</div>
                <div class="text-h6">Rp{{ averageBalance.toFixed(2) }}</div>
              </v-col>
            </v-row>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>

    <!-- Account List Table -->
    <v-row>
      <v-col cols="12">
        <v-data-table :headers="headers" :items="accounts" :items-per-page="5" class="elevation-1" dense>
          <template v-slot:[`item.name`]="{ item }">
            <div class="text-center">{{ item.name }}</div>
          </template>
          <template v-slot:[`item.balance`]="{ item }">
            <div class="text-center">Rp{{ item.balance.toFixed(2) }}</div>
          </template>
          <template v-slot:[`item.actions`]="{ item }">
            <div class="text-center">
              <v-btn color="red darken-1" @click="accounts.splice(accounts.indexOf(item), 1)" rounded small>
                Delete
              </v-btn>
            </div>
          </template>
        </v-data-table>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
.text-center {
  text-align: center;
}
.text-h6 {
  font-weight: bold;
  color: #333;
}
.summary-title {
  font-weight: 500;
  color: #555;
  margin-bottom: 8px;
}
.hidden-sm-and-down {
  display: none;
}
@media (min-width: 600px) {
  .hidden-sm-and-down {
    display: block;
  }
}
</style>
