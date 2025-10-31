export interface User {
  type: string;
  name: string;
  email: string;
  key: string;
  groups: string[];
  hosts: string[];
  keyfile?: string;
  roles?: string[];
  updatedAt?: string;
}

export interface Host {
  host: string;
  user: string;
  key: string;
  alias: string;
  userlist?: User[];
  groups: string[];
  last_updated?: string;
  checksum?: string;
  modified?: boolean;
}

export interface Group {
  label: string;
  users: string[];
  hosts: string[];
}

export interface LogEntry {
  timestamp: string;
  level: string;
  message: string;
}
