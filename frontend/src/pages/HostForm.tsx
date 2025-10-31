import { useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { getHost, createHost, updateHost, getGroups } from "@/lib/api";
import { toast } from "@/hooks/use-toast";
import { Host } from "@/types/api";
import { Checkbox } from "@/components/ui/checkbox";
import { ApiError } from "@/lib/errors";

export default function HostForm() {
  const { id } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const isEditing = !!id;

  const { data: host } = useQuery({
    queryKey: ["host", id],
    queryFn: () => getHost(id!),
    enabled: isEditing,
  });

  const { data: groups = [] } = useQuery({
    queryKey: ["groups"],
    queryFn: getGroups,
  });

  const { register, handleSubmit, setValue, watch } = useForm<Partial<Host>>({
    defaultValues: {
      alias: "",
      host: "",
      user: "",
      key: "",
      groups: [],
    },
  });

  const selectedGroups = watch("groups") || [];

  useEffect(() => {
    if (host) {
      setValue("alias", host.alias);
      setValue("host", host.host);
      setValue("user", host.user);
      setValue("key", host.key);
      setValue("groups", host.groups);
    }
  }, [host, setValue]);

  const mutation = useMutation({
    mutationFn: (data: Partial<Host>) =>
      isEditing ? updateHost(id!, data) : createHost(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["hosts"] });
      toast({
        title: isEditing ? "Host updated" : "Host created",
        description: `The host has been successfully ${isEditing ? "updated" : "created"}.`,
      });
      navigate("/hosts");
    },
    onError: (error) => {
      let description = `Failed to ${isEditing ? "update" : "create"} host. Please try again.`;
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

  const onSubmit = (data: Partial<Host>) => {
    mutation.mutate(data);
  };

  const toggleGroup = (groupLabel: string) => {
    const current = selectedGroups;
    const updated = current.includes(groupLabel)
      ? current.filter((g) => g !== groupLabel)
      : [...current, groupLabel];
    setValue("groups", updated);
  };

  return (
    <div className="min-h-screen bg-background">
      <header className="border-b border-border bg-card">
        <div className="px-8 py-6">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => navigate("/hosts")}
            className="mb-4"
          >
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Hosts
          </Button>
          <h1 className="text-3xl font-bold text-foreground">
            {isEditing ? "Edit Host" : "Add New Host"}
          </h1>
          <p className="mt-1 text-sm text-muted-foreground">
            {isEditing ? "Update host configuration" : "Configure a new SSH host"}
          </p>
        </div>
      </header>

      <div className="p-8">
        <Card className="mx-auto max-w-2xl bg-gradient-card shadow-md">
          <form onSubmit={handleSubmit(onSubmit)} className="p-6">
            <div className="space-y-6">
              <div className="space-y-2">
                <Label htmlFor="alias">Alias *</Label>
                <Input
                  id="alias"
                  {...register("alias", { required: true })}
                  placeholder="server-1"
                  disabled={isEditing}
                />
                <p className="text-xs text-muted-foreground">
                  Unique identifier for this host
                </p>
              </div>

              <div className="space-y-2">
                <Label htmlFor="host">Host *</Label>
                <Input
                  id="host"
                  {...register("host", { required: true })}
                  placeholder="192.168.1.100 or example.com"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="user">User *</Label>
                <Input
                  id="user"
                  {...register("user", { required: true })}
                  placeholder="root"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="key">Key File *</Label>
                <Input
                  id="key"
                  {...register("key", { required: true })}
                  placeholder="~/.ssh/id_rsa"
                />
                <p className="text-xs text-muted-foreground">
                  Path to the SSH private key file
                </p>
              </div>

              <div className="space-y-3">
                <Label>Groups</Label>
                <div className="space-y-2">
                  {groups.map((group) => (
                    <div key={group.label} className="flex items-center space-x-2">
                      <Checkbox
                        id={`group-${group.label}`}
                        checked={selectedGroups.includes(group.label)}
                        onCheckedChange={() => toggleGroup(group.label)}
                      />
                      <label
                        htmlFor={`group-${group.label}`}
                        className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                      >
                        {group.label}
                      </label>
                    </div>
                  ))}
                  {groups.length === 0 && (
                    <p className="text-sm text-muted-foreground">
                      No groups available. Create groups first to assign them.
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
                    ? "Update Host"
                    : "Create Host"}
                </Button>
                <Button
                  type="button"
                  variant="secondary"
                  onClick={() => navigate("/hosts")}
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
