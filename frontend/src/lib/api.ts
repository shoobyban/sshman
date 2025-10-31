import { Host, User, Group } from "@/types/api";
import { apiFetch } from "./fetch";

const API_BASE = import.meta.env.VITE_API_BASE || "/api";

// Hosts
export const getHosts = async (): Promise<Host[]> => {
  const data = await apiFetch<Host[] | Record<string, Host>>(`${API_BASE}/hosts`);
  // backend returns a map (object) of hosts, convert to array if needed
  if (Array.isArray(data)) return data;
  return Object.values(data);
};

export const getHost = async (id: string): Promise<Host> => {
  return apiFetch<Host>(`${API_BASE}/hosts/${id}`);
};

export const createHost = async (host: Partial<Host>): Promise<Host> => {
  return apiFetch<Host>(`${API_BASE}/hosts`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(host),
  });
};

export const updateHost = async (id: string, host: Partial<Host>): Promise<Host> => {
  return apiFetch<Host>(`${API_BASE}/hosts/${id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(host),
  });
};

export const deleteHost = async (id: string): Promise<void> => {
  await apiFetch<void>(`${API_BASE}/hosts/${id}`, {
    method: "DELETE",
  });
};

export const syncHosts = async (): Promise<void> => {
  await apiFetch<void>(`${API_BASE}/hosts/sync`);
};

// Users
export const getUsers = async (): Promise<User[]> => {
  const data = await apiFetch<User[] | Record<string, User>>(`${API_BASE}/users`);
  if (Array.isArray(data)) return data;
  return Object.values(data);
};

export const getUser = async (id: string): Promise<User> => {
  return apiFetch<User>(`${API_BASE}/users/${id}`);
};

export const createUser = async (user: Partial<User>): Promise<User> => {
  return apiFetch<User>(`${API_BASE}/users`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(user),
  });
};

export const updateUser = async (id: string, user: Partial<User>): Promise<User> => {
  return apiFetch<User>(`${API_BASE}/users/${id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(user),
  });
};

export const deleteUser = async (id: string): Promise<void> => {
  await apiFetch<void>(`${API_BASE}/users/${id}`, {
    method: "DELETE",
  });
};

// Groups
export const getGroups = async (): Promise<Group[]> => {
  const data = await apiFetch<Group[] | Record<string, Group>>(`${API_BASE}/groups`);
  // backend returns a map of label => { label, hosts, users }
  if (Array.isArray(data)) return data;
  return Object.keys(data).map((k) => data[k]);
};

export const getGroup = async (id: string): Promise<Group> => {
  return apiFetch<Group>(`${API_BASE}/groups/${id}`);
};

export const createGroup = async (group: Partial<Group>): Promise<Group> => {
  return apiFetch<Group>(`${API_BASE}/groups`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(group),
  });
};

export const updateGroup = async (id: string, group: Partial<Group>): Promise<Group> => {
  return apiFetch<Group>(`${API_BASE}/groups/${id}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(group),
  });
};

export const deleteGroup = async (id: string): Promise<void> => {
  await apiFetch<void>(`${API_BASE}/groups/${id}`, {
    method: "DELETE",
  });
};
