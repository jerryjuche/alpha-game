import { useState } from "react"
import { useNavigate } from "react-router-dom"

function Register() {
    const navigate = useNavigate()
    const [username, setUsername] = useState("")
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [error, setError] = useState("")

    async function signUp() {
        
        const response = await fetch('http://localhost:8080/auth/register', {
            method: 'POST',
            headers: {
                'Content-Type' : 'application/json'
            },
            body: JSON.stringify({
                username : username,
                email: email,
                password: password

            })
        })

        if (response.status === 201) {
            navigate('/login')
        } else {
            setError("Registration failed, try again later!")
        }

    }

    return(
        <div className="register-container">
            <div className="register">
            <input type="text" className="" placeholder="Username" value={username} onChange={ (e) => setUsername(e.target.value)}/>
            <input type="email" className="" placeholder="Email" value={email} onChange={ (e) => setEmail(e.target.value)}/>
            <input type="password" className="" placeholder="Password" value={password} onChange={ (e) => setPassword(e.target.value)}/>
            {error && <p className="error">{error}</p>}
            <button className="btn" onClick={signUp}>sign up</button>
            </div>
        </div>
    )
}

export default Register