import { ArrowRightIcon } from "lucide-react";
import { useMemo } from "react";
import { Link, useLoaderData } from "react-router";
import { Button } from "~/components/ui/button";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "~/components/ui/table";
import { stringToColor } from "~/lib/color";
import { request } from "~/lib/lavinmq";
import type { Queue } from "~/types/mq-queues";
import type { Route } from "./+types/tasks";

export const handle = {
  title: () => "Tasks",
  group: "tasks",
};

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Tasks" },
    {
      name: "description",
      content: "Zsched is a task scheduler built for Go applications.",
    },
  ];
}

export async function loader() {
  const [tasks, queues] = await Promise.all([
    fetch(process.env.ZSCHED_URL + "/tasks").then((res) => res.json()),
    request<Queue[]>("/api/queues"),
  ]);
  return {
    tasks: tasks,
    queues: queues,
  };
}

function Tag({ tag }: { tag: string }) {
  return (
    <span
      className="px-2 py-0.5 rounded-sm text-sm font-medium text-white"
      style={{ backgroundColor: stringToColor(tag) }}
    >
      {tag}
    </span>
  );
}

function TaskRow({ task, queue }: { task: any; queue: Queue }) {
  const successRate = queue.message_stats.ack_details.rate;
  const errorRate =
    queue.message_stats.redeliver_details.rate +
    queue.message_stats.reject_details.rate;
  const pending = queue.ready + queue.unacked;
  return (
    <TableRow key={task.name}>
      <TableCell className="px-6 py-4 text-blue-500 hover:underline font-bold">
        <div className="flex flex-row items-baseline gap-1.5">
          {queue.state === "running" ? (
            <div className="size-2 bg-green-500 border border-green-400 rounded-full"></div>
          ) : null}
          {queue.state === "paused" ? (
            <div className="size-2 bg-yellow-500 border border-yellow-400 rounded-full"></div>
          ) : null}
          <Link
            to={{
              pathname: `/tasks/${task.name}`,
              search: `?vhost=${queue.vhost}`,
            }}
          >
            {task.name}
          </Link>
        </div>
      </TableCell>
      <TableCell className="px-6 py-4">
        <div className="flex flex-wrap gap-1.5">
          {task.tags &&
            task.tags.map((tag: string) => <Tag key={tag} tag={tag} />)}
        </div>
      </TableCell>
      <TableCell className="px-6 py-4">{task.concurrency}</TableCell>
      <TableCell className="px-6 py-4">{task.max_retries}</TableCell>
      <TableCell className="px-6 py-4">{pending}</TableCell>
      <TableCell className="px-6 py-4">{successRate}</TableCell>
      <TableCell className="px-6 py-4">{errorRate}</TableCell>
      <TableCell className="px-6 py-4 text-right">
        <Link
          to={{
            pathname: `/tasks/${task.name}`,
            search: `?vhost=${queue.vhost}`,
          }}
        >
          <Button variant="outline" size="icon">
            <ArrowRightIcon className="size-4" />
          </Button>
        </Link>
      </TableCell>
    </TableRow>
  );
}

export default function Tasks() {
  const { tasks, queues } = useLoaderData<typeof loader>();

  const getQueue = (taskName: string) => {
    return queues.find((queue: Queue) => queue.name === taskName);
  };

  const sortedTasks = useMemo(() => {
    return tasks.sort((a: any, b: any) => a.name.localeCompare(b.name));
  }, [tasks]);

  return (
    <>
      <Table>
        <TableHeader>
          <TableRow className="bg-foreground/5">
            <TableHead className="px-6 py-4">Name</TableHead>
            <TableHead className="px-6 py-4">Tags</TableHead>
            <TableHead className="px-6 py-4">Concurrency</TableHead>
            <TableHead className="px-6 py-4">Max Retries</TableHead>
            <TableHead className="px-6 py-4">Pending</TableHead>
            <TableHead className="px-6 py-4">Success rate</TableHead>
            <TableHead className="px-6 py-4">Error rate</TableHead>
            <TableHead className="px-6 py-4" />
          </TableRow>
        </TableHeader>
        <TableBody>
          {sortedTasks.map((task: any) => (
            <TaskRow
              key={task.name}
              task={task}
              queue={getQueue(task.name) as Queue}
            />
          ))}
        </TableBody>
      </Table>
    </>
  );
}
