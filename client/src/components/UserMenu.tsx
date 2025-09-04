import { User } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "./ui/dropdown-menu";

import userStore from "@/stores/userStore";
import { Separator } from "./ui/separator";

const UserMenu = () => {
  const { logout, user } = userStore();

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <button className="cursor-pointer">
          <User className="w-4 h-4 text-gray-600" />
        </button>
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-32 cursor-pointer" align="start">
        <DropdownMenuLabel>My Account</DropdownMenuLabel>
        <DropdownMenuGroup>
          <DropdownMenuItem>Profile</DropdownMenuItem>
          <DropdownMenuItem>Settings</DropdownMenuItem>
          <DropdownMenuItem className="mb-2" onClick={logout}>
            Logout
          </DropdownMenuItem>
          <Separator />
          <DropdownMenuItem className="mt-2">{user?.name}</DropdownMenuItem>
        </DropdownMenuGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export default UserMenu;
