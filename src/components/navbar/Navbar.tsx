"use client";

import { useState } from "react";
import NavbarMenu from "./navbar-menu";
import Hamburger from "./hamburger";
import Link from "next/link";
import Image from "next/image";
import { Info, LogOut, Menu, User, X } from "lucide-react";

export default function Navbar() {
  const [open, setOpen] = useState(false);
  const [menuOpen, setMenuOpen] = useState(false);

  return (
    <nav className="bg-(--primary-1000) z-10 h-16 fixed w-full px-4 sm:px-6 md:px-10 flex items-center justify-between text-white border-b border-white/10">
      {/* Logo */}
      <Link href="/">
        <div className="flex items-center gap-2 sm:gap-3">
          <Image
            src="/navbar/logo.webp"
            alt="OmahTI"
            className="w-11 h-11 sm:w-11 sm:h-10 object-contain"
            width={100}
            height={100}
          />

          <Image
            src="/navbar/omahtoosn.webp"
            alt="OmahTOOSN"
            className="h-9 w-auto object-contain"
            width={200}
            height={100}
          />
        </div>
      </Link>

      {/* Desktop Menu */}
      <div className="hidden text-md  lg:block absolute left-1/2 -translate-x-1/2">
        <NavbarMenu />
      </div>

      {/* Right Side */}
      <div className="flex items-center gap-3">
        {/* Desktop User */}
        <div className="hidden lg:block relative">
          <button onClick={() => setOpen(!open)}>
            <User/>
          </button>

          {open && (
            <div className="absolute pt-[12] right-0 mt-3 bg-(--primary-1000) rounded-xl p-4 w-40 shadow-lg border border-white/10">
              <button className="ml-5 flex items-center gap-3 w-full text-left hover:opacity-80">
                <LogOut/>
                Keluar
              </button>

              <button className="flex items-center gap-3 w-full text-left mt-3 mb-[-5] hover:opacity-80">
                <Info/>
                Bantuan
              </button>
            </div>
          )}
        </div>

        {/* Hamburger (mobile) */}
        <button
          className="lg:hidden text-2xl transition-transform duration-300"
          onClick={() => setMenuOpen((prev) => !prev)}
        >
          <span
            className={`inline-block transition-all duration-300 ${menuOpen ? "rotate-90 scale-110" : "rotate-0 scale-100"
              }`}
          >
            {menuOpen ? <X /> : <Menu />}
          </span>
        </button>
      </div>

      {/* Mobile Sidebar */}
      <Hamburger open={menuOpen} onClose={() => setMenuOpen(false)} />
    </nav>
  );
}
