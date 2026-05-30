"use client";

import Image from "next/image";
import { useRouter } from "next/navigation";
import type { FormEvent } from "react";
import { ChevronLeft } from "lucide-react";

export default function Register() {
  const router = useRouter();

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const formData = new FormData(event.currentTarget);
    const res = await fetch("/api/v1/auth/register", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        full_name: formData.get("full_name"),
        email: formData.get("email"),
        password: formData.get("password"),
      }),
    });

    if (!res.ok) {
      const payload = await res.json().catch(() => null);
      alert(payload?.message || "Registrasi gagal");
      return;
    }

    router.push("/login");
  }

  return (
    <main className="flex w-full min-h-screen items-center justify-center bg-linear-to-t from-primary-900 to-primary-1000 text-black font-(Plus Jakarta Sans) px-4 py-16">

      <button
        type="button"
        onClick={() => router.back()}
        className="flex items-center gap-1 absolute left-10 top-15 sm:left-8 sm:top-8 rounded-lg bg-primary-600 px-3 py-2 sm:px-3 sm:py-2 text-sm font-medium text-white transition-colors cursor-pointer"
      >
        <ChevronLeft className="w-4 h-4 text-white" />
        <span>Kembali</span>
      </button>
	  
      <div className="w-full max-w-full sm:max-w-md bg-neutral-100 rounded-3xl px-5 py-6 sm:px-5 sm:py-6 flex flex-col items-center gap-10">

        <div className="flex flex-col items-center gap-3">
          <Image
            src="/logo/omahti-dark.webp"
            alt="Omahti Logo"
            width={72}
            height={72}
            className="object-contain sm:scale-140"
          />
          <h1 className="text-neutral-1000 font-semibold text-3xl sm:text-3xl text-center">
            Daftar Akun Peserta
          </h1>
        </div>

        <form onSubmit={handleSubmit} className="flex flex-col w-full items-center gap-3">
          <div className="flex flex-col gap-1 w-full">
            <label className="font-normal text-neutral-1000 text-md">Nama Lengkap</label>
            <input
              name="full_name"
              type="text"
              placeholder="Masukkan nama lengkap Anda"
              required
              className="w-full py-3 rounded-md bg-transparent border border-neutral-1000 px-3 text-sm text-[#71717A]"
            />
          </div>

          <div className="flex flex-col gap-1 w-full">
            <label className="font-normal text-neutral-1000 text-md">Email</label>
            <input
              name="email"
              type="email"
              placeholder="Masukkan email Anda"
              required
              className="w-full py-3 rounded-md bg-transparent border border-neutral-1000 px-3 text-sm text-[#71717A]"
            />
          </div>

          <div className="flex flex-col gap-1 w-full">
            <label className="font-normal text-neutral-1000 text-md">Password</label>
            <input
              name="password"
              type="password"
              placeholder="Masukkan kata sandi Anda"
              required
              minLength={8}
              className="w-full py-3 rounded-md bg-transparent border border-neutral-1000 px-3 text-sm text-[#71717A]"
            />
          </div>

          <button
            type="submit"
            className="w-10/12 sm:w-full py-2 px-3 rounded-lg bg-primary-600 text-sm font-medium text-white transition-colors mt-5"
          >
            Daftar
          </button>

          <p className="text-neutral-1000 text-xs font-normal text-center">
            Sudah punya akun?{" "}
            <a href="/login" className="text-neutral-900 font-normal hover:underline">
              Masuk
            </a>
          </p>
        </form>
      </div>
    </main>
  );
}
