
<script setup>
import { ref, computed, onMounted } from 'vue'
import {
  Plus,
  Search,
  Download,
  Package,
  ClipboardList,
  Warehouse,
  X,
  Trash2,
  ArrowUpFromLine
} from 'lucide-vue-next'

const API_URL = 'http://localhost:8080/api/stock-out'

const showModal = ref(false)
const searchQuery = ref('')
const transactions = ref([])
const loading = ref(false)
const errorMessage = ref('')

const form = ref({
  transaction_no: '',
  transaction_date: new Date().toISOString().slice(0, 10),
  warehouse_id: '',
  destination: '',
  reference_no: '',
  notes: '',
  items: [{ product_id: '', quantity: 1 }]
})

const totalQuantity = computed(() =>
  transactions.value.reduce(
    (total, item) => total + Number(item.total_quantity || 0),
    0
  )
)

const filteredTransactions = computed(() => {
  const keyword = searchQuery.value.trim().toLowerCase()

  if (!keyword) return transactions.value

  return transactions.value.filter(item =>
    [
      item.transaction_no,
      item.destination,
      item.warehouse_name,
      item.reference_no
    ].some(value => String(value || '').toLowerCase().includes(keyword))
  )
})

async function loadTransactions() {
  loading.value = true
  errorMessage.value = ''

  try {
    const response = await fetch(API_URL)

    if (!response.ok) {
      throw new Error(`HTTP ${response.status}: ${response.statusText}`)
    }

    const result = await response.json()

    if (result.status && result.status !== 'success') {
      throw new Error(result.message || 'API returned an unsuccessful status')
    }

    transactions.value = Array.isArray(result)
      ? result
      : Array.isArray(result.data)
        ? result.data
        : []

  } catch (error) {
    console.error('Failed to load Stock Out transactions:', error)
    errorMessage.value =
      'Unable to load Stock Out data. Make sure the Go backend is running.'
  } finally {
    loading.value = false
  }
}

onMounted(loadTransactions)

function formatDate(value) {
  if (!value) return '-'
  return String(value).split('T')[0]
}

function openModal() {
  form.value = {
    transaction_no: '',
    transaction_date: new Date().toISOString().slice(0, 10),
    warehouse_id: '',
    destination: '',
    reference_no: '',
    notes: '',
    items: [{ product_id: '', quantity: 1 }]
  }

  showModal.value = true
}

function closeModal() {
  showModal.value = false
}

function addItem() {
  form.value.items.push({
    product_id: '',
    quantity: 1
  })
}

function removeItem(index) {
  if (form.value.items.length > 1) {
    form.value.items.splice(index, 1)
  }
}

function saveTransaction() {
  alert('Create Stock Out API belum dihubungkan. Saat ini halaman sudah membaca data dari backend.')
}
</script>

