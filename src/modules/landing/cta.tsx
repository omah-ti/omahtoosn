import Container from "@/components/ui/container";
import Image from "next/image";
import Button from "@/components/ui/button";

const CTA = () => {
  return (
    <section className="w-full bg-white py-16 md:py-24">
      <Container>
        <Container>
          <Container>
            <div className="flex flex-col md:flex-row items-center justify-between gap-6 md:gap-10">
              <div className="flex flex-col items-center md:items-start text-center md:text-left w-full md:w-1/2">
                <h2 className="text-[24px] md:text-[46px] font-bold leading-tight w-full">
                  <span className="text-primary-1000">Mulai</span> <span className="text-primary-700">Langkahmu</span>
                  <span className="hidden md:inline"><br /></span>
                  <span className="md:hidden"> </span>
                  <span className="text-primary-1000">Menuju</span> <span className="text-primary-700">Medali!</span>
                </h2>
                
                <div className="flex justify-center w-full md:hidden">
                  <div className="relative w-[280px] h-[280px]">
                    <Image src="/osn-cta.webp" alt="OSN CTA" fill className="object-contain" />
                  </div>
                </div>

                <p className="text-sm md:text-lg text-primary-1000 max-w-[320px] md:max-w-[450px] mt-0 md:mt-4">
                  Simulasi presisi yang dirancang khusus untuk calon peraih medali OSN Informatika
                </p>

                <div className="w-full flex justify-center md:justify-start mt-6 md:mt-8">
                  <Button className="w-full md:w-fit">
                    Coba Try Out
                  </Button>
                </div>
              </div>
              
              <div className="hidden md:flex justify-center w-full md:w-1/2">
                <div className="relative w-[400px] h-[400px] lg:w-[500px] lg:h-[500px]">
                  <Image src="/osn-cta.webp" alt="OSN CTA" fill className="object-contain" />
                </div>
              </div>
            </div>
          </Container>
        </Container>
      </Container>
    </section>
  );
}

export default CTA;