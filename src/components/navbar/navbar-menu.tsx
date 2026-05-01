"use client";

import { useState, useRef, useEffect } from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";

export default function NavbarMenu() {
  const items = [
    { name: "Dashboard", path: "/" },
    { name: "Try Out", path: "/try-out" },
    { name: "Hasil Try Out", path: "/hasil-try-out" },
  ];

  const pathname = usePathname();
  const [left, setLeft] = useState(0);

  const containerRef = useRef<HTMLDivElement>(null);
  const itemRefs = useRef<(HTMLAnchorElement | null)[]>([]);
  const underlineWidth = 90;

  useEffect(() => {
    const activeIndex = items.findIndex((item) => pathname === item.path);
    const currentIndex = activeIndex !== -1 ? activeIndex : 0;

    const el = itemRefs.current[currentIndex];
    const container = containerRef.current;

    if (el && container) {
      const center = el.offsetLeft + el.offsetWidth / 2;
      setLeft(center - underlineWidth / 2);
    }
  }, [pathname]);

  return (
    <div
      ref={containerRef}
      className="relative flex gap-16 h-full items-center text-white font-medium"
    >
      {items.map((item, index) => {
        const isActive = pathname === item.path;

        return (
          <Link
            key={item.name}
            href={item.path}
            ref={(el) => {
              itemRefs.current[index] = el;
            }}
            className={`transition-colors duration-200 ${
              isActive ? "text-white" : "text-gray-300 hover:text-white"
            }`}
          >
            {item.name}
          </Link>
        );
      })}

      {/* Underline */}
      <div
        className="absolute bottom-[-20.5px] h-[4px] bg-[var(--primary-700)] rounded transition-all duration-300 ease-in-out"
        style={{
          width: underlineWidth,
          left: left + 2,
        }}
      />
    </div>
  );
}
