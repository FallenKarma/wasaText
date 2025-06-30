<!-- src/components/layout/AppLayout.vue -->
<template>
  <div class="app-container">
    <AppSidebar />
    <div class="main-content">
      <AppHeader />
      <div class="content-area">
        <slot></slot>
      </div>
    </div>
  </div>
</template>

<script>
import { defineComponent } from 'vue'
import AppHeader from './AppHeader.vue'
import AppSidebar from './AppSidebar.vue'

export default defineComponent({
  name: 'AppLayout',
  components: {
    AppHeader,
    AppSidebar,
  },
})
</script>

<style scoped>
.app-container {
  display: flex;
  height: 100vh;
  width: 100%;
  overflow: hidden;
  background-color: #f9fafb;
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.content-area {
  flex: 1;
  overflow-y: auto;
  padding: 1rem;
}
</style>
