package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type CountryHandler struct{ DB *sql.DB }

func NewCountryHandler(db *sql.DB) *CountryHandler {
	return &CountryHandler{DB: db}
}

func (h *CountryHandler) Handle(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/countries"), "/")
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

func scanCountry(rows interface{ Scan(dest ...interface{}) error }, c *Country) error {
	var iso3 sql.NullString
	var numcode sql.NullInt64
	if err := rows.Scan(&c.ID, &c.ISO, &c.Name, &c.NiceName, &iso3, &numcode, &c.PhoneCode); err != nil {
		return err
	}
	if iso3.Valid {
		c.ISO3 = &iso3.String
	}
	if numcode.Valid {
		n := int(numcode.Int64)
		c.NumCode = &n
	}
	return nil
}

func (h *CountryHandler) list(w http.ResponseWriter, _ *http.Request) {
	rows, err := h.DB.Query("SELECT id, iso, name, nicename, iso3, numcode, phonecode FROM country LIMIT 100")
	if err != nil {
		log.Printf("[ERROR] listCountries query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	items := []Country{}
	for rows.Next() {
		var c Country
		if err := scanCountry(rows, &c); err != nil {
			log.Printf("[ERROR] listCountries scan failed: %v", err)
			WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		items = append(items, c)
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

func (h *CountryHandler) get(w http.ResponseWriter, r *http.Request) {
	countryID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/countries/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid country ID"})
		return
	}

	var c Country
	err = scanCountry(h.DB.QueryRow("SELECT id, iso, name, nicename, iso3, numcode, phonecode FROM country WHERE id = ?", countryID), &c)
	if err == sql.ErrNoRows {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Country not found"})
		return
	}
	if err != nil {
		log.Printf("[ERROR] getCountry query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: c})
}

func (h *CountryHandler) create(w http.ResponseWriter, r *http.Request) {
	var c Country
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := h.DB.Exec("INSERT INTO country (iso, name, nicename, iso3, numcode, phonecode) VALUES (?, ?, ?, ?, ?, ?)", c.ISO, c.Name, c.NiceName, c.ISO3, c.NumCode, c.PhoneCode)
	if err != nil {
		log.Printf("[ERROR] createCountry insert failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create country"})
		return
	}

	id, _ := result.LastInsertId()
	c.ID = int(id)
	WriteJSON(w, http.StatusCreated, APIResponse{Success: true, Data: c})
}

func (h *CountryHandler) update(w http.ResponseWriter, r *http.Request) {
	countryID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/countries/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid country ID"})
		return
	}

	var c Country
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := h.DB.Exec("UPDATE country SET iso = ?, name = ?, nicename = ?, iso3 = ?, numcode = ?, phonecode = ? WHERE id = ?", c.ISO, c.Name, c.NiceName, c.ISO3, c.NumCode, c.PhoneCode, countryID)
	if err != nil {
		log.Printf("[ERROR] updateCountry query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to update country"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Country not found"})
		return
	}

	c.ID = countryID
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: c})
}

func (h *CountryHandler) delete(w http.ResponseWriter, r *http.Request) {
	countryID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/countries/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid country ID"})
		return
	}

	result, err := h.DB.Exec("DELETE FROM country WHERE id = ?", countryID)
	if err != nil {
		log.Printf("[ERROR] deleteCountry query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to delete country"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Country not found"})
		return
	}

	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"message": "Country deleted"}})
}
