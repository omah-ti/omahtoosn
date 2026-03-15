import React from 'react'

const Petunjuk = () => {
    return (
    <div className='flex justify-center gap-[24px] pt-[90px]'>
        <div className='bg-[#fff] box-border flex flex-col w-[264px] h-[294px] border border-[#A4A4A4] rounded-[20px] border-solid justify-center p-[24px] gap-[12px]'>
            <p className='--font-plus size text-[26px] font-bold text-center'>Konfirmasi Test</p>
            <hr className='w-[216px] h-[1px] bg-black'/>
            <div className='flex-col gap-px'>
                <p className='text-[16px] font-normal'>Waktu Tes</p>
                <div className='flex justify-between'>
                    <p className='text-[16px] font-bold'>06/03/2026</p>
                    <p className='text-[16px] font-bold'>09:00</p>
                </div>
            </div>
            <div className='flex-col gap-px'>
                <p className='text-[16px] font-normal'>Durasi</p>
                <p className='text-[16px] font-bold'>150 menit</p>
            </div>
            <div className='flex-col gap-px'>
                <p className='text-[16px] font-normal'>Jumlah</p>
                <p className='text-[16px] font-bold'>30-50 Soal</p>
            </div>
        </div>
        <div className='bg-[#fff] flex flex-col box-border w-[840px] h-[426px] border border-[#A4A4A4] rounded-[20px] p-[24px] gap-[44px]'>
            <div className='flex flex-col gap-[12px]'>
                <p className='size text-[26px] font-bold text-center'>Petunjuk Try Out OSN-K Informatika</p>
                <hr className='w-[792px] h-[1px] bg-black'/>
                <p className='text-[16px] font-normal'><b>1. Waktu pengerjaan adalah 150 menit. </b> Setelah waktu habis, tes akan ditutup secara otomatis.</p>
                <p className='text-[16px] font-normal'><b>2. Selama waktu masih tersedia, </b> Anda dapat membuka kembali soal mana pun dan mengubah jawaban.</p>
                <p className='text-[16px] font-normal pl-6 -indent-5'><b>3. Gunakan Daftar Soal </b>untuk berpindah ke nomor tertentu, atau gunakan tombol Sebelumnya dan Berikutnya untuk navigasi.</p>
                <p className='text-[16px] font-normal'>4. Gunakan fitur <b>Ragu-ragu</b> untuk menandai soal yang ingin Anda tinjau kembali.</p>
                <p className='text-[16px] font-normal'>5. Jika sudah yakin dengan jawaban Anda, tekan tombol <b>Kirim Sekarang</b> untuk mengirimkan tes.</p>
                <p className='text-[16px] font-normal'>6. Tes ini hanya dapat dikerjakan <b>satu kali.</b> Tes tidak dapat diulang setelah dimulai.</p>
            </div>
            <div>
                <button className=' w-[792px] h-[40px] border-[#535353] bg-[#535353] border-[1px] p-[8px] rounded-[8px] text-white flex justify-center shadow-xs cursor-pointer'>Mulai</button>
            </div>
        </div>
    </div>
    )
}

export default Petunjuk