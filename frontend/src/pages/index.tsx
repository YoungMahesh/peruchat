import { Button } from "@nextui-org/react";
import Head from "next/head";
import { useRouter } from "next/router";

export default function Home() {
  const router = useRouter();
  return (
    <>
      <Head>
        <title>Realtime Chat</title>
        <meta name="description" content="Realtime chat app" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main>
        <h3 className="text-center">Welcome to realtime-chat-app</h3>
        <div>
          <Button color="primary" onClick={() => router.push("/register")}>
            Register
          </Button>
        </div>
      </main>
    </>
  );
}
