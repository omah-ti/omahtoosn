"use client";

import Link from "next/link";
import { useState } from "react";

type Props = {
  open: boolean;
  onClose: () => void;
};

export default function Hamburger({ open, onClose }: Props) {
  const [userOpen, setUserOpen] = useState(false);

  const handleClose = () => {
    setUserOpen(false);
    onClose();
  };

  return (
    <div
      className={`fixed inset-0 z-50 
      ${open ? "pointer-events-auto" : "pointer-events-none"}`}
    >
      {/* overlay */}
      <div
        className={`absolute inset-0 bg-black/30 transition-opacity duration-300
        ${open ? "opacity-100" : "opacity-0"}`}
        onClick={handleClose}
      />

      {/* panel */}
      <div
        onClick={(e) => e.stopPropagation()}
        className={`absolute top-0 left-0 w-full 
        ${userOpen ? "h-[50vh]" : "h-[35vh]"} 
        bg-[#0b1a33] text-white p-6
        transform transition-all duration-300 ease-in-out
        ${open ? "translate-y-0" : "-translate-y-full"}`}
      >
        {/* Header */}
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center gap-3">
            <img src="/navbar/logo.webp" className="h-6" />
            <img src="/navbar/omahtoosn.webp" className="h-6" />
          </div>

          <button onClick={handleClose}>✕</button>
        </div>

        {/* Menu */}
        <div className="pt-3 flex flex-col gap-5 text-base">
          <Link href="/" className="hover:opacity-80" onClick={handleClose}>
            Dashboard
          </Link>
          <Link
            href="/try-out"
            className="hover:opacity-80"
            onClick={handleClose}
          >
            Try Out
          </Link>
          <Link
            href="/hasil-try-out"
            className="hover:opacity-80"
            onClick={handleClose}
          >
            Hasil Try Out
          </Link>
        </div>

        {/* USER ICON */}
        <div className="absolute right-5 top-48">
          <button
            onClick={(e) => {
              e.stopPropagation();
              setUserOpen((prev) => !prev);
            }}
          >
            <img
              src="/navbar/user.webp"
              className="w-6 h-6 mr-[3] object-cover"
            />
          </button>
        </div>

        {/* DROPDOWN */}
        {userOpen && (
          <div
            onClick={(e) => e.stopPropagation()}
            className="absolute left-10 top-[255px] flex flex-col gap-3"
          >
            <button className="flex items-center gap-3 hover:opacity-80">
              <img src="/navbar/log_out.webp" className="w-4 h-4" />
              Keluar
            </button>

            <button className="flex items-center gap-3 hover:opacity-80">
              <img src="/navbar/help_circle.webp" className="w-4 h-4" />
              Bantuan
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
