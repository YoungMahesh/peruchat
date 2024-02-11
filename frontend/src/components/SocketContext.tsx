import { createContext, useContext } from "react";

const SocketContext = createContext<WebSocket | null>(null)

export function useSocket() {
    return useContext(SocketContext)
}

export default SocketContext