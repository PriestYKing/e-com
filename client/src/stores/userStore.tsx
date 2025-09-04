import { UserStoreActionsType, UserStoreStateType } from "@/types";
import { toast } from "sonner";

import { create } from "zustand";
import { persist } from "zustand/middleware";

const userStore = create<UserStoreStateType & UserStoreActionsType>()(
  persist(
    (set) => ({
      user: null,
      isAuthenticated: false,
      setUser: (user) => set((state) => ({ ...state, user })),
      setIsAuthenticated: (authStatus: boolean) =>
        set((state) => ({ ...state, isAuthenticated: authStatus })),
      logout: async () => {
        await fetch("http://localhost:8080/logout", {
          method: "POST",
          credentials: "include",
        })
          .then((res) => {
            if (res.ok) {
              set({ user: null, isAuthenticated: false });
              toast.info("Logout successful");
            }
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
