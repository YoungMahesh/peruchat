/* eslint-disable @typescript-eslint/no-misused-promises */
/* eslint-disable @typescript-eslint/no-unsafe-argument */
import { Button } from "@nextui-org/react";
import Head from "next/head";
import { useRouter } from "next/router";
import {  useState } from "react";
import ChatBox from "~/components/ChatBox";
import { useUser } from "~/components/UserContext";
import UsersList from "~/components/UsersList";

export default function Home() {
  const [receiver, setReceiver] = useState<string>("");
  const currUser = useUser() 

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
            else if (currUser.username.length)
              return (
                <ChatWindow
                  currUsername={currUser.username}
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
