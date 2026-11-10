
<script setup>
import { ref, computed, onMounted } from 'vue'

import {
  Package,
  Warehouse,
  TrendingUp,
  TrendingDown,
  AlertTriangle,
  RefreshCw,
  ArrowUpRight,
  ArrowDownRight,
  Activity,
  Boxes
} from 'lucide-vue-next'

import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Tooltip,
  Legend,
  Filler
} from 'chart.js'

import { Line, Doughnut, Bar } from 'vue-chartjs'

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Tooltip,
  Legend,
  Filler
)

const API_URL = 'http://localhost:8080/api/dashboard'

const loading = ref(false)
const errorMessage = ref('')
const dashboardData = ref(null)

const stats = computed(() => {
  const d = dashboardData.value || {}

  return [
    {
      title: 'Total Products',
      value: d.total_products ?? 0,
      icon: Package,
      color: 'teal'
    },
    {
      title: 'Total Warehouses',
      value: d.total_warehouses ?? 0,
      icon: Warehouse,
      color: 'blue'
    },
    {
      title: 'Stock In Quantity',
      value: formatNumber(d.total_stock_in_quantity),
      icon: TrendingUp,
      color: 'green'
    },
    {
      title: 'Stock Out Quantity',
      value: formatNumber(d.total_stock_out_quantity),
      icon: TrendingDown,
      color: 'orange'
    }
  ]
})

function formatNumber(value) {
  return Number(value || 0).toLocaleString('id-ID')
}

function formatRupiah(value) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    maximumFractionDigits: 0
  }).format(Number(value || 0))
}

const monthlyLabels = computed(() =>
  (dashboardData.value?.monthly_movement || [])
    .map(item => item.month)
)

const monthlyIn = computed(() =>
  (dashboardData.value?.monthly_movement || [])
    .map(item => Number(item.stock_in_amount || 0))
)

const monthlyOut = computed(() =>
  (dashboardData.value?.monthly_movement || [])
    .map(item => Number(item.stock_out_amount || 0))
)

const monthlyChartData = computed(() => ({
  labels: monthlyLabels.value,
  datasets: [
    {
      label: 'Stock In',
      data: monthlyIn.value,
      borderColor: '#F97316',
      backgroundColor: 'rgba(45, 212, 191, 0.12)',
      fill: true,
      tension: 0,
      borderWidth: 3
    },
    {
      label: 'Stock Out',
      data: monthlyOut.value,
      borderColor: '#60a5fa',
      backgroundColor: 'rgba(96, 165, 250, 0.08)',
      fill: true,
      tension: 0,
      borderWidth: 3
    }
  ]
}))

const monthlyChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      labels: { color: '#cbd5e1' }
    },
    tooltip: {
      callbacks: {
        label(context) {
          return `${context.dataset.label}: ${formatRupiah(context.raw)}`
        }
      }
    }
  },
  scales: {
    x: {
      ticks: { color: '#94a3b8' },
      grid: { color: 'rgba(148,163,184,0.08)' }
    },
    y: {
      ticks: {
        color: '#94a3b8',
        callback(value) {
          return 'Rp ' + Number(value).toLocaleString('id-ID')
        }
      },
      grid: { color: 'rgba(148,163,184,0.08)' }
    }
  }
}

const categoryChartData = computed(() => {
  const categories = dashboardData.value?.category_values || []

  return {
    labels: categories.map(item => item.category),
    datasets: [{
      data: categories.map(item => Number(item.stock_value || 0)),

      backgroundColor: [
        '#F97316',
        '#60a5fa',
        '#a78bfa',
        '#fbbf24',
        '#fb7185',
        '#34d399'
      ],
      borderWidth: 0
    }]
  }
})

const categoryChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'bottom',
      labels: {
        color: '#cbd5e1',
        padding: 18
      }
    },
    tooltip: {
      callbacks: {
        label(context) {
          return `${context.label}: ${formatRupiah(context.raw)}`
        }
      }
    }
  }
}

const movementChartData = computed(() => ({
  labels: monthlyLabels.value,
  datasets: [
    {
      label: 'Stock In',
      data: monthlyIn.value,
      backgroundColor: '#F97316',
      borderRadius: 5
    },
    {
      label: 'Stock Out',
      data: monthlyOut.value,
      backgroundColor: '#60a5fa',
      borderRadius: 5
    }
  ]
}))

const movementChartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      labels: { color: '#cbd5e1' }
    }
  },
  scales: {
    x: {
      ticks: { color: '#94a3b8' },
      grid: { display: false }
    },
    y: {
      ticks: { color: '#94a3b8' },
      grid: { color: 'rgba(148,163,184,0.08)' }
    }
  }
}

const recentTransactions = computed(() =>
  dashboardData.value?.recent_transactions || []
)

const lowStockItems = computed(() =>
  dashboardData.value?.low_stock_items || []
)

const inventoryValue = computed(() =>
  formatRupiah(dashboardData.value?.total_inventory_value)
)

const totalInAmount = computed(() =>
  formatRupiah(dashboardData.value?.total_stock_in_amount)
)

const totalOutAmount = computed(() =>
  formatRupiah(dashboardData.value?.total_stock_out_amount)
)

async function fetchDashboardData() {
  loading.value = true
  errorMessage.value = ''

  try {
    const response = await fetch(API_URL)

    if (!response.ok) {
      throw new Error(`HTTP Error: ${response.status}`)
    }

    const result = await response.json()

    if (result.status !== 'success' || !result.data) {
      throw new Error('Invalid Dashboard API response')
    }

    dashboardData.value = result.data

  } catch (error) {
    console.error('Dashboard API Error:', error)
    errorMessage.value = error.message

  } finally {
    loading.value = false
  }
}

function formatDate(date) {
  if (!date) return '-'

  return new Date(date).toLocaleDateString('id-ID', {
    day: '2-digit',
    month: 'short',
    year: 'numeric'
  })
}

onMounted(fetchDashboardData)
</script>

