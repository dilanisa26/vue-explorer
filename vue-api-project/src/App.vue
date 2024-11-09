<script setup>
import { RouterView } from 'vue-router';
</script>

<template>
  <!-- Conditional layout component, if layout is specified in the route meta -->
  <component :is="$route.meta.layout || 'div'">
    <RouterView />
  </component>
</template>

<style>
/* Add global styles here if needed */
</style>

