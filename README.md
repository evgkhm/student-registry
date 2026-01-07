# 📚 Реестр учеников - Система документооборота

Веб-приложение для управления реестром учеников с возможностью добавления, фильтрации, отчисления и генерации справок.

## ✨ Возможности

- ✅ Просмотр списка учеников с подробной информацией
- ✅ Добавление новых учеников через удобную форму
- ✅ Отчисление учеников без удаления из базы данных
- ✅ Фильтрация списка по классам (1-11)
- ✅ Генерация и скачивание справок об обучении
- ✅ Валидация данных на стороне сервера
- ✅ Современный адаптивный интерфейс
- ✅ Запуск одной командой через Docker Compose

## 🛠 Технологический стек

### Backend
- **Go 1.25** - основной язык программирования
- **Gorilla Mux** - HTTP router
- **PostgreSQL 15** - база данных
- **lib/pq** - PostgreSQL драйвер

### Frontend
- **Vue 3** - фреймворк с Composition API
- **Vite 5** - сборщик и dev сервер
- **Tailwind CSS** - утилитарный CSS фреймворк
- **Axios** - HTTP клиент

### DevOps
- **Docker** - контейнеризация
- **Docker Compose** - оркестрация сервисов
- **Nginx** - веб-сервер для фронтенда

## 🚀 Быстрый старт

### Требования

- Docker Desktop для Mac/Windows/Linux
- (Опционально) Go 1.25+ и Node.js 18+ для разработки без Docker

### Автоматическая установка (рекомендуется)

```bash
# 1. Клонируйте репозиторий
git clone https://github.com/evgkhm/student-registry.git
cd student-registry

# 2. Запустите скрипт установки
chmod +x setup.sh
./setup.sh

# Приложение автоматически запустится и будет доступно на:
# http://localhost:8080
```

### Ручная установка

```bash
# 1. Клонируйте репозиторий
git clone https://github.com/evgkhm/student-registry.git
cd student-registry

# 2. Запустите Docker Compose
docker-compose up --build

# 3. Откройте в браузере
# Frontend: http://localhost:8080
# Backend API: http://localhost:3000
```

### Остановка приложения

```bash
# Остановка без удаления данных
docker-compose stop

# Полная остановка и удаление контейнеров
docker-compose down

# Удаление контейнеров и данных БД
docker-compose down -v
```

## 📖 API Документация

### Получить всех учеников

```bash
GET /api/students

# С фильтром по классу
GET /api/students?class=5
```

**Ответ:**
```json
[
  {
    "id": 1,
    "full_name": "Иванов Иван Иванович",
    "class": 5,
    "payment_status": "paid",
    "enrollment_date": "2024-01-15",
    "is_expelled": false,
    "created_at": "2024-01-15T10:00:00Z"
  }
]
```

### Добавить ученика

```bash
POST /api/students
Content-Type: application/json

{
  "full_name": "Петров Петр Петрович",
  "class": 7,
  "payment_status": "paid",
  "enrollment_date": "2024-01-15"
}
```

### Отчислить ученика

```bash
PATCH /api/students/:id/expel
```

### Скачать справку

```bash
GET /api/students/:id/certificate
```

## 🏗 Структура проекта

```
student-registry/
├── backend/                 # Go Backend
│   ├── main.go             # REST API endpoints
│   ├── go.mod              # Зависимости Go
│   └── Dockerfile          # Docker образ
│
├── frontend/               # Vue.js Frontend
│   ├── src/
│   │   ├── App.vue        # Главный компонент
│   │   ├── main.js        # Точка входа
│   │   └── style.css      # Tailwind стили
│   ├── index.html
│   ├── package.json
│   ├── vite.config.js
│   ├── tailwind.config.js
│   └── Dockerfile
│
├── docker-compose.yml     # Оркестрация
├── setup.sh               # Скрипт установки
└── README.md              # Эта документация
```

### Быстрая установка на Mac

```bash
# Установка инструментов
brew install go node postgresql@15

# Backend
cd backend
go mod download
export DB_HOST=localhost DB_USER=$(whoami) DB_NAME=students_db
createdb students_db
go run main.go

# Frontend (в новом терминале)
cd frontend
npm install
npm run dev
```

## 🧪 Тестирование

### Проверка работоспособности

```bash
# Health check
curl http://localhost:3000/health

# Получить список учеников
curl http://localhost:3000/api/students

# Добавить тестового ученика
curl -X POST http://localhost:3000/api/students \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Тестов Тест Тестович",
    "class": 5,
    "payment_status": "paid",
    "enrollment_date": "2024-01-15"
  }'
```

## 📊 Модель данных

### Таблица `students`

| Поле | Тип | Описание |
|------|-----|----------|
| id | SERIAL | Первичный ключ |
| full_name | VARCHAR(255) | ФИО ученика |
| class | INTEGER | Класс (1-11) |
| payment_status | VARCHAR(50) | Статус оплаты |
| enrollment_date | DATE | Дата зачисления |
| is_expelled | BOOLEAN | Отчислен или нет |
| created_at | TIMESTAMP | Дата создания |

## 🔧 Переменные окружения

### Backend

- `DB_HOST` - хост PostgreSQL (default: localhost)
- `DB_PORT` - порт БД (default: 5432)
- `DB_USER` - пользователь БД (default: postgres)
- `DB_PASSWORD` - пароль БД (default: postgres)
- `DB_NAME` - имя БД (default: students_db)
- `PORT` - порт API (default: 3000)

## 🐛 Решение проблем

### Порт уже занят

```bash
# Найти процесс на порту
lsof -ti:3000
lsof -ti:8080

# Остановить процесс
kill -9 <PID>
```

### База данных не подключается

```bash
# Проверить логи
docker-compose logs postgres
docker-compose logs backend

# Перезапустить
docker-compose restart postgres backend
```

## 🔄 Полезные команды

```bash
# Просмотр логов
docker-compose logs -f

# Логи конкретного сервиса
docker-compose logs -f backend
docker-compose logs -f frontend

# Перезапуск сервиса
docker-compose restart backend

# Пересборка образа
docker-compose build backend
docker-compose up -d backend

# Запуск в фоне
docker-compose up -d

# Просмотр запущенных контейнеров
docker-compose ps
```