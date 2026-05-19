const pageRoles = {
  productos: ["cliente", "proveedor", "empleado", "gerente", "superadmin"],
  categorias: ["cliente", "proveedor", "empleado", "gerente", "superadmin"],
  proveedores: ["proveedor", "empleado", "gerente", "superadmin"],
  compras: ["empleado", "gerente", "superadmin"],
  clientes: ["empleado", "gerente", "superadmin"],
  empleados: ["empleado", "gerente", "superadmin"],
}

export function canViewPage(role, pageID) {
  return pageRoles[pageID]?.includes(role) ?? false
}

export function visiblePagesForRole(role, pages) {
  return pages.filter((page) => canViewPage(role, page.id))
}

export function canManageProducts(role) {
  return role === "gerente" || role === "superadmin"
}

export function canDeleteProducts(role) {
  return role === "gerente"
}

export function canManageCategories(role) {
  return role === "gerente" || role === "superadmin"
}

export function canDeleteCategories(role) {
  return role === "gerente"
}

export function canManageProviders(role) {
  return role === "gerente" || role === "superadmin"
}

export function canDeleteProviders(role) {
  return role === "gerente"
}

export function canCreateClients(role) {
  return role === "empleado" || role === "gerente" || role === "superadmin"
}

export function canEditClients(role) {
  return role === "empleado" || role === "superadmin"
}

export function canDeleteClients(role) {
  return role === "gerente"
}

export function canEditPurchases(role) {
  return role === "empleado" || role === "superadmin"
}

export function canDeletePurchases(role) {
  return role === "gerente"
}

export function canManageEmployees(role) {
  return role === "gerente" || role === "superadmin"
}

export function canDeleteEmployees(role) {
  return role === "gerente"
}
