import { useEffect, useRef, useState } from "react"
import "./TopBar.css"

export function TopBar({ employee, pages, currentPage, onNavigate, onLogout }) {
  const [open, setOpen] = useState(false)
  const menuRef = useRef(null)

  useEffect(() => {
    function closeMenu(event) {
      if (menuRef.current && !menuRef.current.contains(event.target)) {
        setOpen(false)
      }
    }

    document.addEventListener("mousedown", closeMenu)
    return () => document.removeEventListener("mousedown", closeMenu)
  }, [])

  return (
    <header className="topbar">
      <div className="topbar-section topbar-left" ref={menuRef}>
        <button
          className="topbar-menu-button"
          type="button"
          onClick={() => setOpen((current) => !current)}
          aria-label="Abrir menú"
        >
          ☰
        </button>
        {open && (
          <nav className="topbar-dropdown">
            {pages.map((page) => (
              <button
                key={page.id}
                type="button"
                className={page.id === currentPage ? "active" : ""}
                onClick={() => {
                  onNavigate(page.id)
                  setOpen(false)
                }}
              >
                {page.label}
              </button>
            ))}
          </nav>
        )}
      </div>
      <div className="topbar-section topbar-center">
        <h1>PCFast</h1>
      </div>
      <div className="topbar-section topbar-right">
        <span>{employee.nombre}</span>
        <button type="button" onClick={onLogout}>Cerrar Sesión</button>
      </div>
    </header>
  )
}
