import type { ReactNode } from "react";

export function Card({
  children,
  className = "",
}: {
  children: ReactNode;
  className?: string;
}) {
  return (
    <div
      className={`bg-card border border-border-custom rounded-3xl shadow-sm overflow-hidden ${className}`}
    >
      {children}
    </div>
  );
}
