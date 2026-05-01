import Image from "next/image";
import Link from "next/link";

type LogoProps = {
  href?: string;
  className?: string;
  logoSize?: string;
  textSize?: string;
};

export default function Logo({
  href = "/",
  className = "",
  logoSize = "w-11 h-11 sm:w-11 sm:h-10",
  textSize = "h-9 w-auto",
}: LogoProps) {
  const content = (
    <div className={`flex items-center gap-2 sm:gap-3 ${className}`}>
      <Image
        src="/navbar/logo.webp"
        alt="OmahTI"
        className={`${logoSize} object-contain`}
        width={100}
        height={100}
      />

      <div className="flex flex-col justify-center">
        <p className="font-bold text-base leading-none">Omah<span className="text-[#F0861A]">TOOSN</span></p>
        <p className="text-[10px]">Computer Science UGM</p>
      </div>
    </div>
  );

  if (href) {
    return <Link href={href}>{content}</Link>;
  }

  return content;
}
