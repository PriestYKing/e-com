"use client";
import React, {
  createContext,
  useContext,
  useState,
  useEffect,
  ReactNode,
  useMemo,
  useCallback,
} from "react";
import { CartItemType as CartItem, ShippingFormInputs } from "@/types";

interface CartContextType {
  cart: CartItem[];
  hasHydrated: boolean;
  shippingForm: ShippingFormInputs | null;
  addToCart: (product: CartItem) => void;
  removeFromCart: (
    product: Pick<CartItem, "id" | "selectedSize" | "selectedColor">
  ) => void;
  clearCart: () => void;
  setShippingForm: (data: ShippingFormInputs) => void;
  clearShippingForm: () => void;
  getCartTotal: () => number;
  getCartItemsCount: () => number;
  getCartSubtotal: () => number;
  getDiscount: () => number;
  getShippingFee: () => number;
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
  const [cart, setCart] = useState<CartItem[]>([]);
  const [hasHydrated, setHasHydrated] = useState(false);
  const [shippingForm, setShippingFormState] =
    useState<ShippingFormInputs | null>(null);

  // Safe localStorage loading
  const loadCartFromStorage = (): CartItem[] => {
    try {
      const savedCart = localStorage.getItem("cart");
      if (savedCart) {
        const parsedCart = JSON.parse(savedCart);
        if (Array.isArray(parsedCart)) {
          return parsedCart;
        } else {
          console.warn("Cart data in localStorage is not an array, resetting");
          localStorage.removeItem("cart");
        }
      }
    } catch (error) {
      console.error("Error parsing cart from localStorage:", error);
      localStorage.removeItem("cart");
    }
    return [];
  };

  // Load cart from localStorage on mount
  useEffect(() => {
    if (typeof window !== "undefined") {
      const loadedCart = loadCartFromStorage();
      setCart(loadedCart);
      setHasHydrated(true);
    }
  }, []);

  // Save cart to localStorage
  useEffect(() => {
    if (hasHydrated && Array.isArray(cart)) {
      try {
        localStorage.setItem("cart", JSON.stringify(cart));
      } catch (error) {
        console.error("Error saving cart to localStorage:", error);
      }
    }
  }, [cart, hasHydrated]);

  const addToCart = useCallback((product: CartItem) => {
    setCart((prevCart) => {
      if (!Array.isArray(prevCart)) {
        return [product];
      }

      const existingIndex = prevCart.findIndex(
        (p) =>
          p.id === product.id &&
          p.selectedSize === product.selectedSize &&
          p.selectedColor === product.selectedColor
      );

      if (existingIndex !== -1) {
        const updatedCart = [...prevCart];
        updatedCart[existingIndex] = {
          ...updatedCart[existingIndex],
          quantity: updatedCart[existingIndex].quantity + product.quantity,
        };
        return updatedCart;
      }

      return [...prevCart, product];
    });
  }, []);

  const removeFromCart = useCallback(
    (product: Pick<CartItem, "id" | "selectedSize" | "selectedColor">) => {
      setCart((prevCart) => {
        if (!Array.isArray(prevCart)) {
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
    },
    []
  );

  const clearCart = useCallback(() => {
    setCart([]);
  }, []);

  const setShippingForm = useCallback((data: ShippingFormInputs) => {
    setShippingFormState(data);
  }, []);

  const clearShippingForm = useCallback(() => {
    setShippingFormState(null);
  }, []);

  const getCartSubtotal = useCallback(() => {
    if (!Array.isArray(cart)) return 0;
    return cart.reduce((total, item) => total + item.price * item.quantity, 0);
  }, [cart]);

  const getDiscount = useCallback(() => {
    // 10% discount logic or fixed amount
    return 10;
  }, []);

  const getShippingFee = useCallback(() => {
    return 10;
  }, []);

  const getCartTotal = useCallback(() => {
    const subtotal = getCartSubtotal();
    const discount = getDiscount();
    const shipping = getShippingFee();
    return subtotal - discount + shipping;
  }, [getCartSubtotal, getDiscount, getShippingFee]);

  const getCartItemsCount = useCallback(() => {
    if (!Array.isArray(cart)) return 0;
    return cart.reduce((total, item) => total + item.quantity, 0);
  }, [cart]);

  const value = useMemo(
    () => ({
      cart,
      hasHydrated,
      shippingForm,
      addToCart,
      removeFromCart,
      clearCart,
      setShippingForm,
      clearShippingForm,
      getCartTotal,
      getCartItemsCount,
      getCartSubtotal,
      getDiscount,
      getShippingFee,
    }),
    [
      cart,
      hasHydrated,
      shippingForm,
      addToCart,
      removeFromCart,
      clearCart,
      setShippingForm,
      clearShippingForm,
      getCartTotal,
      getCartItemsCount,
      getCartSubtotal,
      getDiscount,
      getShippingFee,
    ]
  );

  return <CartContext.Provider value={value}>{children}</CartContext.Provider>;
};
