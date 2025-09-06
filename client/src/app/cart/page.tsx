"use client";

import Login from "@/components/Login";
import PaymentForm from "@/components/PaymentForm";
import ShippingForm from "@/components/ShippingForm";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/contexts/AuthContext";
import { useCart } from "@/contexts/cartContext";
import { ArrowRight, Trash2 } from "lucide-react";
import Image from "next/image";
import { useRouter, useSearchParams } from "next/navigation";
import { useState } from "react";

const steps = [
  {
    id: 1,
    title: "Shopping Cart",
  },
  {
    id: 2,
    title: "Shipping Address",
  },
  {
    id: 3,
    title: "Payment Method",
  },
];

const CartPage = () => {
  const searchParams = useSearchParams();
  const router = useRouter();
  const activeStep = parseInt(searchParams.get("step") || "1");

  // Get authentication state
  const { isAuthenticated } = useAuth();

  // Use context instead of local state
  const {
    cart,
    removeFromCart,
    shippingForm,
    getCartSubtotal,
    getDiscount,
    getShippingFee,
    getCartTotal,
  } = useCart();

  const handleContinueToShipping = () => {
    if (!isAuthenticated) {
      return; // Login component will be shown instead
    }
    router.push("/cart?step=2", { scroll: false });
  };

  // Show login modal/component if user is not authenticated
  if (!isAuthenticated) {
    return (
      <div className="flex flex-col gap-8 items-center justify-center mt-12 min-h-[60vh]">
        <div className="text-center">
          <h1 className="text-2xl font-medium mb-4">Authentication Required</h1>
          <p className="text-gray-600 mb-8">
            Please log in to access your shopping cart
          </p>
          <Button variant="outline">
            <Login />
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-8 items-center justify-center mt-12 ">
      {/* TITLE */}
      <h1 className="text-2xl font-medium">Your Shopping Cart</h1>

      {/* STEPS */}
      <div className="flex flex-col lg:flex-row items-center gap-8 lg:gap-16">
        {steps.map((step) => (
          <div
            key={step.id}
            className={`flex items-center gap-2 border-b-2 pb-4 ${
              step.id === activeStep ? "border-gray-800" : "border-gray-200"
            }`}
          >
            <div
              className={`w-6 h-6 rounded-full text-white p-4 flex items-center justify-center ${
                step.id === activeStep ? "bg-gray-800" : "bg-gray-400"
              }`}
            >
              {step.id}
            </div>
            <p
              className={`text-sm font-medium ${
                step.id === activeStep ? "text-gray-800" : "text-gray-400"
              }`}
            >
              {step.title}
            </p>
          </div>
        ))}
      </div>

      {/* STEPS & DETAILS */}
      <div className="w-full flex flex-col lg:flex-row gap-16">
        {/* STEPS */}
        <div className="w-full lg:w-7/12 shadow-lg border-1 border-gray-100 p-8 rounded-lg flex flex-col gap-8">
          {activeStep === 1 ? (
            cart.map((item) => (
              <div
                key={item.id + item.selectedColor + item.selectedSize}
                className="flex items-center justify-between"
              >
                <div className="flex gap-8">
                  <div className="relative w-32 h-32 bg-gray-50 rounded-lg overflow-hidden">
                    <Image
                      src={item.images[item.selectedColor]}
                      alt={item.name}
                      fill
                      className="object-contain"
                    />
                  </div>
                  <div className="flex flex-col justify-between">
                    <div className="flex flex-col gap-1">
                      <p className="text-sm font-medium">{item.name}</p>
                      <p className="text-xs text-gray-500">
                        Quantity : {item.quantity}
                      </p>
                      <p className="text-xs text-gray-500">
                        Size: {item.selectedSize.toUpperCase()}
                      </p>
                      <p className="text-xs text-gray-500">
                        Color: {item.selectedColor.toUpperCase()}
                      </p>
                    </div>
                    <p className="font-medium">${item.price.toFixed(2)}</p>
                  </div>
                </div>
                <button
                  onClick={() => removeFromCart(item)}
                  className="w-8 h-8 rounded-full bg-red-100 text-red-400 flex items-center justify-center cursor-pointer hover:bg-red-200 transition-all duration-300"
                >
                  <Trash2 className="w-3 h-3" />
                </button>
              </div>
            ))
          ) : activeStep === 2 ? (
            <ShippingForm />
          ) : activeStep === 3 && shippingForm ? (
            <PaymentForm />
          ) : (
            <p className="text-sm text-gray-500">
              Please fill in the shipping form to continue.
            </p>
          )}
        </div>

        {/* DETAILS */}
        <div className="w-full lg:w-5/12 shadow-lg border-1 border-gray-100 p-8 rounded-lg flex flex-col gap-8 h-max">
          <h2 className="font-semibold">Cart Details</h2>
          <div className="flex flex-col gap-4">
            <div className="flex justify-between text-sm">
              <p className=" text-gray-500">Subtotal</p>
              <p className=" font-medium">${getCartSubtotal().toFixed(2)}</p>
            </div>
            <div className="flex justify-between text-sm">
              <p className=" text-gray-500">Discount</p>
              <p className=" font-medium">-${getDiscount().toFixed(2)}</p>
            </div>
            <div className="flex justify-between text-sm">
              <p className=" text-gray-500">Shipping Fee</p>
              <p className=" font-medium">${getShippingFee().toFixed(2)}</p>
            </div>
            <hr className="border-gray-200" />
            <div className="flex justify-between">
              <p className=" text-gray-800 font-semibold">Total</p>
              <p className=" font-medium">${getCartTotal().toFixed(2)}</p>
            </div>
          </div>

          {/* Show shipping info if completed */}
          {/* {shippingForm && activeStep > 2 && (
            <div className="border-t pt-4">
              <h3 className="font-medium mb-2">Shipping Address</h3>
              <div className="text-sm text-gray-600">
                <p>{shippingForm.name}</p>
                <p>{shippingForm.email}</p>
                <p>{shippingForm.address}</p>
                <p>{shippingForm.city}</p>
                <p>{shippingForm.phone}</p>
              </div>
            </div>
          )} */}

          {activeStep === 1 && (
            <button
              className="w-full bg-gray-800 text-white p-2 rounded-lg cursor-pointer flex items-center justify-center gap-2 hover:bg-gray-900 transition-all duration-300"
              onClick={handleContinueToShipping}
            >
              Continue <ArrowRight className="w-3 h-3" />
            </button>
          )}
        </div>
      </div>
    </div>
  );
};

export default CartPage;
