import { useState } from "react";
import { useNavigate } from "react-router-dom";

function Login() {
    const navigate = useNavigate()
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [error, setError] = useState("")

    async function login() {
        const response = await fetch('http://localhost:8080/auth/login', {
            method: 'POST' ,
            headers: {
                'Content-Type' : 'application/json'
            },
            body: JSON.stringify ({
                email: email,
                password: password
            })
        })

        if (response.status === 200 ) {
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
            <input type="email" className="" placeholder="yourcompany@gmail.com" value={email} onChange={(e) => setEmail(e.target.value)}/>
            <input type="password" className="" placeholder="password" value={password} onChange={(e) => setPassword(e.target.value)} />
            {error && <p className="error">{error}</p>}
            <button className="" onClick={login}>login</button>
        </div>
       </div>
    )
}

export default Login