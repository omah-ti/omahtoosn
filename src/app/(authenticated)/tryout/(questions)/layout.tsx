import Navbar from "@/components/authenticated/navbar/navbar";
import Footer from "@/components/ui/footer";

export default function AuthenticatedLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <>
      {children}
    </>
  );
}
