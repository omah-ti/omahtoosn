"use client";

import { useState, useEffect } from "react";
import { Grid, ChevronLeft, ChevronRight, CheckSquare, Square, X } from "lucide-react";
import Container from "@/components/ui/container";
import Button from "@/components/ui/button";

const dummySoal = Array.from({ length: 40 }).map((_, i) => ({
  id: i + 1,
  type: i === 1 ? "isian" : "pilihan_ganda",
  text_soal: "Pak Dengklek memiliki N ekor bebek yang dinomori dari Bebek ke-1 hingga Bebek ke-N. Karena hari ini adalah Hari Mandi Nasional, Pak Dengklek ingin memandikan seluruh bebek-bebeknya. Diketahui bahwa Pak Dengklek memiliki dua buah gayung. Dengan gayung pertama, Pak Dengklek dapat secara tepat mengambil X liter; sedangkan dengan gayung kedua, Pak Dengklek dapat secara tepat mengambil Y liter. Perhatikan bahwa Pak Dengklek tidak bisa mengira-ngira bagian dari gayung tersebut, sehingga Pak Dengklek tidak dapat mengambil tepat separuh dari X liter atau sepertiga dari Y liter misalnya. Bebek-bebek Pak Dengklek sebenarnya tidak suka mandi. Sehingga, mereka memberikan dua persyaratan kepada Pak Dengklek agar mereka mau mandi sebagai berikut: 1. Sekali Pak Dengklek mengambil air dengan gayung, maka air tersebut harus langsung digunakan.",
  pertanyaan: "Berapakah banyak kemungkinan posisi awal Kwak sedemikian sehingga tidak ada satu pun instruksi yang menyebabkan Kwak keluar pekarangan?",
  options: ["Option A", "Option B", "Option C", "Option D", "Option E"]
}));

