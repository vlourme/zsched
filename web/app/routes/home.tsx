import { useMemo } from "react";
import { useLoaderData } from "react-router";
import { DataRateCard } from "~/components/data-rate-card";
import { GlobalMessagesCard } from "~/components/global-messages-card";
import { MessageQueueCard } from "~/components/message-queue-card";
import {
  Card,
  CardContent,
  CardDescription,
  CardTitle,
} from "~/components/ui/card";
import { pool } from "~/lib/db";
import { formatDuration } from "~/lib/formatters";
import { request } from "~/lib/lavinmq";
import type { MQOverview } from "~/types/mq-overview";
import type { Route } from "./+types/home";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Zsched - Task Scheduler" },
    {
      name: "description",
      content: "Zsched is a task scheduler built for Go applications.",
    },
  ];
}

export async function loader() {
  const [overview, executions, errors] = await Promise.all([
    request<MQOverview>("/api/overview"),
    pool.query(`
      SELECT count() as c 
      FROM tasks 
      WHERE started_at > dateadd('h', -24, now())
    `),
    pool.query(`
      SELECT count() as c 
      FROM tasks 
      WHERE last_error != ''
        AND started_at > dateadd('h', -24, now())
    `),
  ]);

  return {
    overview: overview,
    executions: executions.rows[0].c || 0,
    errors: errors.rows[0].c || 0,
  };
}

export const handle = {
  title: () => "Home",
};

export default function Home() {
  const { overview, executions, errors } = useLoaderData<typeof loader>();

  const cards = useMemo(() => {
    return [
      {
        title: "Executions (24h)",
        value: executions,
      },
      {
        title: "Errors (24h)",
        value: errors,
      },
      {
        title: "Queues / Consumers",
        value: `${overview.object_totals.queues} / ${overview.object_totals.consumers}`,
      },
      {
        title: "Uptime",
        value: formatDuration(overview.uptime),
      },
    ];
  }, [executions, errors, overview]);

  return (
    <div className="flex flex-col">
      <div className="grid md:divide-x grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 border-b">
        {cards.map((card) => (
          <Card key={card.title} className="rounded-none bg-background">
            <CardContent>
              <CardTitle className="text-4xl font-bold">{card.value}</CardTitle>
              <CardDescription>{card.title}</CardDescription>
            </CardContent>
          </Card>
        ))}
      </div>

      <div className="grid md:divide-x grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4 border-b">
        <GlobalMessagesCard data={overview} />
        <MessageQueueCard data={overview.message_stats} />
        <DataRateCard data={overview} />
      </div>
    </div>
  );
}
