<template>
  <form @submit.prevent="handleSubmit" class="login-form">
    <div class="form-group" :class="{ error: errors.username }">
      <label for="username">Username</label>
      <input
        type="text"
        id="username"
        v-model="form.username"
        placeholder="Enter your username"
        autocomplete="username"
        @blur="validateField('username')"
      />
      <span v-if="errors.username" class="error-message">{{ errors.username }}</span>
    </div>

    <div class="form-options">
      <label class="remember-me">
        <input type="checkbox" v-model="form.rememberMe" />
        <span>Remember me</span>
      </label>
    </div>

    <div class="form-submit">
      <button type="submit" class="login-button" :disabled="isSubmitting || !isFormValid">
        <span v-if="isSubmitting">Signing in...</span>
        <span v-else>Sign in / Register</span>
      </button>
    </div>

    <div v-if="errorMessage" class="auth-error">
      {{ errorMessage }}
    </div>

    <div v-if="successMessage" class="auth-success">
      {{ successMessage }}
    </div>
  </form>
</template>

<script>
import { useAuthStore } from '@/store/auth'

export default {
  name: 'LoginForm',
  data() {
    return {
      form: {
        username: '',
        rememberMe: false,
      },
      errors: {
        username: '',
      },
      isSubmitting: false,
      errorMessage: '',
      successMessage: '',
    }
  },
  computed: {
    isFormValid() {
      return this.form.username && !this.errors.username
    },
  },
  methods: {
    validateField(field) {
      if (field === 'username') {
        this.errors.username = !this.form.username ? 'Username is required' : ''
      }
    },

    validateForm() {
      this.validateField('username')
      return !this.errors.username
    },

    async handleSubmit() {
      if (!this.validateForm()) {
        return
      }

      this.isSubmitting = true
      this.errorMessage = ''
      this.successMessage = ''

      try {
        const authStore = useAuthStore()
        const response = await authStore.login({
          username: this.form.username,
          rememberMe: this.form.rememberMe,
        })
        this.successMessage = `Welcome ${this.form.username}!`

        this.$emit('login-success', response.id)

        // Redirect to dashboard or home page after successful login
        // this.$router.push('/dashboard')
      } catch (error) {
        this.errorMessage = error.response?.data?.message || 'Login failed. Please try again later.'
      } finally {
        this.isSubmitting = false
      }
    },
  },
}
</script>

<style scoped>
.login-form {
  width: 400px;
  margin: 0 auto;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-group label {
  display: block;
  margin-bottom: 0.5rem;
  font-size: 0.9rem;
  font-weight: 500;
  color: #333;
}

.form-group input {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
  transition: border-color 0.3s;
}

.form-group input:focus {
  outline: none;
  border-color: #4a6cf7;
  box-shadow: 0 0 0 2px rgba(74, 108, 247, 0.2);
}

.form-group.error input {
  border-color: #dc3545;
}

.error-message {
  display: block;
  color: #dc3545;
  font-size: 0.8rem;
  margin-top: 0.25rem;
}

.form-options {
  display: flex;
  align-items: center;
  margin-bottom: 1.5rem;
  font-size: 0.85rem;
}

.remember-me {
  display: flex;
  align-items: center;
  cursor: pointer;
}

.remember-me input {
  margin-right: 0.5rem;
}

.login-button {
  display: block;
  width: 100%;
  padding: 0.75rem;
  background-color: #4a6cf7;
  color: white;
  border: none;
  border-radius: 4px;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: background-color 0.3s;
}

.login-button:hover:not(:disabled) {
  background-color: #3a5cdc;
}

.login-button:disabled {
  background-color: #a0aeff;
  cursor: not-allowed;
}

.auth-error {
  margin-top: 1rem;
  padding: 0.75rem;
  background-color: #ffebee;
  border-radius: 4px;
  color: #c62828;
  font-size: 0.9rem;
  text-align: center;
}

.auth-success {
  margin-top: 1rem;
  padding: 0.75rem;
  background-color: #e8f5e9;
  border-radius: 4px;
  color: #2e7d32;
  font-size: 0.9rem;
  text-align: center;
}
</style>
