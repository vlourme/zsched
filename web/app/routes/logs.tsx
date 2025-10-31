import { redirect, useLoaderData } from "react-router";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "~/components/ui/card";
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
  const task = await pool.query(`SELECT state FROM tasks WHERE task_id = $1`, [
    params.task_id,
  ]);
  const logs = await pool.query(
    `
    SELECT task_id, state_id, level, message, data, logged_at FROM logs WHERE task_id = $1 ORDER BY logged_at DESC
    `,
    [params.task_id]
  );

  if (task.rows.length === 0) {
    return redirect("/tasks");
  }

  return {
    parameters: JSON.parse(task.rows[0].state),
    logs: logs.rows,
  };
}

export const handle = {
  title: () => "Logs",
  group: "tasks",
};

export default function Logs() {
  const { logs, parameters } = useLoaderData<typeof loader>();

  return (
    <>
      <Card className="rounded-none bg-background">
        <CardHeader>
          <CardTitle>Task parameters</CardTitle>
          <CardDescription>
            This task was started with the following parameters.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <pre className="text-sm font-mono">
            {JSON.stringify(parameters, null, 2)}
          </pre>
        </CardContent>
      </Card>

      <Table className="border-t">
        <TableHeader>
          <TableRow className="bg-foreground/5">
            <TableHead className="px-6 py-4">State ID</TableHead>
            <TableHead className="px-6 py-4">Level</TableHead>
            <TableHead className="px-6 py-4">Message</TableHead>
            <TableHead className="px-6 py-4">Data</TableHead>
            <TableHead className="px-6 py-4">Logged At</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {logs.length === 0 ? (
            <TableRow>
              <TableCell colSpan={5} className="h-24 text-center">
                No logs found.
              </TableCell>
            </TableRow>
          ) : (
            logs.map((log: any, idx) => (
              <TableRow key={idx}>
                <TableCell className="px-6 py-4">{log.state_id}</TableCell>
                <TableCell className="px-6 py-4">{log.level}</TableCell>
                <TableCell className="px-6 py-4">{log.message}</TableCell>
                <TableCell className="px-6 py-4 font-mono">
                  {log.data}
                </TableCell>
                <TableCell className="px-6 py-4">
                  {new Date(log.logged_at).toLocaleString()}
                </TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </>
  );
}
