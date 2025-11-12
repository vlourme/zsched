import { ClockIcon } from "lucide-react";
import { Button } from "./ui/button";
import { Card, CardContent } from "./ui/card";
import { Dialog, DialogContent, DialogTrigger } from "./ui/dialog";

export function SchedulesDialog({ schedules }: { schedules: any[] }) {
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button size="sm" variant="outline">
          <ClockIcon className="size-4" />
          View schedules
        </Button>
      </DialogTrigger>
      <DialogContent>
        <h2 className="text-lg font-bold">Schedules</h2>
        <div className="flex flex-col mt-2 gap-2 max-h-[60vh] overflow-y-auto">
          {schedules.length > 0 &&
            schedules.map((schedule: any, idx: number) => (
              <Card key={idx} className="border py-3.5">
                <CardContent className="flex flex-col gap-1">
                  <div className="flex flex-col w-36 gap-1">
                    <p className="text-sm text-muted-foreground">Schedule</p>
                    <p className="text-sm">{schedule.schedule}</p>
                  </div>
                  <div className="flex flex-col w-36 gap-1">
                    <p className="text-sm text-muted-foreground">Parameters</p>
                    <p className="text-sm font-mono">
                      {JSON.stringify(schedule.parameters)}
                    </p>
                  </div>
                </CardContent>
              </Card>
            ))}
          {schedules.length === 0 && (
            <p className="text-sm text-muted-foreground">No schedules found.</p>
          )}
        </div>
      </DialogContent>
    </Dialog>
  );
}
