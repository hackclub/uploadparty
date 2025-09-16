import Navbar from "../../components/shared/Navbar";

export const metadata = {
  title: "App | UploadParty",
  description: "Authenticated area of UploadParty",
};

export default function AppLayout({ children }) {
  return (
    <section>
      <Navbar />
      {children}
    </section>
  );
}
