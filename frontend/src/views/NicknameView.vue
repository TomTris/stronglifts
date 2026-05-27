<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const auth = useAuthStore()
const router = useRouter()
const nickname = ref('')
const error = ref('')
const submitting = ref(false)

async function submit() {
  const n = nickname.value.trim()
  if (!n) { error.value = 'Nhập nickname'; return }
  if (n.length > 20) { error.value = 'Tối đa 20 ký tự'; return }

  submitting.value = true
  try {
    await auth.setNickname(n)
    router.push('/')
  } catch (e: any) {
    error.value = e.message
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <div class="nickname-view">
    <h1 class="title">Chào mừng!</h1>
    <p class="desc">Chọn một nickname để hiển thị trên bảng xếp hạng.</p>

    <div class="card">
      <input
        v-model="nickname"
        placeholder="Nickname (1-20 ký tự)"
        maxlength="20"
        @keyup.enter="submit"
        class="nick-input"
      >
      <div v-if="error" class="error">{{ error }}</div>
      <button class="btn btn-primary" @click="submit" :disabled="submitting">
        Tiếp tục
      </button>
    </div>
  </div>
</template>

<style scoped>
.nickname-view {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding-top: 48px;
}

.title { font-size: 28px; font-weight: 800; }
.desc { font-size: 14px; color: var(--text-dim); text-align: center; }

.card { width: 100%; display: flex; flex-direction: column; gap: 12px; }

.nick-input {
  padding: 14px;
  border: 1px solid var(--surface-2);
  border-radius: 8px;
  background: var(--surface-2);
  color: var(--text);
  font-size: 18px;
  font-weight: 600;
  text-align: center;
  outline: none;
}
.nick-input:focus { border-color: var(--primary); }

.error { color: var(--fail); font-size: 13px; text-align: center; }
</style>
