import { ArrowRightIcon } from "lucide-react";
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

function TaskRow({ task, queue }: { task: any; queue: Queue }) {
  const successRate = queue.message_stats.ack_details.rate;
  const errorRate =
    queue.message_stats.redeliver_details.rate +
    queue.message_stats.reject_details.rate;
  const pending = queue.ready + queue.unacked;
  return (
    <TableRow key={task.name}>
      <TableCell className="px-6 py-4 text-blue-500 hover:underline font-bold">
        <Link
          to={{
            pathname: `/tasks/${task.name}`,
            search: `?vhost=${queue.vhost}`,
          }}
        >
          {task.name}
        </Link>
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

  return (
    <>
      <Table>
        <TableHeader>
          <TableRow className="bg-foreground/5">
            <TableHead className="px-6 py-4">Name</TableHead>
            <TableHead className="px-6 py-4">Concurrency</TableHead>
            <TableHead className="px-6 py-4">Max Retries</TableHead>
            <TableHead className="px-6 py-4">Pending</TableHead>
            <TableHead className="px-6 py-4">Success rate</TableHead>
            <TableHead className="px-6 py-4">Error rate</TableHead>
            <TableHead className="px-6 py-4" />
          </TableRow>
        </TableHeader>
        <TableBody>
          {tasks.map((task: any) => (
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
