import { AlertTriangleIcon, CheckIcon } from "lucide-react";
import { Link, redirect, useLoaderData, useSearchParams } from "react-router";
import { MessageQueueCard } from "~/components/message-queue-card";
import { NewTaskDialog } from "~/components/new-task-dialog";
import { Button } from "~/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "~/components/ui/card";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "~/components/ui/table";
import { pool } from "~/lib/db";
import { formatDuration } from "~/lib/formatters";
import { request } from "~/lib/lavinmq";
import type { Route } from "./+types/task";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Zsched - Tasks" },
    {
      name: "description",
      content: "Zsched is a task scheduler built for Go applications.",
    },
  ];
}

export async function action({ request, params }: Route.ActionArgs) {
  const formData = await request.formData();
  const parameters = formData.get("parameters");

  await fetch(process.env.ZSCHED_URL + "/tasks/" + params.name, {
    method: "POST",
    body: parameters as string,
  });

  return;
}

export async function loader({ params, request: req }: Route.LoaderArgs) {
  const searchParams = new URLSearchParams(req.url.split("?")[1]);
  const vhost = searchParams.get("vhost");
  if (!vhost) {
    return redirect("/tasks");
  }

  let after = parseInt(searchParams.get("after") ?? "0");

  const [task, stats, executions, queues] = await Promise.all([
    fetch(process.env.ZSCHED_URL + "/tasks/" + params.name).then((res) =>
      res.json()
    ),
    pool.query(
      `
      SELECT 
        count(*) as total_exec,
        MAX(started_at) as last_exec,
        sum(CASE WHEN last_error != '' THEN 1 ELSE 0 END) as total_err
      FROM tasks
      WHERE task_name = $1
      `,
      [params.name]
    ),
    pool.query(
      `
      SELECT 
        task_id, 
        parent_id, 
        started_at, 
        ended_at,
        started_at::long as cursor, 
        last_error, 
        (ended_at - started_at) / 1000 as duration
      FROM tasks
      WHERE task_name = $1 ${after ? `AND started_at <= $2` : ""}
      ORDER BY started_at DESC
      LIMIT 100
      `,
      [params.name, after]
    ),
    request<any>(`/api/queues/${encodeURIComponent(vhost)}/${params.name}`),
  ]);

  return {
    task: task,
    stats: stats.rows[0],
    executions: executions.rows,
    queue: queues,
  };
}

