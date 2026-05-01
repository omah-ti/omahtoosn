import KonfirmasiCard from '@/components/tryout/konfirmasi-card'
import PetunjukCard from '@/components/tryout/petunjuk-card'
import React from 'react'

const Petunjuk = () => {
    return (
    <div className='flex flex-col md:flex-row md:items-start justify-center items-center gap-4 pt-8 px-4 md:pt-[90px]'>
        <KonfirmasiCard />
        <PetunjukCard />
    </div>
    )
}

export default Petunjuk