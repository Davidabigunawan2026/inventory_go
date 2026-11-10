
<script setup>
import {
  LayoutDashboard,
  Package,
  Tags,
  Warehouse,
  ArrowDownToLine,
  ArrowUpFromLine,
  ClipboardList,
  Settings,
  Boxes
} from 'lucide-vue-next'

defineProps({
  selectedPage: {
    type: String,
    default: 'Dashboard'
  }
})

const emit = defineEmits(['navigate'])

const menuItems = [
  { name: 'Dashboard', icon: LayoutDashboard },
  { name: 'Products', icon: Package },
  { name: 'Categories', icon: Tags },
  { name: 'Warehouses', icon: Warehouse },
  { name: 'Stock In', icon: ArrowDownToLine },
  { name: 'Stock Out', icon: ArrowUpFromLine },
  { name: 'Stock Reports', icon: ClipboardList }
]
</script>

<template>
  <aside class="sidebar">
    <div class="brand">
      <div class="brand-icon">
        <Boxes :size="27" />
      </div>

      <div>
        <h2>GUNAWAN</h2>
        <span>INVENTORY GO</span>
      </div>
    </div>

    <div class="menu-label">MAIN MENU</div>

    <nav>
      <a
        v-for="item in menuItems"
        :key="item.name"
        href="#"
        class="menu-item"
        :class="{ active: selectedPage === item.name }"
        @click.prevent="emit('navigate', item.name)"
      >
        <component :is="item.icon" :size="19" />
        <span>{{ item.name }}</span>
      </a>
    </nav>

    <div class="sidebar-bottom">
      <a
        href="#"
        class="menu-item"
        :class="{ active: selectedPage === 'Settings' }"
        @click.prevent="emit('navigate', 'Settings')"
      >
        <Settings :size="19" />
        <span>Settings</span>
      </a>

      <p>© 2026 Gunawan Brothers</p>
    </div>
  </aside>
</template>
