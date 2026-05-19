package routes

import (
	"net/http"

	database "proyecto2/backend/services"
)

func RegisterRoutes(mux *http.ServeMux, manager *database.Manager) {
	catalogReaders := []string{"cliente", "proveedor", "empleado", "gerente", "superadmin"}
	providerReaders := []string{"proveedor", "empleado", "gerente", "superadmin"}
	operationsReaders := []string{"empleado", "gerente", "superadmin"}
	operationsEditors := []string{"empleado", "superadmin"}
	clientCreators := []string{"empleado", "gerente", "superadmin"}
	managersAndAdmins := []string{"gerente", "superadmin"}
	managersOnly := []string{"gerente"}

	catalogRead := RequireRoles(manager, catalogReaders...)
	providerRead := RequireRoles(manager, providerReaders...)
	operationsRead := RequireRoles(manager, operationsReaders...)
	operationsEdit := RequireRoles(manager, operationsEditors...)
	clientCreate := RequireRoles(manager, clientCreators...)
	managerOrAdmin := RequireRoles(manager, managersAndAdmins...)
	managerOnly := RequireRoles(manager, managersOnly...)

	mux.HandleFunc("GET /api/health", Health(manager))
	mux.HandleFunc("POST /api/login", LoginEmployee(manager))
	mux.HandleFunc("POST /api/logout", Logout(manager))
	mux.HandleFunc("GET /api/session", CurrentSession(manager))
	mux.HandleFunc("POST /api/productos", managerOrAdmin(CreateProduct(manager)))
	mux.HandleFunc("GET /api/productos", catalogRead(IndexProducts(manager)))
	mux.HandleFunc("GET /api/productos/{id}", catalogRead(ShowProduct(manager)))
	mux.HandleFunc("PUT /api/productos/{id}", managerOrAdmin(UpdateProduct(manager)))
	mux.HandleFunc("PATCH /api/productos/{id}", managerOrAdmin(PatchProduct(manager)))
	mux.HandleFunc("DELETE /api/productos/{id}", managerOnly(DestroyProduct(manager)))
	mux.HandleFunc("POST /api/categorias", managerOrAdmin(CreateCategory(manager)))
	mux.HandleFunc("GET /api/categorias", catalogRead(ListCategories(manager)))
	mux.HandleFunc("GET /api/categorias/{id}", catalogRead(ShowCategory(manager)))
	mux.HandleFunc("PUT /api/categorias/{id}", managerOrAdmin(UpdateCategory(manager)))
	mux.HandleFunc("DELETE /api/categorias/{id}", managerOnly(DestroyCategory(manager)))
	mux.HandleFunc("POST /api/proveedores", managerOrAdmin(CreateProvider(manager)))
	mux.HandleFunc("GET /api/proveedores", providerRead(ListProviders(manager)))
	mux.HandleFunc("GET /api/proveedores/{id}", providerRead(ShowProvider(manager)))
	mux.HandleFunc("PUT /api/proveedores/{id}", managerOrAdmin(UpdateProvider(manager)))
	mux.HandleFunc("DELETE /api/proveedores/{id}", managerOnly(DestroyProvider(manager)))
	mux.HandleFunc("POST /api/empleados", managerOrAdmin(CreateEmployee(manager)))
	mux.HandleFunc("GET /api/empleados", operationsRead(ListEmployees(manager)))
	mux.HandleFunc("GET /api/empleados/{id}", operationsRead(ShowEmployee(manager)))
	mux.HandleFunc("PUT /api/empleados/{id}", managerOrAdmin(UpdateEmployee(manager)))
	mux.HandleFunc("DELETE /api/empleados/{id}", managerOnly(DestroyEmployee(manager)))
	mux.HandleFunc("POST /api/clientes", clientCreate(CreateClient(manager)))
	mux.HandleFunc("GET /api/clientes", operationsRead(ListClients(manager)))
	mux.HandleFunc("GET /api/clientes/{id}", operationsRead(ShowClient(manager)))
	mux.HandleFunc("PUT /api/clientes/{id}", operationsEdit(UpdateClient(manager)))
	mux.HandleFunc("DELETE /api/clientes/{id}", managerOnly(DestroyClient(manager)))
	mux.HandleFunc("POST /api/compras", operationsEdit(CreateCompra(manager)))
	mux.HandleFunc("GET /api/compras", operationsRead(IndexCompras(manager)))
	mux.HandleFunc("GET /api/compras/{id}", operationsRead(ShowCompra(manager)))
	mux.HandleFunc("PUT /api/compras/{id}", operationsEdit(UpdateCompra(manager)))
	mux.HandleFunc("DELETE /api/compras/{id}", managerOnly(DestroyCompra(manager)))
}

func Health(manager *database.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusOK
		databaseStatus := "ok"

		if err := manager.Ping(r.Context()); err != nil {
			status = http.StatusServiceUnavailable
			databaseStatus = "error"
		}

		WriteJSON(w, status, map[string]string{
			"status":   "ok",
			"backend":  "ok",
			"database": databaseStatus,
		})
	}
}
