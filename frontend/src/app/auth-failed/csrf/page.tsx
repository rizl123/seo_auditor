import { ShieldX } from "lucide-react";
import { ErrorCard } from "@/components/ErrorCard";

export default function CSRFErrorPage() {
  return (
    <ErrorCard
      title="Security check failed"
      description="The request state is invalid or has expired."
      icon={ShieldX}
      variant="amber"
      actionLabel="Try Again"
      actionHref="/login"
    />
  );
}
