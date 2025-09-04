import { User } from "lucide-react";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuTrigger,
} from "./ui/dropdown-menu";
import { Separator } from "@radix-ui/react-dropdown-menu";
import userStore from "@/stores/userStore";

const UserMenu = () => {
  const { logout } = userStore();

  return (
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
          <DropdownMenuItem className="mb-2" onClick={logout}>
            Logout
          </DropdownMenuItem>
          <Separator />
          <DropdownMenuItem className="mt-2">Yash Patidar</DropdownMenuItem>
        </DropdownMenuGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

export default UserMenu;
