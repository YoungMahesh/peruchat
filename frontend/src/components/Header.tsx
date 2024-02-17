import { Button } from "@nextui-org/react";
import { useRouter } from "next/router";
import authStore from "~/utils/auth";

export default function Header() {
  const router = useRouter();
  return (
    <header className="flex justify-between p-2">
      <h3>PeruChat</h3>
      <Button
        onClick={() => {
          localStorage.removeItem(authStore.USER_IDENTY_KEY);
          void router.push("/login");
        }}
      >
        Logout
      </Button>
    </header>
  );
}
