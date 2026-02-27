import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import './LandingPage.css'



function LandingPage() {
    return (
        <div className="container">
            <h1>ALPHABLITZ</h1>
            <p>AlphaBlitz is a real-time multiplayer word game built for speed, accuracy, and strategy. <br></br>
                Players compete head-to-head in fast-paced typing challenges where quick thinking and precision determine the winner.
                <br></br>Every match is live, every second counts, and every word moves you closer to the top of the global leaderboard.
                <br></br>Designed for instant play in the browser, AlphaBlitz transforms language into competition.</p>
        </div>
    )
}

export default LandingPage


