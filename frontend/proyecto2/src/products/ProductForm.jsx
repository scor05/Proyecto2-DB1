import { useState } from "react"
import "./ProductForm.css"

export function ProductForm({ product, categories, providers, onSubmit, onCancel }) {
  const [form, setForm] = useState({
    nombre: product.nombre,
    descripcion: product.descripcion,
    precio: product.precio,
    stock: product.stock,
    id_categoria: product.id_categoria,
    id_proveedor: product.id_proveedor,
    imagen: product.imagen ?? "",
  })
  const [saving, setSaving] = useState(false)

  function updateField(field, value) {
    setForm((current) => ({ ...current, [field]: value }))
  }

  async function handleSubmit(event) {
    event.preventDefault()
    setSaving(true)

    await onSubmit({
      nombre: form.nombre,
      descripcion: form.descripcion,
      precio: Number(form.precio),
      stock: Number(form.stock),
      id_categoria: Number(form.id_categoria),
      id_proveedor: Number(form.id_proveedor),
      imagen: form.imagen.trim() === "" ? null : form.imagen,
    })

    setSaving(false)
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
      <div className="modal-actions">
        <button className="secondary-button" type="button" onClick={onCancel}>
          Cancelar
        </button>
        <button className="primary-button" type="submit" disabled={saving}>
          {saving ? "Guardando..." : "Guardar cambios"}
        </button>
      </div>
    </form>
  )
}
