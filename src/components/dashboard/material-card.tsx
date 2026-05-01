type Props = {
  title: string;
  desc: string;
  img: string;
};

export default function MaterialCard({ title, desc, img }: Props) {
  return (
    <div className="bg-[var(--primary-100)] rounded-lg p-4 shadow-sm w-full">
      <div className="w-full aspect-[2/1] rounded mb-4 overflow-hidden">
        <img src={img} alt={title} className="w-full h-full object-cover" />
      </div>

      <h3 className="font-semibold border-b pb-2 mb-2">{title}</h3>

      <p className="text-sm text-neutral-1000">{desc}</p>
    </div>
  );
}
