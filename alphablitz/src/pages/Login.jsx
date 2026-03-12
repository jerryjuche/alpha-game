import { useState } from "react";
import { useNavigate } from "react-router-dom";
import './GameRoom.css'

function Login() {
    const navigate = useNavigate()
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [error, setError] = useState("")

    async function login() {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                email: email,
                password: password
            })
        })

        if (response.status === 200) {
            const data = await response.json()
            localStorage.setItem('token', data.token)
            navigate('/lobby')
        } else {
            setError("Invalid Credentials")
        }


    }

    return (
        <div className="login-container">
            <div className="container">
                <label> Email:
                    <input type="email" placeholder="yourcompany@gmail.com" value={email} onChange={(e) => setEmail(e.target.value)} />
                </label>
                <label> Password:
                    <input type="password" className="" placeholder="password" value={password} onChange={(e) => setPassword(e.target.value)} />
                </label>
                {error && <p className="error">{error}</p>}
                <button className="" onClick={login}>login</button>
            </div>
        </div>
    )
}

export default Login