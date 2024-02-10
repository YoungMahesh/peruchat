import { Button } from "@nextui-org/react";
import Head from "next/head";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import ChatBox from "~/components/ChatBox";
import UsersList from "~/components/UsersList";
import authStore from "~/utils/auth";

export default function Home() {
  const [currUser, setCurrUser] = useState<string | null>(null);
  const [receiver, setReceiver] = useState<string>("");

  useEffect(() => {
    void (async () => {
      const _username = await authStore.retrieveUsername();
      setCurrUser(_username);

      const _token = await authStore.retreiveToken();
      const socket = new WebSocket(`ws://localhost:3001/ws?token=${_token}`);
      // const socket = new WebSocket(`ws://localhost:3001/ws`);

      socket.addEventListener("open", function (event) {
        console.log("Connected to the WS Server");
      });
      socket.addEventListener("message", function (event) {
        console.log("Message from server ", event.data);
      })

      socket.onopen = () => {
        console.log("Connected to the server");
      };

      socket.onerror = function (event) {
        console.error("WebSocket error observed:", Object.keys(event));
      };
    })();
  }, []);

  return (
    <>
      <Head>
        <title>Realtime Chat</title>
        <meta name="description" content="Realtime chat app" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main>
        <div>
          <h3 className="text-center">Welcome to realtime-chat-app</h3>
          {(() => {
            if (currUser === null) return <p>Loading...</p>;
            else if (currUser.length)
              return (
                <ChatWindow
                  currUsername={currUser}
                  receiver={receiver}
                  setReceiver={setReceiver}
                />
              );
            else return <LogInBox />;
          })()}
        </div>
      </main>
    </>
  );
}

const LogInBox = () => {
  const router = useRouter();
  return (
    <div className="mt-8 grid place-items-center">
      <Button
        className="m-2"
        color="primary"
        onClick={() => router.push("/login")}
      >
        Login
      </Button>
      <Button
        className="m-2"
        color="primary"
        onClick={() => router.push("/register")}
      >
        Register
      </Button>
    </div>
  );
};

const ChatWindow = ({
  currUsername,
  receiver,
  setReceiver,
}: {
  currUsername: string;
  receiver: string;
  setReceiver: (_: string) => void;
}) => {
  return (
    <div className="grid" style={{ gridTemplateColumns: "25% 75%" }}>
      <UsersList currUsername={currUsername} setReceiver={setReceiver} />
      <ChatBox sender={currUsername} receiver={receiver} />
    </div>
  );
};
