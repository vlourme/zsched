import { Link, useLoaderData } from "react-router";
import { Button } from "~/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "~/components/ui/card";
import { request } from "~/lib/lavinmq";
import type { Queue } from "~/types/mq-queues";
import type { Route } from "./+types/tasks";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Zsched - Tasks" },
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

function TaskCard({ task, queue }: { task: any; queue: Queue }) {
  return (
    <Card className="gap-2">
      <CardHeader className="flex flex-row justify-between">
        <CardTitle>{task.name}</CardTitle>
        <Link
          to={`/tasks/${task.name}?vhost=${encodeURIComponent(queue.vhost)}`}
        >
          <Button variant="outline">View</Button>
        </Link>
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
            <p className="text-sm text-muted-foreground">Pending</p>
            <p className="text-sm">{queue.ready + queue.unacked}</p>
          </div>
          <div className="flex flex-col w-36 gap-1">
            <p className="text-sm text-muted-foreground">Success rate</p>
            <p className="text-sm">{queue.message_stats.ack_details.rate}</p>
          </div>
          <div className="flex flex-col w-36 gap-1">
            <p className="text-sm text-muted-foreground">Error rate</p>
            <p className="text-sm">
              {queue.message_stats.redeliver_details.rate +
                queue.message_stats.reject_details.rate}
            </p>
          </div>
        </div>
      </CardContent>
    </Card>
  );
}

export default function Tasks() {
  const { tasks, queues } = useLoaderData<typeof loader>();

  const getQueue = (taskName: string) => {
    return queues.find((queue: Queue) => queue.name === taskName);
  };

  return (
    <>
      <div className="flex flex-col gap-1">
        <h1 className="text-3xl font-bold">Tasks</h1>
        <p className="text-muted-foreground">
          Here you can manage your tasks and see the status of your tasks.
        </p>
      </div>

      <div className="flex flex-col gap-4 rounded-xl overflow-x-auto">
        {tasks.map((task: any) => (
          <TaskCard
            key={task.name}
            task={task}
            queue={getQueue(task.name) as Queue}
          />
        ))}
      </div>
    </>
  );
}
