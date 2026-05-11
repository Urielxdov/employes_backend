package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type RegionHandler struct{ DB *sql.DB }

func NewRegionHandler(db *sql.DB) *RegionHandler {
	return &RegionHandler{DB: db}
}

func (h *RegionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/regions"), "/")
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

func (h *RegionHandler) list(w http.ResponseWriter, _ *http.Request) {
	rows, err := h.DB.Query("SELECT id, name, nicename, note, country FROM region LIMIT 100")
	if err != nil {
		log.Printf("[ERROR] listRegions query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	defer rows.Close()

	items := []Region{}
	for rows.Next() {
		var rg Region
		if err := rows.Scan(&rg.ID, &rg.Name, &rg.NiceName, &rg.Note, &rg.Country); err != nil {
			log.Printf("[ERROR] listRegions scan failed: %v", err)
			WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to parse results"})
			return
		}
		items = append(items, rg)
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: items})
}

func (h *RegionHandler) get(w http.ResponseWriter, r *http.Request) {
	regionID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/regions/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid region ID"})
		return
	}

	var rg Region
	err = h.DB.QueryRow("SELECT id, name, nicename, note, country FROM region WHERE id = ?", regionID).
		Scan(&rg.ID, &rg.Name, &rg.NiceName, &rg.Note, &rg.Country)
	if err == sql.ErrNoRows {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Region not found"})
		return
	}
	if err != nil {
		log.Printf("[ERROR] getRegion query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Database query failed"})
		return
	}
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: rg})
}

func (h *RegionHandler) create(w http.ResponseWriter, r *http.Request) {
	var rg Region
	if err := json.NewDecoder(r.Body).Decode(&rg); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := h.DB.Exec("INSERT INTO region (name, nicename, note, country) VALUES (?, ?, ?, ?)", rg.Name, rg.NiceName, rg.Note, rg.Country)
	if err != nil {
		log.Printf("[ERROR] createRegion insert failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to create region"})
		return
	}

	id, _ := result.LastInsertId()
	rg.ID = int(id)
	WriteJSON(w, http.StatusCreated, APIResponse{Success: true, Data: rg})
}

func (h *RegionHandler) update(w http.ResponseWriter, r *http.Request) {
	regionID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/regions/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid region ID"})
		return
	}

	var rg Region
	if err := json.NewDecoder(r.Body).Decode(&rg); err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid JSON"})
		return
	}

	result, err := h.DB.Exec("UPDATE region SET name = ?, nicename = ?, note = ?, country = ? WHERE id = ?", rg.Name, rg.NiceName, rg.Note, rg.Country, regionID)
	if err != nil {
		log.Printf("[ERROR] updateRegion query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to update region"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Region not found"})
		return
	}

	rg.ID = regionID
	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: rg})
}

func (h *RegionHandler) delete(w http.ResponseWriter, r *http.Request) {
	regionID, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/api/v1/regions/"))
	if err != nil {
		WriteJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "Invalid region ID"})
		return
	}

	result, err := h.DB.Exec("DELETE FROM region WHERE id = ?", regionID)
	if err != nil {
		log.Printf("[ERROR] deleteRegion query failed: %v", err)
		WriteJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "Failed to delete region"})
		return
	}

	if n, _ := result.RowsAffected(); n == 0 {
		WriteJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "Region not found"})
		return
	}

	WriteJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]string{"message": "Region deleted"}})
}
