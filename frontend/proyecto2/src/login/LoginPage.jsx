import { useState } from "react"
import { loginEmployee } from "../api/client.js"
import "./LoginPage.css"

export function LoginPage({ onLogin }) {
  const [correo, setCorreo] = useState("")
  const [error, setError] = useState("")
  const [loading, setLoading] = useState(false)

  async function handleSubmit(event) {
    event.preventDefault()
    setError("")
    setLoading(true)

    try {
      const employee = await loginEmployee(correo)
      onLogin(employee)
    } catch {
      setError("No se encontró un empleado con ese correo.")
    } finally {
      setLoading(false)
    }
  }

  return (
    <main className="login-page">
      <form className="login-panel" onSubmit={handleSubmit}>
        <h1>PCFast</h1>
        <label htmlFor="correo-empleado">Correo del empleado</label>
        <input
          id="correo-empleado"
          type="email"
          value={correo}
          onChange={(event) => setCorreo(event.target.value)}
          placeholder="empleado@pcmarket.com"
          required
        />
        {error && <p className="login-error">{error}</p>}
        <button className="primary-button" type="submit" disabled={loading}>
          {loading ? "Ingresando..." : "Iniciar sesión"}
        </button>
      </form>
    </main>
  )
}
