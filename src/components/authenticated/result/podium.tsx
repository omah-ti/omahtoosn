import React from "react";

export interface LeaderboardEntry {
  rank: number;
  name: string;
  score: number;
  isCurrentUser?: boolean;
}

interface PodiumProps {
  entries: LeaderboardEntry[];
}

const Pillar1 = ({ entry }: { entry: LeaderboardEntry }) => (
  <div className="relative flex flex-col justify-end">
    <div className="w-full flex flex-col relative items-center">
      <span className="absolute -top-8 md:-top-9 left-1/2 -translate-x-1/2 z-10 whitespace-nowrap rounded-[10px] text-center px-2.5 md:px-3 py-1.5 bg-white text-sm md:text-sm font-bold text-neutral-900 shadow-xs">
        {entry.name}
      </span>
      <div className="h-5 w-27 md:w-34 lg:w-36 xl:w-40 bg-primary-600 [clip-path:polygon(10%_0,90%_0,100%_100%,0%_100%)] border-b-[0.1px] border-white" />
      <div className="w-27 md:w-34 lg:w-36 xl:w-40 h-62.5 md:h-72 lg:h-72 xl:h-80 bg-linear-to-b from-primary-600 to-primary-1000 flex justify-center">
        <span className="bg-primary-background h-fit mt-5 font-bold px-4 py-0.5 border border-black text-lg rounded-[10px]">
          {entry.score}
        </span>
        <span className="absolute bottom-18 text-5xl md:bottom-20 md:text-5xl font-bold text-white">
          1
        </span>
      </div>
    </div>
  </div>
);

const Pillar2 = ({ entry }: { entry: LeaderboardEntry }) => (
  <div className="relative w-full flex flex-col items-center justify-end">
    <div className="w-full flex flex-col relative items-center">
      <span className="absolute -top-8 md:-top-9 left-1/2 -translate-x-1/2 z-10 whitespace-nowrap rounded-[10px] text-center px-2.5 md:px-3 py-1.5 bg-white text-sm md:text-sm font-bold text-neutral-900 shadow-xs">
        {entry.name}
      </span>
      <div className="h-5 w-27 md:w-34 lg:w-36 xl:w-40 bg-primary-400 [clip-path:polygon(10%_0,100%_0,100%_100%,0%_100%)] border-b-[0.1px] border-white" />
      <div className="w-27 md:w-34 lg:w-36 xl:w-40 h-42.5 md:h-56 lg:h-56 xl:h-62.5 bg-linear-to-b from-primary-400 to-primary-800 flex justify-center">
        <span className="bg-primary-background h-fit mt-5 font-bold px-4 py-0.5 border border-black text-lg rounded-[10px]">
          {entry.score}
        </span>
        <span className="absolute bottom-16 text-5xl md:bottom-17 md:text-5xl font-bold text-white">
          2
        </span>
      </div>
    </div>
  </div>
);

const Pillar3 = ({ entry }: { entry: LeaderboardEntry }) => (
  <div className="relative w-full flex flex-col justify-end">
    <div className="w-full flex flex-col relative items-center">
      <span className="absolute -top-8 md:-top-9 left-1/2 -translate-x-1/2 z-10 whitespace-nowrap rounded-[10px] text-center px-2.5 md:px-3 py-1.5 bg-white text-sm md:text-sm font-bold text-neutral-900 shadow-xs">
        {entry.name}
      </span>
      <div className="h-5 w-27 md:w-34 lg:w-36 xl:w-40 bg-primary-400 [clip-path:polygon(0%_0,90%_0,100%_100%,0%_100%)] border-b-[0.1px] border-white" />
      <div className="w-27 md:w-34 lg:w-36 xl:w-40 h-37.5 md:h-50 lg:h-50 xl:h-55 bg-linear-to-b from-primary-400 to-primary-800 flex justify-center">
        <span className="bg-primary-background h-fit mt-5 font-bold px-4 py-0.5 border border-black text-lg rounded-[10px]">
          {entry.score}
        </span>
        <span className="absolute bottom-14 text-5xl md:bottom-15 md:text-5xl font-bold text-white">
          3
        </span>
      </div>
    </div>
  </div>
);

export default function Podium({ entries }: PodiumProps) {
  const p1 = entries.find((e) => e.rank === 1) || { rank: 1, name: "-", score: 0 };
  const p2 = entries.find((e) => e.rank === 2) || { rank: 2, name: "-", score: 0 };
  const p3 = entries.find((e) => e.rank === 3) || { rank: 3, name: "-", score: 0 };

  return (
    <div className="relative flex items-end justify-center w-full max-w-138 lg:max-w-108 xl:max-w-120 h-full mx-auto pt-10 px-4">
      <Pillar2 entry={p2} />
      <Pillar1 entry={p1} />
      <Pillar3 entry={p3} />
    </div>
  );
}
