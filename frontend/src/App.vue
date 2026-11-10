
<script setup>
import { ref, computed } from 'vue'

import Sidebar from './components/Sidebar.vue'
import TopNavbar from './components/TopNavbar.vue'

import Dashboard from './views/Dashboard.vue'
import Products from './views/Products.vue'
import Categories from './views/Categories.vue'
import Warehouses from './views/Warehouses.vue'
import StockIn from './views/StockIn.vue'
import StockOut from './views/StockOut.vue'


const selectedPage = ref('Dashboard')

const pageTitle = computed(() => {
  return selectedPage.value === 'Dashboard'
    ? 'Inventory Management'
    : selectedPage.value
})
</script>

<template>
  <div class="app-layout">

    <Sidebar
      :selected-page="selectedPage"
      @navigate="selectedPage = $event"
    />

    <main class="main-content">

      <TopNavbar :title="pageTitle" />

      <Dashboard
        v-if="selectedPage === 'Dashboard'"
      />

      <Products
        v-else-if="selectedPage === 'Products'"
      />

      <Categories
        v-else-if="selectedPage === 'Categories'"
      />

      <Warehouses
        v-else-if="selectedPage === 'Warehouses'"
      />

      <StockIn v-else-if="selectedPage === 'Stock In'" />

      <StockOut v-else-if="selectedPage === 'Stock Out'" />



      <div v-else class="page-placeholder">
        <h1>{{ selectedPage }}</h1>
        <p>{{ selectedPage }} module is coming soon.</p>
      </div>

    </main>

  </div>
</template>
