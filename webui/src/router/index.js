import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/store/auth'

// Import views
import LoginView from '@/views/LoginView.vue'
import HomeView from '@/views/HomeView.vue'
import ConversationView from '@/views/ConversationView.vue'

// Define routes
const routes = [
  {
    path: '/',
    redirect: '/home',
  },
  {
    path: '/login',
    name: 'login',
    component: LoginView,
    meta: { requiresAuth: false },
  },
  {
    path: '/home',
    name: 'home',
    component: HomeView,
    meta: { requiresAuth: true },
  },
  /*   {
    path: '/settings',
    name: 'settings',
    component: HomeView,
    meta: { requiresAuth: true },
  }, */
  {
    path: '/conversations/:conversationId',
    name: 'conversation',
    component: ConversationView,
    props: true,
    meta: { requiresAuth: true },
  },
  {
    // Catch all / 404 route
    path: '/:pathMatch(.*)*',
    redirect: '/home',
  },
]

// Create router instance
const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Navigation guard for authentication
router.beforeEach(async (to, from, next) => {
  // Check if route requires authentication
  const requiresAuth = to.matched.some((record) => record.meta.requiresAuth)

  // Get auth store and check authentication status
  const authStore = useAuthStore()
  const isAuthenticated = authStore.isAuthenticated
  const token = authStore.token

  // If authentication is required and user is not authenticated, redirect to login
  if (requiresAuth && !isAuthenticated) {
    next('/login')
  }
  // If user is authenticated and trying to access login/register pages, redirect to home
  else if (isAuthenticated && to.name === 'login') {
    next('/home')
  }
  // Otherwise, allow navigation
  else {
    next()
  }
})

export default router
