import { useState } from "react"
import { useNavigate } from "react-router-dom"
import './GameRoom.css'

function Register() {
    const navigate = useNavigate()
    const [username, setUsername] = useState("")
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [error, setError] = useState("")

    async function signUp() {

        const response = await fetch(`${import.meta.env.VITE_API_URL}/auth/register`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                username: username,
                email: email,
                password: password

            })
        })

        if (response.status === 201) {
            const data = await response.json()
            localStorage.setItem('token', data.token)

            navigate('/lobby')
        } else {
            setError("error creating account")
        }
    }

    return (
        <div className="register-container">
            <div className="register">
                <input type="text" className="" placeholder="Username" value={username} onChange={(e) => setUsername(e.target.value)} />
                <input type="email" className="" placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} />
                <input type="password" className="" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} />
                {error && <p className="error">{error}</p>}
                <button onClick={signUp}>sign up</button>
            </div>
        </div>
    )
}

export default Register