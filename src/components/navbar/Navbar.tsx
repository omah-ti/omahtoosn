"use client";

import { useState } from "react";
import Button from "@/components/ui/button";
import NavbarMenu from "./NavbarMenu";

export default function Navbar() {
  const [open, setOpen] = useState(false);

  return (
    <nav className="bg-[#1f2a3a] h-16 px-10 flex items-center justify-between text-white">
      {/* Logo */}
      <div className="flex items-center gap-3 font-semibold text-lg">
        <div className="w-6 h-6 bg-orange-400 rotate-45"></div>
        OmahTOOSN
      </div>

      {/* Menu */}
      <NavbarMenu />

      {/* User */}
      <div className="relative">
        <Button variant="primary" onClick={() => setOpen(!open)}>
          User
        </Button>

        {open && (
          <div className="absolute right-0 mt-3 bg-[#3d475a] rounded-xl p-4 w-44 shadow-lg border border-white/10">
            <button className="flex items-center gap-3 w-full text-left hover:opacity-80">
              <span>⎋</span>
              Keluar
            </button>

            <button className="flex items-center gap-3 w-full text-left mt-3 hover:opacity-80">
              <span>?</span>
              Bantuan
            </button>
          </div>
        )}
      </div>
    </nav>
  );
}
