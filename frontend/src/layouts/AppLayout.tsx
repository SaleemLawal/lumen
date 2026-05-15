import { Outlet, useLocation } from "react-router-dom";
import {
  SidebarProvider,
  SidebarInset,
  SidebarTrigger,
} from "@/components/ui/sidebar";
import SideBar from "@/components/ui/side-bar";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Bell, Plus, Search } from "lucide-react";

const titles: Record<string, string> = {
  "/dashboard": "Dashboard",
  "/transactions": "Transactions",
  "/budgets": "Budgets",
  "/accounts": "Accounts",
  "/investments": "Investments",
  "/connect": "Connect",
};

export default function AppLayout() {
  const { pathname } = useLocation();
  const title = titles[pathname] ?? "Lumen";

  return (
    <SidebarProvider>
        <SideBar />
        <SidebarInset>
          <header className="sticky top-0 z-30 flex h-16 items-center gap-3 rounded-t-xl border-b border-border/60 bg-background/70 px-4 backdrop-blur-xl md:px-6">
            <SidebarTrigger />
            <h1 className="text-lg font-semibold tracking-tight">{title}</h1>
            <div className="ml-auto flex items-center gap-2">
              <div className="relative hidden md:block">
                <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                <Input placeholder="Search transactions…" className="h-9 w-64 pl-9" />
              </div>
              <Button variant="ghost" size="icon" aria-label="Notifications">
                <Bell className="h-4 w-4" />
              </Button>
              <Button size="sm" className="gap-1">
                <Plus className="h-4 w-4" /> Add
              </Button>
            </div>
          </header>
          <main className="flex-1 px-4 py-6 md:px-8 md:py-8">
            <Outlet />
          </main>
        </SidebarInset>
    </SidebarProvider>
  );
}
