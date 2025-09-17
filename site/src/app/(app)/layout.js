import SubmissionForm from "../../components/shared/submissionForm";

export const metadata = {
  title: "App | UploadParty",
  description: "Authenticated area of UploadParty",
};

export default function AppLayout({ children }) {
  return (
    <section>
      <SubmissionForm />
      {children}
    </section>
  );
}
