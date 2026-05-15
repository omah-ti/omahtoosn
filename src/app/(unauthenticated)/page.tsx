import Hero from "@/modules/landing/hero";
import Statistics from "@/modules/landing/statistics";
import Explore from "@/modules/landing/explore";
import Kemampuan from "@/modules/landing/kemampuan";
import Alumni from "@/modules/landing/alumni";
import About from "@/modules/landing/about";
import CTA from "@/modules/landing/cta";

export default function Home() {
  return (
    <>
      <Hero />
      <Statistics />
      <Explore />
      <Kemampuan />
      <Alumni />
      <About />
      <CTA />
    </>
  );
}
