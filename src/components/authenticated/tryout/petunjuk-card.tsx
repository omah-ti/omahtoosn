"use client"

import React from 'react'
import Button from "@/components/ui/button";
import Link from 'next/link';
const PetunjukCard = () => {
	const instruksiList = [
		"<b>1. Waktu pengerjaan adalah 150 menit. </b> Setelah waktu habis, tes akan ditutup secara otomatis.",
		"<b>2. Selama waktu masih tersedia, </b> Anda dapat membuka kembali soal mana pun dan mengubah jawaban.",
		"<b>3. Gunakan Daftar Soal </b>untuk berpindah ke nomor tertentu, atau gunakan tombol Sebelumnya dan Berikutnya untuk navigasi.",
		"<b>4. Gunakan fitur Ragu-ragu </b>untuk menandai soal yang ingin Anda tinjau kembali.",
		"<b>5. Jika sudah yakin dengan jawaban Anda, </b> tekan tombol Kirim Sekarang untuk mengirimkan tes.",
		"<b>6. Tes ini hanya dapat dikerjakan satu kali. </b> Tes tidak dapat diulang setelah dimulai."
	]
	return (
		<div className='bg-primary-background flex flex-col box-border w-full w-max-[398px] h-auto md:w-[840px] rounded-[20px] p-6 gap-[44px]'>
			<div className='flex flex-col gap-3'>
				<p className='size text-[26px] font-bold text-center'>Petunjuk Try Out OSN-K Informatika</p>
				<hr className=' w-full h-px bg-black' />
				{instruksiList.map((instruksi, index) => (
					<p
						key={index} className='text-base font-normal'
						dangerouslySetInnerHTML={{ __html: instruksi }}
					/>
				))}
			</div>
			<div>
				<Link href="/tryout">
					<Button variant='primary' className='w-full bg-primary-600 cursor-pointer'>Mulai</Button>
				</Link>
			</div>
		</div>
	)
}

export default PetunjukCard