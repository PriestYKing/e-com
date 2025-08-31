import { PaymentFormInputs, paymentFormSchema } from "@/types";
import { SubmitHandler, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowRight, ShoppingCart, ShoppingCartIcon } from "lucide-react";
import { useRouter } from "next/navigation";
import Image from "next/image";

const PaymentForm = () => {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<PaymentFormInputs>({
    resolver: zodResolver(paymentFormSchema),
  });
  const router = useRouter();
  const handlePaymentForm: SubmitHandler<PaymentFormInputs> = (data) => {};
  return (
    <form
      className="flex flex-col gap-4"
      onSubmit={handleSubmit(handlePaymentForm)}
    >
      <div className="flex flex-col gap-1">
        <label
          className="text-xs font-medium text-gray-500"
          htmlFor="cardHolder"
        >
          Name on card
        </label>
        <input
          type="text"
          className="border-b border-gray-200 py-2 outline-none text-sm"
          id="name"
          placeholder="John Doe"
          {...register("cardHolder")}
        />
        {errors.cardHolder && (
          <p className="text-xs text-red-500">{errors.cardHolder.message}</p>
        )}
      </div>
      <div className="flex flex-col gap-1">
        <label
          className="text-xs font-medium text-gray-500"
          htmlFor="cardNumber"
        >
          Card Number
        </label>
        <input
          type="text"
          className="border-b border-gray-200 py-2 outline-none text-sm"
          id="cardNumber"
          placeholder="1234 5678 9012 3456"
          {...register("cardNumber")}
        />
        {errors.cardNumber && (
          <p className="text-xs text-red-500">{errors.cardNumber.message}</p>
        )}
      </div>
      <div className="flex flex-col gap-1">
        <label
          className="text-xs font-medium text-gray-500"
          htmlFor="expirationDate"
        >
          Expiration Date
        </label>
        <input
          type="text"
          className="border-b border-gray-200 py-2 outline-none text-sm"
          id="expirationDate"
          placeholder="MM/YY"
          {...register("expirationDate")}
        />
        {errors.expirationDate && (
          <p className="text-xs text-red-500">
            {errors.expirationDate.message}
          </p>
        )}
      </div>
      <div className="flex flex-col gap-1">
        <label className="text-xs font-medium text-gray-500" htmlFor="cvv">
          CVV
        </label>
        <input
          type="text"
          className="border-b border-gray-200 py-2 outline-none text-sm"
          id="cvv"
          placeholder="123"
          {...register("cvv")}
        />
        {errors.cvv && (
          <p className="text-xs text-red-500">{errors.cvv.message}</p>
        )}
      </div>
      <div className="flex items-center gap-2 mt-4">
        <Image
          src="/klarna.png"
          alt="klarna"
          className="rounded-md"
          width={50}
          height={25}
        />
        <Image
          src="/cards.png"
          alt="cards"
          className="rounded-md"
          width={50}
          height={25}
        />
        <Image
          src="/stripe.png"
          alt="stripe"
          className="rounded-md"
          width={50}
          height={25}
        />
      </div>
      <button
        type="submit"
        className="w-full bg-gray-800 text-white p-2 rounded-lg cursor-pointer flex items-center justify-center gap-2 hover:bg-gray-900 transition-all duration-300"
      >
        Checkout <ShoppingCart className="w-3 h-3" />
      </button>
    </form>
  );
};

export default PaymentForm;
