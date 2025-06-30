<template>
  <div class="login-container">
    <div class="login-card">
      <h1 class="login-title">Welcome to WasaText</h1>
      <p class="login-subtitle">Sign in to start chatting</p>

      <LoginForm @login-success="onLoginSuccess" />
    </div>
  </div>
</template>

<script>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import LoginForm from '@/components/auth/LoginForm.vue'
import { useAuthStore } from '@/store/auth'

export default {
  name: 'LoginView',
  components: {
    LoginForm,
  },
  setup() {
    const authStore = useAuthStore()
    const router = useRouter()

    const onLoginSuccess = () => {
      // Redirect to home page after successful login
      router.push('/')
    }

    // Check if user is already authenticated when component is created
    onMounted(() => {
      if (authStore.isAuthenticated) {
        router.push('/home')
      }
    })

    return {
      onLoginSuccess,
    }
  },
}
</script>

<style scoped>
/* Styles optimized for desktop/laptop screens */
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 100vh;
  background-color: #f5f5f5;
  /* Set min-width to prevent layout issues on smaller screens */
  min-width: 1024px;
}

.login-card {
  width: 450px;
  padding: 2.5rem;
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.login-title {
  font-size: 1.75rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 0.5rem;
  text-align: center;
}

.login-subtitle {
  font-size: 1rem;
  color: #666;
  margin-bottom: 2rem;
  text-align: center;
}

.register-link {
  margin-top: 1.5rem;
  text-align: center;
  font-size: 0.9rem;
  color: #666;
}

.register-link a {
  color: #4a6cf7;
  text-decoration: none;
  font-weight: 500;
}

.register-link a:hover {
  text-decoration: underline;
}
</style>
