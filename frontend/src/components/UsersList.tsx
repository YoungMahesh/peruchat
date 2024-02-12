import { useEffect, useState } from "react";
import { useSocket } from "./SocketContext";
import { Event1, FetchedRecipients } from "~/utils/types.utils";

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
  console.log({ usersList });
  useEffect(() => {
    void (async () => {
      if (!socket) return console.error("socket not connected");
      if (socket.readyState !== WebSocket.OPEN)
        return console.error("socket is not ready");

      console.log("sending get_users request");
      socket.send("get_users");

      socket.addEventListener("message", async function (event: MessageEvent) {
        try {
          const _event = (await JSON.parse(event.data)) as Event1;
          const eventName1 = "get_msgs_resp";
          if (_event.type === eventName1) {
            if (_event.payload === null) setUserList([]);
            else setUserList(_event.payload as FetchedRecipients[]);
          }
        } catch (err) {
          console.error("Could not able to parse websocket-message", err);
        }
      });

      // const fetchUsers = async () => {
      //   try {
      //     const res = await fetch("http://localhost:3001/users", {
      //       method: "GET",
      //       headers: { "Content-Type": "application/json" },
      //     });

      //     const _users = (await res.json()) as unknown as {
      //       username: string;
      //     }[];
      //     if (res.status !== 200)
      //       return alert("could not able to fetch users list");
      //     setUserList(_users.filter((user) => user.username !== currUsername));
      //   } catch (error) {
      //     console.error(error);
      //   }
      // };

      // void fetchUsers();
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
