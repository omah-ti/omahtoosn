'use client';
import { useRef } from 'react';
import { ChevronLeft, ChevronRight } from 'lucide-react';
import Container from "@/components/ui/container";
import Image from "next/image";

const Alumni = () => {
  const scrollRef = useRef<HTMLDivElement>(null);

  const scroll = (direction: 'left' | 'right') => {
    if (scrollRef.current) {
      const { current } = scrollRef;
      const scrollAmount = direction === 'left' ? -current.offsetWidth : current.offsetWidth;
      current.scrollBy({ left: scrollAmount, behavior: 'smooth' });
    }
  };

  const alumniData = [
    {
      text: "Saya OSN Nasional sih biasa lah masuk sini ga kaget ok aja si isinya radak2 aneh aja tapi boleh lah buat kuliah jisjeindjbeb bxhjebdhbhb ubdhbehcbbhcbhcbh",
      name: "Ayasha Rahmadinni",
      track: "UTUL wkwk",
      image: "/hero-character.webp"
    },
    {
      text: "Saya OSN Nasional sih biasa lah masuk sini ga kaget ok aja si isinya radak2 aneh aja tapi boleh lah buat kuliah jisjeindjbeb bxhjebdhbhb ubdhbehcbbhcbhcbh",
      name: "Maemunah",
      track: "SNBP",
      image: "/hero-character.webp"
    },
    {
      text: "Saya OSN Nasional sih biasa lah masuk sini ga kaget ok aja si isinya radak2 aneh aja tapi boleh lah buat kuliah jisjeindjbeb bxhjebdhbhb ubdhbehcbbhcbhcbh",
      name: "Bimo",
      track: "SNBP",
      image: "/hero-character.webp"
    },
  ];

  return (
    <section className="relative w-full bg-white overflow-hidden py-20 md:py-32">
      <Container className="relative z-10 md:items-center">
        <h2 className="md:text-[38px] text-[22px] font-bold text-center">
          <span className="text-primary-1000">Jejak</span>{' '}
          <span className="text-primary-700">Alumni OSN</span>{' '}
          <span className="text-primary-1000">Informatika</span>
        </h2>

        <div className="relative group">
          {/* <button
            onClick={() => scroll('left')}
            className="absolute -left-4 md:-left-6 top-1/2 -translate-y-1/2 z-20 bg-primary-200 hover:bg-primary-background p-2 rounded-lg cursor-pointer text-white backdrop-blur-sm transition-all hidden md:block shadow-md"
            aria-label="Scroll left"
          >
            <ChevronLeft className="w-6 h-6" />
          </button>

          <button
            onClick={() => scroll('right')}
            className="absolute -right-4 md:-right-6 top-1/2 -translate-y-1/2 z-20 bg-primary-200 hover:bg-primary-background p-2 rounded-lg cursor-pointer text-white backdrop-blur-sm transition-all hidden md:block shadow-md"
            aria-label="Scroll right"
          >
            <ChevronRight className="w-6 h-6" />
          </button> */}

          <div
            ref={scrollRef}
            className="flex gap-4 md:gap-6 overflow-x-auto snap-x snap-mandatory [&::-webkit-scrollbar]:hidden [-ms-overflow-style:none] [scrollbar-width:none]"
          >
            {alumniData.map((alumni, idx) => {
              const isDark = idx % 2 === 0;
              return (
                <div
                  key={idx}
                  className={`w-full shrink-0 md:w-[350px] lg:w-[400px] rounded-[20px] p-5 flex flex-col justify-between snap-center md:snap-start shadow-sm transition-colors duration-300 min-h-[250px] ${
                    isDark ? 'bg-primary-800 text-white' : 'bg-primary-200 text-primary-1000'
                  }`}
                >
                  <p className="text-sm md:text-base leading-relaxed">
                    {alumni.text}
                  </p>
                  
                  <div>
                    <div className={`w-full h-px mb-5 ${isDark ? 'bg-white/30' : 'bg-primary-1000/20'}`} />
                    <div className="flex items-center gap-4">
                      <div className="w-12 h-12 rounded-full shrink-0 relative overflow-hidden bg-primary-100">
                        <Image src={alumni.image} alt={alumni.name} fill className="object-cover" />
                      </div>
                      <div className="flex flex-col">
                        <span className="font-semibold">{alumni.name}</span>
                        <span className={`text-xs ${isDark ? 'text-white/80' : 'text-primary-1000/80'}`}>{alumni.track}</span>
                      </div>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        </div>

        <div className="flex justify-between items-center md:hidden">
          <button
            onClick={() => scroll('left')}
            className="bg-primary-700 hover:bg-primary-700 p-[10px] rounded-lg cursor-pointer text-white transition-all shadow-md justify-center"
            aria-label="Scroll left"
          >
            <ChevronLeft className="w-6 h-6" />
          </button>

          <button
            onClick={() => scroll('right')}
            className="bg-primary-700 hover:bg-primary-700 p-[10px] rounded-lg text-white transition-all shadow-md justify-center"
            aria-label="Scroll right"
          >
            <ChevronRight className="w-6 h-6" />
          </button>
        </div>
      </Container>
    </section>
  );
}

export default Alumni;