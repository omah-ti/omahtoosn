"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { Menu, X, User, LogOut, Info } from "lucide-react";
import Logo from "@/components/ui/logo";
import Button from "../ui/button";

export default function Navbar() {
  const [menuOpen, setMenuOpen] = useState(false);
  const [scrolled, setScrolled] = useState(false);
  const [userOpen, setUserOpen] = useState(false); // Desktop user dropdown
  const [userOpenMobile, setUserOpenMobile] = useState(false); // Mobile user dropdown

  // TODO: Replace with actual authentication check
  const isLoggedIn = false;

  useEffect(() => {
    const handleScroll = () => setScrolled(window.scrollY > 0);
    handleScroll();
    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  const isSolid = scrolled || menuOpen;

  return (
    <nav
      className={`z-10 h-16 fixed w-full px-4 sm:px-6 md:px-10 flex items-center justify-between text-white transition-colors duration-300
      ${isSolid ? "bg-primary-600" : "bg-transparent"}`}
    >
      {/* Logo */}
      <Logo />

      {/* Desktop Right Side */}
      <div className="hidden lg:flex items-center gap-4">
        <Link
          href="/about"
          className="text-sm font-medium text-white/90 hover:text-white transition-colors"
        >
          Tentang Ilmu Komputer
        </Link>

        {isLoggedIn ? (
          <div className="relative">
            <button onClick={() => setUserOpen(!userOpen)}>
              <User />
            </button>

            {userOpen && (
              <div className="absolute pt-3 right-0 mt-3 bg-primary-1000 rounded-xl p-4 w-48 shadow-lg border border-white/10 flex flex-col gap-4">
                <Link href="/dashboard" className="flex items-center gap-3 w-full text-left hover:opacity-80 text-sm">
                  <User size={18} />
                  Dashboard
                </Link>
                <button className="flex items-center gap-3 w-full text-left hover:opacity-80 text-sm">
                  <LogOut size={18} />
                  Keluar
                </button>
                <button className="flex items-center gap-3 w-full text-left hover:opacity-80 text-sm">
                  <Info size={18} />
                  Bantuan
                </button>
              </div>
            )}
          </div>
        ) : (
          <>
            <Link
              href="/login"
            >
              <Button variant="transparent" className="border border-primary-100 bg-primary-800 font-bold text-primary-100">
                Masuk
              </Button>
            </Link>

            <Link
              href="/register"
            >
              <Button variant="transparent" className="text-black bg-primary-100">
                Coba Sekarang
              </Button>
            </Link>
          </>
        )}
      </div>

      {/* Mobile Hamburger Toggle */}
      <button
        className="lg:hidden text-2xl transition-transform duration-300"
        onClick={() => setMenuOpen((prev) => !prev)}
      >
        <span
          className={`inline-block transition-all duration-300 p-2 bg-primary-100 rounded-lg text-primary-600 ${menuOpen ? "rotate-90 scale-110" : "rotate-0 scale-100"
            }`}
        >
          {menuOpen ? <X size={20} /> : <Menu size={20} />}
        </span>
      </button>

      {/* Mobile Overlay */}
      <div
        className={`fixed inset-0 top-16 z-40 bg-black/30 transition-opacity duration-300 lg:hidden
        ${menuOpen ? "opacity-100 pointer-events-auto" : "opacity-0 pointer-events-none"}`}
        onClick={() => setMenuOpen(false)}
      />

      {/* Mobile Panel */}
      <div
        className={`fixed top-16 left-0 w-full z-50 bg-primary-600 text-white overflow-hidden transition-all duration-300 ease-in-out lg:hidden
        ${menuOpen ? "max-h-80" : "max-h-0"}`}
      >
        <div className="p-6 flex flex-col gap-5">
          <Link
            href="/about"
            className="text-base hover:opacity-80"
            onClick={() => setMenuOpen(false)}
          >
            Tentang Ilmu Komputer
          </Link>

          {isLoggedIn ? (
            <>
              <div className="flex justify-end mt-2">
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    setUserOpenMobile((prev) => !prev);
                  }}
                >
                  <User />
                </button>
              </div>

              <div
                className={`overflow-hidden transition-all duration-300 ease-in-out ${userOpenMobile ? "max-h-40 opacity-100" : "max-h-0 opacity-0"
                  }`}
              >
                <div
                  onClick={(e) => e.stopPropagation()}
                  className="flex flex-col gap-4 pl-4 border-t border-white/10 pt-4"
                >
                  <Link href="/dashboard" className="flex items-center gap-3 hover:opacity-80 text-sm" onClick={() => setMenuOpen(false)}>
                    <User size={18} />
                    Dashboard
                  </Link>
                  <button className="flex items-center gap-3 hover:opacity-80 text-sm">
                    <LogOut size={18} />
                    Keluar
                  </button>
                  <button className="flex items-center gap-3 hover:opacity-80 text-sm">
                    <Info size={18} />
                    Bantuan
                  </button>
                </div>
              </div>
            </>
          ) : (
            <>
              <Link
                href="/login"
                onClick={() => setMenuOpen(false)}
              >
                <Button variant="transparent" className="border border-primary-100 w-full bg-primary-900 font-bold text-primary-100">
                  Masuk
                </Button>
              </Link>

              <Link
                href="/register"
                onClick={() => setMenuOpen(false)}
              >
                <Button variant="transparent" className="bg-primary-100 w-full text-black">
                  Coba Sekarang
                </Button>
              </Link>
            </>
          )}
        </div>
      </div>
    </nav>
  );
}
