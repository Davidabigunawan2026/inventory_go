
<script setup>
import { ref, onMounted } from 'vue'
import {
  Plus,
  Search,
  Pencil,
  Trash2,
  RefreshCw,
  Package,
  X
} from 'lucide-vue-next'

const API_URL = 'http://localhost:8080/api/products'

const products = ref([])
const loading = ref(false)
const search = ref('')
const showForm = ref(false)
const editingId = ref(null)
const errorMessage = ref('')
const successMessage = ref('')

const emptyForm = () => ({
  sku: '',
  name: '',
  description: '',
  category: '',
  unit: 'PCS',
  min_stock: 0,
  is_active: true
})

const form = ref(emptyForm())

// GET PRODUCTS
async function fetchProducts() {
  loading.value = true
  errorMessage.value = ''

  try {
    const response = await fetch(API_URL)
    const result = await response.json()

    if (!response.ok) {
      throw new Error(result.message || 'Failed to load products')
    }

    products.value = result.data || []
  } catch (error) {
    errorMessage.value = error.message
  } finally {
    loading.value = false
  }
}

// SEARCH
const filteredProducts = () => {
  const keyword = search.value.toLowerCase()

  return products.value.filter(product =>
    [
      product.sku,
      product.name,
      product.category
    ].some(value =>
      String(value || '').toLowerCase().includes(keyword)
    )
  )
}

// OPEN ADD FORM
function addProduct() {
  editingId.value = null
  form.value = emptyForm()
  errorMessage.value = ''
  showForm.value = true
}

// OPEN EDIT FORM
function editProduct(product) {
  editingId.value = product.id

  form.value = {
    sku: product.sku,
    name: product.name,
    description: product.description || '',
    category: product.category || '',
    unit: product.unit || 'PCS',
    min_stock: Number(product.min_stock || 0),
    is_active: product.is_active
  }

  errorMessage.value = ''
  showForm.value = true
}

// SAVE PRODUCT: POST / PUT
async function saveProduct() {
  errorMessage.value = ''

  const isEditing = editingId.value !== null

  const url = isEditing
    ? `${API_URL}/${editingId.value}`
    : API_URL

  try {
    const response = await fetch(url, {
      method: isEditing ? 'PUT' : 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(form.value)
    })

    const result = await response.json()

    if (!response.ok) {
      throw new Error(result.message || 'Failed to save product')
    }

    showForm.value = false
    successMessage.value = isEditing
      ? 'Product updated successfully'
      : 'Product created successfully'

    await fetchProducts()
  } catch (error) {
    errorMessage.value = error.message
  }
}

// DELETE PRODUCT
async function deleteProduct(product) {
  const confirmed = window.confirm(
    `Delete product "${product.name}"?`
  )

  if (!confirmed) return

  errorMessage.value = ''

  try {
    const response = await fetch(
      `${API_URL}/${product.id}`,
      { method: 'DELETE' }
    )

    const result = await response.json()

    if (!response.ok) {
      throw new Error(result.message || 'Failed to delete product')
    }

    successMessage.value = 'Product deleted successfully'
    await fetchProducts()
  } catch (error) {
    errorMessage.value = error.message
  }
}

onMounted(fetchProducts)
</script>

