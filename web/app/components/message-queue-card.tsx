import { useMemo } from "react";
import { Area, AreaChart, CartesianGrid, XAxis, YAxis } from "recharts";
import { zip } from "~/lib/zip";
import type { MessageStats } from "~/types/mq-overview";
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
  type ChartConfig,
} from "./ui/chart";

export function MessageQueueCard({ data }: { data: MessageStats }) {
  const chartData = useMemo(() => {
    return zip(
      data.ack_details.log,
      data.deliver_details.log,
      data.get_details.log,
      data.deliver_get_details.log,
      data.confirm_details.log,
      data.publish_details.log,
      data.redeliver_details.log,
      data.reject_details.log
    ).map(
      (
        [ack, confirm, deliver, get, deliver_get, publish, redeliver, reject],
        i
      ) => ({
        timestamp: new Date(
          new Date().getTime() - (data.ack_details.log.length - 1 - i) * 5000
        ),
        ack,
        deliver,
        confirm,
        get,
        deliver_get,
        publish,
        redeliver,
        reject,
      })
    );
  }, [data]);

  const chartConfig = {
    ack: {
      label: "Ack",
      color: "#54be7e",
    },
    deliver: {
      label: "Deliver",
      color: "#4589ff",
    },
    get: {
      label: "Get",
      color: "#d12771",
    },
    deliver_get: {
      label: "Deliver Get",
      color: "#d2a106",
    },
    publish: {
      label: "Publish",
      color: "#08bdba",
    },
    confirm: {
      label: "Confirm",
      color: "#bae6ff",
    },
    redeliver: {
      label: "Redeliver",
      color: "#ba4e00",
    },
    reject: {
      label: "Reject",
      color: "#d4bbff",
    },
  } satisfies ChartConfig;

  return (
    <Card>
      <CardHeader>
        <CardTitle>Message rates</CardTitle>
        <CardDescription>Message rates for the queue.</CardDescription>
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
              dataKey="ack"
              stroke={chartConfig.ack.color}
            />
            <Area
              type="monotone"
              dataKey="confirm"
              stroke={chartConfig.confirm.color}
            />
            <Area
              type="monotone"
              dataKey="deliver"
              stroke={chartConfig.deliver.color}
            />
            <Area
              type="monotone"
              dataKey="get"
              stroke={chartConfig.get.color}
            />
            <Area
              type="monotone"
              dataKey="deliver_get"
              stroke={chartConfig.deliver_get.color}
            />
            <Area
              type="monotone"
              dataKey="publish"
              stroke={chartConfig.publish.color}
            />
            <Area
              type="monotone"
              dataKey="redeliver"
              stroke={chartConfig.redeliver.color}
            />
            <Area
              type="monotone"
              dataKey="reject"
              stroke={chartConfig.reject.color}
            />
            <ChartTooltip content={<ChartTooltipContent />} />
            <ChartLegend
              className="mt-8 flex-wrap px-4"
              content={<ChartLegendContent />}
            />
          </AreaChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}
