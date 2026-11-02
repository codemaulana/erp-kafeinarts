import { ref, computed } from 'vue'
import { uid } from 'quasar'

export function useScDynamicGrid(initialItems = []) {
  const items = ref(initialItems.length > 0 ? initialItems : [])

  const addItem = (defaultFields = {}) => {
    items.value.push({
      id: uid(),
      item_id: null,
      quantity: 1,
      estimated_price: 0,
      unit_price: 0,
      subtotal: 0,
      ...defaultFields,
    })
  }

  const removeItem = (index) => {
    items.value.splice(index, 1)
  }

  const updateItem = (index, field, value) => {
    items.value[index][field] = value

    // Auto-calculate subtotal for PO
    if (field === 'quantity' || field === 'unit_price') {
      const qty = Number(items.value[index].quantity) || 0
      const price = Number(items.value[index].unit_price) || 0
      items.value[index].subtotal = qty * price
    }
  }

  // Specifically for PR where they use estimated_price
  const updatePrItem = (index, field, value) => {
    items.value[index][field] = value

    if (field === 'quantity' || field === 'estimated_price') {
      const qty = Number(items.value[index].quantity) || 0
      const price = Number(items.value[index].estimated_price) || 0
      items.value[index].subtotal = qty * price
    }
  }

  const subtotal = computed(() => {
    return items.value.reduce((sum, item) => sum + (Number(item.subtotal) || 0), 0)
  })

  // Total can include taxes or just be the subtotal for now
  const total = computed(() => {
    return subtotal.value
  })

  const resetItems = (newItems = []) => {
    items.value = newItems.length > 0 ? newItems : []
  }

  return {
    items,
    addItem,
    removeItem,
    updateItem,
    updatePrItem,
    subtotal,
    total,
    resetItems,
  }
}

export default useScDynamicGrid

