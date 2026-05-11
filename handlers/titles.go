package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type TitleHandler struct{ DB *sql.DB }

func NewTitleHandler(db *sql.DB) *TitleHandler {
	return &TitleHandler{DB: db}
}

func (h *TitleHandler) Handle(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/titles"), "/")
	hasID := len(parts) > 1 && parts[1] != ""

	switch r.Method {
	case "GET":
		if hasID {
			h.get(w, r)
		} else {
			h.list(w, r)
		}
	case "POST":
		h.create(w, r)
	case "PUT":
		if !hasID {
			WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "ID required for PUT"})
			return
		}
		h.update(w, r)
	case "DELETE":
		if !hasID {
			WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "ID required for DELETE"})
			return
		}
		h.delete(w, r)
	default:
		WriteJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "Method not allowed"})
	}
}

func (h *TitleHandler) list(w http.ResponseWriter, _ *http.Request) {
	rows, err := h.DB.Query("SELECT emp_no, title, from_date, to_date FROM titles LIMIT 100")
	if err != nil {
		log.Printf("[ERROR] listTitles query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	titles := []Title{}
	for rows.Next() {
		var t Title
		if err := rows.Scan(&t.EmpNo, &t.Title, &t.FromDate, &t.ToDate); err != nil {
			log.Printf("[ERROR] listTitles scan failed: %v", err)
			WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		titles = append(titles, t)
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: titles})
}

func (h *TitleHandler) get(w http.ResponseWriter, r *http.Request) {
	empNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/titles/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	rows, err := h.DB.Query("SELECT emp_no, title, from_date, to_date FROM titles WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] getTitle query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	titles := []Title{}
	for rows.Next() {
		var t Title
		if err := rows.Scan(&t.EmpNo, &t.Title, &t.FromDate, &t.ToDate); err != nil {
			log.Printf("[ERROR] getTitle scan failed: %v", err)
			WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		titles = append(titles, t)
	}

	if len(titles) == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "No titles found"})
		return
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: titles})
}

func (h *TitleHandler) create(w http.ResponseWriter, r *http.Request) {
	var t Title
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	_, err := h.DB.Exec("INSERT INTO titles (emp_no, title, from_date, to_date) VALUES (?, ?, ?, ?)", t.EmpNo, t.Title, t.FromDate, t.ToDate)
	if err != nil {
		log.Printf("[ERROR] createTitle insert failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create title"})
		return
	}
	WriteJSON(w, http.StatusCreated, APIResponse{Success: true, Data: t})
}

func (h *TitleHandler) update(w http.ResponseWriter, r *http.Request) {
	empNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/titles/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	var t Title
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := h.DB.Exec(
		"UPDATE titles SET title = ?, from_date = ?, to_date = ? WHERE emp_no = ? AND from_date = ?",
		t.Title, t.FromDate, t.ToDate, empNo, t.FromDate,
	)
	if err != nil {
		log.Printf("[ERROR] updateTitle query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to update title"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Title record not found"})
		return
	}

	t.EmpNo = empNo
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: t})
}

func (h *TitleHandler) delete(w http.ResponseWriter, r *http.Request) {
	empNo, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/titles/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid employee ID"})
		return
	}

	result, err := h.DB.Exec("DELETE FROM titles WHERE emp_no = ?", empNo)
	if err != nil {
		log.Printf("[ERROR] deleteTitle query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to delete title"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Title record not found"})
		return
	}

	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"message": "Title deleted"}})
}
