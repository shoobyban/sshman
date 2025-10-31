import { useEffect } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { useForm } from "react-hook-form";
import { ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Checkbox } from "@/components/ui/checkbox";
import { getUser, createUser, updateUser, getGroups } from "@/lib/api";
import { toast } from "@/hooks/use-toast";
import { User } from "@/types/api";
import { ApiError } from "@/lib/errors";

export default function UserForm() {
  const { id } = useParams();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const isEditing = !!id;

  const { data: user } = useQuery({
    queryKey: ["user", id],
    queryFn: () => getUser(id!),
    enabled: isEditing,
  });

  const { data: groups = [] } = useQuery({
    queryKey: ["groups"],
    queryFn: getGroups,
  });

  const { register, handleSubmit, setValue, watch } = useForm<Partial<User>>({
    defaultValues: {
      name: "",
      email: "",
      type: "ssh-rsa",
      key: "",
      groups: [],
    },
  });

  const selectedGroups = watch("groups") || [];

  useEffect(() => {
    if (user) {
      setValue("name", user.name);
      setValue("email", user.email);
      setValue("type", user.type);
      setValue("key", user.key);
      setValue("groups", user.groups);
    }
  }, [user, setValue]);

  const mutation = useMutation({
    mutationFn: (data: Partial<User>) =>
      isEditing ? updateUser(id!, data) : createUser(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
      toast({
        title: isEditing ? "User updated" : "User created",
        description: `The user has been successfully ${isEditing ? "updated" : "created"}.`,
      });
      navigate("/users");
    },
    onError: (error) => {
      let description = `Failed to ${isEditing ? "update" : "create"} user. Please try again.`;
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

  const onSubmit = (data: Partial<User>) => {
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
            onClick={() => navigate("/users")}
            className="mb-4"
          >
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Users
          </Button>
          <h1 className="text-3xl font-bold text-foreground">
            {isEditing ? "Edit User" : "Add New User"}
          </h1>
          <p className="mt-1 text-sm text-muted-foreground">
            {isEditing ? "Update user information" : "Create a new SSH user"}
          </p>
        </div>
      </header>

      <div className="p-8">
        <Card className="mx-auto max-w-2xl bg-gradient-card shadow-md">
          <form onSubmit={handleSubmit(onSubmit)} className="p-6">
            <div className="space-y-6">
              <div className="space-y-2">
                <Label htmlFor="name">Name *</Label>
                <Input
                  id="name"
                  {...register("name", { required: true })}
                  placeholder="John Doe"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="email">Email *</Label>
                <Input
                  id="email"
                  type="email"
                  {...register("email", { required: true })}
                  placeholder="john@example.com"
                  disabled={isEditing}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="type">Key Type *</Label>
                <Input
                  id="type"
                  {...register("type", { required: true })}
                  placeholder="ssh-rsa, ssh-ed25519, etc."
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="key">Public SSH Key *</Label>
                <Textarea
                  id="key"
                  {...register("key", { required: true })}
                  placeholder="Paste the user's public SSH key here (e.g., ssh-rsa AAAAB3NzaC1...)"
                  rows={6}
                  className="font-mono text-sm"
                />
                <p className="text-xs text-muted-foreground">
                  This should be the public key from the user's ~/.ssh/id_rsa.pub or similar
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
                    ? "Update User"
                    : "Create User"}
                </Button>
                <Button
                  type="button"
                  variant="secondary"
                  onClick={() => navigate("/users")}
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