<template>
  <div class="stock-out-page">

    <!-- Page Header -->
    <div class="page-header">
      <div>
        <h1>Stock Out</h1>
        <p>Manage outgoing inventory transactions</p>
      </div>

      <button class="btn-primary" @click="openModal">
        <Plus :size="19" />
        New Stock Out
      </button>
    </div>

    <!-- Summary Cards -->
    <div class="summary-grid">

      <div class="summary-card">
        <div>
          <span>Total Transactions</span>
          <h2>{{ transactions.length }}</h2>
        </div>
        <div class="summary-icon teal">
          <ClipboardList :size="25" />
        </div>
      </div>

      <div class="summary-card">
        <div>
          <span>Total Quantity Out</span>
          <h2>{{ totalQuantity.toLocaleString() }}</h2>
        </div>
        <div class="summary-icon blue">
          <Package :size="25" />
        </div>
      </div>

      <div class="summary-card">
        <div>
          <span>Warehouses</span>
          <h2>3</h2>
        </div>
        <div class="summary-icon purple">
          <Warehouse :size="25" />
        </div>
      </div>

    </div>

    <!-- Transaction Table -->
    <div class="table-card">

      <div class="table-header">
        <div>
          <h2>Stock Out Transactions</h2>
          <p>Outgoing stock transaction history</p>
        </div>

        <div class="table-actions">
          <div class="search-box">
            <Search :size="18" />
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search transactions..."
            />
          </div>

          <button class="btn-export" title="Export">
            <Download :size="18" />
          </button>
        </div>
      </div>

      <div class="table-responsive">
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Transaction No</th>
              <th>Date</th>
              <th>Warehouse</th>
              <th>Destination</th>
              <th>Total Quantity</th>
              <th>Status</th>
            </tr>
          </thead>

          <tbody>
            <tr
              v-for="item in filteredTransactions"
              :key="item.id"
            >
              <td>#{{ item.id }}</td>
              <td class="transaction-number">
                {{ item.transaction_no }}
              </td>
              <td>{{ formatDate(item.transaction_date) }}</td>
              <td>{{ item.warehouse_name }}</td>
              <td>{{ item.destination }}</td>
              <td>{{ item.total_quantity }}</td>
              <td>
                <span class="status-badge">Completed</span>
              </td>
            </tr>

            <tr v-if="loading">
              <td colspan="7" class="empty-state">
                <strong>Loading Stock Out transactions...</strong>
              </td>
            </tr>

            <tr v-else-if="errorMessage">
              <td colspan="7" class="empty-state">
                <strong>{{ errorMessage }}</strong>
                <button class="btn-add-item" type="button" @click="loadTransactions">
                  Retry
                </button>
              </td>
            </tr>

            <tr v-else-if="filteredTransactions.length === 0">
              <td colspan="7" class="empty-state">
                <ArrowUpFromLine :size="38" />
                <strong>No Stock Out transactions found</strong>
                <span>
                  Transactions from the backend will appear here.
                </span>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="table-footer">
        Showing {{ filteredTransactions.length }} transactions
      </div>

    </div>

    <!-- New Stock Out Modal -->
    <div
      v-if="showModal"
      class="modal-overlay"
      @click.self="closeModal"
    >
      <div class="modal-card">

        <div class="modal-header">
          <div>
            <h2>New Stock Out</h2>
            <p>Create outgoing stock transaction</p>
          </div>

          <button class="close-btn" @click="closeModal">
            <X :size="22" />
          </button>
        </div>

        <form @submit.prevent="saveTransaction">

          <div class="form-grid">

            <div class="form-group">
              <label>Transaction No</label>
              <input
                v-model="form.transaction_no"
                placeholder="Example: SOUT-2026-001"
                required
              />
            </div>

            <div class="form-group">
              <label>Transaction Date</label>
              <input
                v-model="form.transaction_date"
                type="date"
                required
              />
            </div>

            <div class="form-group">
              <label>Warehouse</label>
              <select v-model="form.warehouse_id" required>
                <option value="" disabled>Select warehouse</option>
                <option value="1">WH-001 - Main Warehouse</option>
                <option value="2">WH-002 - Raw Material Warehouse</option>
                <option value="3">WH-003 - Finished Goods Warehouse</option>
              </select>
            </div>

            <div class="form-group">
              <label>Destination</label>
              <input
                v-model="form.destination"
                placeholder="Customer / Department / Project"
                required
              />
            </div>

            <div class="form-group full-width">
              <label>Reference No</label>
              <input
                v-model="form.reference_no"
                placeholder="Sales order / delivery note"
              />
            </div>

          </div>

          <div class="items-header">
            <h3>Product Items</h3>

            <button
              type="button"
              class="btn-add-item"
              @click="addItem"
            >
              <Plus :size="16" />
              Add Item
            </button>
          </div>

          <div
            v-for="(item, index) in form.items"
            :key="index"
            class="item-row"
          >
            <div class="form-group">
              <label>Product ID</label>
              <input
                v-model="item.product_id"
                type="number"
                min="1"
                placeholder="Product ID"
                required
              />
            </div>

            <div class="form-group">
              <label>Quantity</label>
              <input
                v-model.number="item.quantity"
                type="number"
                min="0.01"
                step="0.01"
                required
              />
            </div>

            <button
              type="button"
              class="remove-item"
              @click="removeItem(index)"
              title="Remove item"
            >
              <Trash2 :size="18" />
            </button>
          </div>

          <div class="form-group notes-group">
            <label>Notes</label>
            <textarea
              v-model="form.notes"
              rows="3"
              placeholder="Additional notes..."
            ></textarea>
          </div>

          <div class="modal-footer">
            <button
              type="button"
              class="btn-cancel"
              @click="closeModal"
            >
              Cancel
            </button>

            <button type="submit" class="btn-primary">
              Save Transaction
            </button>
          </div>

        </form>
      </div>
    </div>

  </div>
</template>

<style scoped>
.stock-out-page {
  padding: 32px 36px;
  color: #e2e8f0;
}

.page-header,
.table-header,
.modal-header,
.modal-footer,
.items-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.page-header {
  margin-bottom: 32px;
}

.page-header h1 {
  font-size: 30px;
  margin: 0 0 8px;
}

.page-header p,
.table-header p,
.modal-header p {
  color: #94a3b8;
  margin: 0;
}

button {
  cursor: pointer;
  border: none;
  font: inherit;
}

