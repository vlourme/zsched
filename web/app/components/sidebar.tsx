import { HomeIcon, WrenchIcon } from "lucide-react";
import { Link, useMatches } from "react-router";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "./ui/sidebar";

const sidebarItems = [
  {
    label: "Overview",
    icon: HomeIcon,
    to: "/",
    group: "overview",
  },
  {
    label: "Tasks",
    icon: WrenchIcon,
    to: "/tasks",
    group: "tasks",
  },
];

export function AppSidebar() {
  const matches = useMatches();
  const lastMatch = matches[matches.length - 1];
  const isActive = (group: string) => {
    if (!lastMatch?.handle) return false;
    return (lastMatch?.handle as { group: string }).group === group;
  };
  return (
    <Sidebar className="group-data-[side=left]:border-r-0">
      <SidebarHeader>
        <div className="flex items-center gap-3 p-2">
          <div className="size-10 bg-border rounded-full flex items-center justify-center font-bold">
            <span>⚙️</span>
          </div>
          <div className="leading-3">
            <h1 className="text-lg font-bold">Zsched</h1>
            <p className="text-sm text-muted-foreground">Task Scheduler</p>
          </div>
        </div>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              {sidebarItems.map((item) => (
                <SidebarMenuItem key={item.to}>
                  <SidebarMenuButton
                    className="data-[active=true]:bg-primary px-4 gap-3 data-[active=true]:text-primary-foreground"
                    isActive={isActive(item.group)}
                    asChild
                  >
                    <Link to={item.to}>
                      <item.icon />
                      <span>{item.label}</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter />
    </Sidebar>
  );
}
