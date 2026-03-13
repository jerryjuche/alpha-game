import { useEffect, useState } from "react"

function AuditDashboard() {

    const [points, setPoints] = useState({})
    const [pending, setPending] = useState([])
    const [isLoading, setIsLoading] = useState(true)


    useEffect(() => {
        async function getPending() {
            const response = await fetch(`${import.meta.env.VITE_API_URL}/audit/pending`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${localStorage.getItem(`token`)}`,
                },
            })

            const data = await response.json()
            setPending(data)
            setIsLoading(false)
        }
        getPending()
    }, [])

    async function approveSubmission(submission) {

        const response = await fetch(`${import.meta.env.VITE_API_URL}/audit/approve`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
            },
            body: JSON.stringify({
                submission_id: submission.ID,
                points: points[submission.ID],
                category: submission.Category,
                word: submission.Word,
            })
        },)

        if (response.status === 200) {
            setPending(prev => prev.filter(sub => sub.ID !== submission.ID))
        }
    }

    async function rejectSubmission(submissionId) {
        const response = await fetch(`${import.meta.env.VITE_API_URL}/audit/reject`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${localStorage.getItem('token')}`,
            },
            body: JSON.stringify({ submission_id: submissionId })
        })

        if (response.status === 200) {
            setPending(prev => prev.filter(sub => sub.ID !== submissionId))


        }


    }
    return (
        <div>
            {isLoading && <p>Loading...</p>}
            {!isLoading && pending.map(submission => (
                <div key={submission.ID}>
                    <p>{submission.Word} — {submission.Category}</p>
                    <input
                        type="number"
                        value={points[submission.ID] || ""}
                        onChange={e => setPoints(prev => ({ ...prev, [submission.ID]: parseInt(e.target.value) }))}
                    />
                    <button onClick={() => approveSubmission(submission)}>Approve</button>
                    <button onClick={() => rejectSubmission(submission.ID)}>Reject</button>
                </div>
            ))}
        </div>

    )
}

export default AuditDashboard