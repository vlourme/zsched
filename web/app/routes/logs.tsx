import { Editor } from "@monaco-editor/react";
import { Form, redirect, useLoaderData } from "react-router";
import { Button } from "~/components/ui/button";
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
  const [task, logs] = await Promise.all([
    pool.query(`SELECT task_name, state FROM tasks WHERE task_id = $1`, [
      params.task_id,
    ]),
    pool.query(
      `SELECT task_id, state_id, level, message, data, logged_at FROM logs WHERE task_id = $1 ORDER BY logged_at DESC
      `,
      [params.task_id]
    ),
  ]);

  if (task.rows.length === 0) {
    return redirect("/tasks");
  }

  return {
    task: task.rows[0].task_name,
    parameters: JSON.parse(task.rows[0].state),
    logs: logs.rows,
  };
}

export const handle = {
  title: () => "Logs",
  group: "tasks",
};

export async function action({ request, params }: Route.ActionArgs) {
  const formData = await request.formData();
  const taskName = formData.get("task_name");
  const parameters = formData.get("parameters");

  await fetch(process.env.ZSCHED_URL + "/tasks/" + taskName, {
    method: "POST",
    body: parameters as string,
  });
}

export default function Logs() {
  const { task, logs, parameters } = useLoaderData<typeof loader>();

  return (
    <>
      <Card className="rounded-none bg-background">
        <CardHeader className="flex flex-row justify-between items-center">
          <div>
            <CardTitle>{task} task parameters</CardTitle>
            <CardDescription>
              This task was started with the following parameters.
            </CardDescription>
          </div>
          <Form method="post">
            <input
              type="hidden"
              name="parameters"
              value={JSON.stringify(parameters)}
            />
            <input type="hidden" name="task_name" value={task} />
            <Button variant="outline" size="sm" type="submit">
              Dispatch again
            </Button>
          </Form>
        </CardHeader>
        <CardContent>
          <Editor
            className="w-full rounded-sm overflow-hidden"
            language="json"
            defaultValue={JSON.stringify(parameters, null, 2)}
            theme="vs-dark"
            options={{
              minimap: { enabled: false },
              padding: { top: 10, bottom: 10 },
              lineNumbers: "off",
              fontSize: 14,
              readOnly: true,
            }}
            height="20vh"
          />
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
