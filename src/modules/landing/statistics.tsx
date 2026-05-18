import Container from "@/components/ui/container";
import Image from "next/image";

const Statistics = () => {
  return (
    <Container className="pb-9 md:pb-25">
      <div className="w-full flex flex-row py-23 items-center justify-center md:gap-x-10 gap-x-5">
        <div className="flex flex-col items-center w-full max-w-89 gap-1">
          <h2 className="md:text-[68px] text-[26px] text-primary-500 font-bold">75%</h2>
          <h1 className="md:text-lg text-xs text-center text-black">Fresh graduate  mendapat pekerjaan 3 bulan setelah lulus</h1>
        </div>
        <div className="h-20 md:h-35 w-px bg-black" />
        <div className="flex flex-col items-center w-full max-w-89 gap-1">
          <h2 className="md:text-[68px] text-[26px] text-primary-900 font-bold">&gt;10 Jt</h2>
          <h1 className="md:text-lg text-xs text-center text-black">Jalur tercepat ke gaji tinggi awal karir dengan peluang global</h1>
        </div>
        <div className="h-20 md:h-35 w-px bg-black" />
        <div className="flex flex-col items-center w-full max-w-89 gap-1">
          <h2 className="md:text-[68px] text-[26px] text-primary-500 font-bold">35%</h2>
          <h1 className="md:text-lg text-xs text-center text-black">Lowongan IT di Indonesia tumbuh lebih cepat daripada sektor lain </h1>
        </div>
      </div>
      <Container className="relative bg-linear-to-b from-primary-600 to-primary-1000 rounded-[23px] lg:h-105 flex flex-col lg:flex-row items-center lg:justify-between py-20">
        <div className="text-white flex flex-col max-w-140 gap-90 lg:gap-4 text-center lg:text-left">
          <h1 className="text-[22px] font-bold lg:text-[38px]">Selangkah <span className="lg:font-bold lg:italic lg:text-primary-200">Lebih Unggul</span></h1>
          <p className="text-sm lg:text-lg">OSN Informatika adalah <span className="font-bold italic">kesempatan emas</span> untuk membangun prestasi dan memperbesar peluang masuk PTN. OSN bisa menjadi langkah strategis yang membuatmu selangkah lebih unggul dari yang lain.  Buktikan kemampuanmu, amankan keunggulanmu, dan buka jalan menuju kampus impian.</p>
        </div>
        <div className="relative w-full lg:w-132 h-full lg:h-105">
          <Image src="/statistics-characters.webp" alt="statistics-character" width={800} height={1200} className="absolute lg:-bottom-10 w-full h-72 lg:h-full object-contain z-1 bottom-35" />
          <Image src="/statistics-circles.webp" alt="landing-circle" width={800} height={1200} className="absolute w-lg h-70 lg:h-auto object-contain right-0 lg:-top-10 bottom-50 rotate-340" />
          {/* <div className="w-110 h-69 rounded-full top-0 bg-radial from-primary-400 to-transparent" /> */}
        </div>
      </Container>
    </Container>
  );
}

export default Statistics;