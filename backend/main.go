package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Student struct {
	ID             int       `json:"id"`
	FullName       string    `json:"full_name"`
	Class          int       `json:"class"`
	PaymentStatus  string    `json:"payment_status"`
	EnrollmentDate string    `json:"enrollment_date"`
	IsExpelled     bool      `json:"is_expelled"`
	CreatedAt      time.Time `json:"created_at"`
}

type StudentInput struct {
	FullName       string `json:"full_name"`
	Class          int    `json:"class"`
	PaymentStatus  string `json:"payment_status"`
	EnrollmentDate string `json:"enrollment_date"`
}

var db *sql.DB

func main() {
	var err error

	// Подключение к БД
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbName := getEnv("DB_NAME", "students_db")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Ожидание доступности БД
	for i := 0; i < 30; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				break
			}
		}
		log.Printf("Waiting for database... (%d/30)", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}
	defer db.Close()

	log.Println("Connected to database successfully")

	// Создание таблицы
	createTable()

	// Настройка роутера
	r := mux.NewRouter()

	// CORS middleware
	r.Use(corsMiddleware)

	// API routes
	r.HandleFunc("/api/students", getStudents).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/students", createStudent).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/students/{id}/expel", expelStudent).Methods("PATCH", "OPTIONS")
	r.HandleFunc("/api/students/{id}/certificate", downloadCertificate).Methods("GET", "OPTIONS")

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}).Methods("GET")

	port := getEnv("PORT", "3000")
	log.Printf("Server starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS students (
		id SERIAL PRIMARY KEY,
		full_name VARCHAR(255) NOT NULL,
		class INTEGER NOT NULL CHECK (class >= 1 AND class <= 11),
		payment_status VARCHAR(50) NOT NULL,
		enrollment_date DATE NOT NULL,
		is_expelled BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
	log.Println("Table created or already exists")
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	classFilter := r.URL.Query().Get("class")

	var rows *sql.Rows
	var err error

	if classFilter != "" {
		rows, err = db.Query(
			"SELECT id, full_name, class, payment_status, enrollment_date, is_expelled, created_at FROM students WHERE class = $1 ORDER BY created_at DESC",
			classFilter,
		)
	} else {
		rows, err = db.Query(
			"SELECT id, full_name, class, payment_status, enrollment_date, is_expelled, created_at FROM students ORDER BY created_at DESC",
		)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	students := []Student{}
	for rows.Next() {
		var s Student
		err := rows.Scan(&s.ID, &s.FullName, &s.Class, &s.PaymentStatus, &s.EnrollmentDate, &s.IsExpelled, &s.CreatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		students = append(students, s)
	}

	json.NewEncoder(w).Encode(students)
}

func createStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var input StudentInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Валидация
	if input.FullName == "" {
		http.Error(w, "ФИО обязательно для заполнения", http.StatusBadRequest)
		return
	}

	if input.Class < 1 || input.Class > 11 {
		http.Error(w, "Класс должен быть от 1 до 11", http.StatusBadRequest)
		return
	}

	if input.PaymentStatus == "" {
		input.PaymentStatus = "not_paid"
	}

	if input.EnrollmentDate == "" {
		input.EnrollmentDate = time.Now().Format("2006-01-02")
	}

	var student Student
	err = db.QueryRow(
		`INSERT INTO students (full_name, class, payment_status, enrollment_date) 
		 VALUES ($1, $2, $3, $4) 
		 RETURNING id, full_name, class, payment_status, enrollment_date, is_expelled, created_at`,
		input.FullName, input.Class, input.PaymentStatus, input.EnrollmentDate,
	).Scan(&student.ID, &student.FullName, &student.Class, &student.PaymentStatus,
		&student.EnrollmentDate, &student.IsExpelled, &student.CreatedAt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(student)
}

func expelStudent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	id := vars["id"]

	result, err := db.Exec("UPDATE students SET is_expelled = TRUE WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Ученик не найден", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Ученик отчислен"})
}

func downloadCertificate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var student Student
	err := db.QueryRow(
		"SELECT id, full_name, class, payment_status, enrollment_date, is_expelled FROM students WHERE id = $1",
		id,
	).Scan(&student.ID, &student.FullName, &student.Class, &student.PaymentStatus, &student.EnrollmentDate, &student.IsExpelled)

	if err != nil {
		http.Error(w, "Ученик не найден", http.StatusNotFound)
		return
	}

	certificate := fmt.Sprintf(`
СПРАВКА ОБ ОБУЧЕНИИ
═══════════════════════════════════════════

ФИО: %s
Класс: %d
Статус: %s
Дата зачисления: %s
Статус оплаты: %s

Дата выдачи: %s

_________________
   (подпись)
`,
		student.FullName,
		student.Class,
		func() string {
			if student.IsExpelled {
				return "Отчислен"
			}
			return "Обучается"
		}(),
		student.EnrollmentDate,
		student.PaymentStatus,
		time.Now().Format("02.01.2006"),
	)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=certificate_%d.txt", student.ID))
	w.Write([]byte(certificate))
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
