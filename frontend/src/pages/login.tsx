/* eslint-disable @next/next/no-img-element */
import { Button, Input } from "@nextui-org/react";
import Head from "next/head";
import Link from "next/link";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import { toast } from "react-toastify";
import authStore from "~/utils/auth";
import { toastOptions } from "~/utils/constants";

export default function LoginPage() {
  const router = useRouter();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const handleSubmit = async () => {
    try {
      if (!handleValidation())
        return toast.error("Please fill all the fields.", toastOptions);

      const res = await fetch("http://localhost:3001/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, password }),
      });

      const data = (await res.json()) as unknown as {
        message: string;
        token: string;
      };
      if (res.status !== 200) {
        alert(data.message);
        return toast.error(data.message, toastOptions);
      }
      authStore.storeToken(username, data.token);
      void router.push("/");
    } catch (error) {
      console.error(error);
      toast.error("An error occurred, check the console for more information");
    }
  };

  useEffect(() => {
    if (localStorage.getItem(authStore.USER_IDENTY_KEY)) {
      void router.push("/");
    }
  }, [router]);

  const handleValidation = () => {
    if (username.length < 3) {
      toast.error(
        "Username should be greater than 3 characters.",
        toastOptions,
      );
      return false;
    } else if (password.length < 8) {
      toast.error(
        "Password should be equal or greater than 8 characters.",
        toastOptions,
      );
      return false;
    }
    return true;
  };

  return (
    <>
      <Head>
        <title>Login | Peru Chat</title>
        <meta name="description" content="Register page for peru chat" />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <main className="p-2">
        <div className="flex items-center justify-center">
          <img
            style={{ width: "50px", backgroundColor: "white" }}
            src={"/logo1.svg"}
            alt="logo"
          />
          <h3 className="m-2 text-center">PeruChat</h3>
        </div>

        <div className="m-auto mt-4 max-w-xs">
          <Input
            type="text"
            label="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
          <Input
            type="password"
            label="Password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <div className="mt-2 flex justify-end">
            <Button color="primary" onClick={handleSubmit}>
              Login
            </Button>
          </div>
          <span>
            Don&apos;t have an account ? <Link href="/register">Register</Link>
          </span>
        </div>
      </main>
    </>
  );
}
