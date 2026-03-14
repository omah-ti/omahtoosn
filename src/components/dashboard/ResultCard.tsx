import Button from "@/components/ui/button";

export default function ResultCard() {
  return (
    <div className="bg-[#1f2a3a] text-white rounded-xl p-6 flex flex-col justify-between h-full">
      <div>
        <h2 className="text-2xl font-semibold mb-2">Hasil Tes</h2>

        <p className="text-gray-300 text-sm">
          Ketahui peringkatmu dan pahami solusi dari setiap soal
        </p>
      </div>

      <Button variant="primary" className="w-full md:w-1/2">
        Lihat Hasil
      </Button>
    </div>
  );
}
