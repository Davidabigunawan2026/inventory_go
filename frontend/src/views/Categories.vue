
<script setup>
import { ref, computed, onMounted } from 'vue'
import { Search, Plus, Pencil, Trash2, X, RefreshCw } from 'lucide-vue-next'

const API_URL = 'http://localhost:8080/api/categories'

const categories = ref([])
const loading = ref(false)
const errorMessage = ref('')
const searchQuery = ref('')

const showModal = ref(false)
const isEditing = ref(false)
const editingId = ref(null)

const form = ref({
  name: '',
  description: '',
  is_active: true
})

const filteredCategories = computed(() => {
  const keyword = searchQuery.value.toLowerCase()

  return categories.value.filter(category =>
    category.name.toLowerCase().includes(keyword) ||
    (category.description || '').toLowerCase().includes(keyword)
  )
})

const activeCount = computed(() =>
  categories.value.filter(item => item.is_active).length
)

async function loadCategories() {
  loading.value = true
  errorMessage.value = ''

  try {
    const response = await fetch(API_URL)

    if (!response.ok) {
      throw new Error(`HTTP Error: ${response.status}`)
    }

    const result = await response.json()
    categories.value = result.data || []
  } catch (error) {
    errorMessage.value = error.message
    console.error('Load Categories Error:', error)
  } finally {
    loading.value = false
  }
}

function openAddModal() {
  isEditing.value = false
  editingId.value = null

  form.value = {
    name: '',
    description: '',
    is_active: true
  }

  showModal.value = true
}

function openEditModal(category) {
  isEditing.value = true
  editingId.value = category.id

  form.value = {
    name: category.name,
    description: category.description || '',
    is_active: category.is_active
  }

  showModal.value = true
}

async function saveCategory() {
  if (!form.value.name.trim()) {
    alert('Category name is required!')
    return
  }

  const url = isEditing.value
    ? `${API_URL}/${editingId.value}`
    : API_URL

  const method = isEditing.value ? 'PUT' : 'POST'

  try {
    const response = await fetch(url, {
      method,
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        name: form.value.name.trim(),
        description: form.value.description.trim(),
        is_active: form.value.is_active
      })
    })

    const result = await response.json()

    if (!response.ok) {
      throw new Error(result.error || result.message || 'Save failed')
    }

    showModal.value = false
    await loadCategories()
  } catch (error) {
    alert(error.message)
    console.error('Save Category Error:', error)
  }
}

async function deleteCategory(category) {
  const confirmed = confirm(
    `Delete category "${category.name}"?`
  )

  if (!confirmed) return

  try {
    const response = await fetch(`${API_URL}/${category.id}`, {
      method: 'DELETE'
    })

    if (!response.ok) {
      const result = await response.json()
      throw new Error(result.error || 'Delete failed')
    }

    await loadCategories()
  } catch (error) {
    alert(error.message)
    console.error('Delete Category Error:', error)
  }
}

onMounted(loadCategories)
</script>

<template>
  <div class="categories-page">

    <div class="page-heading">
      <div>
        <h1>Categories</h1>
        <p>Organize and manage your product categories</p>
      </div>

      <button class="add-button" @click="openAddModal">
        <Plus :size="18" />
        Add Category
      </button>
    </div>

    <div class="summary-grid">
      <div class="summary-card">
        <span>Total Categories</span>
        <strong>{{ categories.length }}</strong>
      </div>

      <div class="summary-card">
        <span>Active Categories</span>
        <strong>{{ activeCount }}</strong>
      </div>

      <div class="summary-card">
        <span>Inactive Categories</span>
        <strong>{{ categories.length - activeCount }}</strong>
      </div>
    </div>

    <div class="table-card">

      <div class="table-toolbar">
        <h2>Category List</h2>

        <div class="search-box">
          <Search :size="18" />
          <input
            v-model="searchQuery"
            placeholder="Search categories..."
          />
        </div>
      </div>

      <div v-if="errorMessage" class="error-message">
        {{ errorMessage }}
        <button @click="loadCategories">Try Again</button>
      </div>

      <div v-else-if="loading" class="empty-state">
        Loading categories...
      </div>

      <div v-else class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>Category Name</th>
              <th>Description</th>
              <th>Status</th>
              <th>Actions</th>
            </tr>
          </thead>

          <tbody>
            <tr
              v-for="category in filteredCategories"
              :key="category.id"
            >
              <td>#{{ category.id }}</td>
              <td class="category-name">{{ category.name }}</td>
              <td>{{ category.description || '-' }}</td>

              <td>
                <span
                  class="status"
                  :class="category.is_active ? 'active' : 'inactive'"
                >
                  {{ category.is_active ? 'Active' : 'Inactive' }}
                </span>
              </td>

              <td>
                <div class="actions">
                  <button
                    class="edit-button"
                    title="Edit"
                    @click="openEditModal(category)"
                  >
                    <Pencil :size="16" />
                  </button>

                  <button
                    class="delete-button"
                    title="Delete"
                    @click="deleteCategory(category)"
                  >
                    <Trash2 :size="16" />
                  </button>
                </div>
              </td>
            </tr>

            <tr v-if="filteredCategories.length === 0">
              <td colspan="5" class="empty-state">
                No categories found.
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="table-footer">
        Showing {{ filteredCategories.length }} of
        {{ categories.length }} categories
        <button class="refresh-button" @click="loadCategories">
          <RefreshCw :size="15" />
          Refresh
        </button>
      </div>

    </div>

    <!-- ADD / EDIT MODAL -->
    <div
      v-if="showModal"
      class="modal-overlay"
      @click.self="showModal = false"
    >
      <div class="modal-card">

        <div class="modal-heading">
          <h2>{{ isEditing ? 'Edit Category' : 'Add Category' }}</h2>
          <button class="close-button" @click="showModal = false">
            <X :size="20" />
          </button>
        </div>

        <form @submit.prevent="saveCategory">

          <label>Category Name *</label>
          <input
            v-model="form.name"
            placeholder="Enter category name"
            maxlength="100"
            required
          />

          <label>Description</label>
          <textarea
            v-model="form.description"
            placeholder="Enter category description"
            rows="4"
          ></textarea>

          <label class="active-option">
            <input type="checkbox" v-model="form.is_active" />
            Active
          </label>

          <div class="modal-actions">
            <button
              type="button"
              class="cancel-button"
              @click="showModal = false"
            >
              Cancel
            </button>

            <button type="submit" class="save-button">
              {{ isEditing ? 'Update Category' : 'Save Category' }}
            </button>
          </div>

        </form>
      </div>
    </div>

  </div>
