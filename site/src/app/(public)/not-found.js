// 404 handler for the public (pre-auth) area
// Behavior: route users back to the pre-auth landing page
// Route scope: any unknown route within the (public) group

import { redirect } from "next/navigation";

export default function NotFound() {
  // Immediately route back to the public landing page
  redirect("/");
}
