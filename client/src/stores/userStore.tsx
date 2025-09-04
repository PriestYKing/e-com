import {
  LoginFormInputs,
  RegisterFormInputs,
  UserStoreActionsType,
  UserStoreStateType,
} from "@/types";

import { create } from "zustand";
import { persist } from "zustand/middleware";
import Cookies from "js-cookie";
const userStore = create<UserStoreStateType & UserStoreActionsType>()(
  persist(
    (set) => ({
      user: null,
      isAuthenticated: false,
      setUser: (user) => set((state) => ({ ...state, user })),
      setIsAuthenticated: (authStatus: boolean) =>
        set((state) => ({ ...state, isAuthenticated: authStatus })),
      logout: async () => {
        const accessToken = Cookies.get("access_token");

        const headers = new Headers({
          Authorization: `Bearer ${accessToken}`,
        });

        fetch("http://localhost:8080/logout", {
          method: "POST",
          headers,
          credentials: "include",
        })
          .then(() => {
            set({ user: null, isAuthenticated: false });
          })
          .catch((err) => {
            console.error("Failed to logout:", err);
          });
      },
    }),
    {
      name: "user-store",
    }
  )
);
export default userStore;
