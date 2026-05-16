import Link from "next/link"
import Container from "./container"
import Logo from "./logo"
import Image from "next/image"
import { Phone } from "lucide-react"

const Footer = () => {
  return (
    <section className="bg-primary-1000 md:py-30 text-white text-[10px] md:text-lg pb-20 pt-11">
      <Container>
        <div className="flex justify-between md:items-stretch">
          <div className="flex flex-col max-w-105 md:w-full w-[50%] justify-between space-y-4">
            <Logo />
            <p>Gedung Fakultas MIPA UGM Sekip Utara,Bulaksumur, Sinduadi, Mlati, Sleman, DI Yogyakarta</p>
            <p>2026 Copyright - OmahTI</p>
          </div>

          <div className="flex flex-col max-w-105 items-end justify-end md:gap-10 w-[50%] h-full gap-3 pt-3 md:pt-0">
            <p className="text-right w-full">OmahTI - Organisasi Mahasiswa Ahli Teknologi Informasi. Berkomitmen membantu calon mahasiswa dalam persiapan OSN Informatika. Mari wujudkan masa depanmu bersama kami!</p>
            <div className="flex flex-row gap-3">
              <Link href="https://www.instagram.com/omahti_ugm">
                <Image src="/logo/instagram.svg" alt="Instagram" width={30} height={30} className="bg-[#D9D9D9] p-1 rounded-md" />
              </Link>
              <Link href="https://www.instagram.com/omahti_ugm">
                <Image src="/logo/linkedin.svg" alt="Instagram" width={30} height={30} className="bg-[#D9D9D9] p-1 rounded-md" />
              </Link>
              <Link href="https://www.instagram.com/omahti_ugm">
                <Phone className="text-black size-[30px] bg-[#D9D9D9] p-1 rounded-md" />
              </Link>
            </div>
          </div>
        </div>
      </Container>
    </section>
  )
}

export default Footer