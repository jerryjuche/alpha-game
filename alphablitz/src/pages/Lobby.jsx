import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import '../styles/alphablitz.css'
import './Lobby.css'

export default function Lobby() {
  const navigate = useNavigate()
  const [error, setError] = useState('')
  const [inviteCode, setInviteCode] = useState('')
  const [joinCode, setJoinCode] = useState('')
  const [copied, setCopied] = useState(false)

  async function createGame() {
    setError('')
    const res = await fetch(`${import.meta.env.VITE_API_URL}/game/create`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${localStorage.getItem('token')}`,
      },
    })
    if (res.status === 201) {
      const data = await res.json()
      setInviteCode(data.invite_code)
      navigate(`/game/${data.game_id}`)
    } else {
      setError('Failed to create game. Try again.')
    }
  }

  async function joinGame() {
    setError('')
    const res = await fetch(`${import.meta.env.VITE_API_URL}/game/join`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${localStorage.getItem('token')}`,
      },
      body: JSON.stringify({ invite_code: joinCode }),
    })
    if (res.status === 200) {
      const data = await res.json()
      navigate(`/game/${data.game_id}`)
    } else {
      setError('Invalid code or game already started.')
    }
  }

  function copyCode() {
    navigator.clipboard.writeText(inviteCode)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }

  return (
    <div className="lobby grain">
      <div className="scanlines" />

      {/* Nav */}
      <nav className="lobby__nav">
        <span className="lobby__nav-logo">ALPHABLITZ</span>
        <div className="lobby__nav-user">
          <span className="nav-dot" />
          <span>Connected</span>
        </div>
      </nav>

      <div className="lobby__body">
        {/* Left column */}
        <div className="lobby__left">
          <div className="lobby__heading">
            <p className="lobby__heading-sub">Game Lobby</p>
            <h1 className="lobby__heading-title">
              CREATE<br />OR JOIN
              <span>Choose your path</span>
            </h1>
          </div>

          {/* Info cards */}
          <div className="lobby__info">
            <div className="lobby__info-card">
              <div className="lobby__info-val">8s</div>
              <div className="lobby__info-key">Per Letter</div>
            </div>
            <div className="lobby__info-card">
              <div className="lobby__info-val">5</div>
              <div className="lobby__info-key">Categories</div>
            </div>
            <div className="lobby__info-card">
              <div className="lobby__info-val">3m</div>
              <div className="lobby__info-key">Round Time</div>
            </div>
          </div>

          {/* Create Room */}
          <div className="lobby__panel">
            <div className="bracket bracket--tl" />
            <div className="bracket bracket--tr" />
            <div className="bracket bracket--bl" />
            <div className="bracket bracket--br" />

            <p className="lobby__panel-tag">Host</p>
            <h2 className="lobby__panel-title">CREATE ROOM</h2>
            <p className="lobby__panel-desc">
              Start a new session. A unique invite code will be generated
              for other players to join.
            </p>
            <div className="lobby__divider" />
            <button className="lobby__btn-create" onClick={createGame}>
              Create Room
            </button>

            {inviteCode && (
              <div className="invite-reveal" style={{ marginTop: 16 }}>
                <p className="invite-reveal__label">Your Invite Code</p>
                <div className="invite-reveal__row">
                  <span className="invite-reveal__code">{inviteCode}</span>
                  <button
                    className={`invite-reveal__copy ${copied ? 'copied' : ''}`}
                    onClick={copyCode}
                  >
                    {copied ? '✓ Copied' : 'Copy'}
                  </button>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Right column */}
        <div className="lobby__right">
          <div className="lobby__panel">
            <div className="bracket bracket--tl" />
            <div className="bracket bracket--tr" />
            <div className="bracket bracket--bl" />
            <div className="bracket bracket--br" />

            <p className="lobby__panel-tag">Player</p>
            <h2 className="lobby__panel-title">JOIN ROOM</h2>
            <p className="lobby__panel-desc">
              Enter the invite code shared by the host to jump into
              an existing game session.
            </p>
            <div className="lobby__divider" />
            <input
              className="lobby__input"
              type="text"
              placeholder="Enter Code"
              value={joinCode}
              onChange={e => setJoinCode(e.target.value.toUpperCase())}
            />
            <button className="lobby__btn-join" onClick={joinGame}>
              Join Room
            </button>
          </div>

          {error && <p className="lobby__error">{error}</p>}
        </div>
      </div>
    </div>
  )
}