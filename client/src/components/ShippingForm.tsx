import { ShippingFormInputs, shippingFormSchema } from "@/types";
import { SubmitHandler, useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowRight } from "lucide-react";
import { useRouter } from "next/navigation";

const ShippingForm = ({
  setShippingForm,
}: {
  setShippingForm: (data: ShippingFormInputs) => void;
}) => {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<ShippingFormInputs>({
    resolver: zodResolver(shippingFormSchema),
  });
  const router = useRouter();
  const handleShippingForm: SubmitHandler<ShippingFormInputs> = (data) => {
    setShippingForm(data);
    router.push("/cart?step=3", { scroll: false });
  };
  return (
    <form
      className="flex flex-col gap-4"
      onSubmit={handleSubmit(handleShippingForm)}
    >
      <div className="flex flex-col gap-1">
        <label className="text-xs font-medium text-gray-500" htmlFor="name">
          Name
        </label>
        <input
          type="text"
          className="border-b border-gray-200 py-2 outline-none text-sm"
          id="name"
          placeholder="John Doe"
          {...register("name")}
        />
        {errors.name && (
          <p className="text-xs text-red-500">{errors.name.message}</p>
        )}
      </div>
      <div className="flex flex-col gap-1">
        <label className="text-xs font-medium text-gray-500" htmlFor="email">
          Email
        </label>
        <input
          type="email"
          className="border-b border-gray-200 py-2 outline-none text-sm"
          id="email"
          placeholder="John.Doe@example.com"
          {...register("email")}
        />
        {errors.email && (
          <p className="text-xs text-red-500">{errors.email.message}</p>
        )}
      </div>
      <div className="flex flex-col gap-1">
        <label className="text-xs font-medium text-gray-500" htmlFor="Phone">
          Phone
        </label>
        <input
          type="text"
          className="border-b border-gray-200 py-2 outline-none text-sm"
          id="phone"
          placeholder="(123) 456-7890"
          {...register("phone")}
        />
        {errors.phone && (
          <p className="text-xs text-red-500">{errors.phone.message}</p>
        )}
      </div>
      <div className="flex flex-col gap-1">
        <label className="text-xs font-medium text-gray-500" htmlFor="address">
          Address
        </label>
        <input
          type="text"
          className="border-b border-gray-200 py-2 outline-none text-sm"
          id="address"
          placeholder="123 Main St"
          {...register("address")}
        />
        {errors.address && (
          <p className="text-xs text-red-500">{errors.address.message}</p>
        )}
      </div>
      <div className="flex flex-col gap-1">
        <label className="text-xs font-medium text-gray-500" htmlFor="city">
          City
        </label>
        <input
          type="text"
          className="border-b border-gray-200 py-2 outline-none text-sm"
          id="city"
          placeholder="New York"
          {...register("city")}
        />
        {errors.city && (
          <p className="text-xs text-red-500">{errors.city.message}</p>
        )}
      </div>
      <button
        type="submit"
        className="w-full bg-gray-800 text-white p-2 rounded-lg cursor-pointer flex items-center justify-center gap-2 hover:bg-gray-900 transition-all duration-300"
      >
        Continue <ArrowRight className="w-3 h-3" />
      </button>
    </form>
  );
};

export default ShippingForm;
