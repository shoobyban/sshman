import { useState, useEffect, useRef } from "react";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Pause, Play, Search, Download } from "lucide-react";
import { cn } from "@/lib/utils";

interface LogEntry {
  timestamp: string;
  level: string;
  message: string;
}

export default function Logs() {
  const [logs, setLogs] = useState<LogEntry[]>([]);
  const [isPaused, setIsPaused] = useState(false);
  const [searchTerm, setSearchTerm] = useState("");
  const logsEndRef = useRef<HTMLDivElement>(null);
  const eventSourceRef = useRef<EventSource | null>(null);

  useEffect(() => {
    // Connect to SSE endpoint
    const API_BASE = import.meta.env.VITE_API_BASE || "/api";
    const eventSource = new EventSource(`${API_BASE}/logs`);
    eventSourceRef.current = eventSource;

    eventSource.onmessage = (event) => {
      if (!isPaused) {
        try {
          const logEntry: LogEntry = JSON.parse(event.data);
          setLogs((prev) => [...prev, logEntry].slice(-1000)); // Keep last 1000 logs
        } catch (e) {
          // If not JSON, treat as raw message
          setLogs((prev) => [
            ...prev,
            {
              timestamp: new Date().toISOString(),
              level: "INFO",
              message: event.data,
            },
          ].slice(-1000));
        }
      }
    };

    eventSource.onerror = () => {
      console.error("SSE connection error");
    };

    return () => {
      eventSource.close();
    };
  }, [isPaused]);

  useEffect(() => {
    if (!isPaused && logsEndRef.current) {
      logsEndRef.current.scrollIntoView({ behavior: "smooth" });
    }
  }, [logs, isPaused]);

  const filteredLogs = logs.filter(
    (log) =>
      log.message.toLowerCase().includes(searchTerm.toLowerCase()) ||
      log.level.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const getLevelColor = (level: string) => {
    switch (level.toUpperCase()) {
      case "ERROR":
        return "text-destructive";
      case "WARN":
        return "text-yellow-500";
      case "INFO":
        return "text-accent";
      case "DEBUG":
        return "text-muted-foreground";
      default:
        return "text-foreground";
    }
  };

  const exportLogs = () => {
    const content = logs
      .map((log) => `[${log.timestamp}] ${log.level}: ${log.message}`)
      .join("\n");
    const blob = new Blob([content], { type: "text/plain" });
    const url = URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = `sshman-logs-${new Date().toISOString()}.txt`;
    a.click();
    URL.revokeObjectURL(url);
  };

  return (
    <div className="min-h-screen bg-background">
      <header className="border-b border-border bg-card">
        <div className="flex items-center justify-between px-8 py-6">
          <div>
            <h1 className="text-3xl font-bold text-foreground">Live Logs</h1>
            <p className="mt-1 text-sm text-muted-foreground">
              Real-time streaming of backend activity
            </p>
          </div>
          <div className="flex gap-3">
            <Button variant="secondary" onClick={exportLogs}>
              <Download className="mr-2 h-4 w-4" />
              Export
            </Button>
            <Button
              variant={isPaused ? "default" : "secondary"}
              onClick={() => setIsPaused(!isPaused)}
            >
              {isPaused ? (
                <>
                  <Play className="mr-2 h-4 w-4" />
                  Resume
                </>
              ) : (
                <>
                  <Pause className="mr-2 h-4 w-4" />
                  Pause
                </>
              )}
            </Button>
          </div>
        </div>
      </header>

      <div className="p-8">
        <Card className="bg-gradient-card shadow-md">
          <div className="p-6">
            <div className="mb-4 flex items-center gap-4">
              <div className="relative flex-1">
                <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
                <Input
                  placeholder="Filter logs by message or level..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="pl-9"
                />
              </div>
              <div className="text-sm text-muted-foreground">
                {filteredLogs.length} log{filteredLogs.length !== 1 ? 's' : ''}
              </div>
            </div>

            <div className="h-[600px] overflow-y-auto rounded-lg border border-border bg-card p-4 font-mono text-sm">
              {filteredLogs.length === 0 ? (
                <div className="flex h-full items-center justify-center text-muted-foreground">
                  {searchTerm ? "No logs match your filter." : "Waiting for logs..."}
                </div>
              ) : (
                <div className="space-y-1">
                  {filteredLogs.map((log, index) => (
                    <div key={index} className="flex gap-4 py-1">
                      <span className="text-muted-foreground">
                        {new Date(log.timestamp).toLocaleTimeString()}
                      </span>
                      <span className={cn("font-semibold", getLevelColor(log.level))}>
                        [{log.level}]
                      </span>
                      <span className="flex-1 text-foreground">{log.message}</span>
                    </div>
                  ))}
                  <div ref={logsEndRef} />
                </div>
              )}
            </div>
          </div>
        </Card>
      </div>
    </div>
  );
}
