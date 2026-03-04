import { useNavigate } from 'react-router-dom'
import { Outlet } from 'react-router-dom'

function ProtectedRoute() {
    const token = localStorage.getItem('token')

    if (!token){
        navigate('/login')
        return null
    }
     
   
    return <Outlet />
}

export default ProtectedRoute