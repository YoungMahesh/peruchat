import { type AppType } from "next/dist/shared/lib/utils";
import { NextUIProvider, User } from "@nextui-org/react";
import { ToastContainer } from "react-toastify";
import "~/styles/globals.css";
import "react-toastify/dist/ReactToastify.css";
import SocketContext from "~/components/SocketContext";
import { useEffect, useState } from "react";
import authStore from "~/utils/auth";
import UserContext from "~/components/UserContext";
import type { LoggedInUserInfo } from "~/utils/types.utils";

const MyApp: AppType = ({ Component, pageProps }) => {
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [user, setUser] = useState<LoggedInUserInfo | null>(null);

  useEffect(() => {
    void (async () => {
      const _userInfo = await authStore.retreiveUser();
      if (!_userInfo) return;
      setUser(_userInfo);

      const _socket = new WebSocket(
        `ws://localhost:3001/ws?token=${_userInfo.token}`,
      );
      _socket.onopen = () => {
        console.log("Connected to the server");
        setSocket(_socket);
      };
      _socket.onerror = function (event) {
        console.error("WebSocket error observed:", Object.keys(event));
      };
      return () => {
        console.log("closing socket");
        _socket.close();
      };
    })();
  }, []);

  return (
    <NextUIProvider>
      <UserContext.Provider value={user}>
        <SocketContext.Provider value={socket}>
          <Component {...pageProps} />
        </SocketContext.Provider>
      </UserContext.Provider>
      <ToastContainer />
    </NextUIProvider>
  );
};

export default MyApp;
