#!/bin/bash

# Скрипт быстрой настройки проекта Student Registry
# Использование: chmod +x setup.sh && ./setup.sh

set -e

echo "🚀 Установка проекта Student Registry"
echo "======================================"
echo ""

# Проверка наличия Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker не установлен!"
    echo "📥 Установите Docker Desktop: https://www.docker.com/products/docker-desktop"
    exit 1
fi

echo "✅ Docker найден: $(docker --version)"

# Проверка наличия Docker Compose
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose не установлен!"
    exit 1
fi

echo "✅ Docker Compose найден: $(docker-compose --version)"
echo ""

# Проверка структуры проекта
echo "📁 Проверка структуры проекта..."

if [ ! -f "docker-compose.yml" ]; then
    echo "❌ Файл docker-compose.yml не найден!"
    echo "Убедитесь, что вы находитесь в корневой директории проекта"
    exit 1
fi

if [ ! -d "backend" ]; then
    echo "❌ Директория backend не найдена!"
    exit 1
fi

if [ ! -d "frontend" ]; then
    echo "❌ Директория frontend не найдена!"
    exit 1
fi

echo "✅ Структура проекта корректна"
echo ""

# Остановка существующих контейнеров
echo "🛑 Остановка существующих контейнеров (если есть)..."
docker-compose down -v 2>/dev/null || true
echo ""

# Сборка и запуск
echo "🔨 Сборка Docker образов..."
echo "Это может занять несколько минут при первом запуске..."
docker-compose build

echo ""
echo "🚀 Запуск приложения..."
docker-compose up -d

echo ""
echo "⏳ Ожидание запуска сервисов..."
sleep 10

# Проверка здоровья сервисов
echo ""
echo "🏥 Проверка состояния сервисов..."

if docker-compose ps | grep -q "Up"; then
    echo "✅ Контейнеры запущены"
else
    echo "❌ Ошибка запуска контейнеров"
    echo "Логи:"
    docker-compose logs
    exit 1
fi

# Проверка Backend
echo ""
echo "🔍 Проверка Backend API..."
for i in {1..30}; do
    if curl -f -s http://localhost:3000/health > /dev/null 2>&1; then
        echo "✅ Backend API работает"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "❌ Backend API не отвечает"
        echo "Логи Backend:"
        docker-compose logs backend
        exit 1
    fi
    sleep 2
done

# Проверка Frontend
echo ""
echo "🔍 Проверка Frontend..."
for i in {1..15}; do
    if curl -f -s http://localhost:8080 > /dev/null 2>&1; then
        echo "✅ Frontend работает"
        break
    fi
    if [ $i -eq 15 ]; then
        echo "❌ Frontend не отвечает"
        echo "Логи Frontend:"
        docker-compose logs frontend
        exit 1
    fi
    sleep 2
done

echo ""
echo "════════════════════════════════════════"
echo "✅ Установка завершена успешно!"
echo "════════════════════════════════════════"
echo ""
echo "📝 Приложение доступно по адресам:"
echo ""
echo "   🌐 Frontend:  http://localhost:8080"
echo "   🔌 Backend:   http://localhost:3000"
echo "   🗄️  Database:  localhost:5432"
echo ""
echo "════════════════════════════════════════"
echo ""
echo "📋 Полезные команды:"
echo ""
echo "   Просмотр логов:          docker-compose logs -f"
echo "   Остановка:               docker-compose stop"
echo "   Запуск:                  docker-compose start"
echo "   Перезапуск:              docker-compose restart"
echo "   Полная остановка:        docker-compose down"
echo "   Остановка с удалением:   docker-compose down -v"
echo ""
echo "🎉 Готово! Откройте http://localhost:8080 в браузере"
echo ""