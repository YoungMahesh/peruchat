import { Input } from "@nextui-org/react";
import { PaperAirplaneIcon } from "@heroicons/react/24/outline";
import { useEffect, useState } from "react";
import authStore from "~/utils/auth";
import { useRouter } from "next/router";

type FetchedMsg = {
  fromSender: boolean;
  message: string;
};

export default function ChatBox({
  sender,
  receiver,
}: {
  sender: string;
  receiver: string;
}) {
  const router = useRouter();
  const [message, setMessage] = useState("");
  const [fetchedMsgs, setFetchedMsgs] = useState<FetchedMsg[]>([]);

  useEffect(() => {
    const fetchMessages = async () => {
      if (receiver === "") return;

      const _token = await authStore.retreiveToken();
      if (!_token) return console.error("Invalid token");

      try {
        const res = await fetch(
          `http://localhost:3001/get_msgs?to=${receiver}`,
          {
            method: "GET",
            headers: {
              "Content-Type": "application/json",
              Authorization: _token,
            },
          },
        );
        const _messages = (await res.json()) as unknown as {
          from: string;
          to: string;
          message: string;
        }[];
        if (res.status !== 200)
          return alert("could not able to fetch messages");

        console.log({ _messages });
        const _messagesWithSender = _messages.map((msg) => ({
          fromSender: msg.from === sender,
          message: msg.message,
        }));
        setFetchedMsgs(_messagesWithSender);
      } catch (error) {
        console.error(error);
      }
    };

    void fetchMessages();
  }, [sender, receiver, router]);

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
        {fetchedMsgs.map((fm, idx) => (
          <p key={idx} className={`${fm.fromSender ? "self-end" : ""} m-2 p-1`}>
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
