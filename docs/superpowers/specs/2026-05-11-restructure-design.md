# Restructure: separación por entidad + fix coherencia DB

**Date:** 2026-05-11  
**Scope:** Refactor del backend Go — separar handlers por entidad y corregir nombres de tabla

---

## Contexto

El backend actual tiene 4 archivos en raíz (`main.go`, `handlers.go`, `models.go`, `db.go`). `handlers.go` tiene 1451 líneas con CRUD completo para 11 entidades mezclado en un solo archivo. El routing vive en `main.go` junto con `main()`. Además, las queries de Employee usan el nombre de tabla `employees` (plural) pero el schema SQL define la tabla como `employee` (singular), causando errores en runtime.

---

## Cambios

### 1. Corrección de tabla `employee`

Todas las queries que referencian `employees` (plural) se corrigen a `employee`:

| Operación | Antes | Después |
|---|---|---|
| listEmployees | `FROM employees` | `FROM employee` |
| getEmployee | `FROM employees WHERE` | `FROM employee WHERE` |
| createEmployee | `INSERT INTO employees` | `INSERT INTO employee` |
| updateEmployee | `UPDATE employees SET` | `UPDATE employee SET` |
| deleteEmployee | `DELETE FROM employees` | `DELETE FROM employee` |

Las demás tablas (`department`, `salary`, `titles`, `dept_emp`, `dept_manager`, `salary_group`, `sg_emp`, `country`, `region`, `region_emp`) ya coinciden con el schema.

### 2. Estructura de archivos final

```
employes_backend/
├── main.go              # solo: initDB, setupRouter, ListenAndServe
├── db.go                # sin cambios
├── models.go            # sin cambios
├── router.go            # NUEVO: corsMiddleware + setupRouter() registra todos los handlers
└── handlers/
    ├── common.go        # writeJSON() compartido
    ├── employees.go     # EmployeeHandler
    ├── departments.go   # DepartmentHandler
    ├── salaries.go      # SalaryHandler
    ├── titles.go        # TitleHandler
    ├── dept_emp.go      # DeptEmpHandler
    ├── dept_manager.go  # DeptManagerHandler
    ├── salary_groups.go # SalaryGroupHandler
    ├── sg_emp.go        # SgEmpHandler
    ├── countries.go     # CountryHandler
    ├── regions.go       # RegionHandler
    └── region_emp.go    # RegionEmpHandler
```

### 3. Patrón Handler struct

Cada archivo en `handlers/` sigue el mismo patrón:

```go
package handlers

import "database/sql"

type EmployeeHandler struct{ DB *sql.DB }

func NewEmployeeHandler(db *sql.DB) *EmployeeHandler {
    return &EmployeeHandler{DB: db}
}

// Handle despacha a list/get/create/update/delete según Method + presencia de ID en path
func (h *EmployeeHandler) Handle(w http.ResponseWriter, r *http.Request) { ... }
func (h *EmployeeHandler) list(w http.ResponseWriter, r *http.Request)   { ... }
func (h *EmployeeHandler) get(w http.ResponseWriter, r *http.Request)    { ... }
func (h *EmployeeHandler) create(w http.ResponseWriter, r *http.Request) { ... }
func (h *EmployeeHandler) update(w http.ResponseWriter, r *http.Request) { ... }
func (h *EmployeeHandler) delete(w http.ResponseWriter, r *http.Request) { ... }
```

### 4. router.go

`setupRouter(db *sql.DB) http.Handler` instancia todos los handlers e registra rutas en `http.NewServeMux()`. Contiene también `corsMiddleware`.

### 5. main.go final

```go
func main() {
    if err := initDB(); err != nil {
        log.Fatalf("[FATAL] Failed to initialize database: %v", err)
    }
    defer closeDB()

    port := os.Getenv("API_PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("[INFO] Starting server on port %s", port)
    log.Fatal(http.ListenAndServe(":"+port, setupRouter(db)))
}
```

---

## Modelo de paquetes

- `package main` — `main.go`, `db.go`, `router.go`
- `package handlers` — `handlers/models.go` + todos los handlers

`models.go` se mueve a `handlers/models.go` con `package handlers`. `router.go` (package main) importa `handlers` para instanciar los structs. `db.go` permanece en `package main` exponiendo `var db *sql.DB` — `router.go` lo pasa a cada `NewXxxHandler(db)`.

---

## Entrega

- Un solo commit atómico: fix de tabla + reestructura completa
- El proyecto debe compilar (`go build ./...`) y todos los endpoints deben funcionar igual que antes
