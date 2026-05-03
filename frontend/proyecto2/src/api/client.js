const apiUrl = import.meta.env.VITE_API_URL ?? "http://localhost:8080"

export async function apiRequest(path, options = {}) {
  const response = await fetch(`${apiUrl}${path}`, {
    headers: {
      "Content-Type": "application/json",
      ...(options.headers ?? {}),
    },
    ...options,
  })

  if (response.status === 204) {
    return null
  }

  const payload = await response.json()
  if (!response.ok) {
    throw new Error(payload.error ?? "Ocurrió un error inesperado")
  }

  return payload
}

export const loginEmployee = (correo) => apiRequest("/api/login", {
  method: "POST",
  body: JSON.stringify({ correo }),
})

export const fetchProducts = () => apiRequest("/api/productos")
export const fetchCategories = () => apiRequest("/api/categorias")
export const fetchProviders = () => apiRequest("/api/proveedores")
export const fetchEmployees = () => apiRequest("/api/empleados")
export const fetchClients = () => apiRequest("/api/clientes")
export const fetchCompras = () => apiRequest("/api/compras")

export const createProduct = (product) => apiRequest("/api/productos", {
  method: "POST",
  body: JSON.stringify(product),
})

export const updateProduct = (id, product) => apiRequest(`/api/productos/${id}`, {
  method: "PUT",
  body: JSON.stringify(product),
})

export const deleteProduct = (id) => apiRequest(`/api/productos/${id}`, {
  method: "DELETE",
})

export const createCompra = (compra) => apiRequest("/api/compras", {
  method: "POST",
  body: JSON.stringify(compra),
})

export const updateCompra = (id, compra) => apiRequest(`/api/compras/${id}`, {
  method: "PUT",
  body: JSON.stringify(compra),
})

export const deleteCompra = (id) => apiRequest(`/api/compras/${id}`, {
  method: "DELETE",
})

export const createCategory = (category) => apiRequest("/api/categorias", {
  method: "POST",
  body: JSON.stringify(category),
})

export const updateCategory = (id, category) => apiRequest(`/api/categorias/${id}`, {
  method: "PUT",
  body: JSON.stringify(category),
})

export const deleteCategory = (id) => apiRequest(`/api/categorias/${id}`, {
  method: "DELETE",
})

export const createProvider = (provider) => apiRequest("/api/proveedores", {
  method: "POST",
  body: JSON.stringify(provider),
})

export const updateProvider = (id, provider) => apiRequest(`/api/proveedores/${id}`, {
  method: "PUT",
  body: JSON.stringify(provider),
})

export const deleteProvider = (id) => apiRequest(`/api/proveedores/${id}`, {
  method: "DELETE",
})

export const createEmployee = (employee) => apiRequest("/api/empleados", {
  method: "POST",
  body: JSON.stringify(employee),
})

export const updateEmployee = (id, employee) => apiRequest(`/api/empleados/${id}`, {
  method: "PUT",
  body: JSON.stringify(employee),
})

export const deleteEmployee = (id) => apiRequest(`/api/empleados/${id}`, {
  method: "DELETE",
})

export const createClient = (client) => apiRequest("/api/clientes", {
  method: "POST",
  body: JSON.stringify(client),
})

export const updateClient = (id, client) => apiRequest(`/api/clientes/${id}`, {
  method: "PUT",
  body: JSON.stringify(client),
})

export const deleteClient = (id) => apiRequest(`/api/clientes/${id}`, {
  method: "DELETE",
})
