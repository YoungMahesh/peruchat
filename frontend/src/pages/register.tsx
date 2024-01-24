import { Button, Input } from "@nextui-org/react";
import Head from "next/head";
import { useState } from "react";

export default function Home() {

  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  const handleSubmit = async () => {
    try {
      const res = await fetch("http://localhost:3001/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, email, password }),
      });
      console.log('status', res.status)
      const data = (await res.json()) as unknown;
      console.log(data);
    } catch (error) {
      console.error(error);
    }
  };

  return (
    <>
      <Head>
        <title>Register</title>
        <meta name="description" content="Register page for realtime-chat" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main>
        <h3 className="text-center">Register to realtime-chat-app</h3>
        <div className="m-auto mt-4 max-w-xs">
          <Input
            type="text"
            label="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
          <Input
            type="email"
            label="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
          />
          <Input
            type="password"
            label="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <div className="mt-2 flex justify-end">
            <Button color="primary" onClick={handleSubmit}>
              Submit
            </Button>
          </div>
        </div>
      </main>
    </>
  );
}
