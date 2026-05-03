import { useEffect, useMemo, useState } from "react"
import {
  createCategory,
  deleteCategory,
  fetchCategories,
  updateCategory,
} from "../api/client.js"
import { ProductModal } from "../products/ProductModal.jsx"
import { CategoriaForm } from "./CategoriaForm.jsx"
import "./CategoriasPage.css"

export function CategoriasPage() {
  const [categories, setCategories] = useState([])
  const [search, setSearch] = useState("")
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")
  const [success, setSuccess] = useState("")
  const [modal, setModal] = useState(null)

  useEffect(() => {
    let active = true

    fetchCategories()
      .then((data) => {
        if (active) {
          setCategories(data)
        }
      })
      .catch(() => {
        if (active) {
          setError("No se pudieron cargar las categorías.")
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

  const filteredCategories = useMemo(() => {
    const normalizedSearch = search.trim().toLowerCase()
    if (!normalizedSearch) {
      return categories
    }
    return categories.filter((category) => (
      String(category.id_categoria).includes(normalizedSearch) ||
      category.nombre.toLowerCase().includes(normalizedSearch)
    ))
  }, [categories, search])

  async function handleCreate(values) {
    setError("")
    setSuccess("")
    try {
      const createdCategory = await createCategory(values)
      setCategories((current) => [...current, createdCategory].sort((a, b) => a.id_categoria - b.id_categoria))
      setSuccess(`Categoría "${createdCategory.nombre}" agregada correctamente.`)
      setModal(null)
      return true
    } catch (err) {
      setError(friendlyCategoryError(err, "No se pudo agregar la categoría."))
      return false
    }
  }

  async function handleUpdate(category, values) {
    setError("")
    setSuccess("")
    try {
      const updatedCategory = await updateCategory(category.id_categoria, values)
      setCategories((current) => current.map((item) => (
        item.id_categoria === updatedCategory.id_categoria ? updatedCategory : item
      )))
      setSuccess(`Categoría "${updatedCategory.nombre}" actualizada correctamente.`)
      setModal(null)
      return true
    } catch (err) {
      setError(friendlyCategoryError(err, "No se pudo actualizar la categoría."))
      return false
    }
  }

  async function handleDelete(category) {
    setError("")
    setSuccess("")
    try {
      await deleteCategory(category.id_categoria)
      setCategories((current) => current.filter((item) => item.id_categoria !== category.id_categoria))
      setSuccess(`Categoría "${category.nombre}" eliminada correctamente.`)
      setModal(null)
    } catch (err) {
      setError(friendlyCategoryError(err, "No se pudo eliminar la categoría."))
    }
  }

  return (
    <section className="entity-page">
      <div className="entity-toolbar">
        <input
          className="entity-search"
          type="search"
          value={search}
          onChange={(event) => setSearch(event.target.value)}
          placeholder="Buscar categoría"
        />
        <button className="entity-add-button" type="button" onClick={() => setModal({ type: "add" })}>
          +
        </button>
      </div>

      {error && <p className="page-error">{error}</p>}
      {success && <p className="page-success">{success}</p>}
      {loading && <p className="entity-status">Cargando categorías...</p>}

      {!loading && (
        filteredCategories.length > 0 ? (
          <div className="entity-table-wrap">
            <table className="entity-table">
              <thead>
                <tr>
                  <th>id_categoria</th>
                  <th>nombre</th>
                </tr>
              </thead>
              <tbody>
                {filteredCategories.map((category) => (
                  <tr key={category.id_categoria} onClick={() => setModal({ type: "edit", category })}>
                    <td>{category.id_categoria}</td>
                    <td>{category.nombre}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <p className="entity-status">No se encontraron categorías con la búsqueda actual.</p>
        )
      )}

      {modal?.type === "add" && (
        <ProductModal title="Agregar categoría" onClose={() => setModal(null)}>
          <CategoriaForm
            onCancel={() => setModal(null)}
            onSubmit={handleCreate}
            submitLabel="Agregar categoría"
          />
        </ProductModal>
      )}

      {modal?.type === "edit" && (
        <ProductModal title={`Actualizar categoría #${modal.category.id_categoria}`} onClose={() => setModal(null)}>
          <CategoriaForm
            category={modal.category}
            onCancel={() => setModal(null)}
            onSubmit={(values) => handleUpdate(modal.category, values)}
            onDelete={() => setModal({ type: "delete", category: modal.category })}
            submitLabel="Guardar cambios"
          />
        </ProductModal>
      )}

      {modal?.type === "delete" && (
        <ProductModal title="Eliminar categoría" onClose={() => setModal(null)}>
          <p className="delete-message">
            ¿Estás seguro que deseas eliminar la categoría "{modal.category.nombre}"?
          </p>
          <div className="modal-actions">
            <button className="secondary-button" type="button" onClick={() => setModal({ type: "edit", category: modal.category })}>
              Cancelar
            </button>
            <button className="danger-button" type="button" onClick={() => handleDelete(modal.category)}>
              Eliminar
            </button>
          </div>
        </ProductModal>
      )}
    </section>
  )
}

function friendlyCategoryError(error, fallback) {
  const message = error?.message ?? ""
  if (message.includes("foreign key")) {
    return `${fallback} La categoría está asociada a productos existentes.`
  }
  if (message.includes("invalid input")) {
    return `${fallback} Revisa los campos del formulario.`
  }
  if (message.includes("record not found")) {
    return `${fallback} La categoría ya no existe.`
  }
  return fallback
}
