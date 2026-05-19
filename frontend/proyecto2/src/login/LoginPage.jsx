import { useState } from "react"
import { loginUser } from "../api/client.js"
import "./LoginPage.css"

export function LoginPage({ onLogin }) {
    const [correo, setCorreo] = useState("")
    const [password, setPassword] = useState("")
    const [error, setError] = useState("")
    const [loading, setLoading] = useState(false)

    async function handleSubmit(event) {
        event.preventDefault()
        setError("")
        setLoading(true)

        try {
            const user = await loginUser(correo, password)
            onLogin(user)
        } catch {
            setError("Correo o contraseña inválidos.")
        } finally {
            setLoading(false)
        }
    }

    return (
        <main className="login-page">
            <form className="login-panel" onSubmit={handleSubmit}>
                <h1>PCFast</h1>
                <label htmlFor="correo-usuario">Correo</label>
                <input
                    id="correo-usuario"
                    type="email"
                    value={correo}
                    onChange={(event) => setCorreo(event.target.value)}
                    placeholder="usuario@pcfast.com"
                    autoComplete="email"
                    required
                />
                <label htmlFor="password-usuario">Contraseña</label>
                <input
                    id="password-usuario"
                    type="password"
                    value={password}
                    onChange={(event) => setPassword(event.target.value)}
                    placeholder="contraseña"
                    autoComplete="current-password"
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
