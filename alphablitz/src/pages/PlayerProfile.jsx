import { useState, useEffect } from "react";

function PlayerProfile() {

    const [profile, setProfile] = useState({})
    const [isLoading, setIsloading] = useState(true)
    const [error, setError] = useState("")

    useEffect(() => {
        async function fetchProfile() {
            const response = await fetch(`${import.meta.env.VITE_API_URL}/profile`, {
                method: `GET`,
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem(`token`)}`,
                },
            })

            const data = await response.json()

            setProfile(data)
            setIsloading(false)
        }
        fetchProfile()
    }, [])

    return (
        <div>
            {isLoading && (<div>
                <h1>PROFILE IS LOADING</h1>
            </div>)}
            {!isLoading && (<div>
                <h1> LEVEL: {profile.Level} </h1>
                <h1> TOTAL GAMES PLAYED: {profile.GamesPlayed} </h1>
                <h1> TOTAL GAMES WON: {profile.GamesWon} </h1>
                <h1> WIN RATE: {profile.WinRate} </h1>
                <h1> LONGEST WORD: {profile.LongestWord} </h1>
                <h1> SHORTEST WORD: {profile.ShortestWord} </h1>
                <h1> TOTAL CORRECT: {profile.TotalCorrect} </h1>
            </div>)}
            {error && <h3>ERROR LOADING</h3>}
        </div>
    )
}

export default PlayerProfile