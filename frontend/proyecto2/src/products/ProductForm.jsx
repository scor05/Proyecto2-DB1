import { useState } from "react"
import "./ProductForm.css"

function initialForm(product, categories, providers) {
  return {
    nombre: product?.nombre ?? "",
    descripcion: product?.descripcion ?? "",
    precio: product?.precio ?? "",
    stock: product?.stock ?? "",
    id_categoria: product?.id_categoria ?? categories[0]?.id_categoria ?? "",
    id_proveedor: product?.id_proveedor ?? providers[0]?.id_proveedor ?? "",
    imagen: product?.imagen ?? "",
  }
}

function validateProduct(form) {
  if (form.nombre.trim() === "") {
    return "El nombre del producto es obligatorio."
  }
  if (form.descripcion.trim() === "") {
    return "La descripción del producto es obligatoria."
  }
  if (!form.id_categoria) {
    return "Debes seleccionar una categoría."
  }
  if (!form.id_proveedor) {
    return "Debes seleccionar un proveedor."
  }
  if (Number(form.precio) <= 0 || Number.isNaN(Number(form.precio))) {
    return "El precio debe ser mayor a 0."
  }
  if (Number(form.stock) < 0 || Number.isNaN(Number(form.stock))) {
    return "El stock no puede ser negativo."
  }
  return ""
}

export function ProductForm({
  product,
  categories,
  providers,
  onSubmit,
  onCancel,
  submitLabel = "Guardar cambios",
}) {
  const [form, setForm] = useState({
    ...initialForm(product, categories, providers),
  })
  const [saving, setSaving] = useState(false)
  const [formError, setFormError] = useState("")

  function updateField(field, value) {
    setForm((current) => ({ ...current, [field]: value }))
  }

  async function handleSubmit(event) {
    event.preventDefault()
    const validationMessage = validateProduct(form)
    if (validationMessage) {
      setFormError(validationMessage)
      return
    }

    setFormError("")
    setSaving(true)

    const success = await onSubmit({
      nombre: form.nombre,
      descripcion: form.descripcion,
      precio: Number(form.precio),
      stock: Number(form.stock),
      id_categoria: Number(form.id_categoria),
      id_proveedor: Number(form.id_proveedor),
      imagen: form.imagen.trim() === "" ? null : form.imagen,
    })

    if (!success) {
      setSaving(false)
    }
  }

  return (
    <form className="product-form" onSubmit={handleSubmit}>
      <label>
        Nombre
        <input value={form.nombre} onChange={(event) => updateField("nombre", event.target.value)} required />
      </label>
      <label>
        Descripción
        <textarea value={form.descripcion} onChange={(event) => updateField("descripcion", event.target.value)} required />
      </label>
      <label>
        Precio
        <input type="number" min="0.01" step="0.01" value={form.precio} onChange={(event) => updateField("precio", event.target.value)} required />
      </label>
      <label>
        Stock
        <input type="number" min="0" step="1" value={form.stock} onChange={(event) => updateField("stock", event.target.value)} required />
      </label>
      <label>
        Categoría
        <select value={form.id_categoria} onChange={(event) => updateField("id_categoria", event.target.value)} required>
          {categories.map((category) => (
            <option key={category.id_categoria} value={category.id_categoria}>
              {category.nombre}
            </option>
          ))}
        </select>
      </label>
      <label>
        Proveedor
        <select value={form.id_proveedor} onChange={(event) => updateField("id_proveedor", event.target.value)} required>
          {providers.map((provider) => (
            <option key={provider.id_proveedor} value={provider.id_proveedor}>
              {provider.nombre}
            </option>
          ))}
        </select>
      </label>
      <label>
        Imagen
        <input value={form.imagen} onChange={(event) => updateField("imagen", event.target.value)} />
      </label>
      {formError && <p className="form-error">{formError}</p>}
      <div className="modal-actions">
        <button className="secondary-button" type="button" onClick={onCancel}>
          Cancelar
        </button>
        <button className="primary-button" type="submit" disabled={saving}>
          {saving ? "Guardando..." : submitLabel}
        </button>
      </div>
    </form>
  )
}
