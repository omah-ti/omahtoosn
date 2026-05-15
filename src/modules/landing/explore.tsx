'use client';
import { useRef } from 'react';
import { Layers, BarChart3, PenTool, Globe, Code, Cpu, ChevronLeft, ChevronRight, Airplay, Podcast } from 'lucide-react';
import Container from "@/components/ui/container";

const Explore = () => {
  const scrollRef = useRef<HTMLDivElement>(null);

  const scroll = (direction: 'left' | 'right') => {
    if (scrollRef.current) {
      const { current } = scrollRef;
      const scrollAmount = direction === 'left' ? -current.offsetWidth : current.offsetWidth;
      current.scrollBy({ left: scrollAmount, behavior: 'smooth' });
    }
  };

  const cards = [
    {
      title: 'Front End',
      desc: 'Membangun tampilan website yang interaktif, responsif, dan menarik',
      Icon: Airplay,
    },
    {
      title: 'Back End',
      desc: 'Mengelola server, database, dan logika sistem agar stabil.',
      Icon: Layers,
    },
    {
      title: 'Data Science',
      desc: 'Menggabungkan statistik dan machine learning untuk pengambilan keputusan.',
      Icon: BarChart3,
    },
    {
      title: 'UI/UX',
      desc: 'Merancang pengalaman pengguna dengan mengutamakan kemudahan dan estetika.',
      Icon: PenTool,
    },
  ];

  return (
    <section className="relative w-full bg-white overflow-hidden py-20 md:py-50">
      <div className="absolute bottom-0 left-1/2 -translate-x-1/2 translate-y-[30%] md:translate-y-[65%] w-[730.768px] h-[639.077px] md:w-[2042px] md:h-[1785.787px] bg-primary-200/25 pointer-events-none rounded-[50%]" />
      <div className="absolute bottom-0 left-1/2 -translate-x-1/2 translate-y-[30%] md:translate-y-[65%] w-[648.503px] h-[567.134px] md:w-[1812.126px] md:h-[1584.755px] bg-primary-200/50 pointer-events-none rounded-[50%]" />
      <div className="absolute bottom-0 left-1/2 -translate-x-1/2 translate-y-[30%] md:translate-y-[65%] w-[567.307px] h-[496.126px] md:w-[1585.237px] md:h-[1386.335px] bg-primary-200 pointer-events-none rounded-[50%]" />

      <div className="absolute top-13 md:top-12 left-[17%] md:left-[10%] p-2 md:p-4 bg-primary-100 rounded-2xl rotate-[-10deg] shadow-sm">
        <Podcast className="w-6 h-6 md:w-8 md:h-8 text-primary-900" />
      </div>
      <div className="absolute top-5 md:top-32 left-[5%] md:left-[20%] p-2 md:p-4 bg-primary-100 rounded-2xl rotate-15 shadow-sm">
        <Code className="w-6 h-6 md:w-8 md:h-8 text-primary-900" />
      </div>
      <div className="absolute top-12 md:top-20 right-[20%] md:right-[20%] p-2 md:p-4 bg-primary-100 rounded-2xl rotate-[-5deg] shadow-sm">
        <Cpu className="w-6 h-6 md:w-8 md:h-8 text-primary-900" />
      </div>
      <div className="absolute top-4 md:top-10 right-[5%] md:right-[10%] p-2 md:p-4 bg-primary-100 rounded-2xl rotate-10 shadow-sm">
        <Globe className="w-6 h-6 md:w-8 md:h-8 text-primary-900" />
      </div>

      <Container className="relative z-10">
        <h2 className="md:text-[38px] text-[22px] font-bold text-center">
          <span className="text-primary-800">Eksplorasi</span>{' '}
          <span className="text-neutral-900">Dunia</span>{' '}
          <span className="text-primary-800">Informatika</span>
        </h2>

        <div className="relative group">
          <button
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
          </button>

          <div
            ref={scrollRef}
            className="flex gap-4 md:gap-6 overflow-x-auto snap-x snap-mandatory [&::-webkit-scrollbar]:hidden [-ms-overflow-style:none] [scrollbar-width:none]"
          >
            {cards.map((card, idx) => (
              <div
                key={idx}
                className="w-full shrink-0 md:w-[380px] lg:w-[420px] md:h-63 h-45 bg-primary-700 rounded-[20px] space-y-2 p-5 text-white relative overflow-hidden snap-center md:snap-start shadow-lg"
              >
                <h3 className="text-lg md:text-[26px] font-semibold">{card.title}</h3>
                <div className="w-full h-px bg-white z-10" />
                <p className="text-sm md:text-lg text-white/90 relative z-10 leading-relaxed max-w-[85%]">
                  {card.desc}
                </p>

                <div className="absolute -bottom-6 -right-6 md:-bottom-8 md:-right-8 text-primary-200 pointer-events-none rotate-13 z-0">
                  <card.Icon className="w-32 h-32 md:w-40 md:h-40" />
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="flex justify-between items-center md:hidden">
          <button
            onClick={() => scroll('left')}
            className="bg-primary-600 hover:bg-primary-700 p-[10px] rounded-lg cursor-pointer text-white transition-all shadow-md justify-center"
            aria-label="Scroll left"
          >
            <ChevronLeft className="w-6 h-6" />
          </button>

          <button
            onClick={() => scroll('right')}
            className="bg-primary-600 hover:bg-primary-700 p-[10px] rounded-lg text-white transition-all shadow-md justify-center"
            aria-label="Scroll right"
          >
            <ChevronRight className="w-6 h-6" />
          </button>
        </div>
      </Container>
    </section>
  );
}

export default Explore;
