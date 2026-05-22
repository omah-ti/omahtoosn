"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { useRouter } from "next/navigation";
import { Grid, ChevronLeft, ChevronRight, CheckSquare, Square, X } from "lucide-react";
import Container from "@/components/ui/container";
import Button from "@/components/ui/button";

type ApiEnvelope<T> = {
  success: boolean;
  message: string;
  data?: T;
};

type TryoutQuestion = {
  id: string;
  code: string;
  question_type: "multiple_choice" | "short_text";
  prompt_html: string;
  display_order: number;
  options?: {
    key: string;
    text: string;
    display_order: number;
  }[];
};

type Attempt = {
  id: string;
  status: string;
  expires_at: string;
  version: number;
};

type AttemptAnswer = {
  question_id: string;
  selected_option_key?: string;
  answer_text?: string;
  is_flagged: boolean;
};

type AttemptPayload = {
  attempt: Attempt;
  questions: TryoutQuestion[];
  answers: AttemptAnswer[];
};

export default function TryoutPage() {
  const router = useRouter();
  const [currentIndex, setCurrentIndex] = useState(0);
  const [questions, setQuestions] = useState<TryoutQuestion[]>([]);
  const [attempt, setAttempt] = useState<Attempt | null>(null);
  const [answers, setAnswers] = useState<Record<number, string>>({});
  const [raguRagu, setRaguRagu] = useState<Record<number, boolean>>({});
  const [loading, setLoading] = useState(true);
  const [loadError, setLoadError] = useState("");

  const [timeLeft, setTimeLeft] = useState(601);
  const [showToast, setShowToast] = useState(false);
  const [showRaguModal, setShowRaguModal] = useState(false);
  const [showConfirmModal, setShowConfirmModal] = useState(false);
  const [showTimeoutModal, setShowTimeoutModal] = useState(false);
  const [showDaftarSoal, setShowDaftarSoal] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const saveTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const pendingSaveRef = useRef<{ index: number; value: string; isFlagged: boolean } | null>(null);

  const saveAnswer = useCallback(async (index: number, value: string, isFlagged: boolean) => {
    const question = questions[index];
    if (!question || !attempt) return;

    const isShortText = question.question_type === "short_text";
    const res = await fetch("/api/v1/attempts/current/answers", {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        attempt_version: attempt.version,
        answers: [
          {
            question_id: question.id,
            selected_option_key: isShortText || value === "" ? null : value,
            answer_text: isShortText && value !== "" ? value : null,
            is_flagged: isFlagged,
            client_updated_at: new Date().toISOString(),
          },
        ],
      }),
    });

    if (res.ok) {
      const payload = (await res.json().catch(() => null)) as ApiEnvelope<{ attempt: Attempt }> | null;
      if (payload?.data?.attempt) {
        setAttempt(payload.data.attempt);
      }
    }
  }, [questions, attempt]);

  const flushSave = useCallback(async () => {
    if (saveTimerRef.current) {
      clearTimeout(saveTimerRef.current);
      saveTimerRef.current = null;
    }
    if (pendingSaveRef.current) {
      const { index, value, isFlagged } = pendingSaveRef.current;
      pendingSaveRef.current = null;
      await saveAnswer(index, value, isFlagged);
    }
  }, [saveAnswer]);

  const scheduleSave = useCallback((index: number, value: string, isFlagged: boolean) => {
    pendingSaveRef.current = { index, value, isFlagged };
    if (saveTimerRef.current) {
      clearTimeout(saveTimerRef.current);
    }
    saveTimerRef.current = setTimeout(() => {
      saveTimerRef.current = null;
      pendingSaveRef.current = null;
      void saveAnswer(index, value, isFlagged);
    }, 800);
  }, [saveAnswer]);

  useEffect(() => {
    return () => {
      void flushSave();
    };
  }, [flushSave]);

  const loadAttempt = useCallback(async () => {
    setLoading(true);
    setLoadError("");

    let res = await fetch("/api/v1/attempts/current", { cache: "no-store" });
    if (res.status === 404) {
      const startRes = await fetch("/api/v1/tryouts/current/start", { method: "POST" });
      if (startRes.status === 401) {
        router.push("/login");
        return;
      }
      if (!startRes.ok) {
        const payload = await startRes.json().catch(() => null);
        setLoadError(payload?.message || "Gagal memulai tryout");
        setLoading(false);
        return;
      }
      res = await fetch("/api/v1/attempts/current", { cache: "no-store" });
    }

    if (res.status === 401) {
      router.push("/login");
      return;
    }
    if (!res.ok) {
      const payload = await res.json().catch(() => null);
      setLoadError(payload?.message || "Gagal memuat tryout");
      setLoading(false);
      return;
    }

    const payload = (await res.json()) as ApiEnvelope<AttemptPayload>;
    if (!payload.success || !payload.data) {
      setLoadError(payload.message || "Gagal memuat tryout");
      setLoading(false);
      return;
    }

    const answerByQuestion = new Map(payload.data.answers.map((answer) => [answer.question_id, answer]));
    const nextAnswers: Record<number, string> = {};
    const nextFlags: Record<number, boolean> = {};
    payload.data.questions.forEach((question, index) => {
      const answer = answerByQuestion.get(question.id);
      if (!answer) return;
      if (answer.selected_option_key) nextAnswers[index] = answer.selected_option_key;
      if (answer.answer_text) nextAnswers[index] = answer.answer_text;
      if (answer.is_flagged) nextFlags[index] = true;
    });

    setQuestions(payload.data.questions);
    setAttempt(payload.data.attempt);
    setAnswers(nextAnswers);
    setRaguRagu(nextFlags);
    setTimeLeft(Math.max(0, Math.floor((new Date(payload.data.attempt.expires_at).getTime() - Date.now()) / 1000)));
    setLoading(false);
  }, [router]);

  useEffect(() => {
    void Promise.resolve().then(loadAttempt);
  }, [loadAttempt]);

  useEffect(() => {
    if (timeLeft <= 0) {
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

  const currentSoal = questions[currentIndex];

  const handleNext = () => {
    void flushSave();
    if (currentIndex < questions.length - 1) setCurrentIndex(currentIndex + 1);
  };

  const handlePrev = () => {
    void flushSave();
    if (currentIndex > 0) setCurrentIndex(currentIndex - 1);
  };

  const handleSelectOption = (value: string) => {
    setAnswers(prev => ({ ...prev, [currentIndex]: value }));
    const isShortText = questions[currentIndex]?.question_type === "short_text";
    if (isShortText) {
      scheduleSave(currentIndex, value, raguRagu[currentIndex] === true);
    } else {
      void saveAnswer(currentIndex, value, raguRagu[currentIndex] === true);
    }
  };

  const handleToggleRagu = () => {
    void flushSave();
    const nextFlag = !raguRagu[currentIndex];
    setRaguRagu(prev => ({ ...prev, [currentIndex]: nextFlag }));
    void saveAnswer(currentIndex, answers[currentIndex] || "", nextFlag);
  };

  const handleSubmitClick = () => {
    const hasRagu = Object.values(raguRagu).some((val) => val === true);
    if (hasRagu) {
      setShowRaguModal(true);
    } else {
      setShowConfirmModal(true);
    }
  };

  const handleFinalSubmit = async () => {
    if (isSubmitting) return;
    setIsSubmitting(true);
    await flushSave();
    if (!attempt) {
      setIsSubmitting(false);
      return;
    }
    fetch("/api/v1/attempts/current/submit", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        attempt_version: attempt.version,
      }),
    }).then(async (res) => {
      setIsSubmitting(false);
      if (!res.ok) {
        const payload = await res.json().catch(() => null);
        alert(payload?.message || "Gagal mengirim jawaban");
        return;
      }
      setShowConfirmModal(false);
      setShowTimeoutModal(false);
      router.push("/dashboard");
    }).catch(() => {
      setIsSubmitting(false);
      alert("Gagal mengirim jawaban");
    });
  };

  const formatTime = (seconds: number) => {
    const h = Math.floor(seconds / 3600);
    const m = Math.floor((seconds % 3600) / 60);
    const s = seconds % 60;
    return `${h.toString().padStart(2, '0')}:${m.toString().padStart(2, '0')}:${s.toString().padStart(2, '0')}`;
  };

  if (loading) {
    return (
      <main className="relative min-h-[calc(100vh-64px)] w-full bg-white flex flex-col items-center justify-center py-10">
        <p className="text-neutral-1000">Memuat tryout...</p>
      </main>
    );
  }

  if (loadError || !currentSoal) {
    return (
      <main className="relative min-h-[calc(100vh-64px)] w-full bg-white flex flex-col items-center justify-center py-10">
        <p className="text-neutral-1000">{loadError || "Tryout tidak tersedia"}</p>
      </main>
    );
  }

  const shouldShowTimeoutModal = showTimeoutModal || timeLeft <= 0;

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
                {questions.map((soal, idx) => {
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
                        void flushSave();
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
                disabled={isSubmitting}
                className={`flex-1 h-10 bg-[#2563EB] text-white font-semibold hover:bg-blue-700 ${isSubmitting ? 'opacity-70 cursor-not-allowed' : ''}`}
              >
                {isSubmitting ? 'Mengirim...' : 'Kirim Sekarang'}
              </Button>
            </div>
          </div>
        </div>
      )}

      {/* Timeout Modal */}
      {shouldShowTimeoutModal && (
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
                disabled={isSubmitting}
                className={`w-full h-10 bg-[#2563EB] text-white font-semibold hover:bg-blue-700 ${isSubmitting ? 'opacity-70 cursor-not-allowed' : ''}`}
              >
                {isSubmitting ? 'Mengirim...' : 'Kirim Sekarang'}
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
              <div
                className="text-neutral-800 leading-relaxed text-sm md:text-base text-justify whitespace-pre-line"
                dangerouslySetInnerHTML={{ __html: currentSoal.prompt_html }}
              />
            </div>

            <div className="w-full lg:w-[450px] bg-white rounded-xl p-6 md:p-8 shadow-sm flex flex-col gap-6 max-h-[60vh] overflow-y-auto">
              <p className="text-neutral-800 leading-relaxed text-sm md:text-base whitespace-pre-line">
                Jawaban untuk {currentSoal.code}
              </p>

              <div className="flex flex-col gap-4">
                <p className="font-medium text-neutral-1000 text-sm md:text-base">Jawaban</p>
                <div className="flex flex-col gap-3">
                  {currentSoal.question_type === "short_text" ? (
                    <input
                      type="text"
                      value={answers[currentIndex] || ''}
                      onChange={(e) => handleSelectOption(e.target.value)}
                      placeholder="Masukkan jawaban bilangan bulat"
                      className="w-full px-4 py-3 rounded-lg border border-neutral-300 focus:outline-none focus:ring-2 focus:ring-primary-600 bg-white text-neutral-800 text-sm md:text-base"
                    />
                  ) : (
                    currentSoal.options?.map((opt, idx) => (
                      <label key={idx} className="flex items-center gap-3 cursor-pointer group">
                        <input 
                          type="radio" 
                          name={`answer-${currentIndex}`} 
                          className="w-4 h-4 accent-neutral-1000 cursor-pointer" 
                          checked={answers[currentIndex] === opt.key}
                          onChange={() => handleSelectOption(opt.key)}
                        />
                        <span
                          className="text-sm md:text-base text-neutral-800 group-hover:text-neutral-1000 transition-colors"
                          dangerouslySetInnerHTML={{ __html: `${opt.key}. ${opt.text}` }}
                        />
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
            {currentIndex === questions.length - 1 ? (
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
