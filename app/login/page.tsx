import Image from "next/image";

export default function Signup() {
  return (
    <main className="flex w-full min-h-screen items-center justify-center bg-white text-black font-(Plus Jakarta Sans) px-3 py-16">

      <button
        type="button"
        className="flex items-center gap-2 absolute left-10 top-15 sm:left-8 sm:top-8 rounded-lg bg-[#F18519] px-4 py-2.5 sm:px-6 sm:py-3 text-sm font-medium text-white transition-colors"
      >
        <Image
          src="/images/arrow_1x.webp"
          alt="Back arrow"
          width={13}
          height={13}
          className="object-contain"
        />
        <span>Kembali</span>
      </button>

      <div className="w-full max-w-full sm:max-w-md bg-[#1E293B] rounded-3xl px-5 py-6 sm:px-5 sm:py-6 flex flex-col items-center gap-10">

        <div className="flex flex-col items-center gap-3">
          <Image
            src="/images/omahti-dark 2.png"
            alt="Omahti Logo"
            width={72}
            height={72}
            className="object-contain"
          />
          <h1 className="text-white font-semibold text-3xl sm:text-3xl text-center">
            Masuk ke Akun
          </h1>
        </div>

        <div className="flex flex-col w-full items-center gap-3">

          <div className="flex flex-col gap-1 w-full">
            <label className="font-extralight text-white text-md">Email</label>
            <input
              type="email"
              placeholder="Masukkan email Anda"
              className="w-full py-3 rounded-md bg-transparent border border-[#E4E4E7] px-3 text-sm text-[#71717A]"
            />
          </div>

          <div className="flex flex-col gap-1 w-full">
            <div className="flex justify-between items-center w-full">
              <label className="font-extralight text-white text-md">Password</label>
              <button className="text-white text-sm font-bold hover:underline">Lupa Password?</button>
            </div>
            <input
              type="password"
              placeholder="Masukkan kata sandi Anda"
              className="w-full py-3 rounded-md bg-transparent border border-[#E4E4E7] px-3 text-sm text-[#71717A]"
            />
          </div>

          <button
            type="submit"
            className="w-10/12 sm:w-full py-2 px-3 rounded-lg bg-[#F18519] text-sm font-medium text-white transition-colors mt-5"
          >
            Daftar
          </button>

          <p className="text-white text-xs font-normal text-center">
            Sudah punya akun?{" "}
            <a href="/Login" className="text-[#F18519] font-medium">
              Masuk
            </a>
          </p>
        </div>
      </div>
    </main>
  );
}