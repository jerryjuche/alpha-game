import { useFetcher, useParams } from "react-router-dom";
import { useEffect, useRef, useState } from "react";
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


    useEffect(() => {
        const socket = new WebSocket(`ws://localhost:8080/ws/${gameId}?token=${localStorage.getItem('token')}`)

        socket.onmessage = (event) => {
            const data = event.data
            console.log(event.data)

            if (data.startsWith("STATE:")) {
                const state = JSON.parse(data.slice(6))
                setCurrentLetter(state.letter)
                setGamePhase(state.phase)
                setTimer(state.timer)
                setGameTime(state.gameTime)
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

        }

        return () => socket.close()
    }, [gameId])

    useEffect(() => {
        gamePhaseRef.current = gamePhase
    }, [gamePhase])

    // timer for the rounds (10secs) & submits answers automitically
    useEffect(() => {
        const interval = setInterval(() => {
            if (gamePhaseRef.current === "playing") {
                setTimer(prev => {
                    if (prev <= 1) {
                        submitRef.current()
                        return 0
                    }
                    return prev - 1
                })
            }

            if (gamePhaseRef.current !== "waiting") {
                setGameTime(prev => {
                    if (prev <= 1) {
                        submitRef.current()
                        setGamePhase("waiting")
                        return 0
                    }
                    return prev - 1
                })
            }

            if (gamePhaseRef.current === "break") {
                setBreakDuration(prev => {
                    if (prev <= 1) {
                        return 5
                    }
                    return prev - 1
                })
            }
        }, 1000)
        return () => clearInterval(interval)
    }, [])


    async function startGame() {
        await fetch(`http://localhost:8080/game/start`, {
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

        const submissions = [
            { round_id: roundId, word: name, category: "name" },
            { round_id: roundId, word: animal, category: "animal" },
            { round_id: roundId, word: place, category: "place" },
            { round_id: roundId, word: thing, category: "thing" },
            { round_id: roundId, word: food, category: "food" },
        ]

        const results = await Promise.allSettled(
            submissions.map(submission =>
                fetch('http://localhost:8080/game/submit', {
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
        </div>
    )
}

export default GameRoom