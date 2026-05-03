import { useEffect, useMemo, useState } from "react"
import { CategoriasPage } from "./categorias/CategoriasPage.jsx"
import { ComprasPage } from "./compras/ComprasPage.jsx"
import { LoginPage } from "./login/LoginPage.jsx"
import { ProductsPage } from "./products/ProductsPage.jsx"
import { ProveedoresPage } from "./proveedores/ProveedoresPage.jsx"
import { TopBar } from "./topbar/TopBar.jsx"
import "./App.css"

const savedEmployeeKey = "pcfast.employee"

export function App() {
  const [employee, setEmployee] = useState(() => {
    const savedEmployee = localStorage.getItem(savedEmployeeKey)
    return savedEmployee ? JSON.parse(savedEmployee) : null
  })
  const [page, setPage] = useState("productos")

  useEffect(() => {
    if (employee) {
      localStorage.setItem(savedEmployeeKey, JSON.stringify(employee))
    } else {
      localStorage.removeItem(savedEmployeeKey)
    }
  }, [employee])

  const pages = useMemo(() => [
    { id: "productos", label: "Productos" },
    { id: "compras", label: "Compras" },
    { id: "proveedores", label: "Proveedores" },
    { id: "categorias", label: "Categorías" },
  ], [])

  if (!employee) {
    return <LoginPage onLogin={setEmployee} />
  }

  return (
    <div className="app-shell">
      <TopBar
        employee={employee}
        pages={pages}
        currentPage={page}
        onNavigate={setPage}
        onLogout={() => setEmployee(null)}
      />
      <main className="app-content">
        {page === "productos" && <ProductsPage />}
        {page === "compras" && <ComprasPage employee={employee} />}
        {page === "proveedores" && <ProveedoresPage />}
        {page === "categorias" && <CategoriasPage />}
      </main>
    </div>
  )
}
