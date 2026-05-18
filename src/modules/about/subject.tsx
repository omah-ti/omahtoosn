'use client';
import { useRef } from 'react';
import Container from "@/components/ui/container";
import Image from "next/image";
import { ChevronLeft, ChevronRight } from 'lucide-react';

const Subject = () => {
  const scrollRef = useRef<HTMLDivElement>(null);

  const scroll = (direction: 'left' | 'right') => {
    if (scrollRef.current) {
      const { current } = scrollRef;
      const scrollAmount = direction === 'left' ? -current.offsetWidth : current.offsetWidth;
      current.scrollBy({ left: scrollAmount, behavior: 'smooth' });
    }
  };

  const semesters = [
    {
      title: 'Semester 1',
      subjects: [
        'Pemrograman',
        'Praktikum Pemrograman',
        'Logika Informatika',
        'Aljabar Linear Fundamental',
        'Kalkulus 1',
        'Kimia Dasar',
        'Fisika Dasar',
        'Bahasa Indonesia',
        'Agama'
      ]
    },
    {
      title: 'Semester 2',
      subjects: [
        'Algoritma dan Struktur Data',
        'Integral dan Persamaan Diferensial',
        'Matematika Diskrit',
        'Organisasi dan Arsitektur Komputer',
        'Pengantar Statistika',
        'Sistem Digital',
        'Praktikum Algoritma dan Struktur Data',
        'Pancasila',
        'Kewarganegaraan'
      ]
    },
    {
      title: 'Semester 3',
      subjects: [
        'Analisis Algoritma dan Kompleksitas',
        'Basis Data',
        'Matematika Diskrit',
        'Jaringan Komputer',
        'Praktikum Basis Data',
        'Praktikum Sistem Komputer dan Jaringan',
        'Probabilitas dan Proses Stokastika',
        'Sistem Operasi'
      ]
    },
    {
      title: 'Semester 4',
      subjects: [
        'Filsafat Ilmu Komputer',
        'Pengembangan Startup Digital',
        'Metode Rekayasa Perangkat Lunak',
        'Workshop Implementasi Rancangan Perangkat Lunak',
        'Bahasa dan Otomata',
        'Metode Numerik',
        'Kriptografi dan Keamanan Informasi',
        'Mata Kuliah Pilihan'
      ]
    },
    {
      title: 'Semester 5',
      subjects: [
        'Pembelajaran Mesin Mendalam',
        'Filsafat Ilmu Komputer',
        'Mata Kuliah Pilihan'
      ]
    },
    {
      title: 'Semester 6',
      subjects: [
        'Pembelajaran Mesin Mendalam',
        'Filsafat Ilmu Komputer',
        'Mata Kuliah Pilihan'
      ]
    },
    {
      title: 'Semester 7',
      subjects: [
        'Pembelajaran Mesin Mendalam',
        'Filsafat Ilmu Komputer',
        'Mata Kuliah Pilihan'
      ]
    },
    {
      title: 'Semester 8',
      subjects: [
        'Pembelajaran Mesin Mendalam',
        'Filsafat Ilmu Komputer',
        'Mata Kuliah Pilihan'
      ]
    }
  ];

  return (
    <section className="relative bg-white py-20 overflow-hidden min-h-screen flex items-center">
      <div className="absolute bottom-0 w-full h-full pointer-events-none z-0 overflow-hidden">
        <Image
          src="/subject-circles.webp"
          alt="Background Circles"
          width={429}
          height={429}
          className="absolute bottom-40 scale-80 md:bottom-20 -left-60 object-cover object-bottom"
        />
       <Image
          src="/subject-circles.webp"
          alt="Background Circles"
          width={429}
          height={429}
          className="absolute bottom-50 md:bottom-30 scale-80 -right-60 object-cover object-bottom"
        />
      </div>

      <div className="absolute bottom-0 left-1/2 -translate-x-1/2 translate-y-[60%] md:translate-y-[80%] w-[1200px] h-[1200px] md:w-[2500px] md:h-[2500px] bg-primary-600 pointer-events-none rounded-[50%] z-0" />

      <Container className="relative z-10 flex flex-col gap-5 md:gap-6">
        <h2 className="md:text-[38px] text-[28px] font-bold text-center">
          <span className="text-primary-1000">Apa Saja yang Akan </span>
          <span className="text-primary-400">Dipelajari?</span>
        </h2>

        <div className="w-full flex items-center justify-center gap-6 group">
          <button
            onClick={() => scroll('left')}
            className="bg-primary-600 hover:bg-primary-700 p-2 rounded-lg cursor-pointer text-white transition-all hidden md:block shadow-md shrink-0"
            aria-label="Scroll left"
          >
            <ChevronLeft className="w-6 h-6" />
          </button>

          <div
            ref={scrollRef}
            className="flex gap-4 md:gap-5 overflow-x-auto snap-x snap-mandatory [&::-webkit-scrollbar]:hidden [-ms-overflow-style:none] [scrollbar-width:none] items-center min-w-0 w-full"
          >
            {semesters.map((sem, idx) => (
              <div
                key={idx}
                className="w-full shrink-0 max-w-48 md:max-w-66 min-h-107 md:min-h-130 bg-radial from-primary-800 to-primary-900 rounded-[20px] py-6 px-1 md:p-[30px] text-white relative flex flex-col snap-center mx-auto md:mx-0"
              >
                <h3 className="text-[22px] font-bold text-center mb-3">{sem.title}</h3>
                <div className="w-33 mx-auto md:w-full h-px bg-white mb-5" />
                <ul className="flex flex-col gap-[10px] text-center">
                  {sem.subjects.map((subject, sIdx) => (
                    <li key={sIdx} className="text-xs md:text-base text-white/90">
                      {subject}
                    </li>
                  ))}
                </ul>
              </div>
            ))}
          </div>

          <button
            onClick={() => scroll('right')}
            className="bg-primary-600 hover:bg-primary-700 p-2 rounded-lg cursor-pointer text-white transition-all hidden md:block shadow-md shrink-0"
            aria-label="Scroll right"
          >
            <ChevronRight className="w-6 h-6" />
          </button>
        </div>

        <div className="flex justify-between items-center md:hidden">
          <button
            onClick={() => scroll('left')}
            className="bg-primary-200 hover:bg-primary-300 p-[10px] rounded-lg cursor-pointer text-primary-1000 transition-all shadow-md justify-center"
            aria-label="Scroll left"
          >
            <ChevronLeft className="w-6 h-6" />
          </button>

          <button
            onClick={() => scroll('right')}
            className="bg-primary-200 hover:bg-primary-300 p-[10px] rounded-lg text-primary-1000 transition-all shadow-md justify-center"
            aria-label="Scroll right"
          >
            <ChevronRight className="w-6 h-6" />
          </button>
        </div>
      </Container>
    </section>
  )
}

export default Subject;