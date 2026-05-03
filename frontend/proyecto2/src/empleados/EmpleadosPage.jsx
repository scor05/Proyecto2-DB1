import { useEffect, useMemo, useState } from "react"
import {
  createEmployee,
  deleteEmployee,
  fetchEmployees,
  updateEmployee,
} from "../api/client.js"
import { ProductModal } from "../products/ProductModal.jsx"
import "../categorias/CategoriasPage.css"
import { EmpleadoForm } from "./EmpleadoForm.jsx"

export function EmpleadosPage() {
  const [employees, setEmployees] = useState([])
  const [search, setSearch] = useState("")
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")
  const [success, setSuccess] = useState("")
  const [modal, setModal] = useState(null)

  useEffect(() => {
    let active = true

    fetchEmployees()
      .then((data) => {
        if (active) {
          setEmployees(data)
        }
      })
      .catch(() => {
        if (active) {
          setError("No se pudieron cargar los empleados.")
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

  const filteredEmployees = useMemo(() => {
    const normalizedSearch = search.trim().toLowerCase()
    if (!normalizedSearch) {
      return employees
    }

    return employees.filter((employee) => {
      const values = [
        employee.id_empleado,
        employee.nombre,
        employee.estado,
        employee.correo,
      ].map((value) => String(value).toLowerCase())
      return values.some((value) => value.includes(normalizedSearch))
    })
  }, [employees, search])

  async function handleCreate(values) {
    setError("")
    setSuccess("")
    try {
      const createdEmployee = await createEmployee(values)
      setEmployees((current) => [...current, createdEmployee].sort((a, b) => a.id_empleado - b.id_empleado))
      setSuccess(`Empleado "${createdEmployee.nombre}" agregado correctamente.`)
      setModal(null)
      return true
    } catch (err) {
      setError(friendlyEmployeeError(err, "No se pudo agregar el empleado."))
      return false
    }
  }

  async function handleUpdate(employee, values) {
    setError("")
    setSuccess("")
    try {
      const updatedEmployee = await updateEmployee(employee.id_empleado, values)
      setEmployees((current) => current.map((item) => (
        item.id_empleado === updatedEmployee.id_empleado ? updatedEmployee : item
      )))
      setSuccess(`Empleado "${updatedEmployee.nombre}" actualizado correctamente.`)
      setModal(null)
      return true
    } catch (err) {
      setError(friendlyEmployeeError(err, "No se pudo actualizar el empleado."))
      return false
    }
  }

  async function handleDelete(employee) {
    setError("")
    setSuccess("")
    try {
      await deleteEmployee(employee.id_empleado)
      setEmployees((current) => current.filter((item) => item.id_empleado !== employee.id_empleado))
      setSuccess(`Empleado "${employee.nombre}" eliminado correctamente.`)
      setModal(null)
    } catch (err) {
      setError(friendlyEmployeeError(err, "No se pudo eliminar el empleado."))
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
          placeholder="Buscar empleado"
        />
        <button className="entity-add-button" type="button" onClick={() => setModal({ type: "add" })}>
          +
        </button>
      </div>

      {error && <p className="page-error">{error}</p>}
      {success && <p className="page-success">{success}</p>}
      {loading && <p className="entity-status">Cargando empleados...</p>}

      {!loading && (
        filteredEmployees.length > 0 ? (
          <div className="entity-table-wrap">
            <table className="entity-table">
              <thead>
                <tr>
                  <th>id_empleado</th>
                  <th>nombre</th>
                  <th>estado</th>
                  <th>correo</th>
                </tr>
              </thead>
              <tbody>
                {filteredEmployees.map((employee) => (
                  <tr key={employee.id_empleado} onClick={() => setModal({ type: "edit", employee })}>
                    <td>{employee.id_empleado}</td>
                    <td>{employee.nombre}</td>
                    <td>{employee.estado}</td>
                    <td>{employee.correo}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <p className="entity-status">No se encontraron empleados con la búsqueda actual.</p>
        )
      )}

      {modal?.type === "add" && (
        <ProductModal title="Agregar empleado" onClose={() => setModal(null)}>
          <EmpleadoForm onCancel={() => setModal(null)} onSubmit={handleCreate} submitLabel="Agregar empleado" />
        </ProductModal>
      )}

      {modal?.type === "edit" && (
        <ProductModal title={`Actualizar empleado #${modal.employee.id_empleado}`} onClose={() => setModal(null)}>
          <EmpleadoForm
            employee={modal.employee}
            onCancel={() => setModal(null)}
            onSubmit={(values) => handleUpdate(modal.employee, values)}
            onDelete={() => setModal({ type: "delete", employee: modal.employee })}
            submitLabel="Guardar cambios"
          />
        </ProductModal>
      )}

      {modal?.type === "delete" && (
        <ProductModal title="Eliminar empleado" onClose={() => setModal(null)}>
          <p className="delete-message">
            ¿Estás seguro que deseas eliminar el empleado "{modal.employee.nombre}"?
          </p>
          <div className="modal-actions">
            <button className="secondary-button" type="button" onClick={() => setModal({ type: "edit", employee: modal.employee })}>
              Cancelar
            </button>
            <button className="danger-button" type="button" onClick={() => handleDelete(modal.employee)}>
              Eliminar
            </button>
          </div>
        </ProductModal>
      )}
    </section>
  )
}

function friendlyEmployeeError(error, fallback) {
  const message = error?.message ?? ""
  if (message.includes("foreign key")) {
    return `${fallback} El empleado está asociado a compras existentes.`
  }
  if (message.includes("invalid input")) {
    return `${fallback} Revisa los campos del formulario.`
  }
  if (message.includes("record not found")) {
    return `${fallback} El empleado ya no existe.`
  }
  return fallback
}
