"use client";
import React, {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
} from "react";
import { CartItemType as CartItem } from "@/types";

interface CartContextType {
  cart: CartItem[];
  hasHydrated: boolean;
  addToCart: (product: CartItem) => void;
  removeFromCart: (
    product: Pick<CartItem, "id" | "selectedSize" | "selectedColor">
  ) => void;
  clearCart: () => void;
  getCartTotal: () => number;
  getCartItemsCount: () => number;
}

const CartContext = createContext<CartContextType | undefined>(undefined);

export const useCart = () => {
  const context = useContext(CartContext);
  if (context === undefined) {
    throw new Error("useCart must be used within a CartProvider");
  }
  return context;
};

interface CartProviderProps {
  children: ReactNode;
}

export const CartProvider = ({ children }: CartProviderProps) => {
  // Initialize cart as an empty array
  const [cart, setCart] = useState<CartItem[]>([]);
  const [hasHydrated, setHasHydrated] = useState(false);

  // Function to safely load cart from localStorage
  const loadCartFromStorage = (): CartItem[] => {
    try {
      const savedCart = localStorage.getItem("cart");
      if (savedCart) {
        const parsedCart = JSON.parse(savedCart);

        // Validate that parsed data is actually an array
        if (Array.isArray(parsedCart)) {
          return parsedCart;
        } else {
          console.warn(
            "Cart data in localStorage is not an array, resetting to empty array"
          );
          localStorage.removeItem("cart");
        }
      }
    } catch (error) {
      console.error("Error parsing cart from localStorage:", error);
      localStorage.removeItem("cart");
    }
    return []; // Always return an array as fallback
  };

  // Load cart from localStorage on mount
  useEffect(() => {
    if (typeof window !== "undefined") {
      const loadedCart = loadCartFromStorage();
      setCart(loadedCart);
      setHasHydrated(true);
    }
  }, []);

  // Save cart to localStorage whenever cart changes
  useEffect(() => {
    if (hasHydrated && Array.isArray(cart)) {
      try {
        localStorage.setItem("cart", JSON.stringify(cart));
      } catch (error) {
        console.error("Error saving cart to localStorage:", error);
      }
    }
  }, [cart, hasHydrated]);

  const addToCart = (product: CartItem) => {
    setCart((prevCart) => {
      // Ensure prevCart is an array
      if (!Array.isArray(prevCart)) {
        console.warn("Previous cart state is not an array, resetting");
        return [product];
      }

      const existingIndex = prevCart.findIndex(
        (p) =>
          p.id === product.id &&
          p.selectedSize === product.selectedSize &&
          p.selectedColor === product.selectedColor
      );

      if (existingIndex !== -1) {
        // Product already exists, update quantity
        const updatedCart = [...prevCart];
        updatedCart[existingIndex] = {
          ...updatedCart[existingIndex],
          quantity: updatedCart[existingIndex].quantity + product.quantity,
        };
        return updatedCart;
      }

      // Product doesn't exist, add new item
      return [...prevCart, product];
    });
  };

  const removeFromCart = (
    product: Pick<CartItem, "id" | "selectedSize" | "selectedColor">
  ) => {
    setCart((prevCart) => {
      // Ensure prevCart is an array
      if (!Array.isArray(prevCart)) {
        console.warn("Previous cart state is not an array, resetting");
        return [];
      }

      return prevCart.filter(
        (p) =>
          !(
            p.id === product.id &&
            p.selectedSize === product.selectedSize &&
            p.selectedColor === product.selectedColor
          )
      );
    });
  };

  const clearCart = () => {
    setCart([]);
  };

  const getCartTotal = () => {
    // Ensure cart is an array before calling reduce
    if (!Array.isArray(cart)) {
      console.warn("Cart is not an array in getCartTotal");
      return 0;
    }
    return cart.reduce((total, item) => total + item.price * item.quantity, 0);
  };

  const getCartItemsCount = () => {
    // Ensure cart is an array before calling reduce
    if (!Array.isArray(cart)) {
      console.warn("Cart is not an array in getCartItemsCount");
      return 0;
    }
    return cart.reduce((total, item) => total + item.quantity, 0);
  };

  const value = {
    cart,
    hasHydrated,
    addToCart,
    removeFromCart,
    clearCart,
    getCartTotal,
    getCartItemsCount,
  };

  return <CartContext.Provider value={value}>{children}</CartContext.Provider>;
};
