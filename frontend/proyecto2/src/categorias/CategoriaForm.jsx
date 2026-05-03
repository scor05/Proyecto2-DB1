import { useState } from "react"
import "./CategoriaForm.css"

function initialForm(category) {
  return {
    nombre: category?.nombre ?? "",
  }
}

function validateCategory(form) {
  if (form.nombre.trim() === "") {
    return "El nombre de la categoría es obligatorio."
  }
  return ""
}

export function CategoriaForm({ category, onSubmit, onCancel, onDelete, submitLabel }) {
  const [form, setForm] = useState({ ...initialForm(category) })
  const [formError, setFormError] = useState("")
  const [saving, setSaving] = useState(false)

  async function handleSubmit(event) {
    event.preventDefault()
    const validationMessage = validateCategory(form)
    if (validationMessage) {
      setFormError(validationMessage)
      return
    }

    setFormError("")
    setSaving(true)

    const success = await onSubmit({
      nombre: form.nombre,
    })

    if (!success) {
      setSaving(false)
    }
  }

  return (
    <form className="category-form" onSubmit={handleSubmit}>
      <label>
        Nombre
        <input value={form.nombre} onChange={(event) => setForm({ nombre: event.target.value })} required />
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
