/* eslint-disable @typescript-eslint/no-unsafe-argument */
/* eslint-disable @typescript-eslint/no-misused-promises */

import { useEffect, useState } from "react";
import { useSocket } from "./SocketContext";
import type { Event1, FetchedRecipients } from "~/utils/types.utils";

export default function UsersList({
  currUsername,
  receiver,
  setReceiver,
}: {
  currUsername: string;
  receiver: string;
  setReceiver: (_: string) => void;
}) {
  const socket = useSocket();
  const [usersList, setUserList] = useState<{ username: string }[]>([]);

  useEffect(() => {
    void (async () => {
      if (!socket) return;
      if (socket.readyState !== WebSocket.OPEN)
        return console.error("socket is not ready");

      console.log("sending get_users request");
      socket.send(JSON.stringify({ type: "get_users", payload: null }));

      socket.addEventListener("message", async function (event: MessageEvent) {
        try {
          const _event = (await JSON.parse(event.data)) as Event1;
          const eventName1 = "get_users_resp";
          if (_event.type === eventName1) {
            if (_event.payload === null) setUserList([]);
            else setUserList(_event.payload as FetchedRecipients[]);
          }
        } catch (err) {
          console.error("Could not able to parse websocket-message", err);
        }
      });
    })();
  }, [currUsername, socket]);

  return (
    <section className="m-2 flex flex-col">
      {usersList.map((user, idx) => (
        <p
          className={`m-1 rounded border p-2 text-center ${user.username === receiver ? "bg-white text-black" : ""}`}
          key={idx}
          onClick={() => setReceiver(user.username)}
        >
          {user.username}
        </p>
      ))}
    </section>
  );
}
