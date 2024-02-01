import { Button } from "@nextui-org/react";
import Head from "next/head";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import ChatBox from "~/components/ChatBox";
import UsersList from "~/components/UsersList";
import { USER_IDENTY_KEY } from "~/utils/constants";

export default function Home() {
  const [currUser, setCurrUser] = useState<string | null>(null);

  useEffect(() => {
    const _currUser = localStorage.getItem(USER_IDENTY_KEY);
    if (_currUser) setCurrUser(_currUser);
    else setCurrUser("");
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
              return <ChatWindow currUsername={currUser} />;
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

const ChatWindow = ({ currUsername }: { currUsername: string }) => {
  return (
    <div className="grid" style={{ gridTemplateColumns: "25% 75%" }}>
      <UsersList currUsername={currUsername} />
      <ChatBox />
    </div>
  );
};
