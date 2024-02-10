/* eslint-disable @next/next/no-img-element */
import { Button, Input } from "@nextui-org/react";
import Head from "next/head";
import Link from "next/link";
import { useRouter } from "next/router";
import { useEffect, useState } from "react";
import {  toast } from "react-toastify";
import authStore from "~/utils/auth";
import {  toastOptions } from "~/utils/constants";

export default function RegisterPage() {
  const router = useRouter();
  const [username, setUsername] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmPass, setConfirmPass] = useState("");

  const handleSubmit = async () => {
    try {
      if (!handleValidation())
        return toast.error("Please fill all the fields.", toastOptions);

      const res = await fetch("http://localhost:3001/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username, email, password }),
      });

      const data = (await res.json()) as unknown as { message: string };
      if (res.status !== 201) {
        alert(data.message);
        return toast.error(data.message, toastOptions);
      }

      localStorage.setItem(authStore.USER_IDENTY_KEY, JSON.stringify(username));
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
    if (password !== confirmPass) {
      toast.error(
        "Password and confirm password should be same.",
        toastOptions,
      );
      return false;
    } else if (username.length < 3) {
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
    } else if (email === "") {
      toast.error("Email is required.", toastOptions);
      return false;
    }
    return true;
  };

  return (
    <>
      <Head>
        <title>Register | Peru Chat</title>
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
          <Input
            type="password"
            label="Confirm Password"
            value={confirmPass}
            onChange={(e) => setConfirmPass(e.target.value)}
          />
          <div className="mt-2 flex justify-end">
            <Button color="primary" onClick={handleSubmit}>
              Register
            </Button>
          </div>
          <span>
            Already have an account ? <Link href="/login">Login.</Link>
          </span>
        </div>
      </main>
    </>
  );
}
