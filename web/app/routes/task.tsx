import {
  AlertCircleIcon,
  ArrowRightIcon,
  CheckIcon,
  ClockIcon,
  Loader2Icon,
} from "lucide-react";
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
        sum(CASE WHEN status = 'success' THEN 1 ELSE 0 END) as total_success,
        sum(CASE WHEN status = 'failed' THEN 1 ELSE 0 END) as total_err
      FROM tasks
      WHERE task_name = $1 ${after ? `AND published_at <= $2` : ""}
      `,
      [params.name, after]
    ),
    pool.query(
      `
      SELECT 
        task_id, 
        status,
        parent_id, 
        published_at,
        started_at, 
        ended_at,
        published_at::long as cursor, 
        last_error, 
        iterations,
        (ended_at - started_at) / 1000 as duration
      FROM tasks
      WHERE task_name = $1 ${after ? `AND published_at <= $2` : ""}
      ORDER BY published_at DESC
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

export const handle = {
  title: () => "Task Details",
  group: "tasks",
};

export function StatusIcon({ status }: { status: string }) {
  switch (status) {
    case "pending":
      return <ClockIcon className="size-4 text-gray-400" />;
    case "running":
      return <Loader2Icon className="size-4 animate-spin text-yellow-500" />;
    case "success":
      return <CheckIcon className="size-4 text-green-500" />;
    case "failed":
      return <AlertCircleIcon className="size-4 text-red-500" />;
  }
}

export function ExecutionRow({ execution }: { execution: any }) {
  return (
    <TableRow key={execution.task_id}>
      <TableCell className="px-6 py-4 w-8">
        <StatusIcon status={execution.status ?? "pending"} />
      </TableCell>
      <TableCell className="px-6 py-4">
        <Link
          className="text-blue-500 hover:underline font-bold"
          to={`/logs/${execution.task_id}`}
        >
          {execution.task_id}
        </Link>
      </TableCell>
      <TableCell className="px-6 py-4">
        {execution.parent_id === "00000000-0000-0000-0000-000000000000"
          ? "None"
          : execution.parent_id}
      </TableCell>
      {execution.status === "pending" ? (
        <TableCell colSpan={2} className="px-6 py-4">
          {new Date(execution.published_at).toLocaleString()}
        </TableCell>
      ) : (
        <>
          <TableCell className="px-6 py-4">
            {new Date(execution.started_at).toLocaleString()}
          </TableCell>
          <TableCell className="px-6 py-4">
            {execution.ended_at
              ? new Date(execution.ended_at).toLocaleString()
              : ""}
          </TableCell>
        </>
      )}
      <TableCell className="px-6 py-4">
        {execution.status === "pending"
          ? "Planned"
          : execution.duration
            ? formatDuration(execution.duration / 1000)
            : "Running"}
      </TableCell>
      <TableCell className="px-6 py-4">{execution.iterations ?? 0}</TableCell>
      <TableCell className="px-6 py-4 w-10">
        <Link to={`/logs/${execution.task_id}`}>
          <Button variant="outline" size="icon">
            <ArrowRightIcon className="size-4" />
          </Button>
        </Link>
      </TableCell>
    </TableRow>
  );
}

export default function Task() {
  const { task, stats, executions, queue } = useLoaderData<typeof loader>();
  const [searchParams, setSearchParams] = useSearchParams();

  return (
    <div className="flex flex-col">
      <div className="grid md:divide-x grid-cols-1 md:grid-cols-2">
        <Card className="rounded-none bg-background">
          <CardHeader>
            <CardTitle>{task.name}</CardTitle>
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
                <p className="text-sm text-muted-foreground">Total successes</p>
                <p className="text-sm">{stats.total_success ?? 0}</p>
              </div>
              <div className="flex flex-col w-36 gap-1">
                <p className="text-sm text-muted-foreground">Total errors</p>
                <p className="text-sm">{stats.total_err ?? 0}</p>
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
            {task.schedules.length > 0 ? (
              <div className="mt-4">
                <h2 className="text-lg font-bold">Schedules</h2>
                <div className="flex flex-col mt-2 gap-2">
                  {task.schedules.map((schedule: any) => (
                    <Card key={schedule.schedule} className="border py-3.5">
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
            ) : null}
          </CardContent>
        </Card>
        <MessageQueueCard data={queue.message_stats} />
      </div>

      <div className="px-4 border-t py-2.5 gap-2 flex flex-row justify-between items-center">
        <p>Executions</p>
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
              setSearchParams(
                {
                  ...Object.fromEntries(searchParams.entries()),
                  after: executions[executions.length - 1]?.cursor.toString(),
                },
                { replace: true }
              );
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
              <TableHead className="px-6 py-4">Status</TableHead>
              <TableHead className="px-6 py-4">Task ID</TableHead>
              <TableHead className="px-6 py-4">Parent ID</TableHead>
              <TableHead className="px-6 py-4">Started At</TableHead>
              <TableHead className="px-6 py-4">Ended At</TableHead>
              <TableHead className="px-6 py-4">Duration</TableHead>
              <TableHead className="px-6 py-4">Iterations</TableHead>
              <TableHead className="px-6 py-4 w-10"></TableHead>
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

            {executions.map((execution: any, idx: number) => (
              <ExecutionRow key={idx} execution={execution} />
            ))}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
