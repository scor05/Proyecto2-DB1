import { useEffect, useMemo, useState } from "react"
import { CategoriasPage } from "./categorias/CategoriasPage.jsx"
import { ClientesPage } from "./clientes/ClientesPage.jsx"
import { ComprasPage } from "./compras/ComprasPage.jsx"
import { EmpleadosPage } from "./empleados/EmpleadosPage.jsx"
import { LoginPage } from "./login/LoginPage.jsx"
import { ProductsPage } from "./products/ProductsPage.jsx"
import { ProveedoresPage } from "./proveedores/ProveedoresPage.jsx"
import { TopBar } from "./topbar/TopBar.jsx"
import { fetchSession, logoutUser } from "./api/client.js"
import { canViewPage, visiblePagesForRole } from "./auth/permissions.js"
import "./App.css"

export function App() {
  const [user, setUser] = useState(null)
  const [checkingSession, setCheckingSession] = useState(true)
  const [page, setPage] = useState("productos")

  useEffect(() => {
    let active = true

    fetchSession()
      .then((sessionUser) => {
        if (active) {
          setUser(sessionUser)
        }
      })
      .catch(() => {
        if (active) {
          setUser(null)
        }
      })
      .finally(() => {
        if (active) {
          setCheckingSession(false)
        }
      })

    return () => {
      active = false
    }
  }, [])

  const pages = useMemo(() => [
    { id: "productos", label: "Productos" },
    { id: "compras", label: "Compras" },
    { id: "proveedores", label: "Proveedores" },
    { id: "categorias", label: "Categorías" },
    { id: "clientes", label: "Clientes" },
    { id: "empleados", label: "Empleados" },
  ], [])

  const visiblePages = useMemo(() => visiblePagesForRole(user?.rol, pages), [pages, user?.rol])

  useEffect(() => {
    if (user && !canViewPage(user.rol, page)) {
      setPage(visiblePages[0]?.id ?? "productos")
    }
  }, [page, user, visiblePages])

  async function handleLogout() {
    try {
      await logoutUser()
    } finally {
      setUser(null)
    }
  }

  if (checkingSession) {
    return <main className="login-page"><p className="login-status">Verificando sesión...</p></main>
  }

  if (!user) {
    return <LoginPage onLogin={setUser} />
  }

  return (
    <div className="app-shell">
      <TopBar
        employee={user}
        pages={visiblePages}
        currentPage={page}
        onNavigate={setPage}
        onLogout={handleLogout}
      />
      <main className="app-content">
        {page === "productos" && <ProductsPage user={user} />}
        {page === "compras" && <ComprasPage employee={user} />}
        {page === "proveedores" && <ProveedoresPage user={user} />}
        {page === "categorias" && <CategoriasPage user={user} />}
        {page === "clientes" && <ClientesPage user={user} />}
        {page === "empleados" && <EmpleadosPage user={user} />}
      </main>
    </div>
  )
}
