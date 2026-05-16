import Button from "@/components/ui/button";

type Variant = "before" | "after";

type Props = {
  variant?: Variant;
};

export default function ResultCard({ variant = "before" }: Props) {
  const isAfter = variant === "after";

  return (
    <div
      className={`relative rounded-xl p-6 h-full overflow-hidden
      ${isAfter ? "bg-linear-to-b from-primary-600 to-primary-1000" : "bg-primary-background"}`}
    >
      <div className="relative z-10">
        <h2
          className={`text-2xl font-semibold mb-2
          ${isAfter ? "text-neutral-100" : "text-neutral-1000"}`}
        >
          {isAfter ? "Hasil Tes" : "Hasil Tes"}
        </h2>

        <p
          className={`text-md mb-6
          ${isAfter ? "text-neutral-100" : "text-neutral-1000"}`}
        >
          {isAfter
            ? "Ketahui peringkatmu dan pahami solusi dari setiap soal"
            : "Ketahui peringkatmu dan pahami solusi dari setiap soal"}
        </p>

        <Button
          variant={isAfter ? "primary" : "primary"}
          className="w-full mt-23 max-w-[180px]"
        >
          {isAfter ? "Lihat Hasil" : "Lihat Hasil"}
        </Button>
      </div>

      <img
        src={isAfter ? "/result/roundresult.webp" : "/result/roundresult.webp"}
        alt="decoration"
        className={`absolute bottom-[-150px] right-[-80px] w-[300px] pointer-events-none
        ${isAfter ? "opacity-60" : "opacity-60"}`}
      />
    </div>
  );
}
