
<script setup>
import { Search, Bell } from 'lucide-vue-next'

defineProps({
  title: {
    type: String,
    default: 'Inventory Management'
  }
})
</script>

<template>
  <header class="top-navbar">
    <div class="page-heading">
      <h2>{{ title }}</h2>
      <p>Manage your inventory efficiently</p>
    </div>

    <div class="navbar-actions">
      <div class="search-box">
        <Search :size="18" />
        <input placeholder="Search..." />
      </div>

      <button class="icon-button">
        <Bell :size="20" />
      </button>

      <div class="user-profile">
        <div class="avatar">GD</div>
        <div>
          <strong>Administrator</strong>
          <span>System Admin</span>
        </div>
      </div>
    </div>
  </header>
</template>
