import { useState } from "react"
import "./ProveedorForm.css"

function initialForm(provider) {
  return {
    nombre: provider?.nombre ?? "",
    telefono: provider?.telefono ?? "",
    correo: provider?.correo ?? "",
    direccion: provider?.direccion ?? "",
  }
}

function validateProvider(form) {
  if (form.nombre.trim() === "") {
    return "El nombre del proveedor es obligatorio."
  }
  if (form.telefono.trim() === "") {
    return "El teléfono del proveedor es obligatorio."
  }
  if (form.correo.trim() === "") {
    return "El correo del proveedor es obligatorio."
  }
  if (!form.correo.includes("@")) {
    return "El correo del proveedor no es válido."
  }
  if (form.direccion.trim() === "") {
    return "La dirección del proveedor es obligatoria."
  }
  return ""
}

export function ProveedorForm({ provider, onSubmit, onCancel, onDelete, submitLabel }) {
  const [form, setForm] = useState({ ...initialForm(provider) })
  const [formError, setFormError] = useState("")
  const [saving, setSaving] = useState(false)

  function updateField(field, value) {
    setForm((current) => ({ ...current, [field]: value }))
  }

  async function handleSubmit(event) {
    event.preventDefault()
    const validationMessage = validateProvider(form)
    if (validationMessage) {
      setFormError(validationMessage)
      return
    }

    setFormError("")
    setSaving(true)

    const success = await onSubmit({
      nombre: form.nombre,
      telefono: form.telefono,
      correo: form.correo,
      direccion: form.direccion,
    })

    if (!success) {
      setSaving(false)
    }
  }

  return (
    <form className="provider-form" onSubmit={handleSubmit}>
      <label>
        Nombre
        <input value={form.nombre} onChange={(event) => updateField("nombre", event.target.value)} required />
      </label>
      <label>
        Teléfono
        <input value={form.telefono} onChange={(event) => updateField("telefono", event.target.value)} required />
      </label>
      <label>
        Correo
        <input type="email" value={form.correo} onChange={(event) => updateField("correo", event.target.value)} required />
      </label>
      <label>
        Dirección
        <input value={form.direccion} onChange={(event) => updateField("direccion", event.target.value)} required />
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
