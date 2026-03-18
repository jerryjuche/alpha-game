import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import '../styles/alphablitz.css'
import './Landing.css'

const LETTERS = ['A', 'B', 'C', 'M', 'S', 'T', 'P', 'R', 'W']
const CATS = ['Name', 'Animal', 'Place', 'Thing', 'Food']

export default function LandingPage() {
  const navigate = useNavigate()
  const [letterIdx, setLetterIdx] = useState(0)

  useEffect(() => {
    const id = setInterval(() => {
      setLetterIdx(i => (i + 1) % LETTERS.length)
    }, 2500)
    return () => clearInterval(id)
  }, [])

  return (
    <div className="landing grain">
      <div className="scanlines" />

      {/* Ambient background */}
      <div className="landing__bg">
        <div className="landing__bg-orb landing__bg-orb--1" />
        <div className="landing__bg-orb landing__bg-orb--2" />
      </div>
      <div className="landing__stripe" />

      {/* Main content */}
      <div className="landing__inner">

        {/* Live status */}
        <div className="landing__status">
          <span className="landing__status-dot" />
          <span>Servers Online</span>
          <span>·</span>
          <span>Real-time Multiplayer</span>
        </div>

        {/* Animated letter board */}
        <div className="landing__board">
          <div className="board__frame">
            <div className="board__inner">
              <span className="board__letter" key={letterIdx}>
                {LETTERS[letterIdx]}
              </span>
              <span className="board__label">Current Letter</span>
            </div>
          </div>
        </div>

        {/* Title with glitch */}
        <div className="landing__title-wrap">
          <h1
            className="landing__title"
            data-text="ALPHABLITZ"
          >
            <span className="word-alpha">ALPHA</span>
            <span className="word-blitz">BLITZ</span>
          </h1>
        </div>

        {/* Tagline */}
        <p className="landing__tagline">
          One Letter · Five Categories · Eight Seconds
        </p>

        {/* Category pills */}
        <div className="landing__cats">
          {CATS.map(cat => (
            <span key={cat} className="cat-pill">{cat}</span>
          ))}
        </div>

        {/* CTA */}
        <div className="landing__cta">
          <button className="btn-primary" onClick={() => navigate('/register')}>
            Play Now
          </button>
          <button className="btn-ghost" onClick={() => navigate('/login')}>
            Sign In
          </button>
        </div>
      </div>

      {/* Stat strip */}
      <div className="landing__stats">
        <div className="stat-item">
          <span className="stat-value">8s</span>
          <span className="stat-label">Per<br/>Letter</span>
        </div>
        <div className="stat-item">
          <span className="stat-value">5</span>
          <span className="stat-label">Categor-<br/>ies</span>
        </div>
        <div className="stat-item">
          <span className="stat-value">3min</span>
          <span className="stat-label">Round<br/>Time</span>
        </div>
        <div className="stat-item">
          <span className="stat-value">∞</span>
          <span className="stat-label">Players<br/>Online</span>
        </div>
      </div>
    </div>
  )
}