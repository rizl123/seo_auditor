import { ShieldAlert } from "lucide-react";
import { ErrorCard } from "@/components/ErrorCard";

export const metadata = {
  title: "Unexpected Error",
  description: "An unknown error occurred during the authentication process.",
};

export default function UnknownErrorPage() {
  return (
    <ErrorCard
      title="Something went wrong"
      description="We encountered an unexpected error while processing your request. This could be due to a temporary connection issue or a system glitch."
      icon={ShieldAlert}
      variant="rose"
      actionLabel="Try Again"
      actionHref="/login"
    />
  );
}