<template>
  <div class="dashboard">

    <div class="welcome">
      <div>
        <span class="eyebrow">GUNAWAN BROTHERS ERP</span>
        <h1>Management Dashboard</h1>
        <p>Inventory overview and warehouse performance.</p>
      </div>

      <button
        class="refresh-btn"
        @click="fetchDashboardData"
        :disabled="loading"
      >
        <RefreshCw :size="16" :class="{ spinning: loading }" />
        {{ loading ? 'Loading...' : 'Refresh Data' }}
      </button>
    </div>

    <div v-if="errorMessage" class="error-message">
      <AlertTriangle :size="18" />
      <span>{{ errorMessage }}</span>
      <button @click="fetchDashboardData">Retry</button>
    </div>

    <div v-if="loading && !dashboardData" class="loading-state">
      Loading dashboard data...
    </div>

    <template v-if="dashboardData">

      <!-- KPI CARDS -->
      <div class="stats-grid">
        <div
          v-for="stat in stats"
          :key="stat.title"
          class="stat-card"
        >
          <div :class="['stat-icon', stat.color]">
            <component :is="stat.icon" :size="23" />
          </div>

          <div class="stat-info">
            <p>{{ stat.title }}</p>
            <h2>{{ stat.value }}</h2>
          </div>
        </div>
      </div>

      <!-- VALUE SUMMARY -->
      <div class="value-grid">
        <div class="value-card inventory-value">
          <div class="value-heading">
            <Boxes :size="20" />
            <span>Total Inventory Value</span>
          </div>
          <h2>{{ inventoryValue }}</h2>
          <small>Current stock valuation</small>
        </div>

        <div class="value-card stock-in-value">
          <div class="value-heading">
            <ArrowUpRight :size="20" />
            <span>Total Stock In</span>
          </div>
          <h2>{{ totalInAmount }}</h2>
          <small>Purchase value</small>
        </div>

        <div class="value-card stock-out-value">
          <div class="value-heading">
            <ArrowDownRight :size="20" />
            <span>Total Stock Out</span>
          </div>
          <h2>{{ totalOutAmount }}</h2>
          <small>Outgoing stock cost</small>
        </div>
      </div>

      <!-- MONTHLY LINE CHART -->
      <div class="chart-grid">
        <section class="panel chart-panel large-panel">
          <div class="panel-heading">
            <div>
              <h3>Monthly Stock In vs Stock Out</h3>
              <p>Monthly transaction value in Rupiah</p>
            </div>
            <Activity :size="21" />
          </div>

          <div class="chart-container">
            <Line
              :data="monthlyChartData"
              :options="monthlyChartOptions"
            />
            <div
              v-if="monthlyLabels.length === 0"
              class="empty-chart"
            >
              No monthly data available
            </div>
          </div>
        </section>

        <!-- CATEGORY PIE -->
        <section class="panel chart-panel">
          <div class="panel-heading">
            <div>
              <h3>Stock Value by Category</h3>
              <p>Current inventory distribution</p>
            </div>
            <Package :size="21" />
          </div>

          <div class="category-chart">
            <Doughnut
              :data="categoryChartData"
              :options="categoryChartOptions"
            />

            <div
              v-if="categoryChartData.labels.length === 0"
              class="empty-chart"
            >
              No category data available
            </div>
          </div>
        </section>
      </div>

      <!-- MOVEMENT BAR CHART -->
      <section class="panel chart-panel movement-panel">
        <div class="panel-heading">
          <div>
            <h3>Monthly Stock Movement</h3>
            <p>Comparison of monthly stock value</p>
          </div>
          <TrendingUp :size="21" />
        </div>

        <div class="chart-container bar-container">
          <Bar
            :data="movementChartData"
            :options="movementChartOptions"
          />
        </div>
      </section>

      <!-- LOW STOCK AND RECENT TRANSACTIONS -->
      <div class="bottom-grid">

        <section class="panel">
          <div class="panel-heading">
            <div>
              <h3>Stock Alerts</h3>
              <p>Products requiring attention</p>
            </div>
            <AlertTriangle :size="20" />
          </div>

          <div v-if="lowStockItems.length" class="alert-list">
            <div
              v-for="item in lowStockItems"
              :key="item.product_id"
              class="alert-item"
            >
              <div>
                <strong>{{ item.product_name }}</strong>
                <small>{{ item.sku }}</small>
              </div>
              <span class="stock-warning">
                {{ formatNumber(item.stock_balance) }}
              </span>
            </div>
          </div>

          <div v-else class="empty-state">
            No stock alerts
          </div>
        </section>

        <section class="panel recent-panel">
          <div class="panel-heading">
            <div>
              <h3>Recent Transactions</h3>
              <p>Latest inventory activities</p>
            </div>
            <Activity :size="20" />
          </div>

          <div v-if="recentTransactions.length" class="transaction-list">
            <div
              v-for="item in recentTransactions"
              :key="item.transaction_no"
              class="transaction-item"
            >
              <div
                :class="[
                  'transaction-icon',
                  item.type === 'IN' ? 'in' : 'out'
                ]"
              >
                {{ item.type }}
              </div>

              <div class="transaction-info">
                <strong>{{ item.transaction_no }}</strong>
                <small>
                  {{ item.warehouse_name }} ·
                  {{ formatDate(item.transaction_date) }}
                </small>
              </div>

              <div class="transaction-amount">
                {{ formatRupiah(item.total_amount) }}
              </div>
            </div>
          </div>

          <div v-else class="empty-state">
            No recent transactions
          </div>
        </section>

      </div>

    </template>

  </div>
</template>

<style scoped>
.dashboard {
  color: #e2e8f0;
  display: flex;
  flex-direction: column;
  gap: 22px;
  padding: 4px 2px 28px;
}

.welcome {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.eyebrow {
  color: #2dd4bf;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 2px;
}

.welcome h1 {
  color: #f1f5f9;
  font-size: 27px;
  margin: 7px 0;
}

.welcome p,
.panel-heading p {
  color: #94a3b8;
  font-size: 13px;
  margin: 0;
}

.refresh-btn {
  display: flex;
  align-items: center;
  gap: 9px;
  padding: 11px 16px;
  background: #0f766e;
  border: 1px solid #14b8a6;
  color: white;
  border-radius: 9px;
  cursor: pointer;
}

.refresh-btn:disabled {
  opacity: .6;
  cursor: not-allowed;
}

.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 15px;
}

.stat-card,
.panel,
.value-card {
  background: #111e30;
  border: 1px solid #23364d;
  border-radius: 13px;
}

.stat-card {
  padding: 19px;
  display: flex;
  align-items: center;
  gap: 15px;
}

.stat-icon {
  width: 47px;
  height: 47px;
  border-radius: 12px;
  display: grid;
  place-items: center;
}

