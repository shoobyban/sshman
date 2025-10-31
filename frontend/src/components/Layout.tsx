import { Link, useLocation } from "react-router-dom";
import { Server, Users, FolderTree, LayoutDashboard, FileText, Terminal } from "lucide-react";
import { cn } from "@/lib/utils";

const navigation = [
  { name: "Dashboard", href: "/", icon: LayoutDashboard },
  { name: "Hosts", href: "/hosts", icon: Server },
  { name: "Users", href: "/users", icon: Users },
  { name: "Groups", href: "/groups", icon: FolderTree },
  { name: "Logs", href: "/logs", icon: FileText },
];

interface LayoutProps {
  children: React.ReactNode;
}

export function Layout({ children }: LayoutProps) {
  const location = useLocation();

  return (
    <div className="flex min-h-screen w-full bg-background">
      {/* Sidebar */}
      <aside className="w-64 border-r border-sidebar-border bg-sidebar">
        <div className="flex h-16 items-center gap-2 border-b border-sidebar-border px-6">
          <Terminal className="h-6 w-6 text-sidebar-primary" />
          <h1 className="text-xl font-bold text-sidebar-foreground">SSHMan</h1>
        </div>
        <nav className="space-y-1 p-4">
          {navigation.map((item) => {
            const isActive = location.pathname === item.href || 
              (item.href !== "/" && location.pathname.startsWith(item.href));
            
            return (
              <Link
                key={item.name}
                to={item.href}
                className={cn(
                  "flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors",
                  isActive
                    ? "bg-sidebar-accent text-sidebar-primary"
                    : "text-sidebar-foreground hover:bg-sidebar-accent/50"
                )}
              >
                <item.icon className="h-4 w-4" />
                {item.name}
              </Link>
            );
          })}
        </nav>
      </aside>

      {/* Main content */}
      <main className="flex-1">
        {children}
      </main>
    </div>
  );
}
