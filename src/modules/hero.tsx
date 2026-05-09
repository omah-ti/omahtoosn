import Image from "next/image";
import Link from "next/link";
import Button from "@/components/ui/button";

const Hero = () => {
  return (
    <section className="hero h-screen relative flex flex-col lg:flex-row items-center justify-between px-6 md:px-20 overflow-hidden">
      <div className="max-w-2xl pb-15 flex flex-col lg:justify-center justify-between h-full lg:gap-8 text-white pt-32 lg:pt-0 z-10 w-full lg:w-1/2 gap-2">
        <div>
          <h1 className="font-bold text-4xl md:text-[46px] text-white leading-tight">Asah Logika, Taklukkan Algoritma, Menangkan <span className="italic text-primary-400">OSN Informatika</span></h1>
          <p className="text-lg mt-2">Di mana kemampuan berpikir kritis dan logika akan diuji.</p>
        </div>
        <div className="flex flex-col lg:flex-row gap-3">
          <Link href="/login">
            <Button className="w-full lg:w-auto">
              Gabung Sekarang
            </Button>
          </Link>
          <Link href="/about">
            <Button variant="outline" className="w-full lg:w-auto">
              Tentang Ilmu Komputer
            </Button>
          </Link>
        </div>
      </div>
      <div className="absolute md:relative flex justify-center md:justify-end w-full h-full lg:w-2/3 mt-0 bottom-0 z-5">
        <Image src="/hero-character.webp" alt="hero-character" width={790} height={1200} className="w-full h-auto object-contain absolute md:-bottom-10 lg:scale-100 bottom-30 scale-90 lg:bottom-0 md:-right-10" />
      </div>

      <div className="absolute w-236 h-236 opacity-70 md:-right-15 md:-bottom-85 -bottom-120">
        <div className="bg-primary-1000/54 w-236 h-236 rounded-full absolute inset-0"/>
        <div className="bg-primary-1000/52 w-210 h-210 rounded-full absolute inset-13"/>
        <div className="bg-primary-1000 w-184 h-184 rounded-full absolute inset-26"/>
      </div>
    </section>
  );
};
export default Hero;