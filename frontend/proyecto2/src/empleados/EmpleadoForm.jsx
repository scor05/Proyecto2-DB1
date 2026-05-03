import { useState } from "react"
import "./EmpleadoForm.css"

function initialForm(employee) {
  return {
    nombre: employee?.nombre ?? "",
    estado: employee?.estado ?? "activo",
    correo: employee?.correo ?? "",
  }
}

function validateEmployee(form) {
  if (form.nombre.trim() === "") {
    return "El nombre del empleado es obligatorio."
  }
  if (form.estado !== "activo" && form.estado !== "inactivo") {
    return "El estado debe ser activo o inactivo."
  }
  if (form.correo.trim() === "") {
    return "El correo del empleado es obligatorio."
  }
  if (!form.correo.includes("@")) {
    return "El correo del empleado no es válido."
  }
  return ""
}

export function EmpleadoForm({ employee, onSubmit, onCancel, onDelete, submitLabel }) {
  const [form, setForm] = useState({ ...initialForm(employee) })
  const [formError, setFormError] = useState("")
  const [saving, setSaving] = useState(false)

  function updateField(field, value) {
    setForm((current) => ({ ...current, [field]: value }))
  }

  async function handleSubmit(event) {
    event.preventDefault()
    const validationMessage = validateEmployee(form)
    if (validationMessage) {
      setFormError(validationMessage)
      return
    }

    setFormError("")
    setSaving(true)

    const success = await onSubmit({
      nombre: form.nombre,
      estado: form.estado,
      correo: form.correo,
    })

    if (!success) {
      setSaving(false)
    }
  }

  return (
    <form className="employee-form" onSubmit={handleSubmit}>
      <label>
        Nombre
        <input value={form.nombre} onChange={(event) => updateField("nombre", event.target.value)} required />
      </label>
      <label>
        Estado
        <select value={form.estado} onChange={(event) => updateField("estado", event.target.value)} required>
          <option value="activo">activo</option>
          <option value="inactivo">inactivo</option>
        </select>
      </label>
      <label>
        Correo
        <input type="email" value={form.correo} onChange={(event) => updateField("correo", event.target.value)} required />
      </label>
      {formError && <p className="form-error">{formError}</p>}
      <div className="modal-actions entity-actions">
        <button className="secondary-button" type="button" onClick={onCancel}>
          Cancelar
        </button>
        {onDelete && (
          <button className="danger-button" type="button" onClick={onDelete}>
            Eliminar
          </button>
        )}
        <button className="primary-button" type="submit" disabled={saving}>
          {saving ? "Guardando..." : submitLabel}
        </button>
      </div>
    </form>
  )
}
