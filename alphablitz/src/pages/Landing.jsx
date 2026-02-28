import { useNavigate } from 'react-router-dom'
import './LandingPage.css'



function LandingPage() {
    const navigate = useNavigate()
    return (
        <div className="container">
            <h1>ALPHABLITZ</h1>
            <div className="description">
                <p>One letter. Five categories. Eight seconds.</p>
            </div>
            <button className="btn" onClick={() => navigate('/register')}>Play Now</button>
        </div>
    )
}

export default LandingPage


