
<script setup>
import { ref, computed, onMounted } from 'vue'

const API_URL = 'http://localhost:8080/api/warehouses'

const warehouses = ref([])
const loading = ref(false)
const saving = ref(false)
const search = ref('')
const showModal = ref(false)
const editingId = ref(null)
const errorMessage = ref('')
const successMessage = ref('')

const emptyForm = () => ({
  warehouse_code: '',
  warehouse_name: '',
  address: '',
  city: '',
  contact_person: '',
  phone: '',
  email: '',
  is_active: true
})

const form = ref(emptyForm())

const totalWarehouses = computed(() => warehouses.value.length)

const activeWarehouses = computed(() =>
  warehouses.value.filter(item => item.is_active).length
)

const inactiveWarehouses = computed(() =>
  warehouses.value.filter(item => !item.is_active).length
)

const filteredWarehouses = computed(() => {
  const keyword = search.value.toLowerCase().trim()

  if (!keyword) return warehouses.value

  return warehouses.value.filter(item =>
    [
      item.warehouse_code,
      item.warehouse_name,
      item.city,
      item.contact_person
    ].some(value => String(value || '').toLowerCase().includes(keyword))
  )
})

async function fetchWarehouses() {
  loading.value = true
  errorMessage.value = ''

  try {
    const response = await fetch(API_URL)
    const result = await response.json()

    if (!response.ok) {
      throw new Error(result.message || 'Failed to load warehouses')
    }

    warehouses.value = result.data || []
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    loading.value = false
  }
}

function openAddModal() {
  editingId.value = null
  form.value = emptyForm()
  errorMessage.value = ''
  showModal.value = true
}

function openEditModal(warehouse) {
  editingId.value = warehouse.id

  form.value = {
    warehouse_code: warehouse.warehouse_code || '',
    warehouse_name: warehouse.warehouse_name || '',
    address: warehouse.address || '',
    city: warehouse.city || '',
    contact_person: warehouse.contact_person || '',
    phone: warehouse.phone || '',
    email: warehouse.email || '',
    is_active: warehouse.is_active
  }

  errorMessage.value = ''
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  editingId.value = null
  form.value = emptyForm()
}

async function saveWarehouse() {
  if (!form.value.warehouse_code.trim() ||
      !form.value.warehouse_name.trim()) {
    errorMessage.value = 'Warehouse code and name are required.'
    return
  }

  saving.value = true
  errorMessage.value = ''

  const isEditing = editingId.value !== null

  const payload = {
    warehouse_code: form.value.warehouse_code.trim(),
    warehouse_name: form.value.warehouse_name.trim(),
    address: form.value.address.trim() || null,
    city: form.value.city.trim() || null,
    contact_person: form.value.contact_person.trim() || null,
    phone: form.value.phone.trim() || null,
    email: form.value.email.trim() || null,
    is_active: form.value.is_active
  }

  try {
    const response = await fetch(
      isEditing ? `${API_URL}/${editingId.value}` : API_URL,
      {
        method: isEditing ? 'PUT' : 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(payload)
      }
    )

    const result = await response.json()

    if (!response.ok) {
      throw new Error(result.message || 'Failed to save warehouse')
    }

    closeModal()
    successMessage.value = isEditing
      ? 'Warehouse updated successfully.'
      : 'Warehouse created successfully.'

    await fetchWarehouses()
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    saving.value = false
  }
}

async function deleteWarehouse(warehouse) {
  const confirmed = window.confirm(
    `Delete warehouse "${warehouse.warehouse_name}"?`
  )

  if (!confirmed) return

  errorMessage.value = ''
  successMessage.value = ''

  try {
    const response = await fetch(`${API_URL}/${warehouse.id}`, {
      method: 'DELETE'
    })

    const result = await response.json()

    if (!response.ok) {
      throw new Error(result.message || 'Failed to delete warehouse')
    }

    successMessage.value = 'Warehouse deleted successfully.'
    await fetchWarehouses()
  } catch (error) {
    errorMessage.value = error.message
  }
}

onMounted(fetchWarehouses)
</script>

