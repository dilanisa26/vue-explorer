<template>
  <v-app id="app">
    <v-navigation-drawer v-model="drawer" app>
      <!-- Sidebar -->
      <v-sheet class="pa-4 d-flex flex-column align-center" color="green darken-3" dark>
        <v-avatar size="64" class="mb-3">
          <v-img :src="logo" alt="Logo"></v-img>
        </v-avatar>
        <div class="font-weight-bold text-h6">Hello...</div>
        <div>Welcome to Dompetku</div>
      </v-sheet>

      <v-divider></v-divider>

      <v-list nav dense>
        <RouterLink
          v-for="{ icon, text, route, click } in sidebarLinks"
          :key="route"
          :to="route"
          custom
        >
          <v-list-item
            class="v-list-item--active"
            :prepend-icon="icon"
            :title="text"
            link
            @click="click"
            v-if="click"
          ></v-list-item>
          <v-list-item :prepend-icon="icon" :title="text" link :to="route" v-else></v-list-item>
        </RouterLink>
      </v-list>
    </v-navigation-drawer>

    <v-app-bar app color="green" dark>
      <v-app-bar-nav-icon @click="drawer = !drawer"></v-app-bar-nav-icon>
      <v-toolbar-title>Dompetku</v-toolbar-title>

      <v-spacer></v-spacer>
      <v-btn icon>
        <v-icon>mdi-account-circle</v-icon>
      </v-btn>
    </v-app-bar>

    <v-main>
      <v-container class="pa-4">
        <RouterView />
      </v-container>
    </v-main>
  </v-app>
</template>

<script setup>
import { ref } from 'vue'
import logo from '@/assets/logo-digi.png'
import { mdiHome, mdiLogout } from '@mdi/js'
import { RouterView, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user.js'

const userStore = useUserStore()
const router = useRouter()

const drawer = ref(false)

function logout() {
  userStore.logout()
  router.push({ name: 'login' })
}

const sidebarLinks = [
  { icon: mdiHome, text: 'Home', route: '/' },
  { icon: mdiLogout, text: 'Logout', route: '/login', click: logout },
]
</script>

<style scoped>
.v-list-item--active {
  background-color: rgba(255, 255, 255, 0.1);
}
.v-navigation-drawer {
  border-right: 1px solid rgba(0, 0, 0, 0.1);
}
</style>
