"use client";

import { FormEvent, useState } from "react";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Label } from "@radix-ui/react-label";
import { Button } from "./ui/button";
import { Input } from "./ui/input";
import { LoginFormInputs, loginFormSchema } from "@/types";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { User } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "./ui/dropdown-menu";
import { Separator } from "./ui/separator";
import userStore from "@/stores/userStore";
import { toast } from "react-toastify";
import UserMenu from "./UserMenu";
const Login = () => {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormInputs>({
    resolver: zodResolver(loginFormSchema),
  });

  const { user, setUser, setIsAuthenticated, isAuthenticated } = userStore();

  const handleLoginForm = async (data: LoginFormInputs) => {
    try {
      const res = await fetch("http://localhost:8080/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
        credentials: "include", // <-- Important for cookies!
      });

      if (res.ok) {
        const result = await res.json();
        console.log("User logged in successfully:", result);
        setIsAuthenticated(true);
        setUser(result.user);
      } else {
        const error = await res.json();
        toast.error(error.message || "Login failed");
        setIsAuthenticated(false);
        setUser(null);
      }
    } catch (err) {
      console.error("Failed to login user:", err);
      setIsAuthenticated(false);
      setUser(null);
    }
  };

  return (
    <>
      {isAuthenticated ? (
        <UserMenu />
      ) : (
        <Dialog>
          <DialogTrigger asChild>
            <p className="cursor-pointer text-gray-600 font-medium text-sm">
              Login
            </p>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[425px]">
            <DialogHeader>
              <DialogTitle>Login</DialogTitle>
              <DialogDescription>
                Please enter your credentials to sign in.
              </DialogDescription>
            </DialogHeader>
            <form onSubmit={handleSubmit(handleLoginForm)}>
              <div className="grid gap-4 mb-4">
                <div className="grid gap-3">
                  <Label htmlFor="email">Email</Label>
                  <Input
                    id="email"
                    type="text"
                    placeholder="John@Doe.com"
                    {...register("email")}
                  />
                  {errors.email && (
                    <p className="text-xs text-red-500">
                      {errors.email.message}
                    </p>
                  )}
                </div>
                <div className="grid gap-3">
                  <Label htmlFor="password">Password</Label>
                  <Input
                    id="password"
                    placeholder="******"
                    type="password"
                    {...register("password")}
                  />
                  {errors.password && (
                    <p className="text-xs text-red-500">
                      {errors.password.message}
                    </p>
                  )}
                </div>
              </div>
              <DialogFooter>
                <DialogClose asChild>
                  <Button variant="outline">Cancel</Button>
                </DialogClose>
                <Button type="submit">Login</Button>
              </DialogFooter>
            </form>
          </DialogContent>
        </Dialog>
      )}
    </>
  );
};

export default Login;