.stat-icon.teal {
  background: #0f766e33;
  color: #2dd4bf;
}

.stat-icon.blue {
  background: #2563eb33;
  color: #60a5fa;
}

.stat-icon.green {
  background: #16a34a33;
  color: #4ade80;
}

.stat-icon.orange {
  background: #ea580c33;
  color: #fb923c;
}

.stat-info p {
  color: #94a3b8;
  font-size: 12px;
  margin: 0 0 7px;
}

.stat-info h2 {
  color: #f8fafc;
  font-size: 23px;
  margin: 0;
}

.value-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 15px;
}

.value-card {
  padding: 20px;
}

.value-heading {
  display: flex;
  align-items: center;
  gap: 9px;
  color: #94a3b8;
  font-size: 13px;
}

.value-card h2 {
  font-size: clamp(19px, 2vw, 25px);
  color: #f8fafc;
  margin: 16px 0 7px;
  overflow-wrap: anywhere;
}

.value-card small {
  color: #64748b;
  font-size: 11px;
}

.inventory-value .value-heading { color: #2dd4bf; }
.stock-in-value .value-heading { color: #4ade80; }
.stock-out-value .value-heading { color: #60a5fa; }

.chart-grid {
  display: grid;
  grid-template-columns: 1.6fr 1fr;
  gap: 16px;
}

.panel {
  padding: 20px;
  min-width: 0;
}

.panel-heading {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  margin-bottom: 20px;
  color: #2dd4bf;
}

.panel-heading h3 {
  color: #f1f5f9;
  font-size: 15px;
  margin: 0 0 7px;
}

.chart-container {
  height: 300px;
  position: relative;
}

.category-chart {
  height: 300px;
  position: relative;
}

.bar-container {
  height: 270px;
}

.movement-panel {
  margin-top: 0;
}

.bottom-grid {
  display: grid;
  grid-template-columns: 1fr 1.5fr;
  gap: 16px;
}

.alert-list,
.transaction-list {
  display: flex;
  flex-direction: column;
}

.alert-item,
.transaction-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 13px 0;
  border-bottom: 1px solid #23364d;
}

.alert-item strong,
.transaction-info strong {
  display: block;
  font-size: 12px;
  color: #e2e8f0;
}

.alert-item small,
.transaction-info small {
  display: block;
  color: #64748b;
  font-size: 11px;
  margin-top: 5px;
}

.stock-warning {
  color: #fb923c;
  font-weight: 700;
}

.transaction-icon {
  min-width: 38px;
  height: 35px;
  border-radius: 9px;
  display: grid;
  place-items: center;
  font-size: 11px;
  font-weight: 700;
}

.transaction-icon.in {
  color: #2dd4bf;
  background: #0f766e33;
}

.transaction-icon.out {
  color: #60a5fa;
  background: #2563eb33;
}

.transaction-info {
  flex: 1;
  min-width: 0;
}

.transaction-amount {
  color: #e2e8f0;
  font-size: 12px;
  font-weight: 600;
  white-space: nowrap;
}

.empty-chart,
.empty-state,
.loading-state {
  color: #64748b;
  font-size: 13px;
  text-align: center;
  padding: 25px 10px;
}

.empty-chart {
  position: absolute;
  inset: 0;
  display: grid;
  place-items: center;
  pointer-events: none;
}

.error-message {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 13px;
  border: 1px solid #b45309;
  background: #78350f33;
  color: #fdba74;
  border-radius: 9px;
  font-size: 13px;
}

.error-message button {
  margin-left: auto;
  border: 0;
  background: #b45309;
  color: white;
  padding: 7px 12px;
  border-radius: 6px;
  cursor: pointer;
}

@media (max-width: 1100px) {
  .stats-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .chart-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 700px) {
  .welcome {
    align-items: flex-start;
    flex-direction: column;
  }

  .welcome h1 {
    font-size: 22px;
  }

  .stats-grid,
  .value-grid,
  .bottom-grid {
    grid-template-columns: 1fr;
  }

  .stat-card {
    padding: 15px;
  }

  .panel {
    padding: 15px;
  }

  .chart-container {
    height: 250px;
  }

  .category-chart {
    height: 260px;
  }

  .transaction-amount {
    font-size: 10px;
  }
}
</style>
