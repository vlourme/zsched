import { useMemo, useState } from "react";
import { Outlet, useMatches } from "react-router";
import { AppSidebar } from "~/components/sidebar";
import { SidebarProvider, SidebarTrigger } from "~/components/ui/sidebar";
import { cn } from "~/lib/utils";

export default function Layout() {
  const matches = useMatches();
  const handle = useMemo(
    () =>
      matches[matches.length - 1]?.handle as
        | { title: () => string }
        | undefined,
    [matches]
  );
  const [open, setOpen] = useState(true);

  return (
    <SidebarProvider open={open} onOpenChange={setOpen}>
      <AppSidebar />
      <div
        className={cn(
          "w-[calc(100%-var(--sidebar-width))] flex-1 pr-2 py-2 bg-sidebar flex flex-col",
          !open && "pl-2"
        )}
      >
        <main className="flex flex-col flex-1 bg-background rounded-md max-w-full">
          <div className="flex items-center divide-x border-b p-3">
            <div className="pr-1">
              <SidebarTrigger />
            </div>
            <h1 className="text-lg font-bold pl-3">
              {handle?.title() || "Zsched"}
            </h1>
          </div>

          <Outlet />
        </main>
      </div>
    </SidebarProvider>
  );
}
