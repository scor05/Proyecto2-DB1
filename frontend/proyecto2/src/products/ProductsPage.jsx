import { useEffect, useMemo, useState } from "react"
import {
  createProduct,
  deleteProduct,
  fetchCategories,
  fetchProducts,
  fetchProviders,
  updateProduct,
} from "../api/client.js"
import { ProductCard } from "./ProductCard.jsx"
import { ProductForm } from "./ProductForm.jsx"
import { ProductModal } from "./ProductModal.jsx"
import "./ProductsPage.css"

export function ProductsPage() {
  const [products, setProducts] = useState([])
  const [categories, setCategories] = useState([])
  const [providers, setProviders] = useState([])
  const [search, setSearch] = useState("")
  const [categoryFilter, setCategoryFilter] = useState("")
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")
  const [success, setSuccess] = useState("")
  const [modal, setModal] = useState(null)

  useEffect(() => {
    let active = true

    Promise.all([
      fetchProducts(),
      fetchCategories(),
      fetchProviders(),
    ])
      .then(([productsData, categoriesData, providersData]) => {
        if (!active) {
          return
        }
        setProducts(productsData)
        setCategories(categoriesData)
        setProviders(providersData)
      })
      .catch(() => {
        if (active) {
          setError("No se pudieron cargar los productos. Verifica que el backend y la base de datos estén activos.")
        }
      })
      .finally(() => {
        if (active) {
          setLoading(false)
        }
      })

    return () => {
      active = false
    }
  }, [])

  const filteredProducts = useMemo(() => {
    const normalizedSearch = search.trim().toLowerCase()
    return products.filter((product) => {
      const matchesSearch = !normalizedSearch || product.nombre.toLowerCase().includes(normalizedSearch)
      const matchesCategory = !categoryFilter || product.id_categoria === Number(categoryFilter)
      return matchesSearch && matchesCategory
    })
  }, [products, search, categoryFilter])

  async function handleCreate(values) {
    setError("")
    setSuccess("")
    try {
      const createdProduct = await createProduct(values)
      setProducts((current) => [...current, createdProduct].sort((a, b) => a.id_producto - b.id_producto))
      setSuccess(`Producto "${createdProduct.nombre}" agregado correctamente.`)
      setModal(null)
      return true
    } catch (err) {
      setError(friendlyProductError(err, "No se pudo agregar el producto."))
      return false
    }
  }

  async function handleUpdate(product, values) {
    setError("")
    setSuccess("")
    try {
      const updatedProduct = await updateProduct(product.id_producto, values)
      setProducts((current) => current.map((item) => (
        item.id_producto === updatedProduct.id_producto ? updatedProduct : item
      )))
      setSuccess(`Producto "${updatedProduct.nombre}" actualizado correctamente.`)
      setModal(null)
      return true
    } catch (err) {
      setError(friendlyProductError(err, "No se pudo actualizar el producto."))
      return false
    }
  }

  async function handleDelete(product) {
    setError("")
    setSuccess("")
    try {
      await deleteProduct(product.id_producto)
      setProducts((current) => current.filter((item) => item.id_producto !== product.id_producto))
      setSuccess(`Producto "${product.nombre}" eliminado correctamente.`)
      setModal(null)
    } catch (err) {
      setError(friendlyProductError(err, "No se pudo eliminar el producto."))
    }
  }

  function applyCategoryFilter(value) {
    setCategoryFilter(value)
    setModal(null)
    setSuccess(value ? "Filtro de categoría aplicado." : "Filtro de categoría removido.")
  }

  const selectedCategory = categories.find((category) => category.id_categoria === Number(categoryFilter))

  return (
    <section className="products-page">
      <div className="products-toolbar">
        <input
          className="products-search"
          type="search"
          value={search}
          onChange={(event) => setSearch(event.target.value)}
          placeholder="Buscar producto"
        />
        <button className="add-product-button" type="button" onClick={() => setModal({ type: "add" })}>
          +
        </button>
        <button className="filter-button" type="button" onClick={() => setModal({ type: "filter" })} aria-label="Filtrar por categoría">
          <svg viewBox="0 0 24 24" aria-hidden="true">
            <path d="M4 6h16l-6 7v5l-4 2v-7L4 6z" />
          </svg>
        </button>
      </div>

      {error && <p className="page-error">{error}</p>}
      {success && <p className="page-success">{success}</p>}
      {categoryFilter && (
        <p className="active-filter">
          <strong>Categoría</strong>: {selectedCategory?.nombre}
          <button type="button" onClick={() => applyCategoryFilter("")}>Quitar filtro</button>
        </p>
      )}
      {loading && <p className="products-status">Cargando productos...</p>}

      {!loading && (
        filteredProducts.length > 0 ? (
          <div className="products-grid">
            {filteredProducts.map((product) => (
              <ProductCard
                key={product.id_producto}
                product={product}
                onUpdate={(selectedProduct) => setModal({ type: "update", product: selectedProduct })}
                onDelete={(selectedProduct) => setModal({ type: "delete", product: selectedProduct })}
              />
            ))}
          </div>
        ) : (
          <p className="products-status">No se encontraron productos con los filtros actuales.</p>
        )
      )}

      {modal?.type === "add" && (
        <ProductModal title="Agregar producto" onClose={() => setModal(null)}>
          <ProductForm
            categories={categories}
            providers={providers}
            onCancel={() => setModal(null)}
            onSubmit={handleCreate}
            submitLabel="Agregar producto"
          />
        </ProductModal>
      )}

      {modal?.type === "filter" && (
        <ProductModal title="Filtrar por categoría" onClose={() => setModal(null)}>
          <div className="filter-panel">
            <label>
              Categoría
              <select value={categoryFilter} onChange={(event) => setCategoryFilter(event.target.value)}>
                <option value="">Todas las categorías</option>
                {categories.map((category) => (
                  <option key={category.id_categoria} value={category.id_categoria}>
                    {category.nombre}
                  </option>
                ))}
              </select>
            </label>
            <div className="modal-actions">
              <button className="secondary-button" type="button" onClick={() => applyCategoryFilter("")}>
                Limpiar
              </button>
              <button className="primary-button" type="button" onClick={() => applyCategoryFilter(categoryFilter)}>
                Aplicar filtro
              </button>
            </div>
          </div>
        </ProductModal>
      )}

      {modal?.type === "update" && (
        <ProductModal title={`Actualizar ${modal.product.nombre}`} onClose={() => setModal(null)}>
          <ProductForm
            product={modal.product}
            categories={categories}
            providers={providers}
            onCancel={() => setModal(null)}
            onSubmit={(values) => handleUpdate(modal.product, values)}
            submitLabel="Guardar cambios"
          />
        </ProductModal>
      )}

      {modal?.type === "delete" && (
        <ProductModal title="Eliminar producto" onClose={() => setModal(null)}>
          <p className="delete-message">
            ¿Estás seguro que deseas eliminar {modal.product.nombre}?
          </p>
          <div className="modal-actions">
            <button className="secondary-button" type="button" onClick={() => setModal(null)}>
              Cancelar
            </button>
            <button className="danger-button" type="button" onClick={() => handleDelete(modal.product)}>
              Eliminar
            </button>
          </div>
        </ProductModal>
      )}
    </section>
  )
}

function friendlyProductError(error, fallback) {
  const message = error?.message ?? ""
  if (message.includes("foreign key")) {
    return `${fallback} La categoría o proveedor seleccionado no existe.`
  }
  if (message.includes("invalid input")) {
    return `${fallback} Revisa los campos del formulario.`
  }
  if (message.includes("record not found")) {
    return `${fallback} El producto ya no existe.`
  }
  return fallback
}
