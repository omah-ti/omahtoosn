import Container from "@/components/ui/container";
import { Trophy, BookOpen } from "lucide-react";

const Kemampuan = () => {
  return (
    <section className="relative w-full bg-primary-200 overflow-hidden md:pt-13 md:pb-25 py-20">
      <Container>
        <h2 className="text-[22px] md:text-[38px] font-bold text-center">
          <span className="text-primary-1000">Sudah</span>{' '}
          <span className="text-primary-600">Sejauh</span>{' '}
          <span className="text-primary-1000">Mana</span>{' '}
          <span className="text-primary-600 italic">Kemampuanmu?</span>
        </h2>

        <div className="flex flex-col md:flex-row justify-center gap-5 md:gap-8 w-full">
          <div className="flex-1 w-full max-w-114 bg-primary-800 rounded-[20px] gap-6 p-5 relative overflow-hidden text-white flex flex-col justify-between min-h-[204px] md:min-h-[218px]">
            <div className="space-y-2">
              <h3 className="text-xl md:text-[26px] font-semibold">Ranking</h3>
              <div className="w-full h-px bg-white" />
              <p className="text-sm md:text-lg text-white/90 leading-relaxed relative z-10 max-w-[85%] md:max-w-[75%]">
                Lihat sejauh mana hasil pekerjaanmu dan bagaimana posisimu
              </p>
            </div>
            <button className="bg-primary-200 text-primary-900 font-semibold px-6 py-2 md:py-3 rounded-lg w-max relative z-10 text-sm md:text-base hover:bg-white transition-colors">
              Lihat Rangking
            </button>
            <div className="absolute -bottom-20 -right-6 md:-bottom-25 md:-right-8 text-primary-100 pointer-events-none rotate-12 z-0">
              <Trophy className="w-40 h-40 md:w-56 md:h-56" strokeWidth={1.5} />
            </div>
          </div>

          <div className="flex-1 w-full max-w-114 bg-primary-800 rounded-[20px] gap-6 p-5 relative overflow-hidden text-white flex flex-col justify-between min-h-[204px] md:min-h-[218px]">
            <div className="space-y-2">
              <h3 className="text-xl md:text-[26px] font-semibold">Pembahasan</h3>
              <div className="w-full h-px bg-white" />
              <p className="text-sm md:text-lg text-white/90 leading-relaxed relative z-10 max-w-[85%] md:max-w-[75%]">
                Masih penasaran dengan jawaban dari soal yang kamu kerjakan?
              </p>
            </div>
            <button className="bg-primary-200 text-primary-900 font-semibold px-6 py-2 md:py-3 rounded-lg w-max relative z-10 text-sm md:text-base hover:bg-white transition-colors">
              Daftar Pembahasan
            </button>
            <div className="absolute -bottom-20 -right-6 md:-bottom-25 md:-right-8 text-primary-100 pointer-events-none rotate-12 z-0">
              <BookOpen className="w-40 h-40 md:w-56 md:h-56" strokeWidth={1.5} />
            </div>
          </div>
        </div>
      </Container>
    </section>
  );
}

export default Kemampuan;