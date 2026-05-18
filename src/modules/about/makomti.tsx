import Container from "@/components/ui/container";
import Image from "next/image";

const Makomti = () => {
  return (
    <section className="relative overflow-hidden bg-linear-to-b from-primary-600 to-primary-800 text-white text-center pt-[152px] pb-95 md:pb-24">
      <div className="absolute z-0 w-50 h-auto md:w-65 shadow-xl overflow-hidden md:rotate-5 -rotate-12 left-10 bottom-40 md:top-24 md:bottom-auto md:left-10 aspect-4/3">
        <Image src="/balairung.webp" alt="Background Element" fill className="object-cover" />
      </div>

      <div className="absolute z-0 w-50 h-auto md:w-65 shadow-xl overflow-hidden md:-rotate-5 rotate-12 right-10 bottom-15 md:top-24 md:bottom-auto md:right-10 aspect-4/3">
        <Image src="/dike.webp" alt="Background Element" fill className="object-cover" />
      </div>

      <Container className="relative z-10 flex flex-col items-center justify-center gap-7">
        <h2 className="text-[22px] md:text-[38px] font-bold text-primary-200">Ruang <span className="text-white"> Berkembang </span>Mahasiswa</h2>
        <div className="grid grid-cols-1 gap-2 md:gap-7 md:grid-cols-2 items-center justify-center w-fit">
          <div className="bg-radial from-primary-800 to-primary-900 max-w-[505px] rounded-2xl h-full px-5 py-10 md:p-11 gap-5 md:gap-8 flex flex-col items-center justify-center shadow-lg">
            <Image src="/logo/omahti-full.png" alt="Logo OmahTI" width={1000} height={70} className="h-13 md:h-[70px] w-auto object-contain"/>
            <p className="text-sm md:text-base">OmahTI merupakan organisasi kemahasiswaan yang berfokus pada peningkatan hard skills di bidang Ilmu Komputer melalui pembelajaran berbasis proyek dan pengembangan portofolio.</p>
          </div>
          <div className="bg-radial from-primary-800 to-primary-900 max-w-[505px] rounded-2xl p-5 pb-10 md:p-11 md:pt-[22px] gap-3 flex flex-col items-center justify-center shadow-lg">
            <Image src="/logo/himakom.png" alt="Logo Himakom" width={116} height={116} className="md:h-[116px] h-21 w-auto object-contain"/>
            <p className="text-sm md:text-base">Himakom adalah organisasi kemahasiswaan yang mewadahi representasi mahasiswa serta mengelola hubungan internal & eksternal melalui program akademik serta penguatan solidaritas.</p>
          </div>
        </div>
      </Container>
    </section>
  )
}

export default Makomti;