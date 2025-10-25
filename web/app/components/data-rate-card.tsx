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

export function DataRateCard({ data }: { data: MQOverview }) {
  const chartData = useMemo(() => {
    return zip(data.recv_oct_details.log, data.send_oct_details.log).map(
      ([recv_oct, send_oct], i) => ({
        timestamp: new Date(
          new Date().getTime() -
            (data.recv_oct_details.log.length - 1 - i) * 5000
        ),
        recv_oct,
        send_oct,
      })
    );
  }, [data]);

  const chartConfig = {
    recv_oct: {
      label: "Recv Oct",
      color: "#54be7e",
    },
    send_oct: {
      label: "Send Oct",
      color: "#4589ff",
    },
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle>Data Rate</CardTitle>
        <CardDescription>Data rate for the queue.</CardDescription>
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
              dataKey="send_oct"
              stroke={chartConfig.send_oct.color}
            />
            <Area
              type="monotone"
              dataKey="recv_oct"
              stroke={chartConfig.recv_oct.color}
            />
            <ChartTooltip content={<ChartTooltipContent />} />
            <ChartLegend className="mt-8" content={<ChartLegendContent />} />
          </AreaChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}
