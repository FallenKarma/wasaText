<template>
  <header class="app-header">
    <div class="header-content">
      <div class="left-section">
        <h1 class="app-title"></h1>
      </div>

      <div class="right-section">
        <div class="user-menu" ref="userMenuContainer">
          <button class="user-menu-button" @click="toggleUserMenu">
            <div class="user-avatar">
              <img
                v-if="usersStore.currentUser?.photo"
                :src="getFullPhotoUrl(usersStore.currentUser?.photo)"
                class="avatar"
              />
              <div v-else class="avatar-placeholder">{{ userInitials }}</div>
            </div>
            <span class="username">{{ username }}</span>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="chevron-icon"
              :class="{ rotated: isUserMenuOpen }"
            >
              <polyline points="6 9 12 15 18 9"></polyline>
            </svg>
          </button>

          <div v-if="isUserMenuOpen" class="user-dropdown">
            <div class="user-info">
              <div class="user-avatar">
                <img
                  v-if="usersStore.currentUser?.photo"
                  :src="getFullPhotoUrl(usersStore.currentUser?.photo)"
                  class="avatar"
                />
                <div v-else class="avatar-placeholder">{{ userInitials }}</div>
              </div>
              <div>
                <div class="user-fullname">{{ username }}</div>
              </div>
            </div>

            <div class="menu-divider"></div>

            <ul class="menu-items">
              <li @click="triggerProfilePhotoUpload">
                <input
                  type="file"
                  ref="profilePhotoInput"
                  @change="handleProfilePhotoChange"
                  accept="image/*"
                  style="display: none"
                />
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="16"
                  height="16"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"></path>
                  <circle cx="12" cy="7" r="4"></circle>
                </svg>
                <span>Set profile photo</span>
              </li>
              <li class="logout-item" @click="logout">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  width="16"
                  height="16"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  stroke-linecap="round"
                  stroke-linejoin="round"
                >
                  <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"></path>
                  <polyline points="16 17 21 12 16 7"></polyline>
                  <line x1="21" y1="12" x2="9" y2="12"></line>
                </svg>
                <span>Logout</span>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </header>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '@/store/auth'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/users'
import { getFullPhotoUrl } from '@/utilities/helpers'

export default {
  name: 'AppHeader',
  setup() {
    const authStore = useAuthStore()
    const router = useRouter()
    const usersStore = useUserStore()

    // User menu state
    const isUserMenuOpen = ref(false)
    const userMenuContainer = ref(null)
    const profilePhotoInput = ref(null) // Ref for the hidden file input

    // User data from store
    const username = computed(() => authStore.name)
    const userInitials = computed(() => {
      const name = username.value
      if (!name) return ''
      return name.charAt(0).toUpperCase()
    })

    // Toggle user menu
    const toggleUserMenu = () => {
      isUserMenuOpen.value = !isUserMenuOpen.value
    }

    // Close user menu when clicking outside
    const handleClickOutside = (event) => {
      if (userMenuContainer.value && !userMenuContainer.value.contains(event.target)) {
        isUserMenuOpen.value = false
      }
    }

    // Trigger file input click
    const triggerProfilePhotoUpload = () => {
      profilePhotoInput.value.click()
      isUserMenuOpen.value = false // Close the menu after triggering
    }

    // Handle file selection and upload
    const handleProfilePhotoChange = async (event) => {
      const file = event.target.files[0]
      if (file) {
        try {
          await usersStore.uploadProfilePhoto(file)
          // Optionally, show a success message or update the UI
          console.log('Profile photo uploaded successfully!')
        } catch (error) {
          // Handle error, e.g., show an error message to the user
          console.error('Error uploading profile photo:', error)
          alert('Failed to upload profile photo. Please try again.')
        }
      }
    }

    // Logout
    const logout = async () => {
      await authStore.logout()
      router.push('/login')
    }

    // Setup event listeners
    onMounted(() => {
      usersStore.fetchCurrentUser()
      document.addEventListener('click', handleClickOutside)
    })

    onUnmounted(() => {
      document.removeEventListener('click', handleClickOutside)
    })

    return {
      isUserMenuOpen,
      userMenuContainer,
      profilePhotoInput,
      username,
      usersStore,
      userInitials,
      getFullPhotoUrl,
      toggleUserMenu,
      logout,
      triggerProfilePhotoUpload,
      handleProfilePhotoChange,
    }
  },
}
</script>

<style scoped>
.app-header {
  height: 64px;
  border-bottom: 1px solid #e5e7eb;
  background-color: white;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
}

.header-content {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 1.5rem;
}

.left-section {
  display: flex;
  align-items: center;
}

.app-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: #3b82f6;
  margin: 0;
}

.right-section {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.user-menu {
  position: relative;
}

.user-menu-button {
  display: flex;
  align-items: center;
  background: none;
  border: none;
  padding: 0.25rem;
  cursor: pointer;
  gap: 0.5rem;
  border-radius: 0.375rem;
  transition: background-color 0.2s ease;
}

.user-menu-button:hover {
  background-color: #f3f4f6;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.user-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-placeholder {
  width: 100%;
  height: 100%;
  background-color: #3b82f6;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 1rem;
}

.username {
  font-size: 0.875rem;
  font-weight: 500;
  color: #111827;
}

.chevron-icon {
  color: #6b7280;
  transition: transform 0.2s ease;
}

.chevron-icon.rotated {
  transform: rotate(180deg);
}

.user-dropdown {
  position: absolute;
  top: calc(100% + 0.5rem);
  right: 0;
  width: 260px;
  background-color: white;
  border: 1px solid #e5e7eb;
  border-radius: 0.375rem;
  box-shadow:
    0 4px 6px -1px rgba(0, 0, 0, 0.1),
    0 2px 4px -1px rgba(0, 0, 0, 0.06);
  z-index: 30;
  overflow: hidden;
}

.user-info {
  padding: 1rem;
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.user-fullname {
  font-size: 0.875rem;
  font-weight: 600;
  color: #111827;
}

.user-email {
  font-size: 0.75rem;
  color: #6b7280;
  margin-top: 0.125rem;
}

.menu-divider {
  height: 1px;
  background-color: #e5e7eb;
  margin: 0;
}

.menu-items {
  list-style: none;
  padding: 0.5rem 0;
  margin: 0;
}

.menu-items li {
  display: flex;
  align-items: center;
  padding: 0.625rem 1rem;
  cursor: pointer;
  gap: 0.75rem;
  transition: background-color 0.2s ease;
}

.menu-items li:hover {
  background-color: #f9fafb;
}

.menu-items li svg {
  color: #6b7280;
}

.menu-items li span {
  font-size: 0.875rem;
  color: #111827;
}

.logout-item {
  color: #ef4444 !important;
}

.logout-item svg {
  color: #ef4444 !important;
}
</style>
