/* eslint-disable @typescript-eslint/no-unsafe-argument */
/* eslint-disable @typescript-eslint/no-misused-promises */
import { Input } from "@nextui-org/react";
import { PaperAirplaneIcon } from "@heroicons/react/24/outline";
import { useEffect, useState } from "react";
import authStore from "~/utils/auth";
import type { FetchedMessage } from "~/utils/types.utils";
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

  console.log({ msgList1 });
  useEffect(() => {
    if (!socket) return;
    if (socket.readyState !== WebSocket.OPEN) return;

    console.log({ receiver });
    const getMsgReq = {
      type: "get_msgs",
      payload: {
        to_user: receiver,
      },
    };
    socket.send(JSON.stringify(getMsgReq));

    socket.addEventListener("message", async function (event: MessageEvent) {
      console.log("Message from server ", event);
      const _msgList = (await JSON.parse(event.data)) as
        | FetchedMessage[]
        | null;
      console.log(_msgList);
      if (_msgList) setMsgList1(_msgList);
      else setMsgList1([]);
    });
  }, [socket, receiver]);

  const sendMessage = async () => {
    const _token = await authStore.retreiveToken();
    if (!_token) return console.error("Invalid token");
    try {
      const res = await fetch("http://localhost:3001/send_msg", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: _token,
        },
        body: JSON.stringify({ to: receiver, message }),
      });
      if (res.status !== 201) return alert("could not able to send message");
      setMessage("");
      alert("message sent");
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <section>
      <div className="flex flex-col">
        {msgList1.map((fm, idx) => (
          <p key={idx} className={`${fm.is_sender ? "self-end" : ""} m-2 p-1`}>
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
    </section>
  );
}
