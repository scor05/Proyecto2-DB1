import { useEffect, useMemo, useState } from "react"
import {
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
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")
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
          setError("No se pudieron cargar los productos.")
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
    if (!normalizedSearch) {
      return products
    }
    return products.filter((product) => product.nombre.toLowerCase().includes(normalizedSearch))
  }, [products, search])

  async function handleUpdate(product, values) {
    setError("")
    try {
      const updatedProduct = await updateProduct(product.id_producto, values)
      setProducts((current) => current.map((item) => (
        item.id_producto === updatedProduct.id_producto ? updatedProduct : item
      )))
      setModal(null)
    } catch {
      setError("No se pudo actualizar el producto.")
    }
  }

  async function handleDelete(product) {
    setError("")
    try {
      await deleteProduct(product.id_producto)
      setProducts((current) => current.filter((item) => item.id_producto !== product.id_producto))
      setModal(null)
    } catch {
      setError("No se pudo eliminar el producto.")
    }
  }

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
      </div>

      {error && <p className="page-error">{error}</p>}
      {loading && <p className="products-status">Cargando productos...</p>}

      {!loading && (
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
      )}

      {modal?.type === "add" && (
        <ProductModal title="Agregar producto" onClose={() => setModal(null)}>
          <p className="pending-message">La creación de productos se agregará en el siguiente paso.</p>
          <div className="modal-actions">
            <button className="primary-button" type="button" onClick={() => setModal(null)}>
              Entendido
            </button>
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
