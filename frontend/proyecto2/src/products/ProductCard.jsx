import "./ProductCard.css"

function isImageUrl(value) {
  return typeof value === "string" && /^(https?:\/\/|data:image\/)/i.test(value)
}

export function ProductCard({ product, canEdit, canDelete, onUpdate, onDelete }) {
  return (
    <article className="product-card">
      {(canEdit || canDelete) && (
        <div className="product-card-actions">
          {canEdit && (
            <button
              className="product-action-button"
              type="button"
              aria-label={`Actualizar ${product.nombre}`}
              onClick={() => onUpdate(product)}
            >
              ✎
            </button>
          )}
          {canDelete && (
            <button
              className="product-action-button product-delete-button"
              type="button"
              aria-label={`Eliminar ${product.nombre}`}
              onClick={() => onDelete(product)}
            >
              🗑
            </button>
          )}
        </div>
      )}

      <div className="product-image-slot">
        {isImageUrl(product.imagen) ? (
          <img src={product.imagen} alt={product.nombre} />
        ) : (
          <span>Sin imagen</span>
        )}
      </div>

      <h2>{product.nombre}</h2>
      <p><strong>Descripción</strong>: {product.descripcion}</p>
      <p><strong>Precio</strong>: Q{Number(product.precio).toFixed(2)}</p>
      <p><strong>Stock</strong>: {product.stock}</p>
      <p><strong>Categoría</strong>: {product.categoria}</p>
    </article>
  )
}
