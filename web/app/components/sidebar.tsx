import { HomeIcon, WrenchIcon } from "lucide-react";
import { Link, useLocation } from "react-router";
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
  },
  {
    label: "Tasks",
    icon: WrenchIcon,
    to: "/tasks",
  },
];

export function AppSidebar() {
  const location = useLocation();
  const isActive = (to: string) => {
    return location.pathname === to;
  };
  return (
    <Sidebar>
      <SidebarHeader>
        <div className="flex items-center gap-3 p-2">
          <div className="size-10 rounded-full bg-border flex items-center justify-center font-bold">
            <span>🔧</span>
          </div>
          <div className="leading-3">
            <h1 className="text-lg font-bold">Zsched</h1>
            <p className="text-xs text-muted-foreground">Task Scheduler</p>
          </div>
        </div>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              {sidebarItems.map((item) => (
                <SidebarMenuItem key={item.to}>
                  <SidebarMenuButton isActive={isActive(item.to)} asChild>
                    <Link to={item.to}>
                      <item.icon className="size-4" />
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
