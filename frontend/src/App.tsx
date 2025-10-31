import { Toaster } from "@/components/ui/toaster";
import { Toaster as Sonner } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { Layout } from "./components/Layout";
import Dashboard from "./pages/Dashboard";
import Hosts from "./pages/Hosts";
import HostForm from "./pages/HostForm";
import Users from "./pages/Users";
import UserForm from "./pages/UserForm";
import Groups from "./pages/Groups";
import GroupForm from "./pages/GroupForm";
import Logs from "./pages/Logs";
import NotFound from "./pages/NotFound";

const queryClient = new QueryClient();

const App = () => (
  <QueryClientProvider client={queryClient}>
    <TooltipProvider>
      <Toaster />
      <Sonner />
      <BrowserRouter>
        <Layout>
          <Routes>
            <Route path="/" element={<Dashboard />} />
            <Route path="/hosts" element={<Hosts />} />
            <Route path="/hosts/new" element={<HostForm />} />
            <Route path="/hosts/:id/edit" element={<HostForm />} />
            <Route path="/users" element={<Users />} />
            <Route path="/users/new" element={<UserForm />} />
            <Route path="/users/:id/edit" element={<UserForm />} />
            <Route path="/groups" element={<Groups />} />
            <Route path="/groups/new" element={<GroupForm />} />
            <Route path="/groups/:id/edit" element={<GroupForm />} />
            <Route path="/logs" element={<Logs />} />
            <Route path="*" element={<NotFound />} />
          </Routes>
        </Layout>
      </BrowserRouter>
    </TooltipProvider>
  </QueryClientProvider>
);

export default App;
