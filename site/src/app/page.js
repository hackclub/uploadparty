// Root route for the public (pre‑auth) experience.
// The actual component lives in ./(public)/page.
// Keeping this file as a simple re-export makes the separation explicit.
import PublicLanding from "./(public)/page";
import { cookies } from 'next/headers';
import { redirect } from 'next/navigation';

export default async function RootPage() {
  // Check if authentication bypass is enabled (development only)
  const bypassAuth = process.env.BYPASS_AUTH === 'true';
  
  if (!bypassAuth) {
    // Check if user is already logged in and redirect to app
    const cookieStore = await cookies();
    const session = cookieStore.get('auth_session');
    
    if (session && session.value === 'logged_in') {
      redirect('/app');
    }
  }
  
  return <PublicLanding />;
}
