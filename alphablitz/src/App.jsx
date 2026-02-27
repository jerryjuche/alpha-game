import { BrowserRouter, Routes, Route } from "react-router-dom"
import LandingPage from "./pages/Landing"
import AdminDashboard from "./pages/AdminDashboard"
import AuditDashboard from "./pages/AuditDashboard"
import GameRoom from "./pages/GameRoom"
import Lobby from "./pages/Lobby"
import PlayerProfile from "./pages/PlayerProfile"
import Login from "./pages/Login"
import Register from "./pages/Register"

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element = {<LandingPage />} />
        <Route path="/admin" element = {<AdminDashboard />} />
        <Route path="/audit" element = {<AuditDashboard />} />        
        <Route path="/game" element = {<GameRoom />} />
        <Route path="/lobby" element = {<Lobby />} />
        <Route path="/profile" element = {<PlayerProfile />} />
        <Route path="/login" element = {<Login />} />
        <Route path="/register" element = {<Register />} />
              </Routes>
    </BrowserRouter>
  )
}

export default App