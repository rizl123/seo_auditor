import { MainClientContainer } from "@/components/MainClientContainer";

export default async function Home() {
  return (
    <div className="min-h-screen bg-background">
      <main className="max-w-5xl mx-auto px-6 pb-20">
        <MainClientContainer />
      </main>
    </div>
  );
}
