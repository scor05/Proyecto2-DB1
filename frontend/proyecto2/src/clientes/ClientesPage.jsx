import { useEffect, useMemo, useState } from "react"
import {
  createClient,
  deleteClient,
  fetchClients,
  updateClient,
} from "../api/client.js"
import { ProductModal } from "../products/ProductModal.jsx"
import "../categorias/CategoriasPage.css"
import { ClienteForm } from "./ClienteForm.jsx"

export function ClientesPage() {
  const [clients, setClients] = useState([])
  const [search, setSearch] = useState("")
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")
  const [success, setSuccess] = useState("")
  const [modal, setModal] = useState(null)

  useEffect(() => {
    let active = true

    fetchClients()
      .then((data) => {
        if (active) {
          setClients(data)
        }
      })
      .catch(() => {
        if (active) {
          setError("No se pudieron cargar los clientes.")
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

  const filteredClients = useMemo(() => {
    const normalizedSearch = search.trim().toLowerCase()
    if (!normalizedSearch) {
      return clients
    }

    return clients.filter((client) => {
      const values = [
        client.id_cliente,
        client.nombre,
        client.telefono,
        client.correo,
      ].map((value) => String(value).toLowerCase())
      return values.some((value) => value.includes(normalizedSearch))
    })
  }, [clients, search])

  async function handleCreate(values) {
    setError("")
    setSuccess("")
    try {
      const createdClient = await createClient(values)
      setClients((current) => [...current, createdClient].sort((a, b) => a.id_cliente - b.id_cliente))
      setSuccess(`Cliente "${createdClient.nombre}" agregado correctamente.`)
      setModal(null)
      return true
    } catch (err) {
      setError(friendlyClientError(err, "No se pudo agregar el cliente."))
      return false
    }
  }

  async function handleUpdate(client, values) {
    setError("")
    setSuccess("")
    try {
      const updatedClient = await updateClient(client.id_cliente, values)
      setClients((current) => current.map((item) => (
        item.id_cliente === updatedClient.id_cliente ? updatedClient : item
      )))
      setSuccess(`Cliente "${updatedClient.nombre}" actualizado correctamente.`)
      setModal(null)
      return true
    } catch (err) {
      setError(friendlyClientError(err, "No se pudo actualizar el cliente."))
      return false
    }
  }

  async function handleDelete(client) {
    setError("")
    setSuccess("")
    try {
      await deleteClient(client.id_cliente)
      setClients((current) => current.filter((item) => item.id_cliente !== client.id_cliente))
      setSuccess(`Cliente "${client.nombre}" eliminado correctamente.`)
      setModal(null)
    } catch (err) {
      setError(friendlyClientError(err, "No se pudo eliminar el cliente."))
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
          placeholder="Buscar cliente"
        />
        <button className="entity-add-button" type="button" onClick={() => setModal({ type: "add" })}>
          +
        </button>
      </div>

      {error && <p className="page-error">{error}</p>}
      {success && <p className="page-success">{success}</p>}
      {loading && <p className="entity-status">Cargando clientes...</p>}

      {!loading && (
        filteredClients.length > 0 ? (
          <div className="entity-table-wrap">
            <table className="entity-table">
              <thead>
                <tr>
                  <th>id_cliente</th>
                  <th>nombre</th>
                  <th>telefono</th>
                  <th>correo</th>
                </tr>
              </thead>
              <tbody>
                {filteredClients.map((client) => (
                  <tr key={client.id_cliente} onClick={() => setModal({ type: "edit", client })}>
                    <td>{client.id_cliente}</td>
                    <td>{client.nombre}</td>
                    <td>{client.telefono}</td>
                    <td>{client.correo}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <p className="entity-status">No se encontraron clientes con la búsqueda actual.</p>
        )
      )}

      {modal?.type === "add" && (
        <ProductModal title="Agregar cliente" onClose={() => setModal(null)}>
          <ClienteForm onCancel={() => setModal(null)} onSubmit={handleCreate} submitLabel="Agregar cliente" />
        </ProductModal>
      )}

      {modal?.type === "edit" && (
        <ProductModal title={`Actualizar cliente #${modal.client.id_cliente}`} onClose={() => setModal(null)}>
          <ClienteForm
            client={modal.client}
            onCancel={() => setModal(null)}
            onSubmit={(values) => handleUpdate(modal.client, values)}
            onDelete={() => setModal({ type: "delete", client: modal.client })}
            submitLabel="Guardar cambios"
          />
        </ProductModal>
      )}

      {modal?.type === "delete" && (
        <ProductModal title="Eliminar cliente" onClose={() => setModal(null)}>
          <p className="delete-message">
            ¿Estás seguro que deseas eliminar el cliente "{modal.client.nombre}"?
          </p>
          <div className="modal-actions">
            <button className="secondary-button" type="button" onClick={() => setModal({ type: "edit", client: modal.client })}>
              Cancelar
            </button>
            <button className="danger-button" type="button" onClick={() => handleDelete(modal.client)}>
              Eliminar
            </button>
          </div>
        </ProductModal>
      )}
    </section>
  )
}

function friendlyClientError(error, fallback) {
  const message = error?.message ?? ""
  if (message.includes("foreign key")) {
    return `${fallback} El cliente está asociado a compras existentes.`
  }
  if (message.includes("invalid input")) {
    return `${fallback} Revisa los campos del formulario.`
  }
  if (message.includes("record not found")) {
    return `${fallback} El cliente ya no existe.`
  }
  return fallback
}
