import { useEffect, useRef, useState } from "react";
import { useParams } from "react-router-dom"

import './GameRoom.css'

function GameRoom() {
    const { gameId } = useParams()
    const [currentLetter, setCurrentLetter] = useState("")
    const [name, setName] = useState("")
    const [animal, setAnimal] = useState("")
    const [place, setPlace] = useState("")
    const [thing, setThing] = useState("")
    const [food, setFood] = useState("")
    const [roundId, setRoundId] = useState("")
    const [gameTime, setGameTime] = useState(180)
    const [timer, setTimer] = useState(10)
    const [gamePhase, setGamePhase] = useState("waiting")
    const submitRef = useRef()
    submitRef.current = submitAnswers
    const gamePhaseRef = useRef(gamePhase)
    const [breakDuration, setBreakDuration] = useState(5)
    const lastSubmittedRoundRef = useRef("x")
    const [scores, setScores] = useState({})



    useEffect(() => {
        const socket = new WebSocket(`${import.meta.env.VITE_WS_URL}/ws/${gameId}?token=${localStorage.getItem('token')}`)

        socket.onopen = () => console.log("WS connected")
        socket.onclose = () => console.log("WS closed")
        socket.onerror = (e) => console.log("WS error", e)

        socket.onmessage = (event) => {
            console.log(event.data)
            const data = event.data

            if (data.startsWith("ROUND:")) {
                setRoundId(data.split(":")[1])

            }

            if (data.startsWith("STATE:")) {
                const state = JSON.parse(data.slice(6))
                setCurrentLetter(state.letter)
                setGamePhase(state.phase)
                setTimer(state.timer)
                setGameTime(state.gameTime)
                setRoundId(state.roundID)
            }

            if (data.startsWith("LETTER:")) {
                setCurrentLetter(data.split(":")[1])

                setGamePhase("playing")
                setTimer(10)

                // Clear inputs on new letter
                setName("")
                setAnimal("")
                setPlace("")
                setThing("")
                setFood("")
                setTimer(10)
            }

            if (data.startsWith("BREAK:")) {
                const duration = parseInt(data.split(":")[1])
                setBreakDuration(duration)
                setGamePhase("break")
                submitRef.current()
            }

            if (data === "GAME:FINISHED") {
                setGamePhase("waiting")
            }

            if (data.startsWith("SCORES:")) {
                const scores = JSON.parse(data.slice(7))
                setScores(scores)
                setGamePhase("break")
            }

        }

        return () => socket.close()
    }, [gameId])

    useEffect(() => {
        gamePhaseRef.current = gamePhase
    }, [gamePhase])

    useEffect(() => {
        if (gameTime === 0) {
            submitRef.current()
            setGamePhase("waiting")
        }
    }, [gameTime])

    useEffect(() => {
        if (timer === 0) {
            submitRef.current()
        }
    }, [timer])

    // timer for the rounds (10secs) & submits answers automitically
    useEffect(() => {
        const interval = setInterval(() => {
            if (gamePhaseRef.current === "playing") {
                setTimer(prev => {
                    if (prev <= 1) {
                        return 0
                    }
                    return prev - 1
                })
            }

            if (gamePhaseRef.current !== "waiting") {
                setGameTime(prev => {
                    if (prev <= 1) {
                        return 0
                    }
                    return prev - 1
                })
            }

            if (gamePhaseRef.current === "break") {
                setBreakDuration(prev => {
                    if (prev === 1) return 0
                    
                    return prev - 1
                })
            }
        }, 1000)
        return () => clearInterval(interval)
    }, [])


    async function startGame() {
        await fetch(`${import.meta.env.VITE_API_URL}/game/start`, {

            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
            },
            body: JSON.stringify({ game_id: gameId })
        })

        setGamePhase("playing")

    }

    async function submitAnswers() {
        if (lastSubmittedRoundRef.current === roundId) return
        lastSubmittedRoundRef.current = roundId

        const submissions = [
            { round_id: roundId, game_id: gameId, word: name, category: "name" },
            { round_id: roundId, game_id: gameId, word: animal, category: "animal" },
            { round_id: roundId, game_id: gameId, word: place, category: "place" },
            { round_id: roundId, game_id: gameId, word: thing, category: "thing" },
            { round_id: roundId, game_id: gameId, word: food, category: "food" },
        ]

        const results = await Promise.allSettled(
            submissions.map(submission =>
                fetch(`${import.meta.env.VITE_API_URL}/game/submit`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'Authorization': `Bearer ${localStorage.getItem('token')}`
                    },
                    body: JSON.stringify(submission)
                }).then(res => res.json())
            )
        )

        results.forEach((result, index) => {
            if (result.status === "fulfilled") {
                console.log("Submitted:", result.value)
            } else {
                console.log("Failed:", submissions[index].category)
            }
        })
    }

    return (
        <div className="room-conainer">
            <div className="game-room">
                <h3>Current Letter: {currentLetter} </h3>
                <div>
                    {gamePhase === "playing" && <h2> Round id: {roundId}</h2>}
                    {gamePhase === "playing" && <h2>Timer: {timer}</h2>}
                    {gamePhase === "break" && <h2>Break: {breakDuration}</h2>}
                </div>

                <div>
                    <h2>Game Time: {gameTime} </h2>
                </div>

                {gamePhase === "waiting" && <button className="button" onClick={startGame}>Start Game</button>}                <div>
                    <input
                        type="text"
                        placeholder="Name"
                        value={name}
                        autoComplete="off"
                        onChange={(e) => setName(e.target.value)}
                        disabled={gamePhase !== "playing"}
                        className="input"

                    />
                    <input
                        type="text"
                        placeholder="Animal"
                        value={animal}
                        autoComplete="off"
                        onChange={(e) => setAnimal(e.target.value)}
                        disabled={gamePhase !== "playing"}
                        className="input"

                    />
                    <input
                        type="text"
                        placeholder="Place"
                        value={place}
                        autoComplete="off"
                        onChange={(e) => setPlace(e.target.value)}
                        disabled={gamePhase !== "playing"}
                        className="input"

                    />
                    <input
                        type="text"
                        placeholder="Thing"
                        value={thing}
                        autoComplete="off"
                        onChange={(e) => setThing(e.target.value)}
                        disabled={gamePhase !== "playing"}
                        className="input"

                    />
                    <input
                        type="text"
                        placeholder="Food"
                        value={food}
                        autoComplete="off"
                        onChange={(e) => setFood(e.target.value)}
                        disabled={gamePhase !== "playing"}
                        className="input"

                    />
                    <button className="button" disabled={gamePhase !== "playing"} type="button" onClick={submitAnswers}>Submit</button>
                </div>
            </div>
            {gamePhase === "break" && (
                <div>
                    <h2>Get Ready</h2>
                    <p>Next Round In {breakDuration}</p>
                    {Object.entries(scores).map(([userId, points]) => (
                        <div key={userId}>
                            {userId} - {points}

                        </div>

                    ))}
                </div>)}
        </div>
    )
}

export default GameRoom