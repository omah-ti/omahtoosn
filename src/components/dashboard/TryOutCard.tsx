import Button from "@/components/ui/button";

export default function TryOutCard() {
  return (
    <div className="bg-gray-100 rounded-xl p-6 flex justify-between items-center min-h-[260px]">
      <div className="max-w-md">
        <h2 className="text-3xl font-bold mb-2 text-[#0a0a0a]">Coba Try Out</h2>

        <p className="mb-20 text-[#0a0a0a]">
          Uji kemampuanmu melalui simulasi berbasis silabus OSN terbaru dengan
          sistem real-time
        </p>

        <Button variant="outline">Sudah Mengerjakan</Button>
      </div>

      <div className="w-40 h-32 bg-gray-300 rounded-lg" />
    </div>
  );
}
