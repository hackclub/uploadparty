// Overview section (public pre-auth)
// Minimal content; server component. Keep styling simple per guidelines.

import LoginButton from '../components/LoginButton';

export default function Overview() {
  return (
    <section className="max-w-4xl mx-auto px-6 py-12">
      <div className="text-center mb-8">
        <h2 className="text-3xl font-bold mb-4">What we're building</h2>
        <p className="text-lg mb-4">
          This month, I'm aiming to ship an early, public slice of UploadParty to help
          contribute toward the $250 goal—even if the full platform isn't ready by
          month's end.
        </p>
        <p className="text-lg mb-8">
          The core question we're exploring is: how can music be "open‑source" in
          practice? Think teen‑built virtual instruments, effects, and shared building
          blocks that anyone can learn from, remix, and improve.
        </p>
      </div>
      
      <div className="bg-gray-50 rounded-lg p-8 text-center">
        <LoginButton />
      </div>
    </section>
  );
}
