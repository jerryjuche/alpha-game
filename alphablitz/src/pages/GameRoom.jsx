import { useEffect, useRef, useState } from 'react'
import { useParams } from 'react-router-dom'
import '../styles/alphablitz.css'
import './GameRoom.css'

const CATS = [
  { key: 'name',   label: 'Name'   },
  { key: 'animal', label: 'Animal' },
  { key: 'place',  label: 'Place'  },
  { key: 'thing',  label: 'Thing'  },
  { key: 'food',   label: 'Food'   },
]

export default function GameRoom() {
  const { gameId } = useParams()

  const [letter,    setLetter]    = useState('')
  const [roundId,   setRoundId]   = useState('')
  const [gameTime,  setGameTime]  = useState(180)
  const [timer,     setTimer]     = useState(10)
  const [phase,     setPhase]     = useState('waiting')
  const [breakSecs, setBreakSecs] = useState(5)
  const [scores,    setScores]    = useState({})
  const [inputs,    setInputs]    = useState({ name:'', animal:'', place:'', thing:'', food:'' })

  const submitRef   = useRef()
  const phaseRef    = useRef(phase)
  const lastRound   = useRef('x')

  submitRef.current = submitAnswers

  useEffect(() => { phaseRef.current = phase }, [phase])

  // ── WebSocket ──────────────────────────────────────────────
  useEffect(() => {
    const ws = new WebSocket(
      `${import.meta.env.VITE_WS_URL}/ws/${gameId}?token=${localStorage.getItem('token')}`
    )
    ws.onopen  = () => console.log('WS connected')
    ws.onclose = () => console.log('WS closed')
    ws.onerror = e  => console.error('WS error', e)

    ws.onmessage = ({ data }) => {
      if (data.startsWith('ROUND:'))  { setRoundId(data.split(':')[1]); return }
      if (data === 'GAME:FINISHED')   { setPhase('waiting'); return }

      if (data.startsWith('STATE:')) {
        const s = JSON.parse(data.slice(6))
        setLetter(s.letter); setPhase(s.phase)
        setTimer(s.timer);   setGameTime(s.gameTime)
        setRoundId(s.roundID)
        return
      }
      if (data.startsWith('LETTER:')) {
        setLetter(data.split(':')[1])
        setPhase('playing'); setTimer(10)
        setInputs({ name:'', animal:'', place:'', thing:'', food:'' })
        return
      }
      if (data.startsWith('BREAK:')) {
        setBreakSecs(parseInt(data.split(':')[1]))
        setPhase('break')
        submitRef.current()
        return
      }
      if (data.startsWith('SCORES:')) {
        setScores(JSON.parse(data.slice(7)))
      }
    }
    return () => ws.close()
  }, [gameId])

  // ── Game over watcher ──────────────────────────────────────
  useEffect(() => {
    if (gameTime === 0) { submitRef.current(); setPhase('waiting') }
  }, [gameTime])

  // ── Letter timer watcher ───────────────────────────────────
  useEffect(() => {
    if (timer === 0) submitRef.current()
  }, [timer])

  // ── Master interval ────────────────────────────────────────
  useEffect(() => {
    const id = setInterval(() => {
      const p = phaseRef.current
      if (p === 'playing')
        setTimer(t => t <= 1 ? 0 : t - 1)
      if (p !== 'waiting')
        setGameTime(t => t <= 1 ? 0 : t - 1)
      if (p === 'break')
        setBreakSecs(t => t <= 1 ? 0 : t - 1)
    }, 1000)
    return () => clearInterval(id)
  }, [])

  // ── Submit ─────────────────────────────────────────────────
  async function submitAnswers() {
    if (lastRound.current === roundId) return
    lastRound.current = roundId

    await Promise.allSettled(
      CATS.map(cat =>
        fetch(`${import.meta.env.VITE_API_URL}/game/submit`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${localStorage.getItem('token')}`,
          },
          body: JSON.stringify({
            round_id: roundId, game_id: gameId,
            word: inputs[cat.key], category: cat.key,
          }),
        })
      )
    )
  }

  async function startGame() {
    await fetch(`${import.meta.env.VITE_API_URL}/game/start`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${localStorage.getItem('token')}`,
      },
      body: JSON.stringify({ game_id: gameId }),
    })
    setPhase('playing')
  }

  function fmt(s) {
    return `${Math.floor(s/60).toString().padStart(2,'0')}:${(s%60).toString().padStart(2,'0')}`
  }

  const timerClass  = timer <= 3 ? 'danger' : timer <= 6 ? 'warn' : ''
  const clockUrgent = gameTime <= 30
  const sorted      = Object.entries(scores).sort((a,b) => b[1]-a[1])

  // ── WAITING ──────────────────────────────────────────────────
  if (phase === 'waiting') return (
    <div className="gameroom grain">
      <div className="scanlines" />
      <div className="gameroom__grid" />
      <div className="hud">
        <span className="hud__brand">ALPHABLITZ</span>
        <span className="hud__phase hud__phase--waiting">Waiting</span>
      </div>
      <div className="waiting">
        <h1 className="waiting__title">
          ALPHA<span>BLITZ</span>
        </h1>
        <p className="waiting__sub">Waiting for host to start</p>
        <button className="waiting__btn" onClick={startGame}>
          Start Game
        </button>
      </div>
    </div>
  )

  // ── BREAK ────────────────────────────────────────────────────
  if (phase === 'break') return (
    <div className="gameroom grain">
      <div className="scanlines" />
      <div className="gameroom__grid" />
      <div className="hud">
        <span className="hud__brand">ALPHABLITZ</span>
        <span className="hud__clock-time">{fmt(gameTime)}</span>
        <span className="hud__phase hud__phase--break">Break</span>
      </div>
      <div className="break-screen">
        <h2 className="break-screen__title">GET READY</h2>
        <p className="break-screen__count">Next round in {breakSecs}s</p>

        {sorted.length > 0 && (
          <div className="scoreboard">
            <div className="scoreboard__head">
              <span>#</span>
              <span>Player</span>
              <span>Pts</span>
            </div>
            {sorted.map(([uid, pts], i) => (
              <div className="score-row" key={uid}>
                <span className="score-row__rank">{i + 1}</span>
                <span className="score-row__id">{uid.slice(0,12)}…</span>
                <span className="score-row__pts">{pts}</span>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )

  // ── PLAYING ──────────────────────────────────────────────────
  return (
    <div className="gameroom grain">
      <div className="scanlines" />
      <div className="gameroom__grid" />

      {/* HUD */}
      <div className="hud">
        <span className="hud__brand">ALPHABLITZ</span>
        <div className="hud__bar">
          <div
            className={`hud__bar-fill ${clockUrgent ? 'warn' : ''}`}
            style={{ width: `${(gameTime / 180) * 100}%` }}
          />
        </div>
        <div className="hud__clock">
          <span className="hud__clock-label">Round</span>
          <span className={`hud__clock-time ${clockUrgent ? 'urgent' : ''}`}>
            {fmt(gameTime)}
          </span>
        </div>
        <span className="hud__phase hud__phase--playing">Playing</span>
      </div>

      {/* Main playing area */}
      <div className="playing">

        {/* Letter display */}
        <div className="letter-stage">
          <span className="letter-stage__label">Current Letter</span>
          <div className="letter-stage__card">
            <div className="letter-stage__face">
              <span className="letter-stage__glyph" key={letter}>{letter}</span>
            </div>
          </div>
          {/* Per-letter timer */}
          <div className="letter-timer">
            <div className="letter-timer__top">
              <span className={`letter-timer__secs ${timerClass}`}>{timer}</span>
              <span className="letter-timer__of">/ 10s</span>
            </div>
            <div className="letter-timer__track">
              <div
                className={`letter-timer__fill ${timerClass}`}
                style={{ width: `${(timer / 10) * 100}%` }}
              />
            </div>
          </div>
        </div>

        {/* Inputs */}
        <div className="fields">
          {CATS.map(cat => (
            <div key={cat.key} className="field">
              <span className="field__label">{cat.label}</span>
              <input
                className="field__input"
                type="text"
                value={inputs[cat.key]}
                autoComplete="off"
                autoCorrect="off"
                spellCheck="false"
                disabled={phase !== 'playing'}
                placeholder={`${cat.label} starting with ${letter || '…'}`}
                onChange={e =>
                  setInputs(p => ({ ...p, [cat.key]: e.target.value }))
                }
              />
            </div>
          ))}
        </div>

        <button
          className="submit-btn"
          disabled={phase !== 'playing'}
          onClick={submitAnswers}
        >
          Submit
        </button>
      </div>
    </div>
  )
}