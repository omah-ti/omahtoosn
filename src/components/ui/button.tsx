import { ButtonHTMLAttributes, ReactNode } from "react";

type Variant = "primary" | "secondary" | "outline" | "disabled";
type Size = "sm" | "md" | "lg" | "xl" | "2xl";

type Props = {
  children: ReactNode;
  variant?: Variant;
  size?: Size;
  className?: string;
} & ButtonHTMLAttributes<HTMLButtonElement>;

export default function Button({
  children,
  variant = "primary",
  size = "md",
  className = "",
  ...props
}: Props) {
  const base =
    "rounded-lg font-medium flex items-center justify-center transition-colors duration-200";

  const variants = {
    primary: "bg-[var(--foreground)] text-[var(--background)] hover:opacity-90",

    secondary: "bg-gray-700 text-white hover:bg-gray-800",

    outline:
      "bg-transparent border border-gray-400 text-gray-700 hover:bg-gray-100 hover:border-gray-500",

    disabled: "bg-gray-200 text-gray-400 cursor-not-allowed",
  };

  const sizes = {
    sm: "px-4 py-2 text-sm",
    md: "px-6 py-2.5 text-sm",
    lg: "px-8 py-3 text-base",
    xl: "px-10 py-4 text-lg",
    "2xl": "px-12 py-5 text-xl",
  };

  return (
    <button
      className={`${base} ${variants[variant]} ${sizes[size]} ${className}`}
      disabled={variant === "disabled"}
      {...props}
    >
      {children}
    </button>
  );
}
