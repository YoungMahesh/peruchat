/* eslint-disable @typescript-eslint/no-misused-promises */
/* eslint-disable @typescript-eslint/no-unsafe-argument */
import { Button } from "@nextui-org/react";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import ChatBox from "~/components/ChatBox";
import Header from "~/components/Header";
import PageTitle from "~/components/PageTitle";
import { useUser } from "~/components/UserContext";
import UsersList from "~/components/UsersList";
import authStore from "~/utils/auth";

export default function Home() {
  const router = useRouter();
  const [receiver, setReceiver] = useState<string>("");
  const currUser = useUser();

  useEffect(() => {
    if (!localStorage.getItem(authStore.USER_IDENTY_KEY)) {
      void router.push("/login");
    }
  }, [router]);

  return (
    <>
      <PageTitle pageName="Chat" pageDescription="" />
      <main>
        <Header />
        <div>
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
    <div
      className="grid"
      style={{ gridTemplateColumns: "minmax(0, max-content) 1fr" }}
    >
      <UsersList
        currUsername={currUsername}
        receiver={receiver}
        setReceiver={setReceiver}
      />
      <ChatBox sender={currUsername} receiver={receiver} />
    </div>
  );
};
