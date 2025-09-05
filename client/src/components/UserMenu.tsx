import { Cog, LogOut, User } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "./ui/dropdown-menu";

import { Separator } from "./ui/separator";
import { useContext } from "react";
import { useAuth } from "@/contexts/AuthContext";

const UserMenu = () => {
  const { user, logout } = useAuth();

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <button className="cursor-pointer">
          <User className="w-4 h-4 text-gray-600" />
        </button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-32 " align="start">
        <DropdownMenuLabel>My Account</DropdownMenuLabel>
        <DropdownMenuGroup>
          <DropdownMenuItem className="cursor-pointer">
            <User /> Profile
          </DropdownMenuItem>
          <DropdownMenuItem className="cursor-pointer">
            <Cog /> Settings
          </DropdownMenuItem>
          <DropdownMenuItem className="mb-2 cursor-pointer" onClick={logout}>
            <LogOut /> Logout
          </DropdownMenuItem>
          <Separator />
          <DropdownMenuItem className="mt-2">
            {user?.name.toUpperCase()}
          </DropdownMenuItem>
        </DropdownMenuGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export default UserMenu;
