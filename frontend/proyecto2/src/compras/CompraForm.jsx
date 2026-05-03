import { useState } from "react"
import "./CompraForm.css"

function today() {
  return new Date().toISOString().slice(0, 10)
}

function initialForm(compra, employee, clients, products) {
  return {
    id_empleado: compra?.id_empleado ?? employee?.id_empleado ?? "",
    id_cliente: compra?.id_cliente ?? clients[0]?.id_cliente ?? "",
    fecha_compra: compra?.fecha_compra ?? today(),
    id_producto: compra?.id_producto ?? products[0]?.id_producto ?? "",
    total_compra: compra?.total_compra ?? "",
  }
}

function validateCompra(form) {
  if (!form.id_empleado) {
    return "Debes seleccionar un empleado."
  }
  if (!form.id_cliente) {
    return "Debes seleccionar un cliente."
  }
  if (!form.fecha_compra) {
    return "Debes seleccionar una fecha de compra."
  }
  if (!form.id_producto) {
    return "Debes seleccionar un producto."
  }
  if (Number(form.total_compra) < 0 || Number.isNaN(Number(form.total_compra))) {
    return "El total de la compra no puede ser negativo."
  }
  return ""
}

export function CompraForm({
  compra,
  employee,
  employees,
  clients,
  products,
  onSubmit,
  onCancel,
  onDelete,
  submitLabel,
}) {
  const [form, setForm] = useState({ ...initialForm(compra, employee, clients, products) })
  const [formError, setFormError] = useState("")
  const [saving, setSaving] = useState(false)

  function updateField(field, value) {
    setForm((current) => ({ ...current, [field]: value }))
  }

  async function handleSubmit(event) {
    event.preventDefault()
    const validationMessage = validateCompra(form)
    if (validationMessage) {
      setFormError(validationMessage)
      return
    }

    setFormError("")
    setSaving(true)

    const success = await onSubmit({
      id_empleado: Number(form.id_empleado),
      id_cliente: Number(form.id_cliente),
      fecha_compra: form.fecha_compra,
      id_producto: Number(form.id_producto),
      total_compra: Number(form.total_compra),
    })

    if (!success) {
      setSaving(false)
    }
  }

  return (
    <form className="compra-form" onSubmit={handleSubmit}>
      <label>
        Empleado
        <select value={form.id_empleado} onChange={(event) => updateField("id_empleado", event.target.value)} required>
          <option value="">Selecciona un empleado</option>
          {employees.map((item) => (
            <option key={item.id_empleado} value={item.id_empleado}>
              {item.nombre}
            </option>
          ))}
        </select>
      </label>
      <label>
        Cliente
        <select value={form.id_cliente} onChange={(event) => updateField("id_cliente", event.target.value)} required>
          <option value="">Selecciona un cliente</option>
          {clients.map((item) => (
            <option key={item.id_cliente} value={item.id_cliente}>
              {item.nombre}
            </option>
          ))}
        </select>
      </label>
      <label>
        Fecha de compra
        <input type="date" value={form.fecha_compra} onChange={(event) => updateField("fecha_compra", event.target.value)} required />
      </label>
      <label>
        Producto comprado
        <select value={form.id_producto} onChange={(event) => updateField("id_producto", event.target.value)} required>
          <option value="">Selecciona un producto</option>
          {products.map((item) => (
            <option key={item.id_producto} value={item.id_producto}>
              {item.nombre}
            </option>
          ))}
        </select>
      </label>
      <label>
        Total de la compra
        <input
          type="number"
          min="0"
          step="0.01"
          value={form.total_compra}
          onChange={(event) => updateField("total_compra", event.target.value)}
          required
        />
      </label>
      {formError && <p className="form-error">{formError}</p>}
      <div className="modal-actions compra-actions">
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