<template>
  <section class="warehouse-page">

    <div class="page-heading">
      <div>
        <h1>Warehouses</h1>
        <p>Manage your warehouse locations and information</p>
      </div>

      <button class="btn-primary" @click="openAddModal">
        <span>＋</span> Add Warehouse
      </button>
    </div>

    <div v-if="errorMessage && !showModal" class="alert alert-error">
      {{ errorMessage }}
    </div>

    <div v-if="successMessage" class="alert alert-success">
      {{ successMessage }}
    </div>

    <!-- SUMMARY CARDS -->
    <div class="summary-grid">

      <div class="summary-card">
        <p>Total Warehouses</p>
        <h2>{{ totalWarehouses }}</h2>
      </div>

      <div class="summary-card">
        <p>Active Warehouses</p>
        <h2>{{ activeWarehouses }}</h2>
      </div>

      <div class="summary-card">
        <p>Inactive Warehouses</p>
        <h2>{{ inactiveWarehouses }}</h2>
      </div>

    </div>

    <!-- WAREHOUSE TABLE -->
    <div class="table-card">

      <div class="table-heading">
        <h2>Warehouse List</h2>

        <input
          v-model="search"
          type="text"
          placeholder="Search warehouses..."
          class="search-input"
        />
      </div>

      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Code</th>
              <th>Warehouse Name</th>
              <th>City</th>
              <th>Contact Person</th>
              <th>Phone</th>
              <th>Status</th>
              <th>Actions</th>
            </tr>
          </thead>

          <tbody>
            <tr v-if="loading">
              <td colspan="8" class="empty-row">
                Loading warehouses...
              </td>
            </tr>

            <tr v-else-if="filteredWarehouses.length === 0">
              <td colspan="8" class="empty-row">
                No warehouses found.
              </td>
            </tr>

            <tr
              v-for="warehouse in filteredWarehouses"
              :key="warehouse.id"
            >
              <td>#{{ warehouse.id }}</td>

              <td class="code-cell">
                {{ warehouse.warehouse_code }}
              </td>

              <td class="name-cell">
                {{ warehouse.warehouse_name }}
              </td>

              <td>{{ warehouse.city || '-' }}</td>

              <td>{{ warehouse.contact_person || '-' }}</td>

              <td>{{ warehouse.phone || '-' }}</td>

              <td>
                <span
                  class="status-badge"
                  :class="warehouse.is_active ? 'active' : 'inactive'"
                >
                  {{ warehouse.is_active ? 'Active' : 'Inactive' }}
                </span>
              </td>

              <td>
                <div class="action-buttons">
                  <button
                    class="btn-edit"
                    title="Edit warehouse"
                    @click="openEditModal(warehouse)"
                  >
                    ✎
                  </button>

                  <button
                    class="btn-delete"
                    title="Delete warehouse"
                    @click="deleteWarehouse(warehouse)"
                  >
                    🗑
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="table-footer">
        Showing {{ filteredWarehouses.length }} of
        {{ totalWarehouses }} warehouses
      </div>

    </div>

    <!-- ADD / EDIT MODAL -->
    <div
      v-if="showModal"
      class="modal-overlay"
      @click.self="closeModal"
    >
      <div class="modal-card">

        <div class="modal-heading">
          <div>
            <h2>
              {{ editingId ? 'Edit Warehouse' : 'Add Warehouse' }}
            </h2>
            <p>Enter warehouse information below</p>
          </div>

          <button class="close-button" @click="closeModal">
            ×
          </button>
        </div>

        <form @submit.prevent="saveWarehouse">

          <div class="form-grid">

            <label>
              Warehouse Code *
              <input
                v-model="form.warehouse_code"
                type="text"
                maxlength="30"
                placeholder="WH-004"
                required
              />
            </label>

            <label>
              Warehouse Name *
              <input
                v-model="form.warehouse_name"
                type="text"
                maxlength="150"
                placeholder="Warehouse Name"
                required
              />
            </label>

            <label class="full-width">
              Address
              <textarea
                v-model="form.address"
                placeholder="Warehouse address"
                rows="2"
              ></textarea>
            </label>

            <label>
              City
              <input
                v-model="form.city"
                type="text"
                placeholder="City"
              />
            </label>

            <label>
              Contact Person
              <input
                v-model="form.contact_person"
                type="text"
                placeholder="Contact person"
              />
            </label>

            <label>
              Phone
              <input
                v-model="form.phone"
                type="text"
                placeholder="Phone number"
              />
            </label>

            <label>
              Email
              <input
                v-model="form.email"
                type="email"
                placeholder="Email address"
              />
            </label>

            <label v-if="editingId" class="status-control full-width">
              <input
                v-model="form.is_active"
                type="checkbox"
              />
              Active Warehouse
            </label>

          </div>

          <p v-if="errorMessage" class="modal-error">
            {{ errorMessage }}
          </p>

          <div class="modal-actions">
            <button
              type="button"
              class="btn-secondary"
              @click="closeModal"
            >
              Cancel
            </button>

            <button
              type="submit"
              class="btn-primary"
              :disabled="saving"
            >
              {{ saving ? 'Saving...' : 'Save Warehouse' }}
            </button>
          </div>

        </form>
      </div>
    </div>

  </section>
</template>

<style scoped>
.warehouse-page {
  padding: 32px 36px;
  color: #e2e8f0;
}

.page-heading,
.table-heading,
.modal-heading,
.modal-actions,
.table-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.page-heading {
  margin-bottom: 28px;
}

.page-heading h1 {
  margin: 0 0 8px;
  font-size: 28px;
}

.page-heading p,
.modal-heading p {
  margin: 0;
  color: #94a3b8;
}

