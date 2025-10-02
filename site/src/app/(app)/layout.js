import SubmissionForm from "../components/shared/submissionForm";
import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';

export const metadata = {
  title: "App | UploadParty",
  description: "Authenticated area of UploadParty",
};

export default async function AppLayout({ children }) {
  // Check if authentication bypass is enabled (development only)
  const bypassAuth = process.env.BYPASS_AUTH === 'true';
  
  if (!bypassAuth) {
    const cookieStore = await cookies();
    const session = cookieStore.get('auth_session');

    if (!session || session.value !== 'logged_in') {
      // Server-side redirect preserves HttpOnly cookie security and avoids client flashes
      redirect('/api/auth/login?returnTo=/app');
    }
  }

  return (
    <main>
      <section>
        <SubmissionForm />
        {children}
      </section>
    </main>
  );
}
