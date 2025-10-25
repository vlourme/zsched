import { useLoaderData } from "react-router";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "~/components/ui/table";
import { pool } from "~/lib/db";
import type { Route } from "./+types/logs";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Zsched - Tasks" },
    {
      name: "description",
      content: "Zsched is a task scheduler built for Go applications.",
    },
  ];
}

export async function loader({ params }: Route.LoaderArgs) {
  const logs = await pool.query(
    `
    SELECT task_id, state_id, level, message, data, logged_at FROM logs WHERE task_id = $1 ORDER BY logged_at DESC
    `,
    [params.task_id]
  );

  return {
    logs: logs.rows,
  };
}

export default function Logs() {
  const { logs } = useLoaderData<typeof loader>();

  return (
    <>
      <div className="flex flex-col gap-1">
        <h1 className="text-3xl font-bold">Logs</h1>
        <p className="text-muted-foreground">Logs for the task.</p>
      </div>

      <div className="bg-foreground/5 rounded-lg">
        <div className="flex flex-wrap items-center p-4 gap-4">
          <h2 className="text-lg font-bold">Logs</h2>
        </div>
        <div className="overflow-hidden border-t">
          <Table>
            <TableHeader>
              <TableRow className="bg-foreground/5">
                <TableHead className="px-4 py-2">State ID</TableHead>
                <TableHead className="px-4 py-2">Level</TableHead>
                <TableHead className="px-4 py-2">Message</TableHead>
                <TableHead className="px-4 py-2">Data</TableHead>
                <TableHead className="px-4 py-2">Logged At</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {logs.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={4} className="h-24 text-center">
                    No logs found.
                  </TableCell>
                </TableRow>
              ) : (
                logs.map((log: any, idx) => (
                  <TableRow key={idx}>
                    <TableCell className="px-4 py-2">{log.state_id}</TableCell>
                    <TableCell className="px-4 py-2">{log.level}</TableCell>
                    <TableCell className="px-4 py-2">{log.message}</TableCell>
                    <TableCell className="px-4 py-2 font-mono">
                      {log.data}
                    </TableCell>
                    <TableCell className="px-4 py-2">
                      {new Date(log.logged_at).toLocaleString()}
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </div>
      </div>
    </>
  );
}
