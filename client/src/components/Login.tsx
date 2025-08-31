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
const Login = () => {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormInputs>({
    resolver: zodResolver(loginFormSchema),
  });
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const handleLoginForm = (data: LoginFormInputs) => {
    setIsAuthenticated(true);
  };

  return (
    <>
      {isAuthenticated ? (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <button className="cursor-pointer">
              <User className="w-4 h-4 text-gray-600" />
            </button>
          </DropdownMenuTrigger>
          <DropdownMenuContent className="w-32" align="start">
            <DropdownMenuLabel>My Account</DropdownMenuLabel>
            <DropdownMenuGroup>
              <DropdownMenuItem>Profile</DropdownMenuItem>
              <DropdownMenuItem>Settings</DropdownMenuItem>
              <DropdownMenuItem className="mb-2">Logout</DropdownMenuItem>
              <Separator />
              <DropdownMenuItem className="mt-2">Yash Patidar</DropdownMenuItem>
            </DropdownMenuGroup>
          </DropdownMenuContent>
        </DropdownMenu>
      ) : (
        <Dialog>
          <DialogTrigger asChild>
            <p className="cursor-pointer text-gray-600 font-medium text-sm">
              Sign in
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
                  <Label htmlFor="username">Username</Label>
                  <Input
                    id="username"
                    type="text"
                    placeholder="John Doe"
                    {...register("username")}
                  />
                  {errors.username && (
                    <p className="text-xs text-red-500">
                      {errors.username.message}
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
