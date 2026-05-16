import Image from "next/image";

export default function Signup() {
  //comment untuk commit
  return (
    <main className="flex w-full min-h-screen items-center justify-center bg-gradient-to-t from-primary-900 to-primary-1000 text-black font-(Plus Jakarta Sans) px-4 py-16">

      <button
        type="button"
        className="flex items-center gap-1 absolute left-10 top-15 sm:left-8 sm:top-8 rounded-lg bg-primary-600 px-3 py-2.5 sm:px-3 sm:py-2.5 text-sm font-medium text-white transition-colors"
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
    
      <div className="w-full max-w-full sm:max-w-md bg-neutral-100 rounded-3xl px-5 py-6 sm:px-5 sm:py-6 flex flex-col items-center gap-10">

        <div className="flex flex-col items-center gap-3">
          <Image
            src="/images/omahti-dark_2x.webp"
            alt="Omahti Logo"
            width={72}
            height={72}
            className="object-contain sm:scale-140"
          />
          <h1 className="text-neutral-1000 font-semibold text-3xl sm:text-3xl text-center">
            Lupa Passwordmu?
          </h1>
        </div>

        <div className="flex flex-col w-full items-center gap-3">
          <div className="flex flex-col gap-1 w-full">
            <label className="font-normal text-neutral-1000 text-md">Email</label>
            <input
              type="email"
              placeholder="Masukkan email Anda"
              className="w-full py-3 rounded-md bg-transparent border border-neutral-1000 px-3 text-sm text-[#71717A]"
            />
          </div>

          <button
            type="submit"
            className="w-10/12 sm:w-full py-2 px-3 rounded-lg bg-primary-600 text-sm font-medium text-white transition-colors mt-5"
          >
            Masuk
          </button>
        </div>
      </div>
    </main>
  );
}