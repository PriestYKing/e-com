"use client";

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
import { useAuth } from "@/contexts/AuthContext";
import { toast } from "sonner";
import UserMenu from "./UserMenu";

const Login = () => {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormInputs>({
    resolver: zodResolver(loginFormSchema),
  });

  // Replace zustand with Auth context
  const { login, isAuthenticated } = useAuth();

  const handleLoginForm = async (data: LoginFormInputs) => {
    try {
      const res = await fetch("http://localhost:8080/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
        credentials: "include", // Important for cookies!
      });

      if (res.ok) {
        const result = await res.json();
        // The JWT token is now in cookies, so just call login() to decode it
        login();
        toast.success("Login successful");
      } else {
        const error = await res.json();
        toast.error(error.message || "Login failed");
      }
    } catch (err) {
      console.error("Failed to login user:", err);
      toast.error("Network error occurred");
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
