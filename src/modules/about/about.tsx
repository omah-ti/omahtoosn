import Container from "@/components/ui/container";
import Image from "next/image";
import Link from "next/link";
import Button from "@/components/ui/button";

const About = () => {
  return (
    <section className="bg-linear-to-b from-primary-1000 to-primary-600 relative pt-20 md:pb-49 pb-20 text-white">
      <Container className="gap-5 md:gap-3 flex flex-col">
        <h2 className="font-bold md:text-[38px] text-[26px] text-primary-200 text-center">Ilmu Komputer UGM <span className="font-normal text-white">untuk</span><span className="text-white"> Masa Depan Digital</span></h2>
        <div className="flex md:flex-row flex-col gap-5 md:gap-7">
          <Image src="/dike.webp" alt="" width={200} height={180} className="rounded-[20px] md:rounded-none h-61 md:h-98 w-auto md:w-115 object-cover" />
          <div className="md:items-end justify-between flex flex-col gap-5">
            <div className="space-y-3 text-xs md:text-base">
              <p>Program Studi Ilmu Komputer UGM menghadirkan lingkungan akademik yang adaptif terhadap perkembangan ilmu pengetahuan dan teknologi yang dinamis. Kurikulum dirancang secara berkelanjutan untuk menyesuaikan diri dengan kemajuan terbaru di bidang komputasi, sehingga mahasiswa memperoleh kompetensi yang relevan dan aplikatif.</p>
              <p>Di samping itu, suasana akademik yang kompetitif secara sehat mendorong mahasiswa untuk senantiasa meningkatkan kualitas diri. Diskusi ilmiah, kegiatan penelitian, serta partisipasi dalam kompetisi dan proyek kolaboratif membentuk karakter yang tangguh, kritis, dan berorientasi pada pencapaian.</p>
              <p>Lingkungan ini tidak hanya menekankan pada pencapaian akademik, tetapi juga pada pengembangan profesionalisme, integritas, dan kemampuan bekerja sama, sehingga lulusan siap berkontribusi secara signifikan di tingkat nasional maupun global.</p>
            </div>
            <Link href="https://dcse.cs.ugm.ac.id/" className="justify-end">
              <Button variant="transparent" className="bg-primary-100 text-black">
                Buka Website Resmi
              </Button>
            </Link>
          </div>
        </div>
      </Container>
      <div className="bg-white w-full h-9 md:h-18 bottom-0 absolute [clip-path:polygon(50%_0%,0%_100%,100%_100%)]" />
    </section>
  )
}

export default About;