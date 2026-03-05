import { useParams } from "react-router-dom";
import { useEffect, useRef, useState } from "react";

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
    const [timer, setTimer] = useState(12)
    const [gamePhase, setGamePhase] = useState(false)
    const submitRef = useRef()
    submitRef.current = submitAnswers




    useEffect(() => {
        const socket = new WebSocket(`ws://localhost:8080/ws/${gameId}?token=${localStorage.getItem('token')}`)

        socket.onmessage = (event) => {
            const data = event.data
            console.log(event.data)

            if (data.startsWith("LETTER:")) {
                setCurrentLetter(data.split(":")[1])
                // Clear inputs on new letter
                setName("")
                setAnimal("")
                setPlace("")
                setThing("")
                setFood("")
                setTimer(12)
            }

            if (data.startsWith("ROUND:")) {
                setRoundId(data.split(":")[1])
            }

            if (data === "GAME:FINISHED") {
                alert("Game Over!")
            }
        }

        return () => socket.close()
    }, [gameId])


    useEffect(() => {
        if (gamePhase === true) {
            const interval = setInterval(() => {

                setTimer(prev => {
                    if (prev <= 1) {
                        submitRef.current()
                        return 12
                    }
                    return prev - 1
                })

            }, 1000)
            return () => clearInterval(interval)

        }

    }, [gamePhase])

    useEffect(() => {
        if (gamePhase === true) {
            const interval = setInterval(() => {
                setGameTime(prev => {
                    if (prev <= 1) return 0
                    return prev - 1
                })
            }, 1000)
            return () => clearInterval(interval)
        }

        if (gamePhase === false) {
            console.log("Game Has Not Started")
        }
    }, [gamePhase])

    useEffect(() => {
        if (gameTime === 0) {
            submitRef.current()
            setGamePhase(false)

        }
    }, [gameTime])


    async function startGame() {
        await fetch(`http://localhost:8080/game/start`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
            },
            body: JSON.stringify({ game_id: gameId })
        })

        setGamePhase(true)

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
        <div>
            <div>
                <h3>Current Letter: {currentLetter} </h3>
                <div>
                    <h2>Timer: {timer} </h2>
                </div>

                <div>
                    <h2>Game Time: {gameTime} </h2>
                </div>

                {!gamePhase && (<button disabled={gamePhase} onClick={startGame}>Start Game</button>)}
                <div>
                    <input
                        type="text"
                        placeholder="Name"
                        value={name}
                        autoComplete="off"
                        onChange={(e) => setName(e.target.value)}
                        disabled={!gamePhase}
                    />
                    <input
                        type="text"
                        placeholder="Animal"
                        value={animal}
                        autoComplete="off"
                        onChange={(e) => setAnimal(e.target.value)}
                        disabled={!gamePhase}
                    />
                    <input
                        type="text"
                        placeholder="Place"
                        value={place}
                        autoComplete="off"
                        onChange={(e) => setPlace(e.target.value)}
                        disabled={!gamePhase}
                    />
                    <input
                        type="text"
                        placeholder="Thing"
                        value={thing}
                        autoComplete="off"
                        onChange={(e) => setThing(e.target.value)}
                        disabled={!gamePhase}
                    />
                    <input
                        type="text"
                        placeholder="Food"
                        value={food}
                        autoComplete="off"
                        onChange={(e) => setFood(e.target.value)}
                        disabled={!gamePhase}
                    />
                    <button disabled={!gamePhase} type="button" onClick={submitAnswers}>Submit</button>
                </div>
            </div>
        </div>
    )
}

export default GameRoom