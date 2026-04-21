"use client";

import { useState } from "react";
import NavbarMenu from "./NavbarMenu";
import Hamburger from "./Hamburger";

export default function Navbar() {
  const [open, setOpen] = useState(false);
  const [menuOpen, setMenuOpen] = useState(false);

  return (
    <nav className="bg-[var(--primary-1000)] h-16 px-4 sm:px-6 md:px-10 flex items-center justify-between text-white relative border-b border-white/10">
      {/* Logo */}
      <div className="flex items-center gap-2 sm:gap-3">
        <img
          src="/navbar/logo.webp"
          alt="OmahTI"
          className="w-11 mt-2 h-11 sm:w-11 sm:h-10 object-contain"
        />

        <img
          src="/navbar/omahtoosn.webp"
          alt="OmahTOOSN"
          className="h-8 sm:h-8 object-contain"
        />
      </div>

      {/* Desktop Menu */}
      <div className="hidden text-md  lg:block absolute left-1/2 -translate-x-1/2">
        <NavbarMenu />
      </div>

      {/* Right Side */}
      <div className="flex items-center gap-3">
        {/* Desktop User */}
        <div className="hidden lg:block relative">
          <button onClick={() => setOpen(!open)}>
            <img
              src="/navbar/user.webp"
              alt="User"
              className="w-7 h-7 mt-2 gobject-cover cursor-pointer hover:opacity-80 transition"
            />
          </button>

          {open && (
            <div className="absolute pt-[12] right-0 mt-3 bg-[var(--primary-1000)] rounded-xl p-4 w-40 shadow-lg border border-white/10">
              <button className="ml-5 flex items-center gap-3 w-full text-left hover:opacity-80">
                <img
                  src="/navbar/log_out.webp"
                  alt="User"
                  className="w-4 h-4 gobject-cover cursor-pointer hover:opacity-80 transition"
                />
                Keluar
              </button>

              <button className="flex items-center gap-3 w-full text-left mt-3 mb-[-5] hover:opacity-80">
                <img
                  src="/navbar/help_circle.webp"
                  alt="User"
                  className="ml-5 w-4 h-4 gobject-cover cursor-pointer hover:opacity-80 transition"
                />
                Bantuan
              </button>
            </div>
          )}
        </div>

        {/* Hamburger (mobile) */}
        <button
          className="lg:hidden text-2xl"
          onClick={() => setMenuOpen(true)}
        >
          ☰
        </button>
      </div>

      {/* Mobile Sidebar */}
      <Hamburger open={menuOpen} onClose={() => setMenuOpen(false)} />
    </nav>
  );
}
