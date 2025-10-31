import { useState } from "react";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Link } from "react-router-dom";
import { Plus, Pencil, Trash2, Search } from "lucide-react";
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
import { getGroups, deleteGroup } from "@/lib/api";
import { toast } from "@/hooks/use-toast";
import { ApiError } from "@/lib/errors";

export default function Groups() {
  const [searchTerm, setSearchTerm] = useState("");
  const [deleteId, setDeleteId] = useState<string | null>(null);
  const queryClient = useQueryClient();

  const { data: groups = [], isLoading } = useQuery({
    queryKey: ["groups"],
    queryFn: getGroups,
  });

  const deleteMutation = useMutation({
    mutationFn: deleteGroup,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["groups"] });
      toast({
        title: "Group deleted",
        description: "The group has been successfully removed.",
      });
      setDeleteId(null);
    },
    onError: (error) => {
      let description = "Failed to delete group. Please try again.";
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

  const filteredGroups = groups.filter((group) =>
    group.label.toLowerCase().includes(searchTerm.toLowerCase())
  );

  return (
    <div className="min-h-screen bg-background">
      <header className="border-b border-border bg-card">
        <div className="flex items-center justify-between px-8 py-6">
          <div>
            <h1 className="text-3xl font-bold text-foreground">Groups</h1>
            <p className="mt-1 text-sm text-muted-foreground">
              Organize users and hosts into access groups
            </p>
          </div>
          <Button asChild>
            <Link to="/groups/new">
              <Plus className="mr-2 h-4 w-4" />
              Add Group
            </Link>
          </Button>
        </div>
      </header>

      <div className="p-8">
        <Card className="bg-gradient-card shadow-md">
          <div className="p-6">
            <div className="mb-6 flex items-center gap-4">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                <Input
                  placeholder="Search groups by name..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-9"
                />
              </div>
            </div>

            {isLoading ? (
              <div className="py-8 text-center text-muted-foreground">Loading groups...</div>
            ) : filteredGroups.length === 0 ? (
              <div className="py-8 text-center">
                <p className="text-muted-foreground">
                  {searchTerm ? "No groups match your search." : "No groups yet. Create your first group!"}
                </p>
              </div>
            ) : (
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Group Name</TableHead>
                    <TableHead>Users</TableHead>
                    <TableHead>Hosts</TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filteredGroups.map((group) => (
                    <TableRow key={group.label}>
                      <TableCell className="font-medium">{group.label}</TableCell>
                      <TableCell>
                        <span className="text-sm text-muted-foreground">
                          {group.users.length} user{group.users.length !== 1 ? 's' : ''}
                        </span>
                      </TableCell>
                      <TableCell>
                        <span className="text-sm text-muted-foreground">
                          {group.hosts.length} host{group.hosts.length !== 1 ? 's' : ''}
                        </span>
                      </TableCell>
                      <TableCell className="text-right">
                        <div className="flex justify-end gap-2">
                          <Button asChild variant="ghost" size="sm">
                            <Link to={`/groups/${group.label}/edit`}>
                              <Pencil className="h-4 w-4" />
                            </Link>
                          </Button>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => setDeleteId(group.label)}
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
            <AlertDialogTitle>Delete Group</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete this group? This action cannot be undone.
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
