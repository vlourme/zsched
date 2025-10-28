import { useMemo } from "react";
import { Outlet, useMatches } from "react-router";
import { AppSidebar } from "~/components/sidebar";
import { SidebarProvider, SidebarTrigger } from "~/components/ui/sidebar";

export default function Layout() {
  const matches = useMatches();
  const handle = useMemo(
    () =>
      matches[matches.length - 1]?.handle as
        | { title: () => string }
        | undefined,
    [matches]
  );

  return (
    <SidebarProvider>
      <AppSidebar />
      <div className="p-2 bg-sidebar flex flex-col w-full">
        <main className="flex flex-col flex-1 bg-background rounded-md">
          <div className="flex items-center divide-x border-b p-3">
            <div className="pr-1">
              <SidebarTrigger />
            </div>
            <h1 className="text-lg font-bold pl-3">
              {handle?.title() || "Zsched"}
            </h1>
          </div>

          <div>
            <Outlet />
          </div>
        </main>
      </div>
    </SidebarProvider>
  );
}
