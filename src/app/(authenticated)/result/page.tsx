import React from "react";
import Button from "@/components/ui/button";
import Podium, { type LeaderboardEntry } from "@/components/authenticated/result/podium";
import { displayName, type BackendUser } from "@/lib/auth";
import { backendJson } from "@/lib/backend";

async function getUser() {
  return backendJson<BackendUser>("/api/v1/me").catch(() => null);
}

interface BackendLeaderboardResponse {
  entries: {
    rank: number;
    user_id: string;
    full_name: string;
    final_score: number;
    submitted_at?: string;
  }[];
}

async function getLeaderboard() {
  return backendJson<BackendLeaderboardResponse>("/api/v1/leaderboard/current?limit=100").catch(() => null);
}

export default async function ResultPage() {
  const [user, leaderboardData] = await Promise.all([getUser(), getLeaderboard()]);

  const entries = leaderboardData?.entries || [];

  const updatedLeaderboard = entries.map((entry) => {
    return {
      rank: entry.rank,
      name: entry.full_name,
      score: entry.final_score,
      isCurrentUser: user ? entry.user_id === user.id : false,
    };
  }) as LeaderboardEntry[];

  const top3 = updatedLeaderboard.slice(0, 3);
  const otherParticipants = updatedLeaderboard.slice(3);

  return (
    <div className="w-full min-h-screen bg-white">
      <div className="max-w-6xl mx-auto px-4 py-6 sm:px-6 sm:py-8 lg:px-8 lg:py-16 xl:py-20">
        <div className="grid lg:grid-cols-[minmax(0,1fr)_minmax(360px,0.9fr)] lg:items-center lg:gap-x-8 xl:gap-x-12">
          <section className="order-2 flex flex-col gap-6 lg:order-1 lg:row-span-2 lg:self-end h-full justify-between">
            <div className="lg:mb-6">
              <h1 className="text-2xl font-bold text-neutral-1000 mb-1 lg:text-xl">
                Leaderboard Top 3
              </h1>
              <p className="text-xs text-neutral-600 lg:text-sm">
                Lihat peringkat peserta lainnya
              </p>
            </div>

            <Podium entries={top3} />
          </section>

          <section className="order-1 relative overflow-hidden rounded-xl bg-gradient-to-b from-primary-600 to-primary-1000 p-6 text-white shadow-md border border-white/5 flex flex-col gap-4 sm:flex-row mb-6 sm:items-center sm:justify-between lg:order-2 lg:min-h-32 lg:flex-col lg:items-stretch lg:justify-between">
            <div>
              <h2 className="text-xl sm:text-2xl lg:text-xl font-bold mb-2">
                Pembahasan
              </h2>
              <p className="text-xs text-primary-100/90 max-w-xs">
                Dapatkan solusi dan strategi penyelesaian soal
              </p>
            </div>
            <div className="flex justify-end">
              <Button size="sm" className="w-full sm:w-auto lg:min-w-32">
                Akses Pembahasan
              </Button>
            </div>
          </section>

          <section className="order-3 bg-white rounded-xl shadow-md border border-neutral-200 p-4 flex flex-col lg:order-3 lg:max-h-86">
            <div className="flex items-center justify-between px-3 pb-2.5 text-xs font-bold text-neutral-900 border-b border-neutral-100">
              <div className="w-[15%] text-center">Rank</div>
              <div className="w-[65%] text-left">Peserta</div>
              <div className="w-[20%] text-right">Skor</div>
            </div>

            <div className="flex flex-col mt-2 gap-1.5 max-h-[400px] overflow-y-auto pr-1 lg:max-h-76">
              {otherParticipants.map((entry, idx) => (
                <div
                  key={`participant-${entry.rank}-${entry.name}-${idx}`}
                  className={`flex items-center justify-between px-3.5 py-3 rounded-lg text-xs transition duration-150 lg:py-2 ${entry.isCurrentUser
                    ? "bg-primary-200 text-primary-950 font-bold border border-primary-300/40"
                    : "bg-primary-background/80 hover:bg-primary-background text-neutral-800 font-medium"
                    }`}
                >
                  <div className="w-[15%] text-center">{entry.rank}</div>
                  <div className="w-[65%] text-left truncate pr-2">{entry.name}</div>
                  <div className="w-[20%] text-right">{entry.score}</div>
                </div>
              ))}
            </div>
          </section>
        </div>
      </div>
    </div>
  );
}
