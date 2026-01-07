<template>
  <div class="min-h-screen bg-gray-100">
    <div class="container mx-auto px-4 py-8">
      <!-- Заголовок -->
      <div class="mb-8">
        <h1 class="text-4xl font-bold text-gray-800 mb-2">Реестр учеников</h1>
        <p class="text-gray-600">Система управления учениками ЦПСО</p>
      </div>

      <!-- Форма добавления -->
      <div class="bg-white rounded-lg shadow-md p-6 mb-8">
        <h2 class="text-2xl font-semibold mb-4 text-gray-800">Добавить ученика</h2>
        
        <form @submit.prevent="addStudent" class="space-y-4">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                ФИО <span class="text-red-500">*</span>
              </label>
              <input
                v-model="newStudent.full_name"
                type="text"
                required
                class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                placeholder="Иванов Иван Иванович"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Класс <span class="text-red-500">*</span>
              </label>
              <select
                v-model.number="newStudent.class"
                required
                class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option v-for="c in 11" :key="c" :value="c">{{ c }} класс</option>
              </select>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Статус оплаты
              </label>
              <select
                v-model="newStudent.payment_status"
                class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="paid">Оплачено</option>
                <option value="not_paid">Не оплачено</option>
                <option value="partial">Частично</option>
              </select>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Дата зачисления
              </label>
              <input
                v-model="newStudent.enrollment_date"
                type="date"
                class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>

          <button
            type="submit"
            :disabled="loading"
            class="w-full md:w-auto px-6 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {{ loading ? 'Добавление...' : 'Добавить ученика' }}
          </button>
        </form>

        <div v-if="error" class="mt-4 p-4 bg-red-100 border border-red-400 text-red-700 rounded">
          {{ error }}
        </div>
      </div>

      <!-- Фильтр -->
      <div class="bg-white rounded-lg shadow-md p-6 mb-8">
        <h2 class="text-xl font-semibold mb-4 text-gray-800">Фильтр</h2>
        <div class="flex flex-wrap gap-2">
          <button
            @click="filterByClass(null)"
            :class="[
              'px-4 py-2 rounded-md transition-colors',
              selectedClass === null
                ? 'bg-blue-600 text-white'
                : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
            ]"
          >
            Все классы
          </button>
          <button
            v-for="c in 11"
            :key="c"
            @click="filterByClass(c)"
            :class="[
              'px-4 py-2 rounded-md transition-colors',
              selectedClass === c
                ? 'bg-blue-600 text-white'
                : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
            ]"
          >
            {{ c }} класс
          </button>
        </div>
      </div>

      <!-- Список учеников -->
      <div class="bg-white rounded-lg shadow-md overflow-hidden">
        <div class="px-6 py-4 bg-gray-50 border-b border-gray-200">
          <h2 class="text-xl font-semibold text-gray-800">
            Список учеников
            <span class="text-gray-500 text-sm ml-2">({{ students.length }})</span>
          </h2>
        </div>

        <div v-if="loading && students.length === 0" class="p-8 text-center text-gray-500">
          Загрузка...
        </div>

        <div v-else-if="students.length === 0" class="p-8 text-center text-gray-500">
          Нет учеников для отображения
        </div>

        <div v-else class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  ФИО
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Класс
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Статус
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Оплата
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Дата зачисления
                </th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  Действия
                </th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="student in students" :key="student.id" :class="{'bg-red-50': student.is_expelled}">
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="text-sm font-medium text-gray-900">{{ student.full_name }}</div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="text-sm text-gray-900">{{ student.class }} класс</div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span
                    :class="[
                      'px-2 inline-flex text-xs leading-5 font-semibold rounded-full',
                      student.is_expelled
                        ? 'bg-red-100 text-red-800'
                        : 'bg-green-100 text-green-800'
                    ]"
                  >
                    {{ student.is_expelled ? 'Отчислен' : 'Обучается' }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span
                    :class="[
                      'px-2 inline-flex text-xs leading-5 font-semibold rounded-full',
                      student.payment_status === 'paid'
                        ? 'bg-green-100 text-green-800'
                        : student.payment_status === 'partial'
                        ? 'bg-yellow-100 text-yellow-800'
                        : 'bg-red-100 text-red-800'
                    ]"
                  >
                    {{ getPaymentStatusText(student.payment_status) }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {{ formatDate(student.enrollment_date) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium space-x-2">
                  <button
                    @click="downloadCertificate(student.id)"
                    class="text-blue-600 hover:text-blue-900"
                    title="Скачать справку"
                  >
                    📄 Справка
                  </button>
                  <button
                    v-if="!student.is_expelled"
                    @click="expelStudent(student.id)"
                    class="text-red-600 hover:text-red-900"
                    title="Отчислить"
                  >
                    ❌ Отчислить
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'

const API_URL = window.location.hostname === 'localhost' 
  ? 'http://localhost:3000/api'
  : '/api'

const students = ref([])
const loading = ref(false)
const error = ref('')
const selectedClass = ref(null)

const newStudent = ref({
  full_name: '',
  class: 1,
  payment_status: 'not_paid',
  enrollment_date: new Date().toISOString().split('T')[0]
})

onMounted(() => {
  fetchStudents()
})

const fetchStudents = async () => {
  loading.value = true
  error.value = ''
  try {
    const url = selectedClass.value 
      ? `${API_URL}/students?class=${selectedClass.value}`
      : `${API_URL}/students`
    const response = await axios.get(url)
    students.value = response.data || []
  } catch (err) {
    error.value = 'Ошибка загрузки данных: ' + (err.response?.data || err.message)
  } finally {
    loading.value = false
  }
}

const addStudent = async () => {
  loading.value = true
  error.value = ''
  try {
    await axios.post(`${API_URL}/students`, newStudent.value)
    
    // Сброс формы
    newStudent.value = {
      full_name: '',
      class: 1,
      payment_status: 'not_paid',
      enrollment_date: new Date().toISOString().split('T')[0]
    }
    
    await fetchStudents()
  } catch (err) {
    error.value = 'Ошибка добавления: ' + (err.response?.data || err.message)
  } finally {
    loading.value = false
  }
}

const expelStudent = async (id) => {
  if (!confirm('Вы уверены, что хотите отчислить этого ученика?')) {
    return
  }
  
  try {
    await axios.patch(`${API_URL}/students/${id}/expel`)
    await fetchStudents()
  } catch (err) {
    error.value = 'Ошибка отчисления: ' + (err.response?.data || err.message)
  }
}

const downloadCertificate = async (id) => {
  try {
    const response = await axios.get(`${API_URL}/students/${id}/certificate`, {
      responseType: 'blob'
    })
    
    const url = window.URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `certificate_${id}.txt`)
    document.body.appendChild(link)
    link.click()
    link.remove()
  } catch (err) {
    error.value = 'Ошибка скачивания справки: ' + (err.response?.data || err.message)
  }
}

const filterByClass = (classNum) => {
  selectedClass.value = classNum
  fetchStudents()
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('ru-RU')
}

const getPaymentStatusText = (status) => {
  const statuses = {
    'paid': 'Оплачено',
    'not_paid': 'Не оплачено',
    'partial': 'Частично'
  }
  return statuses[status] || status
}
</script>