import TryOutCard from "@/components/dashboard/tryout-card";
import ResultCard from "@/components/dashboard/result-card";
import MaterialCard from "@/components/dashboard/material-card";

async function getUser() {
  const res = await fetch("http://localhost:3000/api/me", {
    cache: "no-store",
  });

  if (!res.ok) {
    return { name: "User" };
  }

  return res.json();
}

export default async function DashboardPage() {
  const user = await getUser();

  const testStatus = "NOT_STARTED"; // "NOT_STARTED" | "COMPLETED"
  const variant = testStatus === "COMPLETED" ? "after" : "before";

  return (
    <div className="w-full min-h-screen pt-5 bg-white">
      <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8">
        <div className="mb-6 sm:mb-8">
          <h1 className="text-2xl sm:text-3xl font-bold text-[#0a0a0a] mb-3">
            Hi, {user.name || "User"}!
          </h1>

          <p className="text-sm mb-[-15] sm:text-base text-[#0a0a0a]">
            Setiap latihan membawamu selangkah lebih dekat ke kompetisi
            sebenarnya
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-4 sm:gap-6 mb-6">
          <div className="lg:col-span-2">
            <TryOutCard variant={variant} />
          </div>

          <div>
            <ResultCard variant={variant} />
          </div>
        </div>

        <div className="bg-primary-background rounded-xl p-4 sm:p-6">
          <h2 className="text-lg sm:text-xl font-semibold mb-1 text-neutral-1000">
            Materi
          </h2>

          <p className="text-xs sm:text-sm text-neutral-1000 mb-4 sm:mb-6">
            Materi try out disusun berdasarkan silabus TOKI untuk OSN-K
            Informatika
          </p>

          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 text-[#0a0a0a] gap-4 sm:gap-6">
            <MaterialCard
              img="/materi/materi1.webp"
              title="Abstraksi Komputasional"
              desc="Soal berbasis cerita untuk menguji konsep dasar informatika dan pola berpikir komputasional"
            />
            <MaterialCard
              img="/materi/materi2.webp"
              title="Pemecahan Masalah"
              desc="Studi kasus komputasional yang diselesaikan melalui penalaran logis"
            />
            <MaterialCard
              img="/materi/materi3.webp"
              title="Algoritma C++"
              desc="Analisis potongan kode C++ untuk memahami logika dan hasil program"
            />
          </div>
        </div>
      </div>
    </div>
  );
}
