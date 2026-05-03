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
