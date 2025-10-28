import { useMemo } from "react";
import { Area, AreaChart, CartesianGrid, XAxis, YAxis } from "recharts";
import { zip } from "~/lib/zip";
import type { MQOverview } from "~/types/mq-overview";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "./ui/card";
import {
  ChartContainer,
  ChartLegend,
  ChartLegendContent,
  ChartTooltip,
  ChartTooltipContent,
} from "./ui/chart";

export function GlobalMessagesCard({ data }: { data: MQOverview }) {
  const chartData = useMemo(() => {
    return zip(
      data.queue_totals.messages_ready_log,
      data.queue_totals.messages_unacknowledged_log
    ).map(([ready, unacknowledged], i) => ({
      timestamp: new Date(
        new Date().getTime() -
          (data.queue_totals.messages_ready_log.length - 1 - i) * 5000
      ),
      ready,
      unacknowledged,
    }));
  }, [data]);

  const chartConfig = {
    ready: {
      label: "Ready",
      color: "#54be7e",
    },
    unacknowledged: {
      label: "Unacknowledged",
      color: "#4589ff",
    },
  };

  return (
    <Card className="rounded-none bg-background">
      <CardHeader>
        <CardTitle>Messages</CardTitle>
        <CardDescription>Total messages in the queue.</CardDescription>
      </CardHeader>
      <CardContent className="pl-0">
        <ChartContainer config={chartConfig}>
          <AreaChart data={chartData}>
            <CartesianGrid />
            <XAxis
              dataKey="timestamp"
              tickFormatter={(value) => {
                const date = new Date(value);
                return date.toLocaleTimeString("en-US", {
                  hour: "2-digit",
                  minute: "2-digit",
                  second: "2-digit",
                });
              }}
              angle={-50}
              tickMargin={5}
              textAnchor="end"
              textRendering="geometricPrecision"
            />
            <YAxis />
            <Area
              type="monotone"
              dataKey="ready"
              stroke={chartConfig.ready.color}
            />
            <Area
              type="monotone"
              dataKey="unacknowledged"
              stroke={chartConfig.unacknowledged.color}
            />
            <ChartTooltip content={<ChartTooltipContent />} />
            <ChartLegend className="mt-8" content={<ChartLegendContent />} />
          </AreaChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}
