import Image from "next/image";

export default function Signup() {
	return (
    <main className="flex w-full min-h-screen items-center justify-center bg-white text-black font-(Plus Jakarta Sans)">
			<button
				type="button"
				className=" flex absolute justify-center items-center left-10 top-10 rounded-md bg-[#F18519] px-10 py-4.5 text-sm font-medium text-white transition-colors">
				<Image
						src="/images/arrow.png"
						alt="Omahti Logo"
						width={13}
						height={13}
						className="object-cover absolute transform items-center justify-start left-1/4 -translate-x-1/2"
					/>
				<p className="flex absolute transform right-1/6">Back</p>
			</button>
			<div className="py-50 px-50 bg-[#1E293B] overflow-hidden rounded-xl relative flex items-center justify-center">
				<div className="absolute items-end justify-center translate-x-0 left-0 translate-y-0 top-14 flex w-full h-fit">
					<Image
						src="/images/omahti-dark 2.png"
						alt="Omahti Logo"
						width={80}
						height={80}
						className="object-cover absolute"
					/>
					<h1 className="absolute text-white font-semibold text-2xl translate-y-10">Daftar Akun Peserta</h1>
				</div>
				<div className="flex-col gap-0.5 items-center justify-start absolute flex w-fit h-fit translate-y-0 top-32">
					<p className="font-extralight text-white text-xs flex mr-auto">Nama Lengkap</p>
					<input
						type="text"
						placeholder="Masukkan nama lengkap"
						className="w-90 h-8 rounded-sm bg-transparent border-1 border-[#E4E4E7] px-2 py-2 text-xs text-[#71717A] mt-0.5"
					/>
					<p className="font-extralight text-white text-xs flex mr-auto mt-2">Email</p>
					<input
						type="email"
						placeholder="Masukkan email"
						className="w-90 h-8 rounded-sm bg-transparent border-1 border-[#E4E4E7] px-2 py-2 text-xs text-[#71717A] mt-0.5"
					/>
					<p className="font-extralight text-white text-xs flex mr-auto mt-2">Password</p>
					<input
						type="password"
						placeholder="Masukkan password"
						className="w-90 h-8 rounded-sm bg-transparent border-1 border-[#E4E4E7] px-2 py-2 text-xs text-[#71717A] mt-0.5"
					/>
					<button
						type="submit"
						className="items-center justify-center w-90 h-8 rounded-sm bg-[#F18519] px-3 text-xs font-medium text-white transition-colors mt-5">
						Daftar
					</button>
					<p className="text-white text-xs font-normal">
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