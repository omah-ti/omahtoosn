import React from 'react'

const KonfirmasiCard = () => {
    return (
        <div className='bg-[#DCE5F9] box-border flex flex-col w-full w-max-[398px] h-[190px] md:w-[264px] md:h-[294px] rounded-[20px] shadow-[0_2px_10px_0_rgba(0,0,0,0.25)] justify-center p-6 gap-3'>
            <p className='--font-plus size text-[26px] font-bold text-center'>Konfirmasi Test</p>
            <hr className='w-full h-[1px] bg-neutral-100'/>
            <div className='flex-col gap-px'>
                <p className='text-base font-normal'>Waktu Tes</p>
                <div className='flex justify-between'>
                    <p className='text-base font-bold'>06/03/2026</p>
                    <p className='text-base font-bold'>09:00 WIB</p>
                </div>
            </div>
            <div className='flex flex-row justify-between md:flex-col gap-3'>
                <div className='flex-col gap-px'>
                    <p className='text-base font-normal'>Durasi</p>
                    <p className='text-base font-bold'>150 menit</p>
                </div>
                <div className='flex-col gap-px'>
                    <p className='text-base font-normal text-right md:text-left'>Jumlah</p>
                    <p className='text-base font-bold'>30-50 Soal</p>
                </div>
            </div>
        </div>
    )
}

export default KonfirmasiCard