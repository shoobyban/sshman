import { useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Checkbox } from "@/components/ui/checkbox";
import { getGroup, createGroup, updateGroup, getUsers, getHosts } from "@/lib/api";
import { toast } from "@/hooks/use-toast";
import { Group } from "@/types/api";
import { ApiError } from "@/lib/errors";

export default function GroupForm() {
  const { id } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const isEditing = !!id;

  const { data: group } = useQuery({
    queryKey: ["group", id],
    queryFn: () => getGroup(id!),
    enabled: isEditing,
  });

  const { data: users = [] } = useQuery({
    queryKey: ["users"],
    queryFn: getUsers,
  });

  const { data: hosts = [] } = useQuery({
    queryKey: ["hosts"],
    queryFn: getHosts,
  });

  const { register, handleSubmit, setValue, watch } = useForm<Partial<Group>>({
    defaultValues: {
      label: "",
      users: [],
      hosts: [],
    },
  });

  const selectedUsers = watch("users") || [];
  const selectedHosts = watch("hosts") || [];

  useEffect(() => {
    if (group) {
      setValue("label", group.label);
      setValue("users", group.users);
      setValue("hosts", group.hosts);
    }
  }, [group, setValue]);

  const mutation = useMutation({
    mutationFn: (data: Partial<Group>) =>
      isEditing ? updateGroup(id!, data) : createGroup(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["groups"] });
      toast({
        title: isEditing ? "Group updated" : "Group created",
        description: `The group has been successfully ${isEditing ? "updated" : "created"}.`,
      });
      navigate("/groups");
    },
    onError: (error) => {
      let description = `Failed to ${isEditing ? "update" : "create"} group. Please try again.`;
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

  const onSubmit = (data: Partial<Group>) => {
    mutation.mutate(data);
  };

  const toggleUser = (userEmail: string) => {
    const current = selectedUsers;
    const updated = current.includes(userEmail)
      ? current.filter((u) => u !== userEmail)
      : [...current, userEmail];
    setValue("users", updated);
  };

  const toggleHost = (hostAlias: string) => {
    const current = selectedHosts;
    const updated = current.includes(hostAlias)
      ? current.filter((h) => h !== hostAlias)
      : [...current, hostAlias];
    setValue("hosts", updated);
  };

  return (
    <div className="min-h-screen bg-background">
      <header className="border-b border-border bg-card">
        <div className="px-8 py-6">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => navigate("/groups")}
            className="mb-4"
          >
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Groups
          </Button>
          <h1 className="text-3xl font-bold text-foreground">
            {isEditing ? "Edit Group" : "Add New Group"}
          </h1>
          <p className="mt-1 text-sm text-muted-foreground">
            {isEditing ? "Update group configuration" : "Create a new access group"}
          </p>
        </div>
      </header>

      <div className="p-8">
        <Card className="mx-auto max-w-2xl bg-gradient-card shadow-md">
          <form onSubmit={handleSubmit(onSubmit)} className="p-6">
            <div className="space-y-6">
              <div className="space-y-2">
                <Label htmlFor="label">Group Name *</Label>
                <Input
                  id="label"
                  {...register("label", { required: true })}
                  placeholder="developers"
                  disabled={isEditing}
                />
                <p className="text-xs text-muted-foreground">
                  Unique identifier for this group
                </p>
              </div>

              <div className="space-y-3">
                <Label>Users</Label>
                <div className="max-h-60 space-y-2 overflow-y-auto rounded-lg border border-border p-4">
                  {users.map((user) => (
                    <div key={user.email} className="flex items-center space-x-2">
                      <Checkbox
                        id={`user-${user.email}`}
                        checked={selectedUsers.includes(user.email)}
                        onCheckedChange={() => toggleUser(user.email)}
                      />
                      <label
                        htmlFor={`user-${user.email}`}
                        className="flex-1 text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                      >
                        {user.name} ({user.email})
                      </label>
                    </div>
                  ))}
                  {users.length === 0 && (
                    <p className="text-sm text-muted-foreground">
                      No users available. Create users first to add them to groups.
                    </p>
                  )}
                </div>
              </div>

              <div className="space-y-3">
                <Label>Hosts</Label>
                <div className="max-h-60 space-y-2 overflow-y-auto rounded-lg border border-border p-4">
                  {hosts.map((host) => (
                    <div key={host.alias} className="flex items-center space-x-2">
                      <Checkbox
                        id={`host-${host.alias}`}
                        checked={selectedHosts.includes(host.alias)}
                        onCheckedChange={() => toggleHost(host.alias)}
                      />
                      <label
                        htmlFor={`host-${host.alias}`}
                        className="flex-1 text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                      >
                        {host.alias} ({host.user}@{host.host})
                      </label>
                    </div>
                  ))}
                  {hosts.length === 0 && (
                    <p className="text-sm text-muted-foreground">
                      No hosts available. Create hosts first to add them to groups.
                    </p>
                  )}
                </div>
              </div>

              <div className="flex gap-3 pt-4">
                <Button type="submit" disabled={mutation.isPending}>
                  {mutation.isPending
                    ? isEditing
                      ? "Updating..."
                      : "Creating..."
                    : isEditing
                    ? "Update Group"
                    : "Create Group"}
                </Button>
                <Button
                  type="button"
                  variant="secondary"
                  onClick={() => navigate("/groups")}
                >
                  Cancel
                </Button>
              </div>
            </div>
          </form>
        </Card>
      </div>
    </div>
  );
}
