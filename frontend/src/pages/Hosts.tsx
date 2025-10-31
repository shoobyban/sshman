import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Link } from "react-router-dom";
import { Plus, Pencil, Trash2, RefreshCw, Search } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { getHosts, deleteHost, syncHosts } from "@/lib/api";
import { toast } from "@/hooks/use-toast";
import { ApiError } from "@/lib/errors";

export default function Hosts() {
  const [searchTerm, setSearchTerm] = useState("");
  const [deleteId, setDeleteId] = useState<string | null>(null);
  const queryClient = useQueryClient();

  const { data: hosts = [], isLoading } = useQuery({
    queryKey: ["hosts"],
    queryFn: getHosts,
  });

  const deleteMutation = useMutation({
    mutationFn: deleteHost,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["hosts"] });
      toast({
        title: "Host deleted",
        description: "The host has been successfully removed.",
      });
      setDeleteId(null);
    },
    onError: (error) => {
      let description = "Failed to delete host. Please try again.";
      if (error instanceof ApiError && error.details) {
        description = error.details;
      } else if (error instanceof Error) {
        description = error.message;
      }
      toast({
        title: "Error",
        description,
        variant: "destructive",
      });
    },
  });

  const syncMutation = useMutation({
    mutationFn: syncHosts,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["hosts"] });
      toast({
        title: "Sync complete",
        description: "All hosts have been synchronized.",
      });
    },
    onError: (error) => {
      let description = "Could not sync hosts. Please try again.";
      if (error instanceof ApiError && error.details) {
        description = error.details;
      } else if (error instanceof Error) {
        description = error.message;
      }
      toast({
        title: "Sync failed",
        description,
        variant: "destructive",
      });
    },
  });

  const filteredHosts = hosts.filter(
    (host) =>
      host.alias.toLowerCase().includes(searchTerm.toLowerCase()) ||
      host.host.toLowerCase().includes(searchTerm.toLowerCase()) ||
      host.user.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div className="min-h-screen bg-background">
      <header className="border-b border-border bg-card">
        <div className="flex items-center justify-between px-8 py-6">
          <div>
            <h1 className="text-3xl font-bold text-foreground">Hosts</h1>
            <p className="mt-1 text-sm text-muted-foreground">
              Manage your SSH host configurations
            </p>
          </div>
          <div className="flex gap-3">
            <Button
              variant="secondary"
              onClick={() => syncMutation.mutate()}
              disabled={syncMutation.isPending}
            >
              <RefreshCw className={`mr-2 h-4 w-4 ${syncMutation.isPending ? "animate-spin" : ""}`} />
              Sync All
            </Button>
            <Button asChild>
              <Link to="/hosts/new">
                <Plus className="mr-2 h-4 w-4" />
                Add Host
              </Link>
            </Button>
          </div>
        </div>
      </header>

      <div className="p-8">
        <Card className="bg-gradient-card shadow-md">
          <div className="p-6">
            <div className="mb-6 flex items-center gap-4">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                <Input
                  placeholder="Search hosts by alias, host, or user..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-9"
                />
              </div>
            </div>

            {isLoading ? (
              <div className="py-8 text-center text-muted-foreground">Loading hosts...</div>
            ) : filteredHosts.length === 0 ? (
              <div className="py-8 text-center">
                <p className="text-muted-foreground">
                  {searchTerm ? "No hosts match your search." : "No hosts yet. Add your first host!"}
                </p>
              </div>
            ) : (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Alias</TableHead>
                    <TableHead>Connection</TableHead>
                    <TableHead>Groups</TableHead>
                    <TableHead>Key File</TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filteredHosts.map((host) => (
                    <TableRow key={host.alias}>
                      <TableCell className="font-medium">{host.alias}</TableCell>
                      <TableCell>
                        <code className="rounded bg-muted px-2 py-1 text-sm">
                          {host.user}@{host.host}
                        </code>
                      </TableCell>
                      <TableCell>
                        <div className="flex flex-wrap gap-1">
                          {host.groups.map((group) => (
                            <span
                              key={group}
                              className="rounded-full bg-accent/10 px-2 py-1 text-xs text-accent"
                            >
                              {group}
                            </span>
                          ))}
                        </div>
                      </TableCell>
                      <TableCell>
                        <code className="text-xs text-muted-foreground">{host.key}</code>
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="flex justify-end gap-2">
                          <Button asChild variant="ghost" size="sm">
                            <Link to={`/hosts/${host.alias}/edit`}>
                              <Pencil className="h-4 w-4" />
                            </Link>
                          </Button>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => setDeleteId(host.alias)}
                          >
                            <Trash2 className="h-4 w-4 text-destructive" />
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            )}
          </div>
        </Card>
      </div>

      <AlertDialog open={!!deleteId} onOpenChange={() => setDeleteId(null)}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Host</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete this host? This action cannot be undone.
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={() => deleteId && deleteMutation.mutate(deleteId)}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
