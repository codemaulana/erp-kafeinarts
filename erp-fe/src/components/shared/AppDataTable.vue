<template>
  <div class="adt-wrapper clay-animate-up">
    <!-- Optional filter slot (search bar, selects, etc.) -->
    <div v-if="$slots.filter" class="adt-filter-card q-mb-md">
      <slot name="filter" />
    </div>

    <!-- Main Table Container Card -->
    <div class="adt-card">
      <q-table
        v-bind="$attrs"
        :rows="tableRows"
        :columns="columnsWithNo"
        :row-key="rowKey"
        :loading="loading"
        :pagination="tablePagination"
        flat
        class="adt-table"
        :no-data-label="noDataLabel"
        loading-label="Memuat data..."
        @request="$emit('request', $event)"
      >
        <!-- Auto nomor urut -->
        <template #body-cell-no="props">
          <q-td :props="props" class="text-center adt-col-no">
            <span class="adt-no-badge">{{ rowNo(props.rowIndex) }}</span>
          </q-td>
        </template>

        <!-- Dynamic Column Slot Forwarding (Supports both #body-cell-[col] and #cell-[col]) -->
        <template v-for="col in contentColumns" :key="col.name" #[`body-cell-${col.name}`]="slotProps">
          <!-- 1. If parent provided #body-cell-[name] -->
          <slot
            v-if="$slots[`body-cell-${col.name}`]"
            :name="`body-cell-${col.name}`"
            v-bind="slotProps"
          />
          <!-- 2. If parent provided #cell-[name] -->
          <q-td
            v-else-if="$slots[`cell-${col.name}`]"
            :props="slotProps"
            :class="[col.align ? `text-${col.align}` : '', `adt-col-${col.name}`]"
          >
            <slot
              :name="`cell-${col.name}`"
              v-bind="{ row: slotProps.row, value: slotProps.value, ...slotProps }"
            />
          </q-td>
          <!-- 3. Default column cell rendering -->
          <q-td
            v-else
            :props="slotProps"
            :class="[col.align ? `text-${col.align}` : '', `adt-col-${col.name}`]"
          >
            {{ slotProps.value !== undefined && slotProps.value !== null ? slotProps.value : '-' }}
          </q-td>
        </template>

        <!-- Forward other non-cell slots (e.g. top, bottom, loading) -->
        <template v-for="slotName in otherSlotNames" :key="slotName" #[slotName]="slotProps">
          <slot :name="slotName" v-bind="slotProps ?? {}" />
        </template>

        <!-- Luxury Empty State -->
        <template #no-data="{ message }">
          <div class="adt-empty q-py-xl text-center full-width">
            <div class="adt-empty__icon-box q-mb-md">
              <q-icon name="inventory_2" size="40px" />
            </div>
            <div class="adt-empty__title">{{ message || noDataLabel }}</div>
            <div class="adt-empty__sub">Data belum tersedia atau tidak sesuai dengan kata kunci pencarian.</div>
          </div>
        </template>
      </q-table>
    </div>
  </div>
</template>

<script setup>
import { computed, useSlots } from 'vue'

defineOptions({ inheritAttrs: false })

const props = defineProps({
  rows:         { type: Array,   default: undefined },
  data:         { type: Array,   default: undefined },
  columns:      { type: Array,   required: true },
  rowKey:       { type: String,  default: 'id' },
  loading:      { type: Boolean, default: false },
  pagination:   { type: Object,  default: () => ({ page: 1, rowsPerPage: 10, rowsNumber: 0 }) },
  noDataLabel:  { type: String,  default: 'Belum ada data' },
  /** false = jangan tambahkan kolom No. secara otomatis */
  withRowNo:    { type: Boolean, default: true }
})

defineEmits(['request'])

const slots = useSlots()

const tableRows = computed(() => {
  return props.rows || props.data || []
})

/**
 * Filter non-cell slots to pass through to q-table
 */
const otherSlotNames = computed(() => {
  return Object.keys(slots).filter(name => {
    return !name.startsWith('cell-') && !name.startsWith('body-cell-') && name !== 'filter' && name !== 'no-data'
  })
})

/**
 * Tambahkan kolom "No." di posisi pertama jika withRowNo = true
 * dan kolom "no" belum ada di daftar kolom dari parent.
 */
const columnsWithNo = computed(() => {
  if (!props.withRowNo) return props.columns
  const hasNo = props.columns.some(c => c.name === 'no')
  if (hasNo) return props.columns
  return [
    { name: 'no', label: 'NO.', field: 'no', align: 'center', sortable: false, style: 'width: 60px' },
    ...props.columns
  ]
})

const contentColumns = computed(() => {
  return columnsWithNo.value.filter(c => c.name !== 'no')
})

const tablePagination = computed(() => {
  const p = props.pagination || {}
  const totalRows = tableRows.value.length
  const rowsNumber = (p.rowsNumber && Number(p.rowsNumber) > 0) ? Number(p.rowsNumber) : totalRows
  return {
    ...p,
    rowsNumber
  }
})

/**
 * Hitung nomor urut berdasarkan state pagination.
 */
function rowNo (rowIndex) {
  const page = tablePagination.value.page || 1
  const perPage = tablePagination.value.rowsPerPage || 10
  return (page - 1) * perPage + rowIndex + 1
}
</script>

<style lang="scss">
.adt-wrapper {
  width: 100%;
}

.adt-filter-card {
  background: var(--rb-surface, #ffffff);
  border: 1px solid var(--rb-border, #E2E8F0);
  border-radius: var(--rb-radius-lg, 16px);
  padding: 14px 18px;
  box-shadow: 0 4px 20px -2px rgba(15, 23, 42, 0.05);
  transition: all 0.25s ease;

  body.body--dark & {
    background: #111827;
    border-color: rgba(255, 255, 255, 0.08);
    box-shadow: 0 10px 30px -10px rgba(0, 0, 0, 0.5);
  }
}

.adt-card {
  background: var(--rb-surface, #ffffff);
  border: 1px solid var(--rb-border, #E2E8F0);
  border-radius: var(--rb-radius-lg, 16px);
  box-shadow: 0 4px 20px -2px rgba(15, 23, 42, 0.05);
  overflow: hidden;
  transition: all 0.25s ease;

  body.body--dark & {
    background: #111827;
    border-color: rgba(255, 255, 255, 0.08);
    box-shadow: 0 10px 30px -10px rgba(0, 0, 0, 0.5);
  }
}

.adt-table {
  background: transparent !important;

  .q-table__top {
    padding: 16px 20px;
    border-bottom: 1px solid var(--rb-border, #E2E8F0);

    body.body--dark & {
      border-color: rgba(255, 255, 255, 0.07);
    }
  }

  thead tr th {
    background: rgba(241, 245, 249, 0.7) !important;
    color: var(--rb-text-muted, #64748B) !important;
    font-size: 11.5px !important;
    font-weight: 700 !important;
    text-transform: uppercase !important;
    letter-spacing: 0.75px !important;
    padding: 14px 18px !important;
    border-bottom: 1px solid var(--rb-border, #E2E8F0) !important;
    white-space: nowrap;

    body.body--dark & {
      background: rgba(255, 255, 255, 0.03) !important;
      color: #94A3B8 !important;
      border-color: rgba(255, 255, 255, 0.08) !important;
    }
  }

  tbody tr {
    transition: background 0.15s ease, transform 0.15s ease;

    &:hover td {
      background: rgba(59, 130, 246, 0.04) !important;

      body.body--dark & {
        background: rgba(255, 255, 255, 0.035) !important;
      }
    }

    td {
      padding: 14px 18px !important;
      font-size: 13.5px !important;
      color: var(--rb-text-primary, #1E293B) !important;
      border-bottom: 1px solid rgba(226, 232, 240, 0.7) !important;

      body.body--dark & {
        color: #F3F4F6 !important;
        border-color: rgba(255, 255, 255, 0.05) !important;
      }
    }
  }

  .adt-col-no {
    width: 60px;
  }

  .adt-no-badge {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    min-width: 24px;
    height: 24px;
    padding: 0 6px;
    font-size: 11.5px;
    font-weight: 700;
    color: var(--rb-text-muted, #64748B);
    background: rgba(148, 163, 184, 0.12);
    border-radius: 6px;

    body.body--dark & {
      color: #94A3B8;
      background: rgba(255, 255, 255, 0.06);
    }
  }

  .q-table__bottom {
    padding: 12px 20px;
    font-size: 12.5px;
    color: var(--rb-text-muted, #64748B);
    border-top: 1px solid var(--rb-border, #E2E8F0);

    body.body--dark & {
      border-color: rgba(255, 255, 255, 0.07);
      color: #94A3B8;
    }
  }
}

.adt-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;

  &__icon-box {
    width: 64px;
    height: 64px;
    border-radius: 16px;
    background: rgba(59, 130, 246, 0.08);
    color: #3B82F6;
    display: flex;
    align-items: center;
    justify-content: center;

    body.body--dark & {
      background: rgba(59, 130, 246, 0.15);
      color: #60A5FA;
    }
  }

  &__title {
    font-size: 15px;
    font-weight: 700;
    color: var(--rb-text-primary, #0F172A);
    margin-bottom: 4px;

    body.body--dark & {
      color: #F8FAFC;
    }
  }

  &__sub {
    font-size: 13px;
    color: var(--rb-text-muted, #64748B);
  }
}
</style>
