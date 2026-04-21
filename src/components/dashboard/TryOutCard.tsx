import Button from "@/components/ui/button";

type Variant = "before" | "after";

type Props = {
  variant?: Variant;
};

export default function TryOutCard({ variant = "before" }: Props) {
  const isAfter = variant === "after";

  return (
    <div
      className={`rounded-xl p-6 flex justify-between items-center min-h-[260px] relative overflow-hidden
      ${
        isAfter
          ? "bg-[#DCE5F9]"
          : "bg-gradient-to-b from-[var(--primary-600)] to-[var(--primary-1000)]"
      }`}
    >
      {/* TEXT */}
      <div className="max-w-md pr-28 sm:pr-0">
        <h2
          className={`text-3xl font-bold mb-2
          ${isAfter ? "text-black" : "text-[var(--neutral-100)]"}`}
        >
          {isAfter ? "Coba Try Out" : "Coba Try Out"}
        </h2>

        <p
          className={`mb-20
          ${isAfter ? "text-[var(--neutral-1000)]" : "text-[var(--neutral-100)]"}`}
        >
          {isAfter
            ? "Uji kemampuanmu melalui simulasi berbasis silabus OSN terbaru dengan sistem real-time"
            : "Uji kemampuanmu melalui simulasi berbasis silabus OSN terbaru dengan sistem real-time"}
        </p>

        <Button
          variant={isAfter ? "disabled" : "primary"}
          className="w-full max-w-[180px] mt-7 sm:w-1/2"
        >
          {isAfter ? "Sudah Mengerjakan" : "Mulai Sekarang!"}
        </Button>
      </div>

      {/* IMAGE */}
      <div className="absolute bottom-0 right-0 w-[clamp(280px,28vw,380px)]">
        <img
          src={isAfter ? "/tryout/round.webp" : "/tryout/roundblack.webp"}
          alt="background"
          className="z-[0] absolute top-6 left-10 w-[clamp(320px,38vw,460px)] object-contain"
        />

        <img
          src="/tryout/anaksma.webp"
          alt="student"
          className="relative w-full object-contain z-10"
        />
      </div>
    </div>
  );
}