export default function Task() {
  const { task, stats, executions, queue } = useLoaderData<typeof loader>();
  const [searchParams, setSearchParams] = useSearchParams();

  return (
    <>
      <div className="flex flex-col gap-1">
        <h1 className="text-3xl font-bold">{task.name}</h1>
        <p className="text-muted-foreground">
          Task details and message statistics.
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Card>
          <CardHeader>
            <CardTitle>Task details</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex flex-wrap gap-2">
              <div className="flex flex-col w-36 gap-1">
                <p className="text-sm text-muted-foreground">Concurrency</p>
                <p className="text-sm">{task.concurrency}</p>
              </div>
              <div className="flex flex-col w-36 gap-1">
                <p className="text-sm text-muted-foreground">Max Retries</p>
                <p className="text-sm">{task.max_retries}</p>
              </div>
              <div className="flex flex-col w-36 gap-1">
                <p className="text-sm text-muted-foreground">
                  Total executions
                </p>
                <p className="text-sm">{stats.total_exec}</p>
              </div>
              <div className="flex flex-col w-36 gap-1">
                <p className="text-sm text-muted-foreground">Last execution</p>
                <p className="text-sm">
                  {new Date(stats.last_exec).toLocaleString()}
                </p>
              </div>
              <div className="flex flex-col w-36 gap-1">
                <p className="text-sm text-muted-foreground">Total errors</p>
                <p className="text-sm">{stats.total_err}</p>
              </div>
              <div className="flex flex-col w-36 gap-1">
                <p className="text-sm text-muted-foreground">Pending</p>
                <p className="text-sm">{queue.ready + queue.unacked}</p>
              </div>
              <div className="flex flex-col w-36 gap-1">
                <p className="text-sm text-muted-foreground">Success rate</p>
                <p className="text-sm">
                  {queue.message_stats.ack_details.rate}
                </p>
              </div>
              <div className="flex flex-col w-36 gap-1">
                <p className="text-sm text-muted-foreground">Error rate</p>
                <p className="text-sm">
                  {queue.message_stats.redeliver_details.rate +
                    queue.message_stats.reject_details.rate}
                </p>
              </div>
            </div>
            <div className="mt-4">
              <h2 className="text-lg font-bold">Schedules</h2>
              <div className="flex flex-col mt-2 gap-2">
                {task.schedules.map((schedule: any) => (
                  <Card key={schedule.schedule}>
                    <CardContent className="flex flex-wrap gap-1">
                      <div className="flex flex-col w-36 gap-1">
                        <p className="text-sm text-muted-foreground">
                          Schedule
                        </p>
                        <p className="text-sm">{schedule.schedule}</p>
                      </div>
                      <div className="flex flex-col w-36 gap-1">
                        <p className="text-sm text-muted-foreground">
                          Parameters
                        </p>
                        <p className="text-sm">
                          <pre>{JSON.stringify(schedule.parameters)}</pre>
                        </p>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            </div>
          </CardContent>
        </Card>
        <MessageQueueCard data={queue.message_stats} />
      </div>

      <div className="bg-foreground/5 rounded-lg">
        <div className="flex flex-wrap items-center p-4 gap-4">
          <h2 className="text-lg font-bold">Executions</h2>
          <NewTaskDialog />

          <div className="md:flex-1"></div>

          <p>
            {executions.length} out of {stats.total_exec} executions
          </p>
          {executions.length == 100 ? (
            <Button
              size="sm"
              variant="outline"
              onClick={() => {
                setSearchParams({
                  ...Object.fromEntries(searchParams.entries()),
                  after: executions[executions.length - 1]?.cursor.toString(),
                });
              }}
            >
              View more
            </Button>
          ) : null}
        </div>
        <div className="border-t">
          <Table>
            <TableHeader>
              <TableRow className="bg-foreground/5">
                <TableHead className="px-4 py-2">Status</TableHead>
                <TableHead className="px-4 py-2">Task ID</TableHead>
                <TableHead className="px-4 py-2">Parent ID</TableHead>
                <TableHead className="px-4 py-2">Started At</TableHead>
                <TableHead className="px-4 py-2">Ended At</TableHead>
                <TableHead className="px-4 py-2">Duration</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {executions.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={6} className="h-24 text-center">
                    No executions found.
                  </TableCell>
                </TableRow>
              ) : null}

              {executions.map((execution: any) => (
                <TableRow key={execution.task_id}>
                  <TableCell className="px-4 py-2 w-8">
                    {execution.last_error ? (
                      <AlertTriangleIcon className="size-4 text-red-500" />
                    ) : (
                      <CheckIcon className="size-4 text-green-500" />
                    )}
                  </TableCell>
                  <TableCell className="px-4 py-2">
                    <Link
                      className="text-blue-500 hover:underline"
                      to={`/logs/${execution.task_id}`}
                    >
                      {execution.task_id}
                    </Link>
                  </TableCell>
                  <TableCell className="px-4 py-2">
                    {execution.parent_id ===
                    "00000000-0000-0000-0000-000000000000"
                      ? "None"
                      : execution.parent_id}
                  </TableCell>
                  <TableCell className="px-4 py-2">
                    {new Date(execution.started_at).toLocaleString()}
                  </TableCell>
                  <TableCell className="px-4 py-2">
                    {execution.ended_at
                      ? new Date(execution.ended_at).toLocaleString()
                      : "Running..."}
                  </TableCell>
                  <TableCell className="px-4 py-2">
                    {execution.duration
                      ? formatDuration(execution.duration / 1000)
                      : "Running..."}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      </div>
    </>
  );
}
