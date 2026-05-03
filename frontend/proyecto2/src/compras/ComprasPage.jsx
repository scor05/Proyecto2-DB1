import { useEffect, useMemo, useState } from "react"
import {
  createCompra,
  deleteCompra,
  fetchClients,
  fetchCompras,
  fetchEmployees,
  fetchProducts,
  updateCompra,
} from "../api/client.js"
import { ProductModal } from "../products/ProductModal.jsx"
import { CompraForm } from "./CompraForm.jsx"
import "./ComprasPage.css"

export function ComprasPage({ employee }) {
  const [compras, setCompras] = useState([])
  const [employees, setEmployees] = useState([])
  const [clients, setClients] = useState([])
  const [products, setProducts] = useState([])
  const [search, setSearch] = useState("")
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")
  const [success, setSuccess] = useState("")
  const [modal, setModal] = useState(null)

  useEffect(() => {
    let active = true

    Promise.all([fetchCompras(), fetchEmployees(), fetchClients(), fetchProducts()])
      .then(([comprasData, employeesData, clientsData, productsData]) => {
        if (!active) {
          return
        }
        setCompras(comprasData)
        setEmployees(employeesData)
        setClients(clientsData)
        setProducts(productsData)
      })
      .catch(() => {
        if (active) {
          setError("No se pudieron cargar las compras. Verifica que el backend y la base de datos estén activos.")
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

  const filteredCompras = useMemo(() => {
    const normalizedSearch = search.trim().toLowerCase()
    if (!normalizedSearch) {
      return compras
    }

    return compras.filter((compra) => {
      const values = [
        compra.id_compra,
        compra.id_empleado,
        compra.nombre_empleado,
        compra.id_cliente,
        compra.nombre_cliente,
        compra.fecha_compra,
        compra.productos,
        compra.total_compra,
      ].map((value) => String(value).toLowerCase())

      return values.some((value) => value.includes(normalizedSearch))
    })
  }, [compras, search])

  async function handleCreate(values) {
    setError("")
    setSuccess("")
    try {
      const createdCompra = await createCompra(values)
      setCompras((current) => [...current, createdCompra].sort((a, b) => a.id_compra - b.id_compra))
      setSuccess(`Compra #${createdCompra.id_compra} agregada correctamente.`)
      setModal(null)
      return true
    } catch (err) {
      setError(friendlyCompraError(err, "No se pudo agregar la compra."))
      return false
    }
  }

  async function handleUpdate(compra, values) {
    setError("")
    setSuccess("")
    try {
      const updatedCompra = await updateCompra(compra.id_compra, values)
      setCompras((current) => current.map((item) => (
        item.id_compra === updatedCompra.id_compra ? updatedCompra : item
      )))
      setSuccess(`Compra #${updatedCompra.id_compra} actualizada correctamente.`)
      setModal(null)
      return true
    } catch (err) {
      setError(friendlyCompraError(err, "No se pudo actualizar la compra."))
      return false
    }
  }

  async function handleDelete(compra) {
    setError("")
    setSuccess("")
    try {
      await deleteCompra(compra.id_compra)
      setCompras((current) => current.filter((item) => item.id_compra !== compra.id_compra))
      setSuccess(`Compra #${compra.id_compra} eliminada correctamente.`)
      setModal(null)
    } catch (err) {
      setError(friendlyCompraError(err, "No se pudo eliminar la compra."))
    }
  }

  function exportCSV() {
    setError("")
    setSuccess("")

    if (compras.length === 0) {
      setError("No hay compras disponibles para exportar.")
      return
    }

    const headers = ["id_compra", "empleado", "cliente", "fecha_compra", "productos", "total_compra"]
    const csvRows = compras.map((compra) => ({
      id_compra: compra.id_compra,
      empleado: compra.nombre_empleado,
      cliente: compra.nombre_cliente,
      fecha_compra: compra.fecha_compra,
      productos: compra.productos,
      total_compra: compra.total_compra,
    }))
    const rows = csvRows.map((compra) => headers.map((header) => csvValue(compra[header])).join(","))
    const csv = [headers.join(","), ...rows].join("\n")
    const blob = new Blob([csv], { type: "text/csv;charset=utf-8" })
    const url = URL.createObjectURL(blob)
    const link = document.createElement("a")

    link.href = url
    link.download = "compras.csv"
    link.click()
    URL.revokeObjectURL(url)
    setSuccess("Archivo CSV generado correctamente.")
  }

  return (
    <section className="compras-page">
      <div className="compras-toolbar">
        <input
          className="compras-search"
          type="search"
          value={search}
          onChange={(event) => setSearch(event.target.value)}
          placeholder="Buscar compra"
        />
        <button className="add-compra-button" type="button" onClick={() => setModal({ type: "add" })}>
          +
        </button>
        <button className="csv-button" type="button" onClick={exportCSV}>
          Exportar a CSV
        </button>
      </div>

      {error && <p className="page-error">{error}</p>}
      {success && <p className="page-success">{success}</p>}
      {loading && <p className="compras-status">Cargando compras...</p>}

      {!loading && (
        filteredCompras.length > 0 ? (
          <div className="compras-table-wrap">
            <table className="compras-table">
              <thead>
                <tr>
                  <th>id_compra</th>
                  <th>empleado</th>
                  <th>cliente</th>
                  <th>fecha_compra</th>
                  <th>productos</th>
                  <th>total_compra</th>
                </tr>
              </thead>
              <tbody>
                {filteredCompras.map((compra) => (
                  <tr key={compra.id_compra} onClick={() => setModal({ type: "edit", compra })}>
                    <td>{compra.id_compra}</td>
                    <td>{compra.nombre_empleado}</td>
                    <td>{compra.nombre_cliente}</td>
                    <td>{compra.fecha_compra}</td>
                    <td>{compra.productos}</td>
                    <td>Q{Number(compra.total_compra).toFixed(2)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <p className="compras-status">No se encontraron compras con la búsqueda actual.</p>
        )
      )}

      {modal?.type === "add" && (
        <ProductModal title="Agregar compra" onClose={() => setModal(null)}>
          <CompraForm
            employee={employee}
            employees={employees}
            clients={clients}
            products={products}
            onCancel={() => setModal(null)}
            onSubmit={handleCreate}
            submitLabel="Agregar compra"
          />
        </ProductModal>
      )}

      {modal?.type === "edit" && (
        <ProductModal title={`Actualizar compra #${modal.compra.id_compra}`} onClose={() => setModal(null)}>
          <CompraForm
            compra={modal.compra}
            employee={employee}
            employees={employees}
            clients={clients}
            products={products}
            onCancel={() => setModal(null)}
            onSubmit={(values) => handleUpdate(modal.compra, values)}
            onDelete={() => setModal({ type: "delete", compra: modal.compra })}
            submitLabel="Guardar cambios"
          />
        </ProductModal>
      )}

      {modal?.type === "delete" && (
        <ProductModal title="Eliminar compra" onClose={() => setModal(null)}>
          <p className="delete-message">
            ¿Estás seguro que deseas eliminar la compra #{modal.compra.id_compra}?
          </p>
          <div className="modal-actions">
            <button className="secondary-button" type="button" onClick={() => setModal({ type: "edit", compra: modal.compra })}>
              Cancelar
            </button>
            <button className="danger-button" type="button" onClick={() => handleDelete(modal.compra)}>
              Eliminar
            </button>
          </div>
        </ProductModal>
      )}
    </section>
  )
}

function csvValue(value) {
  const text = String(value ?? "")
  return `"${text.replaceAll('"', '""')}"`
}

function friendlyCompraError(error, fallback) {
  const message = error?.message ?? ""
  if (message.includes("foreign key")) {
    return `${fallback} El empleado, cliente o producto seleccionado no existe.`
  }
  if (message.includes("invalid input")) {
    return `${fallback} Revisa los campos del formulario.`
  }
  if (message.includes("record not found")) {
    return `${fallback} La compra ya no existe.`
  }
  if (message.includes("transaction rollback")) {
    return `${fallback} La transacción fue revertida para proteger los datos.`
  }
  if (message.includes("transaction commit failed")) {
    return `${fallback} No se pudo confirmar la transacción.`
  }
  return fallback
}
