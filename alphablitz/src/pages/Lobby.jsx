import { useState } from "react";
import { useNavigate } from "react-router-dom";

function Lobby() {
    const navigate = useNavigate()
    const [gameId, setGameId] = useState("")
    const [error, setError] = useState("")
    const [createdInviteCode, setCreatedInviteCode] = useState("")
    const [joinInviteCode, setJoinInviteCode] = useState("")

    async function createGame() {
        const response = await fetch('http://localhost:8080/game/create', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            },
        })

        if (response.status === 201) {
            const data = await response.json()
            setGameId(data.game_id)
            setCreatedInviteCode(data.invite_code)

        } else {
            setError("Cannot create game")
        }
    }

    async function joinGame() {
        console.log(joinInviteCode)
        const response = await fetch('http://localhost:8080/game/join', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`
            },
            body: JSON.stringify({
                invite_code: joinInviteCode,
            })
        })

        if (response.status === 200) {
            const data = await response.json()
            setGameId(data.game_id)
            navigate(`/game/${data.game_id}`)

        } else {
            setError("Error joining game")
        }

    }

    function copyInviteCode() {
        navigator.clipboard.writeText(createdInviteCode)
    }

    return (
        <div className="create-game-container">
            <h1 className="h1oth">CREATE OR JOIN A ROOM</h1>
            <div className="game-container">
                <h1 className="h1oth">CREATE ROOM</h1>
                <p>Create Room allows a player to start a new game session and act as the host. <br></br>
                    A unique room is generated where other players can join using a room code. <br></br>
                    The host controls the session setup and begins the game once all players have joined.</p>
                <button className="btn" onClick={createGame}>Create Room</button>
                {createdInviteCode && <p className="invitecode">{createdInviteCode}</p>}
                <button type="button" onClick={copyInviteCode}>Copy invite code</button>
                <button type="button" onClick={() => navigate(`/game/${gameId}`)}>start game</button>
            </div>
            <div className="game-container">
                <h1 className="h1oth">JOIN ROOM</h1>
                <p>Join Room allows a player to enter an existing game session using a valid room code. <br></br>
                    Once joined, the player becomes part of the multiplayer lobby and waits with other participants until the host starts the game.<br></br>
                </p>
                <input type="text" placeholder="invite code" value={joinInviteCode} onChange={(e) => setJoinInviteCode(e.target.value)} />
                <button className="btn" onClick={joinGame}>Join Room</button>
                {error && <p className="error">{error}</p>}
            </div>
        </div>
    )
}

export default Lobby