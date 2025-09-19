// 404 handler for the authenticated (app) area
// Behavior: show a clear message that the page does not exist
// Route scope: any unknown route within the (app) group (rendered without the segment name)

import Link from "next/link";

export default function NotFound() {
  return (
    <main>
      <h1>Page not found</h1>
      <p>The page you’re looking for doesn’t exist.</p>
      <nav>
        <ul>
          <li>
            <Link href="/dashboard">Go to dashboard</Link>
          </li>
          <li>
            <Link href="/">Return to home</Link>
          </li>
        </ul>
      </nav>
    </main>
  );
}
