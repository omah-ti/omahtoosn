import Navbar from "@/components/authenticated/navbar/navbar";

export default function ResultLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex flex-col min-h-screen bg-white">
      <Navbar />
      <main className="flex-1 pt-16">
        {children}
      </main>
    </div>
  );
}
