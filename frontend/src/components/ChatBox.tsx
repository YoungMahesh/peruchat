/* eslint-disable @typescript-eslint/no-unsafe-argument */
/* eslint-disable @typescript-eslint/no-misused-promises */
import { Input } from "@nextui-org/react";
import { PaperAirplaneIcon } from "@heroicons/react/24/outline";
import { useEffect, useState } from "react";
import authStore from "~/utils/auth";
import type { Event1, FetchedMessage } from "~/utils/types.utils";
import { useSocket } from "./SocketContext";

export default function ChatBox({
  receiver,
}: {
  sender: string;
  receiver: string;
}) {
  const [message, setMessage] = useState("");
  const [msgList1, setMsgList1] = useState<FetchedMessage[]>([]);
  const socket = useSocket();

  useEffect(() => {
    if (!socket) return;
    if (socket.readyState !== WebSocket.OPEN)
      return console.error("socket not ready");

    const getMsgReq = {
      type: "get_msgs",
      payload: {
        to_user: receiver,
      },
    };
    socket.send(JSON.stringify(getMsgReq));

    socket.addEventListener("message", async function (event: MessageEvent) {
      try {
        const _event = (await JSON.parse(event.data)) as Event1;
        const eventName1 = "get_msgs_resp";
        if (_event.type === eventName1) {
          if (_event.payload === null) setMsgList1([]);
          else setMsgList1(_event.payload as FetchedMessage[]);
        }
      } catch (err) {
        console.error("Could not able to parse websocket-message", err);
      }
    });
  }, [socket, receiver]);

  const sendMessage = async () => {
    const _token = await authStore.retreiveToken();
    if (!_token) return console.error("Invalid token");
    if (!socket) return console.error("Socket not connected");
    try {
      const sendMsgReq = JSON.stringify({
        type: "send_msg",
        payload: {
          to: receiver,
          message,
        },
      });
      socket.send(sendMsgReq);

      setMessage("");
      alert("message sent");
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <section className="p-3">
      {receiver === "" ? (
        <p className="text-center">Select a user to chat</p>
      ) : (
        <>
          <div className="flex flex-col">
            {msgList1.map((fm, idx) => (
              <p
                key={idx}
                className={`${fm.is_sender ? "self-end" : ""} m-2 p-1`}
              >
                {fm.message}
              </p>
            ))}
          </div>
          <div className="flex items-center">
            <Input
              type="text"
              value={message}
              onChange={(e) => setMessage(e.target.value)}
            />
            <PaperAirplaneIcon
              className="m-2 h-6 w-6 cursor-pointer"
              onClick={sendMessage}
            />
          </div>
        </>
      )}
    </section>
  );
}
