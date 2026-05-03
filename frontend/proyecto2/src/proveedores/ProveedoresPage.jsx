import { useEffect, useMemo, useState } from "react"
import {
  createProvider,
  deleteProvider,
  fetchProviders,
  updateProvider,
} from "../api/client.js"
import { ProductModal } from "../products/ProductModal.jsx"
import { ProveedorForm } from "./ProveedorForm.jsx"
import "../categorias/CategoriasPage.css"

export function ProveedoresPage() {
  const [providers, setProviders] = useState([])
  const [search, setSearch] = useState("")
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")
  const [success, setSuccess] = useState("")
  const [modal, setModal] = useState(null)

  useEffect(() => {
    let active = true

    fetchProviders()
      .then((data) => {
        if (active) {
          setProviders(data)
        }
      })
      .catch(() => {
        if (active) {
          setError("No se pudieron cargar los proveedores.")
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

  const filteredProviders = useMemo(() => {
    const normalizedSearch = search.trim().toLowerCase()
    if (!normalizedSearch) {
      return providers
    }

    return providers.filter((provider) => {
      const values = [
        provider.id_proveedor,
        provider.nombre,
        provider.telefono,
        provider.correo,
        provider.direccion,
      ].map((value) => String(value).toLowerCase())
      return values.some((value) => value.includes(normalizedSearch))
    })
  }, [providers, search])

  async function handleCreate(values) {
    setError("")
    setSuccess("")
    try {
      const createdProvider = await createProvider(values)
      setProviders((current) => [...current, createdProvider].sort((a, b) => a.id_proveedor - b.id_proveedor))
      setSuccess(`Proveedor "${createdProvider.nombre}" agregado correctamente.`)
      setModal(null)
      return true
    } catch (err) {
      setError(friendlyProviderError(err, "No se pudo agregar el proveedor."))
      return false
    }
  }

  async function handleUpdate(provider, values) {
    setError("")
    setSuccess("")
    try {
      const updatedProvider = await updateProvider(provider.id_proveedor, values)
      setProviders((current) => current.map((item) => (
        item.id_proveedor === updatedProvider.id_proveedor ? updatedProvider : item
      )))
      setSuccess(`Proveedor "${updatedProvider.nombre}" actualizado correctamente.`)
      setModal(null)
      return true
    } catch (err) {
      setError(friendlyProviderError(err, "No se pudo actualizar el proveedor."))
      return false
    }
  }

  async function handleDelete(provider) {
    setError("")
    setSuccess("")
    try {
      await deleteProvider(provider.id_proveedor)
      setProviders((current) => current.filter((item) => item.id_proveedor !== provider.id_proveedor))
      setSuccess(`Proveedor "${provider.nombre}" eliminado correctamente.`)
      setModal(null)
    } catch (err) {
      setError(friendlyProviderError(err, "No se pudo eliminar el proveedor."))
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
          placeholder="Buscar proveedor"
        />
        <button className="entity-add-button" type="button" onClick={() => setModal({ type: "add" })}>
          +
        </button>
      </div>

      {error && <p className="page-error">{error}</p>}
      {success && <p className="page-success">{success}</p>}
      {loading && <p className="entity-status">Cargando proveedores...</p>}

      {!loading && (
        filteredProviders.length > 0 ? (
          <div className="entity-table-wrap">
            <table className="entity-table provider-table">
              <thead>
                <tr>
                  <th>id_proveedor</th>
                  <th>nombre</th>
                  <th>telefono</th>
                  <th>correo</th>
                  <th>direccion</th>
                </tr>
              </thead>
              <tbody>
                {filteredProviders.map((provider) => (
                  <tr key={provider.id_proveedor} onClick={() => setModal({ type: "edit", provider })}>
                    <td>{provider.id_proveedor}</td>
                    <td>{provider.nombre}</td>
                    <td>{provider.telefono}</td>
                    <td>{provider.correo}</td>
                    <td>{provider.direccion}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <p className="entity-status">No se encontraron proveedores con la búsqueda actual.</p>
        )
      )}

      {modal?.type === "add" && (
        <ProductModal title="Agregar proveedor" onClose={() => setModal(null)}>
          <ProveedorForm
            onCancel={() => setModal(null)}
            onSubmit={handleCreate}
            submitLabel="Agregar proveedor"
          />
        </ProductModal>
      )}

      {modal?.type === "edit" && (
        <ProductModal title={`Actualizar proveedor #${modal.provider.id_proveedor}`} onClose={() => setModal(null)}>
          <ProveedorForm
            provider={modal.provider}
            onCancel={() => setModal(null)}
            onSubmit={(values) => handleUpdate(modal.provider, values)}
            onDelete={() => setModal({ type: "delete", provider: modal.provider })}
            submitLabel="Guardar cambios"
          />
        </ProductModal>
      )}

      {modal?.type === "delete" && (
        <ProductModal title="Eliminar proveedor" onClose={() => setModal(null)}>
          <p className="delete-message">
            ¿Estás seguro que deseas eliminar el proveedor "{modal.provider.nombre}"?
          </p>
          <div className="modal-actions">
            <button className="secondary-button" type="button" onClick={() => setModal({ type: "edit", provider: modal.provider })}>
              Cancelar
            </button>
            <button className="danger-button" type="button" onClick={() => handleDelete(modal.provider)}>
              Eliminar
            </button>
          </div>
        </ProductModal>
      )}
    </section>
  )
}

function friendlyProviderError(error, fallback) {
  const message = error?.message ?? ""
  if (message.includes("foreign key")) {
    return `${fallback} El proveedor está asociado a productos existentes.`
  }
  if (message.includes("invalid input")) {
    return `${fallback} Revisa los campos del formulario.`
  }
  if (message.includes("record not found")) {
    return `${fallback} El proveedor ya no existe.`
  }
  return fallback
}