export default function TryoutPage() {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [answers, setAnswers] = useState<Record<number, string | number>>({});
  const [raguRagu, setRaguRagu] = useState<Record<number, boolean>>({});

  const [timeLeft, setTimeLeft] = useState(601);
  const [showToast, setShowToast] = useState(false);
  const [showRaguModal, setShowRaguModal] = useState(false);
  const [showConfirmModal, setShowConfirmModal] = useState(false);
  const [showTimeoutModal, setShowTimeoutModal] = useState(false);
  const [showDaftarSoal, setShowDaftarSoal] = useState(false);

  useEffect(() => {
    if (timeLeft <= 0) {
      setShowTimeoutModal(true);
      return;
    }

    const timer = setInterval(() => {
      setTimeLeft((prev) => {
        if (prev === 601) {
          setShowToast(true);
        }
        return prev - 1;
      });
    }, 1000);

    return () => clearInterval(timer);
  }, [timeLeft]);

  const currentSoal = dummySoal[currentIndex];

  const handleNext = () => {
    if (currentIndex < dummySoal.length - 1) setCurrentIndex(currentIndex + 1);
  };

  const handlePrev = () => {
    if (currentIndex > 0) setCurrentIndex(currentIndex - 1);
  };

  const handleSelectOption = (value: string | number) => {
    setAnswers(prev => ({ ...prev, [currentIndex]: value }));
  };

  const handleToggleRagu = () => {
    setRaguRagu(prev => ({ ...prev, [currentIndex]: !prev[currentIndex] }));
  };

  const handleSubmitClick = () => {
    const hasRagu = Object.values(raguRagu).some((val) => val === true);
    if (hasRagu) {
      setShowRaguModal(true);
    } else {
      setShowConfirmModal(true);
    }
  };

  const handleFinalSubmit = () => {
    // Final submission logic
    console.log("Answers Submitted:", answers);
    setShowConfirmModal(false);
    setShowTimeoutModal(false);
    // e.g. redirect to success page
  };

  const formatTime = (seconds: number) => {
    const h = Math.floor(seconds / 3600);
    const m = Math.floor((seconds % 3600) / 60);
    const s = seconds % 60;
    return `${h.toString().padStart(2, '0')}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`;
  };

  return (
    <main className="relative min-h-[calc(100vh-64px)] w-full bg-white flex flex-col items-center justify-center py-10">
      <div className="absolute top-0 left-0 w-full h-32 md:h-48 bg-primary-1000 z-0" />
      
      {/* Toast Notification */}
      {showToast && (
        <div className="fixed top-20 right-4 md:right-10 bg-[#FAE5D3] border border-[#E67E22] text-[#A04000] px-4 py-3 md:px-6 md:py-4 rounded-lg shadow-lg z-50 flex items-start gap-4 transition-all duration-300">
          <div className="flex flex-col">
            <span className="font-bold text-sm">10 Menit Tersisa!</span>
            <span className="text-xs">Waktu pengerjaan hampir habis.</span>
          </div>
          <button onClick={() => setShowToast(false)} className="text-[#A04000] hover:opacity-70 mt-0.5 cursor-pointer">
            <X className="w-4 h-4" />
          </button>
        </div>
      )}

      {/* Daftar Soal Modal */}
      {showDaftarSoal && (
        <div className="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4">
          <div className="bg-white rounded-lg flex flex-col overflow-hidden w-full max-w-[600px] shadow-2xl">
            <div className="bg-primary-1000 flex justify-between items-center p-4">
              <span className="text-white text-lg font-medium">Daftar Soal</span>
              <button onClick={() => setShowDaftarSoal(false)} className="text-white hover:opacity-70 cursor-pointer">
                <X className="w-5 h-5" />
              </button>
            </div>
            <div className="p-10">
              <div className="grid grid-cols-10 max-w-[640px] w-full gap-3 md:gap-4 justify-center items-center">
                {dummySoal.map((soal, idx) => {
                  const isRagu = raguRagu[idx];
                  const isAnswered = answers[idx] !== undefined && answers[idx] !== '';
                  let bgColor = 'bg-white';
                  if (currentIndex === idx) {
                    bgColor = 'bg-primary-400';
                  } else if (isRagu) {
                    bgColor = 'bg-[#FCD34D]';
                  } else if (isAnswered) {
                    bgColor = 'bg-primary-200';
                  }
                  const textColor = isRagu || isAnswered || currentIndex === idx ? 'text-neutral-1000' : 'text-neutral-800';
                  
                  return (
                    <button
                      key={idx}
                      onClick={() => {
                        setCurrentIndex(idx);
                        setShowDaftarSoal(false);
                      }}
                      className={`w-10 h-10 md:w-[42px] md:h-[42px] flex items-center justify-center rounded-md border border-neutral-300 text-sm md:text-base transition-colors ${bgColor} ${textColor} hover:opacity-80 cursor-pointer`}
                    >
                      {idx + 1}
                    </button>
                  );
                })}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Ragu-ragu Modal */}
      {showRaguModal && (
        <div className="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4">
          <div className="bg-white rounded-lg w-[370px] h-[196px] p-6 flex flex-col items-start gap-[10px] shrink-0 shadow-2xl">
            <h3 className="text-xl font-bold text-neutral-1000">Masih Ada Jawaban Ragu-ragu</h3>
            <p className="text-[#71717A] text-sm text-left mb-auto leading-relaxed">
              Anda masih mempunyai jawaban ragu-ragu, silahkan periksa kembali.
            </p>
            <div className="flex gap-3 w-full mt-auto">
              <Button 
                onClick={() => { setShowRaguModal(false); setShowConfirmModal(true); }}
                variant="transparent"
                className="flex-1 h-10 border border-[#2563EB] text-[#2563EB] font-semibold hover:bg-blue-50"
              >
                Tetap Submit
              </Button>
              <Button 
                onClick={() => setShowRaguModal(false)}
                variant="primary"
                className="flex-1 h-10 bg-[#2563EB] text-white font-semibold hover:bg-blue-700"
              >
                Periksa Kembali
              </Button>
            </div>
          </div>
        </div>
      )}

      {/* Confirm Submit Modal */}
      {showConfirmModal && (
        <div className="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4">
          <div className="bg-white rounded-lg w-[370px] h-[196px] p-6 flex flex-col items-start gap-[10px] shrink-0 shadow-2xl">
            <h3 className="text-xl font-bold text-neutral-1000">Konfirmasi Submit</h3>
            <p className="text-[#71717A] text-sm text-left mb-auto leading-relaxed">
              Apakah Anda yakin ingin mengirim jawaban sekarang? Setelah dikirim, jawaban tidak dapat diubah kembali.
            </p>
            <div className="flex gap-3 w-full mt-auto">
              <Button 
                onClick={() => setShowConfirmModal(false)}
                variant="transparent"
                className="flex-1 h-10 border border-[#2563EB] text-[#2563EB] font-semibold hover:bg-blue-50"
              >
                Kembali
              </Button>
              <Button 
                onClick={handleFinalSubmit}
                variant="primary"
                className="flex-1 h-10 bg-[#2563EB] text-white font-semibold hover:bg-blue-700"
              >
                Kirim Sekarang
              </Button>
            </div>
          </div>
        </div>
      )}

      {/* Timeout Modal */}
      {showTimeoutModal && (
        <div className="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4">
          <div className="bg-white rounded-lg w-[370px] h-[196px] p-6 flex flex-col items-start gap-[10px] shrink-0 shadow-2xl">
            <h3 className="text-xl font-bold text-neutral-1000">Waktu Pengerjaan Telah Habis</h3>
            <p className="text-[#71717A] text-sm text-left mb-auto leading-relaxed">
              Jawaban Anda akan otomatis dikumpulkan dan tidak dapat diubah.
            </p>
            <div className="flex justify-start w-[50%] mt-auto ml-auto">
              <Button 
                onClick={handleFinalSubmit}
                variant="primary"
                className="w-full h-10 bg-[#2563EB] text-white font-semibold hover:bg-blue-700"
              >
                Kirim Sekarang
              </Button>
            </div>
          </div>
        </div>
      )}

      <Container className="relative z-10 w-full max-w-[1200px]">
        <div className="relative w-full bg-[#DEE5F4] rounded-2xl p-6 md:p-8 flex flex-col gap-6 shadow-xl">
          <div className="flex flex-col md:flex-row justify-between items-center gap-4 pb-2 border-b border-neutral-300 md:border-b-0 md:pb-0">
            <h2 className="text-[20px] md:text-[24px] font-bold text-neutral-1000">
              Soal Nomor {String(currentIndex + 1).padStart(2, '0')}
            </h2>
            <div className="flex items-center gap-3">
              <div className="bg-white px-4 py-2 rounded-lg font-bold text-sm text-neutral-1000 shadow-sm">
                Sisa Waktu: {formatTime(timeLeft)}
              </div>
              <Button onClick={() => setShowDaftarSoal(true)} variant="transparent" className="bg-white px-4 py-2 rounded-lg flex items-center gap-2 font-bold text-sm text-neutral-1000 shadow-sm border border-neutral-200 hover:bg-neutral-50 transition-colors cursor-pointer">
                <Grid className="w-4 h-4" />
                Daftar Soal
              </Button>
            </div>
          </div>

          <div className="flex flex-col lg:flex-row gap-6">
            <div className="flex-1 bg-white rounded-xl p-6 md:p-8 shadow-sm min-h-[400px] max-h-[60vh] overflow-y-auto">
              <p className="text-neutral-800 leading-relaxed text-sm md:text-base text-justify whitespace-pre-line">
                {currentSoal.text_soal}
              </p>
            </div>

            <div className="w-full lg:w-[450px] bg-white rounded-xl p-6 md:p-8 shadow-sm flex flex-col gap-6 max-h-[60vh] overflow-y-auto">
              <p className="text-neutral-800 leading-relaxed text-sm md:text-base whitespace-pre-line">
                {currentSoal.pertanyaan}
              </p>

              <div className="flex flex-col gap-4">
                <p className="font-medium text-neutral-1000 text-sm md:text-base">Jawaban</p>
                <div className="flex flex-col gap-3">
                  {currentSoal.type === "isian" ? (
                    <input
                      type="text"
                      value={(answers[currentIndex] as string) || ''}
                      onChange={(e) => handleSelectOption(e.target.value)}
                      placeholder="Masukkan jawaban bilangan bulat"
                      className="w-full px-4 py-3 rounded-lg border border-neutral-300 focus:outline-none focus:ring-2 focus:ring-primary-600 bg-white text-neutral-800 text-sm md:text-base"
                    />
                  ) : (
                    currentSoal.options.map((opt, idx) => (
                      <label key={idx} className="flex items-center gap-3 cursor-pointer group">
                        <input 
                          type="radio" 
                          name={`answer-${currentIndex}`} 
                          className="w-4 h-4 accent-neutral-1000 cursor-pointer" 
                          checked={answers[currentIndex] === idx}
                          onChange={() => handleSelectOption(idx)}
                        />
                        <span className="text-sm md:text-base text-neutral-800 group-hover:text-neutral-1000 transition-colors">{opt}</span>
                      </label>
                    ))
                  )}
                </div>
              </div>
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-2">
            <button 
              onClick={handlePrev}
              disabled={currentIndex === 0}
              className={`w-full py-3 rounded-xl border border-primary-600 bg-transparent flex items-center justify-center gap-2 text-primary-600 font-semibold transition-colors cursor-pointer ${currentIndex === 0 ? 'opacity-50 cursor-not-allowed' : 'hover:bg-primary-100'}`}
            >
              <ChevronLeft className="w-5 h-5" />
              Sebelumnya
            </button>
            <button 
              onClick={handleToggleRagu}
              className={`w-full py-3 rounded-xl flex items-center justify-center gap-2 font-semibold transition-colors cursor-pointer ${raguRagu[currentIndex] ? 'bg-yellow-500 text-neutral-1000' : 'bg-[#FCD34D] text-neutral-1000 hover:bg-yellow-400'}`}
            >
              {raguRagu[currentIndex] ? <CheckSquare className="w-5 h-5" /> : <Square className="w-5 h-5" />}
              Ragu-ragu
            </button>
            {currentIndex === dummySoal.length - 1 ? (
              <button 
                onClick={handleSubmitClick}
                className="w-full py-3 rounded-xl bg-[#2563EB] flex items-center justify-center gap-2 text-white font-semibold hover:bg-blue-700 transition-colors cursor-pointer"
              >
                Submit
              </button>
            ) : (
              <button 
                onClick={handleNext}
                className="w-full py-3 rounded-xl bg-[#2563EB] flex items-center justify-center gap-2 text-white font-semibold hover:bg-blue-700 transition-colors cursor-pointer"
              >
                Berikutnya
                <ChevronRight className="w-5 h-5" />
              </button>
            )}
          </div>
        </div>
      </Container>
    </main>
  );
}