.btn-primary {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 9px;
  padding: 15px 22px;
  border-radius: 11px;
  background: #14b8a6;
  color: #062a2b;
  font-weight: 700;
}

.btn-primary:hover {
  background: #2dd4bf;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 24px;
  margin-bottom: 32px;
}

.summary-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #111e30;
  border: 1px solid #26374d;
  border-radius: 17px;
  padding: 28px;
}

.summary-info span,
.summary-card span {
  color: #a8b7ca;
}

.summary-card h2 {
  color: #2dd4bf;
  font-size: 36px;
  margin: 14px 0 0;
}

.summary-icon {
  padding: 15px;
  border-radius: 13px;
}

.teal { color: #2dd4bf; background: #123b40; }
.blue { color: #60a5fa; background: #172f4d; }
.purple { color: #c084fc; background: #35264e; }

.table-card {
  background: #111e30;
  border: 1px solid #26374d;
  border-radius: 17px;
  overflow: hidden;
}

.table-header {
  padding: 26px;
  flex-wrap: wrap;
}

.table-header h2 {
  margin: 0 0 8px;
  font-size: 21px;
}

.table-actions {
  display: flex;
  gap: 12px;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 14px;
  background: #0b1422;
  border: 1px solid #34465e;
  border-radius: 10px;
  color: #94a3b8;
}

.search-box input {
  width: 220px;
  padding: 13px 0;
  background: transparent;
  color: white;
  border: none;
  outline: none;
}

.btn-export,
.btn-add-item {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 7px;
  padding: 12px;
  border-radius: 9px;
  color: #60a5fa;
  background: #193650;
}

.table-responsive {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  white-space: nowrap;
}

th {
  text-align: left;
  padding: 19px 24px;
  background: #17263b;
  color: #a8b7ca;
}

td {
  padding: 20px 24px;
  border-top: 1px solid #26374d;
}

.transaction-number {
  color: #60a5fa;
  font-weight: 600;
}

.status-badge {
  background: #123f3b;
  color: #2dd4bf;
  padding: 7px 13px;
  border-radius: 20px;
  font-size: 13px;
}

.empty-state {
  height: 210px;
  text-align: center;
  color: #94a3b8;
}

.empty-state svg {
  display: block;
  margin: 0 auto 12px;
  color: #64748b;
}

.empty-state strong,
.empty-state span {
  display: block;
}

.empty-state span {
  margin-top: 8px;
  font-size: 13px;
}

.table-footer {
  padding: 18px 24px;
  color: #94a3b8;
  border-top: 1px solid #26374d;
  font-size: 13px;
}

/* Modal */

.modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 1000;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 20px;
  background: #020817cc;
  overflow-y: auto;
}

.modal-card {
  width: 100%;
  max-width: 850px;
  max-height: 90vh;
  overflow-y: auto;
  background: #111e30;
  border: 1px solid #34465e;
  border-radius: 17px;
  padding: 28px;
}

.modal-header {
  margin-bottom: 28px;
}

.modal-header h2 {
  margin: 0 0 7px;
}

.close-btn {
  color: #94a3b8;
  background: transparent;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 9px;
}

.form-group label {
  color: #b6c4d5;
  font-size: 14px;
}

.form-group input,
.form-group select,
.form-group textarea {
  width: 100%;
  box-sizing: border-box;
  padding: 13px;
  border: 1px solid #34465e;
  border-radius: 9px;
  background: #0b1422;
  color: #e2e8f0;
  outline: none;
}

.form-group input:focus,
.form-group select:focus,
.form-group textarea:focus {
  border-color: #14b8a6;
}

.full-width {
  grid-column: 1 / -1;
}

.items-header {
  margin: 30px 0 18px;
}

.items-header h3 {
  margin: 0;
}

.item-row {
  display: grid;
  grid-template-columns: 1.2fr 1fr auto;
  align-items: end;
  gap: 12px;
  margin-bottom: 16px;
}

.remove-item {
  padding: 12px;
  border-radius: 8px;
  color: #fb7185;
  background: #482432;
}

.notes-group {
  margin-top: 22px;
}

.modal-footer {
  justify-content: flex-end;
  margin-top: 28px;
}

.btn-cancel {
  padding: 13px 22px;
  border-radius: 9px;
  color: #cbd5e1;
  background: #26374d;
}

@media (max-width: 900px) {
  .stock-out-page {
    padding: 24px 18px;
  }

  .summary-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 600px) {
  .page-header {
    align-items: flex-start;
    flex-direction: column;
  }

  .form-grid,
  .item-row {
    grid-template-columns: 1fr;
  }

  .full-width {
    grid-column: auto;
  }

  .search-box input {
    width: 150px;
  }
}
</style>
