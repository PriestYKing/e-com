"use client";

import { useCart } from "@/contexts/cartContext";
import useCartStore from "@/stores/cartStore";
import { get } from "http";
import { ShoppingCart } from "lucide-react";
import Link from "next/link";

const ShoppingCartIcon = () => {
  const { hasHydrated, getCartItemsCount } = useCart();
  if (!hasHydrated) return null;

  return (
    <Link href="/cart" className="relative">
      <ShoppingCart className="w-4 h-4 text-gray-600" />
      <span className="absolute -top-3 -right-3 bg-amber-400 text-gray-600 rounded-full w-4 h-4 flex items-center justify-center text-xs font-medium">
        {getCartItemsCount()}
      </span>
    </Link>
  );
};

export default ShoppingCartIcon;
