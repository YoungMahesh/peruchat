import { createContext, useContext } from "react";
import type { LoggedInUserInfo } from "~/utils/types.utils";

const UserContext = createContext<LoggedInUserInfo | null>(null);

export function useUser() {
  return useContext(UserContext);
}

export default UserContext;