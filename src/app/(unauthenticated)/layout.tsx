import LandingNavbar from "@/components/unauthenticated/navbar";
import Footer from "@/components/ui/footer";

export default function UnauthenticatedLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <>
      <LandingNavbar />
      {children}
      <Footer />
    </>
  );
}
