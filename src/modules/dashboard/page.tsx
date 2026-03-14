import TryOutCard from "@/components/dashboard/TryOutCard";
import ResultCard from "@/components/dashboard/ResultCard";
import MaterialCard from "@/components/dashboard/MaterialCard";

export default function DashboardPage() {
  return (
    <div className="w-full min-h-screen bg-white">
      {/* container */}
      <div className="max-w-6xl mx-auto px-8 py-8">
        {/* Greeting */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-[#0a0a0a] mb-2">Hi, User!</h1>

          <p className="text-[#0a0a0a]">
            Setiap latihan membawamu selangkah lebih dekat ke kompetisi
            sebenarnya
          </p>
        </div>

        {/* Top Cards */}
        <div className="grid grid-cols-3 gap-6 mb-5 items-stretch">
          {/* TryOut Card */}
          <div className="col-span-2 h-full">
            <TryOutCard />
          </div>

          {/* Result Card*/}
          <div className="h-full">
            <ResultCard />
          </div>
        </div>

        {/* Materi Section */}
        <div className="bg-gray-200 rounded-xl p-6">
          <h2 className="text-xl font-semibold mb-1 text-[#0a0a0a]">Materi</h2>

          <p className="text-sm text-[#0a0a0a] mb-6">
            Materi try out disusun berdasarkan silabus TOKI untuk OSN-K
            Informatika
          </p>

          {/* Materi Cards */}
          <div className="grid grid-cols-3 gap-6 text-[#0a0a0a]">
            <MaterialCard
              title="Abstraksi Komputasional"
              desc="Soal berbasis cerita untuk menguji konsep dasar informatika"
            />

            <MaterialCard
              title="Pemecahan Masalah"
              desc="Studi kasus komputasional yang diselesaikan melalui penalaran logis"
            />

            <MaterialCard
              title="Algoritma dalam C++"
              desc="Analisis potongan kode C++ untuk memahami logika dan hasil program"
            />
          </div>
        </div>
      </div>
    </div>
  );
}
