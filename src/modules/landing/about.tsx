import Image from "next/image";
import Container from "@/components/ui/container";
import Button from "@/components/ui/button";
import { Code } from "lucide-react";

const About = () => {
  return (
    <section className="relative w-full overflow-hidden bg-linear-to-b from-primary-600 to-primary-1000 py-20 md:py-32">
      <Container>
        <div className="flex md:flex-row flex-col-reverse gap-10 md:gap-16">
          <Image
            src="/statistics-circles.webp"
            alt="circles background"
            width={800}
            height={1200}
            className="absolute w-[500px] md:w-[600px] h-auto object-contain -top-0 -right-20 pointer-events-none opacity-50"
          />

          <div className="absolute top-[10%] right-[15%] w-4 h-4 rounded-full bg-white/40 hidden md:block" />
          <div className="absolute top-[15%] right-[25%] w-2 h-2 rounded-full bg-white/60 hidden md:block" />

          <div className="flex-1 flex flex-col items-start gap-5 w-full text-white relative z-10">
            <h2 className="text-3xl md:text-[42px] font-bold leading-tight">
              <span className="text-white">Menuju</span> <span className="text-primary-300">Ilmu Komputer</span> <span className="text-white">UGM</span>
            </h2>
            <p className="text-sm md:text-lg text-white/80 leading-relaxed max-w-[95%] md:max-w-full">
              Siapkan dirimu menembus salah satu program studi IT paling bergengsi di Indonesia. Temukan kurikulum berstandar internasional, fasilitas riset mutakhir, dan ekosistem yang melahirkan talenta digital masa depan.
            </p>
            <a href="https://dcse.fmipa.ugm.ac.id/" target="_blank" rel="noopener noreferrer" className="mt-4 md:mt-2 w-full md:w-auto">
              <Button className="w-full md:w-auto bg-primary-600 hover:bg-primary-700 text-white border-none">
                Buka Website Resmi
              </Button>
            </a>
          </div>

          <div className="flex-1 w-full relative z-10 overflow-y-visible">
            <div className="w-full h-40 md:h-auto md:aspect-16/10 rounded-[20px]  shadow-2xl">
              <Image
                src="/balairung.webp"
                alt="Balairung UGM"
                fill
                className="object-cover rounded-xl"
              />
            </div>

          </div>
          <div className="absolute bottom-25 right-160  bg-white p-3 md:p-4 hidden md:block rounded-2xl rotate-[-15deg] shadow-xl z-20">
            <Code className="w-6 h-6 md:w-8 md:h-8 text-primary-1000" />
          </div>
        </div>
      </Container>
    </section>
  );
};

export default About;