.btn-primary {
  border: 0;
  border-radius: 10px;
  padding: 14px 22px;
  background: #14b8a6;
  color: #062b2b;
  font-size: 15px;
  font-weight: 700;
  cursor: pointer;
}

.btn-primary:hover {
  background: #2dd4bf;
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 22px;
  margin-bottom: 30px;
}

.summary-card,
.table-card,
.modal-card {
  background: #111e30;
  border: 1px solid #26364c;
  border-radius: 16px;
}

.summary-card {
  padding: 28px;
}

.summary-card p {
  margin: 0 0 14px;
  color: #a5b4c8;
}

.summary-card h2 {
  margin: 0;
  color: #2dd4bf;
  font-size: 36px;
}

.table-card {
  overflow: hidden;
}

.table-heading {
  padding: 24px;
  flex-wrap: wrap;
}

.table-heading h2 {
  margin: 0;
  font-size: 20px;
}

.search-input {
  width: 280px;
  max-width: 100%;
}

.search-input,
.form-grid input,
.form-grid textarea {
  box-sizing: border-box;
  border: 1px solid #33445b;
  border-radius: 9px;
  padding: 12px 14px;
  background: #0b1422;
  color: #e2e8f0;
  font: inherit;
}

.search-input:focus,
.form-grid input:focus,
.form-grid textarea:focus {
  outline: 1px solid #2dd4bf;
  border-color: #2dd4bf;
}

.table-wrapper {
  width: 100%;
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  min-width: 950px;
}

thead {
  background: #17263b;
}

th,
td {
  padding: 19px 24px;
  text-align: left;
  border-bottom: 1px solid #26364c;
  white-space: nowrap;
}

th {
  color: #a5b4c8;
  font-weight: 600;
}

td {
  color: #dbe5f1;
}

.code-cell {
  color: #60a5fa;
  font-weight: 600;
}

.name-cell {
  color: #f1f5f9;
  font-weight: 600;
}

.status-badge {
  display: inline-block;
  padding: 6px 13px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 700;
}

.status-badge.active {
  background: #064e3b;
  color: #2dd4bf;
}

.status-badge.inactive {
  background: #3f2932;
  color: #fda4af;
}

.action-buttons {
  display: flex;
  gap: 9px;
}

.btn-edit,
.btn-delete,
.close-button {
  border: 0;
  border-radius: 8px;
  cursor: pointer;
  font-size: 18px;
}

.btn-edit,
.btn-delete {
  width: 40px;
  height: 40px;
}

.btn-edit {
  background: #153b55;
  color: #60a5fa;
}

.btn-delete {
  background: #482632;
  color: #fb7185;
}

.empty-row {
  padding: 35px;
  text-align: center;
  color: #94a3b8;
}

.table-footer {
  justify-content: flex-start;
  padding: 17px 24px;
  color: #94a3b8;
  font-size: 13px;
}

.alert {
  padding: 13px 16px;
  margin-bottom: 18px;
  border-radius: 9px;
}

.alert-error,
.modal-error {
  color: #fda4af;
  background: #451f2b;
}

.alert-success {
  color: #5eead4;
  background: #123b38;
}

.modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
  background: rgba(2, 8, 20, 0.78);
  overflow-y: auto;
}

.modal-card {
  width: 650px;
  max-width: 100%;
  max-height: 90vh;
  overflow-y: auto;
  padding: 28px;
}

.modal-heading {
  margin-bottom: 25px;
}

.modal-heading h2 {
  margin: 0 0 7px;
}

.close-button {
  width: 36px;
  height: 36px;
  background: #26364c;
  color: #e2e8f0;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 18px;
}

.form-grid label {
  display: flex;
  flex-direction: column;
  gap: 8px;
  color: #a5b4c8;
  font-size: 13px;
  font-weight: 600;
}

.form-grid .full-width {
  grid-column: 1 / -1;
}

.form-grid textarea {
  resize: vertical;
}

.status-control {
  flex-direction: row !important;
  align-items: center;
}

.status-control input {
  width: 17px;
  height: 17px;
  accent-color: #14b8a6;
}

.modal-error {
  padding: 10px;
  border-radius: 8px;
  margin-top: 18px;
}

.modal-actions {
  justify-content: flex-end;
  margin-top: 26px;
}

.btn-secondary {
  padding: 13px 20px;
  border: 1px solid #475569;
  border-radius: 9px;
  background: transparent;
  color: #e2e8f0;
  cursor: pointer;
}

@media (max-width: 900px) {
  .warehouse-page {
    padding: 24px 18px;
  }

  .summary-grid {
    grid-template-columns: 1fr;
    gap: 14px;
  }
}

@media (max-width: 600px) {
  .page-heading {
    align-items: flex-start;
    flex-direction: column;
  }

  .form-grid {
    grid-template-columns: 1fr;
  }

  .form-grid .full-width {
    grid-column: auto;
  }

  .table-heading {
    align-items: stretch;
    flex-direction: column;
  }

  .search-input {
    width: 100%;
  }
}
</style>
