import { useEffect, useState } from "react";

export default function UsersList({
  currUsername,
  setReceiver,
}: {
  currUsername: string;
  setReceiver: (_: string) => void;
}) {
  const [usersList, setUserList] = useState<{ username: string }[]>([]);
  console.log({ usersList });
  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const res = await fetch("http://localhost:3001/users", {
          method: "GET",
          headers: { "Content-Type": "application/json" },
        });

        const _users = (await res.json()) as unknown as {
          username: string;
        }[];
        if (res.status !== 200)
          return alert("could not able to fetch users list");
        setUserList(_users.filter((user) => user.username !== currUsername));
      } catch (error) {
        console.error(error);
      }
    };

    void fetchUsers();
  }, [currUsername]);

  return (
    <section className="flex flex-col">
      {usersList.map((user, idx) => (
        <div key={idx} onClick={() => setReceiver(user.username)}>
          {user.username}
        </div>
      ))}
    </section>
  );
}
