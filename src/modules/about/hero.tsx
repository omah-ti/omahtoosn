import Container from "@/components/ui/container";
import Image from "next/image";

const Hero = () => {
  return (
    <section className="bg-primary-1000 text-white">
      <Image
        src="/balairung-about.webp"
        alt=""
        width={1920}
        height={1080}
        className="w-full h-full object-cover min-h-123 md:min-h-none"
      />
      <Container className="flex md:flex-row flex-col gap-9 justify-between w-full md:py-15 py-5">
        <div className="flex flex-col gap-1 max-w-160">
          <h2 className="md:text-[38px] text-[32px]">Mengakar Kuat, <span className="text-primary-200 font-bold italic"> Menjulang Tinggi</span></h2>
          <p className="md:text-lg">
            Menjadi pusat keunggulan pendidikan di Indonesia dengan pengakuan
            prestasi di kancah global.
          </p>
        </div>
        <div className="grid grid-cols-3 text-center">
          <div className="flex flex-col gap-1 items-center">
            <h3 className="md:text-[38px] text-[32px] font-bold">
              #2
            </h3>
            <p className="md:text-lg text-sm whitespace-nowrap">Kampus Indonesia</p>
          </div>
          <div className="flex flex-col gap-1 items-center">
            <h3 className="md:text-[38px] text-[32px] font-bold">
              #224
            </h3>
            <p className="md:text-lg text-sm whitespace-nowrap">Kampus Global</p>
          </div>
          <div className="flex flex-col gap-1 items-center">
            <h3 className="md:text-[38px] text-[32px] font-bold">
              #3
            </h3>
            <p className="md:text-lg text-sm whitespace-nowrap">Jurusan IT</p>
          </div>
        </div>
      </Container>
    </section>
  )
}

export default Hero;