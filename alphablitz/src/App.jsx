import { BrowserRouter, Routes, Route } from "react-router-dom"
import LandingPage from "./pages/Landing"
import AdminDashboard from "./pages/AdminDashboard"
import AuditDashboard from "./pages/AuditDashboard"
import GameRoom from "./pages/GameRoom"
import Lobby from "./pages/Lobby"
import PlayerProfile from "./pages/PlayerProfile"
import Login from "./pages/Login"
import Register from "./pages/Register"
import ProtectedRoute from "./components/ProtectedRoute"


function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<LandingPage />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />

        <Route element={<ProtectedRoute />}>
          <Route path="/lobby" element={<Lobby />} />
          <Route path="/profile" element={<PlayerProfile />} />
          <Route path="/audit" element={<AuditDashboard />} />
          <Route path="/admin" element={<AdminDashboard />} />
          <Route path="/game/:gameId" element={<GameRoom />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}




export default App