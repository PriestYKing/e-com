import { z } from "zod";
export type ProductType = {
  id: string | number;
  name: string;
  shortDescription: string;
  description: string;
  price: number;
  sizes: string[];
  colors: string[];
  images: Record<string, string>;
};

export type ProductsType = ProductType[];

export type CartItemType = ProductType & {
  quantity: number;
  selectedSize: string;
  selectedColor: string;
};

export type CartItemsType = CartItemType[];

export const shippingFormSchema = z.object({
  name: z.string().min(1, "Name is required"),
  email: z.string().min(1, "Email is required"),
  phone: z
    .string()
    .min(7, "Phone number must be atleast 7 digits")
    .max(10, "Phone number can't be more than 10 digits"),
  //.regex(/^d+$/, "Phone number must contain only numbers"),
  address: z.string().min(1, "Address is required"),
  city: z.string().min(1, "City is required"),
});

export type ShippingFormInputs = z.infer<typeof shippingFormSchema>;

export const paymentFormSchema = z.object({
  cardHolder: z.string().min(1, "Cardholder is required"),
  cardNumber: z
    .string()
    .min(16, "Card Number is required")
    .max(16, "Card Number is required"),
  expirationDate: z
    .string()
    .regex(/^(0[1-9]|1[0-2])\/\d{2}$/, "Expiration date must be MM/YY format"),
  cvv: z.string().min(3, "CVV is required").max(3, "CVV is required"),
});

export type PaymentFormInputs = z.infer<typeof paymentFormSchema>;

export type CartStoreStateType = {
  cart: CartItemsType;
  hasHydrated: boolean;
};

export type CartStoreActionsType = {
  addToCart: (product: CartItemType) => void;
  removeFromCart: (product: CartItemType) => void;
  clearCart: () => void;
};

export const loginFormSchema = z.object({
  email: z.email("Email is invalid"),
  password: z
    .string()
    .min(6, "Password must be at least 6 characters long")
    .max(20, "Password can't be more than 20 characters long"),
});

export type LoginFormInputs = z.infer<typeof loginFormSchema>;

export const registerFormSchema = z.object({
  name: z
    .string()
    .min(2, "Name must be at least 2 characters long")
    .max(100, "Name can't be more than 100 characters long"),
  email: z.email("Email is invalid"),
  password: z
    .string()
    .min(6, "Password must be at least 6 characters long")
    .max(20, "Password can't be more than 20 characters long"),
});

export type RegisterFormInputs = z.infer<typeof registerFormSchema>;
