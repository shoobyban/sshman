### Prompt: Design the Ideal GUI for SSHMan

**Project Vision:**
Create a clean, modern, and responsive web-based graphical user interface (GUI) for SSHMan. The application will serve as the primary control plane for managing SSH infrastructure, allowing administrators to easily perform all CRUD (Create, Read, Update, Delete) operations on hosts, users, and groups. The GUI should be intuitive for both novice and experienced users, providing clear navigation and immediate feedback on all operations.

**Target Audience:**
System administrators, DevOps engineers, and IT managers who need to manage SSH access across multiple servers and for multiple users.

**Key Views and Functionality:**

**1. Dashboard / Home View (`/`)**
*   **Purpose:** Provide an at-a-glance summary of the entire SSHMan-managed infrastructure.
*   **UI Components:**
    *   **Stat Cards:** Display key metrics:
        *   Total number of managed hosts.
        *   Total number of users.
        *   Total number of groups.
    *   **Live Log Stream:** A small, embedded panel showing the last 5-10 log entries from the `/api/logs` stream to show recent activity.
    *   **Quick Actions:** Buttons or links for common tasks like "Add New Host," "Add New User," or "Add New Group."

**2. Hosts Management View (`/hosts`)**
*   **Purpose:** Display, create, and manage all host entries.
*   **UI Components:**
    *   **Hosts List:** A table or card-based list displaying all hosts. Each entry should show:
        *   Host Alias
        *   Connection String (`user@host`)
        *   Groups it belongs to (as tags or a comma-separated list).
        *   Action buttons: `Edit`, `Delete` (with confirmation).
    *   **"Add Host" Button:** Opens a modal or navigates to a form for creating a new host.
    *   **Global Actions:**
        *   A "Sync All Hosts" button to trigger the `/api/hosts/sync` endpoint.
        *   A search/filter bar to quickly find hosts by alias or connection string.
*   **Host Create/Edit Form (`/hosts/new` or `/hosts/:id/edit`)**
    *   Fields for: Alias, Host, User, Key File.
    *   A multi-select component to assign the host to one or more groups.
    *   A view showing which users currently have access to this host.

**3. Users Management View (`/users`)**
*   **Purpose:** Display, create, and manage all user entries.
*   **UI Components:**
    *   **Users List:** A table displaying all users. Each entry should show:
        *   User Name
        *   User Email
        *   Groups the user belongs to.
        *   Action buttons: `Edit`, `Delete` (with confirmation).
    *   **"Add User" Button:** Opens a form for creating a new user.
    *   **Search/Filter Bar:** To find users by name or email.
*   **User Create/Edit Form (`/users/new` or `/users/:id/edit`)**
    *   Fields for: Name, Email.
    *   A text area to paste the user's public SSH key (`type key name`).
    *   A multi-select component to assign the user to one or more groups.

**4. Groups Management View (`/groups`)**
*   **Purpose:** Create and manage groups to associate users with hosts.
*   **UI Components:**
    *   **Groups List:** A table or card list showing all groups. Each entry should display:
        *   Group Name (Label)
        *   Count of users in the group.
        *   Count of hosts in the group.
        *   Action buttons: `Edit`, `Delete` (with confirmation).
    *   **"Add Group" Button:** Opens a form for creating a new group.
*   **Group Create/Edit Form (`/groups/new` or `/groups/:id/edit`)**
    *   A text field for the Group Name (Label).
    *   Two multi-select components:
        *   One to select all users that should be part of this group.
        *   One to select all hosts that this group should grant access to.

**5. Live Logs View (`/logs`)**
*   **Purpose:** Provide a real-time, streaming view of backend activity for monitoring and debugging.
*   **UI Components:**
    *   A full-page, auto-scrolling view that connects to the `/api/logs` Server-Sent Events (SSE) endpoint.
    *   Each log entry should be color-coded based on severity (e.g., INFO, WARN, ERROR).
    *   A toggle to pause/resume the log stream.
    *   A search bar to filter logs by keywords.

**Non-Functional Requirements:**

*   **Responsiveness:** The entire application must be fully responsive and usable on desktop, tablet, and mobile devices.
*   **User Feedback:** Use toast notifications or alerts to provide immediate feedback for all user actions (e.g., "Host 'server-1' was successfully created," "Error: Failed to delete user.").
*   **State Consistency:** The UI must intelligently update itself after CRUD operations. For example, after a new user is created, the user list should automatically refresh to show the new entry without requiring a page reload.
*   **Error Handling:** Gracefully handle API errors and display user-friendly messages. For example, if the backend is unreachable, show a clear error message.
*   **Intuitive Navigation:** A persistent sidebar or top navigation bar should provide easy access to all main views (Dashboard, Hosts, Users, Groups, Logs).
