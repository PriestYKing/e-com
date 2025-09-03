import {
  LoginFormInputs,
  RegisterFormInputs,
  UserStoreActionsType,
  UserStoreStateType,
} from "@/types";
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
    }),
    {
      name: "user-store",
    }
  )
);
export default userStore;