</template>

<style scoped>
.categories-page {
  padding: 28px;
  color: #e2e8f0;
  min-height: 100%;
}

.page-heading {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 20px;
  margin-bottom: 25px;
}

.page-heading h1 {
  font-size: 25px;
  margin: 0 0 7px;
}

.page-heading p {
  color: #94a3b8;
  margin: 0;
}

.add-button,
.save-button {
  background: #14b8a6;
  color: #06211f;
  border: none;
  border-radius: 8px;
  padding: 11px 17px;
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  cursor: pointer;
}

.add-button:hover,
.save-button:hover {
  background: #2dd4bf;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 18px;
  margin-bottom: 24px;
}

.summary-card,
.table-card {
  background: #111e30;
  border: 1px solid #25364b;
  border-radius: 12px;
}

.summary-card {
  padding: 20px;
}

.summary-card span {
  color: #94a3b8;
  font-size: 13px;
}

.summary-card strong {
  display: block;
  margin-top: 10px;
  font-size: 28px;
  color: #2dd4bf;
}

.table-card {
  overflow: hidden;
}

.table-toolbar {
  padding: 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 15px;
  border-bottom: 1px solid #25364b;
}

.table-toolbar h2 {
  font-size: 17px;
  margin: 0;
}

.search-box {
  display: flex;
  align-items: center;
  gap: 9px;
  background: #0b1422;
  border: 1px solid #334155;
  border-radius: 8px;
  padding: 10px 12px;
  color: #94a3b8;
}

.search-box input {
  background: transparent;
  color: white;
  border: none;
  outline: none;
  width: 190px;
}

.table-wrapper {
  overflow-x: auto;
}

table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

th, td {
  padding: 16px 20px;
  border-bottom: 1px solid #25364b;
  font-size: 13px;
}

th {
  color: #94a3b8;
  background: #142238;
  font-weight: 600;
}

td {
  color: #cbd5e1;
}

.category-name {
  color: #f1f5f9;
  font-weight: 600;
}

.status {
  padding: 5px 10px;
  border-radius: 20px;
  font-size: 11px;
  font-weight: 600;
}

.status.active {
  background: #123b39;
  color: #2dd4bf;
}

.status.inactive {
  background: #3b2730;
  color: #fda4af;
}

.actions {
  display: flex;
  gap: 9px;
}

.edit-button,
.delete-button,
.close-button,
.refresh-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  border: none;
  border-radius: 6px;
}

.edit-button {
  background: #183a52;
  color: #60a5fa;
  padding: 8px;
}

.delete-button {
  background: #40232e;
  color: #fb7185;
  padding: 8px;
}

.table-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 15px 20px;
  color: #94a3b8;
  font-size: 12px;
}

.refresh-button {
  background: transparent;
  color: #2dd4bf;
  gap: 7px;
}

.empty-state {
  text-align: center;
  padding: 30px;
  color: #94a3b8;
}

.error-message {
  padding: 20px;
  color: #fda4af;
}

.modal-overlay {
  position: fixed;
  inset: 0;
  z-index: 1000;
  background: #000a;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.modal-card {
  background: #111e30;
  border: 1px solid #334155;
  border-radius: 14px;
  width: 100%;
  max-width: 480px;
  padding: 24px;
}

.modal-heading {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 22px;
}

.modal-heading h2 {
  font-size: 19px;
  margin: 0;
}

.close-button {
  background: transparent;
  color: #94a3b8;
}

.modal-card form > label:not(.active-option) {
  display: block;
  margin: 15px 0 8px;
  font-size: 13px;
  color: #cbd5e1;
}

.modal-card form > input,
.modal-card textarea {
  width: 100%;
  box-sizing: border-box;
  background: #0b1422;
  color: white;
  border: 1px solid #334155;
  border-radius: 8px;
  padding: 12px;
  outline: none;
  font: inherit;
}

.modal-card form > input:focus,
.modal-card textarea:focus {
  border-color: #14b8a6;
}

.active-option {
  display: flex;
  align-items: center;
  gap: 9px;
  margin-top: 18px;
  color: #cbd5e1;
  font-size: 13px;
}

.active-option input {
  accent-color: #14b8a6;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 25px;
}

.cancel-button {
  background: #25364b;
  color: #e2e8f0;
  border: none;
  border-radius: 8px;
  padding: 11px 16px;
  cursor: pointer;
}

@media (max-width: 700px) {
  .categories-page {
    padding: 16px;
  }

  .page-heading {
    align-items: flex-start;
    flex-direction: column;
  }

  .summary-grid {
    grid-template-columns: 1fr;
  }

  .table-toolbar {
    align-items: flex-start;
    flex-direction: column;
  }

  .search-box,
  .search-box input {
    width: 100%;
    box-sizing: border-box;
  }
}
</style>
