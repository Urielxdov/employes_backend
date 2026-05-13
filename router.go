package main

import (
	"database/sql"
	"net/http"

	"github.com/user/employees-api/handlers"
)

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			origin = "*"
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, Origin, X-Requested-With")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.Header().Add("Vary", "Origin")

		if requestHeaders := r.Header.Get("Access-Control-Request-Headers"); requestHeaders != "" {
			w.Header().Set("Access-Control-Allow-Headers", requestHeaders)
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

func setupRouter(database *sql.DB) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", corsMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if err := database.Ping(); err != nil {
			handlers.WriteJSON(w, http.StatusInternalServerError, handlers.APIResponse{Success: false, Error: "Database unavailable"})
			return
		}
		handlers.WriteJSON(w, http.StatusOK, handlers.APIResponse{Success: true, Data: map[string]string{"status": "ok"}})
	}))

	eh := handlers.NewEmployeeHandler(database)
	mux.HandleFunc("/api/v1/employees", corsMiddleware(eh.Handle))
	mux.HandleFunc("/api/v1/employees/", corsMiddleware(eh.Handle))

	dh := handlers.NewDepartmentHandler(database)
	mux.HandleFunc("/api/v1/departments", corsMiddleware(dh.Handle))
	mux.HandleFunc("/api/v1/departments/", corsMiddleware(dh.Handle))

	sh := handlers.NewSalaryHandler(database)
	mux.HandleFunc("/api/v1/salaries", corsMiddleware(sh.Handle))
	mux.HandleFunc("/api/v1/salaries/", corsMiddleware(sh.Handle))

	th := handlers.NewTitleHandler(database)
	mux.HandleFunc("/api/v1/titles", corsMiddleware(th.Handle))
	mux.HandleFunc("/api/v1/titles/", corsMiddleware(th.Handle))

	deh := handlers.NewDeptEmpHandler(database)
	mux.HandleFunc("/api/v1/dept_emp", corsMiddleware(deh.Handle))
	mux.HandleFunc("/api/v1/dept_emp/", corsMiddleware(deh.Handle))

	dmh := handlers.NewDeptManagerHandler(database)
	mux.HandleFunc("/api/v1/dept_manager", corsMiddleware(dmh.Handle))
	mux.HandleFunc("/api/v1/dept_manager/", corsMiddleware(dmh.Handle))

	sgh := handlers.NewSalaryGroupHandler(database)
	mux.HandleFunc("/api/v1/salary_groups", corsMiddleware(sgh.Handle))
	mux.HandleFunc("/api/v1/salary_groups/", corsMiddleware(sgh.Handle))

	seh := handlers.NewSgEmpHandler(database)
	mux.HandleFunc("/api/v1/sg_emp", corsMiddleware(seh.Handle))
	mux.HandleFunc("/api/v1/sg_emp/", corsMiddleware(seh.Handle))

	ch := handlers.NewCountryHandler(database)
	mux.HandleFunc("/api/v1/countries", corsMiddleware(ch.Handle))
	mux.HandleFunc("/api/v1/countries/", corsMiddleware(ch.Handle))

	rh := handlers.NewRegionHandler(database)
	mux.HandleFunc("/api/v1/regions", corsMiddleware(rh.Handle))
	mux.HandleFunc("/api/v1/regions/", corsMiddleware(rh.Handle))

	reh := handlers.NewRegionEmpHandler(database)
	mux.HandleFunc("/api/v1/region_emp", corsMiddleware(reh.Handle))
	mux.HandleFunc("/api/v1/region_emp/", corsMiddleware(reh.Handle))

	return mux
}