<template>
  <div class="products-page">

    <div class="page-intro">
      <div>
        <h1>Product Management</h1>
        <p>Manage your inventory products and information.</p>
      </div>

      <button class="add-button" @click="addProduct">
        <Plus :size="19" />
        Add Product
      </button>
    </div>

    <div class="product-summary">
      <div class="summary-icon">
        <Package :size="25" />
      </div>

      <div>
        <p>Total Products</p>
        <h2>{{ products.length }}</h2>
      </div>
    </div>

    <div v-if="successMessage" class="success-message">
      {{ successMessage }}
      <button @click="successMessage = ''">×</button>
    </div>

    <div v-if="errorMessage && !showForm" class="error-message">
      {{ errorMessage }}
    </div>

    <section class="product-panel">

      <div class="panel-header">
        <div>
          <h3>Products List</h3>
          <p>All registered inventory products</p>
        </div>

        <div class="product-actions">
          <div class="search-box">
            <Search :size="18" />
            <input
              v-model="search"
              placeholder="Search products..."
            />
          </div>

          <button
            class="refresh-button"
            @click="fetchProducts"
            :disabled="loading"
          >
            <RefreshCw :size="18" />
          </button>
        </div>
      </div>

      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th>ID</th>
              <th>SKU</th>
              <th>Product Name</th>
              <th>Category</th>
              <th>Unit</th>
              <th>Min. Stock</th>
              <th>Status</th>
              <th>Actions</th>
            </tr>
          </thead>

          <tbody>
            <tr v-if="loading">
              <td colspan="8" class="empty-row">
                Loading products...
              </td>
            </tr>

            <tr
              v-else
              v-for="product in filteredProducts()"
              :key="product.id"
            >
              <td>{{ product.id }}</td>
              <td class="sku">{{ product.sku }}</td>
              <td>{{ product.name }}</td>
              <td>{{ product.category }}</td>
              <td>{{ product.unit }}</td>
              <td>{{ product.min_stock }}</td>

              <td>
                <span
                  :class="product.is_active ? 'active-status' : 'inactive-status'"
                >
                  {{ product.is_active ? 'Active' : 'Inactive' }}
                </span>
              </td>

              <td>
                <div class="action-buttons">
                  <button
                    class="edit-button"
                    title="Edit"
                    @click="editProduct(product)"
                  >
                    <Pencil :size="16" />
                  </button>

                  <button
                    class="delete-button"
                    title="Delete"
                    @click="deleteProduct(product)"
                  >
                    <Trash2 :size="16" />
                  </button>
                </div>
              </td>
            </tr>

            <tr v-if="!loading && filteredProducts().length === 0">
              <td colspan="8" class="empty-row">
                No products found.
              </td>
            </tr>
          </tbody>
        </table>
      </div>

    </section>

    <!-- PRODUCT FORM -->
    <div v-if="showForm" class="modal-overlay">
      <div class="product-modal">

        <div class="modal-header">
          <div>
            <h2>
              {{ editingId ? 'Edit Product' : 'Add New Product' }}
            </h2>
            <p>Enter product information below.</p>
          </div>

          <button class="close-button" @click="showForm = false">
            <X :size="21" />
          </button>
        </div>

        <form @submit.prevent="saveProduct">

          <div class="form-grid">

            <div class="form-group">
              <label>SKU *</label>
              <input v-model="form.sku" required maxlength="50"
                placeholder="Example: ELC-006" />
            </div>

            <div class="form-group">
              <label>Product Name *</label>
              <input v-model="form.name" required maxlength="150"
                placeholder="Product name" />
            </div>

            <div class="form-group">
              <label>Category</label>
              <input v-model="form.category"
                placeholder="Example: Electronics" />
            </div>

            <div class="form-group">
              <label>Unit</label>
              <input v-model="form.unit" placeholder="PCS" />
            </div>

            <div class="form-group">
              <label>Minimum Stock</label>
              <input v-model.number="form.min_stock"
                type="number" min="0" step="0.01" />
            </div>

            <div class="form-group">
              <label>Status</label>
              <select v-model="form.is_active">
                <option :value="true">Active</option>
                <option :value="false">Inactive</option>
              </select>
            </div>

            <div class="form-group full-width">
              <label>Description</label>
              <textarea
                v-model="form.description"
                rows="3"
                placeholder="Product description..."
              ></textarea>
            </div>

          </div>

          <p v-if="errorMessage" class="error-message">
            {{ errorMessage }}
          </p>

          <div class="modal-footer">
            <button
              type="button"
              class="cancel-button"
              @click="showForm = false"
            >
              Cancel
            </button>

            <button type="submit" class="add-button">
              {{ editingId ? 'Update Product' : 'Save Product' }}
            </button>
          </div>

        </form>
      </div>
    </div>

  </div>
</template>
