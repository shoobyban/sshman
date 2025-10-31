import { useQuery } from "@tanstack/react-query";
import { Server, Users, FolderTree, Plus } from "lucide-react";
import { Link } from "react-router-dom";
import { StatCard } from "@/components/StatCard";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { getHosts, getUsers, getGroups } from "@/lib/api";

export default function Dashboard() {
  const { data: hosts = [] } = useQuery({
    queryKey: ["hosts"],
    queryFn: getHosts,
  });

  const { data: users = [] } = useQuery({
    queryKey: ["users"],
    queryFn: getUsers,
  });

  const { data: groups = [] } = useQuery({
    queryKey: ["groups"],
    queryFn: getGroups,
  });

  return (
    <div className="min-h-screen bg-background">
      <header className="border-b border-border bg-card">
        <div className="px-8 py-6">
          <h1 className="text-3xl font-bold text-foreground">Dashboard</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Overview of your SSH infrastructure
          </p>
        </div>
      </header>

      <div className="p-8">
        {/* Stats */}
        <div className="grid gap-6 md:grid-cols-3">
          <StatCard
            title="Total Hosts"
            value={hosts.length}
            icon={Server}
            description="Managed SSH hosts"
          />
          <StatCard
            title="Total Users"
            value={users.length}
            icon={Users}
            description="Registered users"
          />
          <StatCard
            title="Total Groups"
            value={groups.length}
            icon={FolderTree}
            description="Access groups"
          />
        </div>

        {/* Quick Actions */}
        <Card className="mt-8 bg-gradient-card shadow-md">
          <div className="p-6">
            <h2 className="text-lg font-semibold text-foreground">Quick Actions</h2>
            <div className="mt-4 flex flex-wrap gap-3">
              <Button asChild>
                <Link to="/hosts/new">
                  <Plus className="mr-2 h-4 w-4" />
                  Add Host
                </Link>
              </Button>
              <Button asChild variant="secondary">
                <Link to="/users/new">
                  <Plus className="mr-2 h-4 w-4" />
                  Add User
                </Link>
              </Button>
              <Button asChild variant="secondary">
                <Link to="/groups/new">
                  <Plus className="mr-2 h-4 w-4" />
                  Add Group
                </Link>
              </Button>
            </div>
          </div>
        </Card>

        {/* Recent Activity Preview */}
        <Card className="mt-8 bg-gradient-card shadow-md">
          <div className="p-6">
            <div className="flex items-center justify-between">
              <h2 className="text-lg font-semibold text-foreground">Recent Hosts</h2>
              <Button asChild variant="ghost" size="sm">
                <Link to="/hosts">View All</Link>
              </Button>
            </div>
            <div className="mt-4 space-y-3">
              {hosts.slice(0, 5).map((host) => (
                <div
                  key={host.alias}
                  className="flex items-center justify-between rounded-lg border border-border bg-card p-3"
                >
                  <div>
                    <p className="font-medium text-foreground">{host.alias}</p>
                    <p className="text-sm text-muted-foreground">
                      {host.user}@{host.host}
                    </p>
                  </div>
                  <div className="flex gap-2">
                    {host.groups.map((group) => (
                      <span
                        key={group}
                        className="rounded-full bg-accent/10 px-2 py-1 text-xs text-accent"
                      >
                        {group}
                      </span>
                    ))}
                  </div>
                </div>
              ))}
              {hosts.length === 0 && (
                <p className="text-center text-sm text-muted-foreground">
                  No hosts yet. Add your first host to get started!
                </p>
              )}
            </div>
          </div>
        </Card>
      </div>
    </div>
  );
}
