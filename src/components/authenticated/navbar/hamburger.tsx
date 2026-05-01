"use client";

import { Info, LogOut, User } from "lucide-react";
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
    <>
      {/* overlay */}
      <div
        className={`fixed inset-0 top-16 z-40 bg-black/30 transition-opacity duration-300
        ${open ? "opacity-100 pointer-events-auto" : "opacity-0 pointer-events-none"}`}
        onClick={handleClose}
      />

      {/* panel */}
      <div
        onClick={(e) => e.stopPropagation()}
        className={`fixed top-16 left-0 w-full z-50
        bg-primary-1000 text-white
        overflow-hidden
        transition-all duration-300 ease-in-out
        ${open ? (userOpen ? "max-h-[50vh]" : "max-h-80") : "max-h-0"}
        `}
        style={{ borderBottom: open ? "1px solid rgba(255,255,255,0.1)" : "none" }}
      >
        <div className="p-6">
          {/* Menu */}
          <div className="flex flex-col gap-5 text-base">
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
          <div className="flex justify-end mt-6">
            <button
              onClick={(e) => {
                e.stopPropagation();
                setUserOpen((prev) => !prev);
              }}
            >
              <User/>
            </button>
          </div>

          {/* DROPDOWN */}
          <div
            className={`overflow-hidden transition-all duration-300 ease-in-out
            ${userOpen ? "max-h-40 opacity-100 mt-4" : "max-h-0 opacity-0 mt-0"}`}
          >
            <div
              onClick={(e) => e.stopPropagation()}
              className="flex flex-col gap-3 pl-4"
            >
              <button className="flex items-center gap-3 hover:opacity-80">
                <LogOut />
                Keluar
              </button>

              <button className="flex items-center gap-3 hover:opacity-80">
                <Info />
                Bantuan
              </button>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}
