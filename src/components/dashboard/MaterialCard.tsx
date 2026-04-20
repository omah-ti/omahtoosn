type Props = {
  title: string;
  desc: string;
};

export default function MaterialCard({ title, desc }: Props) {
  return (
    <div className="bg-white rounded-lg p-4 shadow-sm">
      <div className="h-32 bg-gray-200 rounded mb-4" />

      <h3 className="font-semibold border-b pb-2 mb-2">{title}</h3>

      <p className="text-sm text-gray-600">{desc}</p>
    </div>
  );
}
