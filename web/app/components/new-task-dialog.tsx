import { Editor } from "@monaco-editor/react";
import { useState } from "react";
import { Form } from "react-router";
import { Button } from "./ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "./ui/dialog";
import { Label } from "./ui/label";

export function NewTaskDialog() {
  const [parameters, setParameters] = useState<string | undefined>(
    '{\"name\": \"John\"}'
  );
  const [open, setOpen] = useState(false);

  return (
    <Dialog open={open} onOpenChange={(open) => setOpen(open)}>
      <DialogTrigger asChild>
        <Button size="sm" variant="outline">
          New Task
        </Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>New Task</DialogTitle>
          <DialogDescription>
            Create a new task to be executed.
          </DialogDescription>
        </DialogHeader>

        <Form
          onSubmit={() => {
            setOpen(false);
          }}
          className="flex flex-col gap-3"
          method="post"
        >
          <div className="grid w-full gap-3">
            <Label htmlFor="message-2">Parameters</Label>
            <input type="hidden" name="parameters" value={parameters} />
            <Editor
              className="rounded-md overflow-hidden"
              language="json"
              defaultValue={parameters}
              onChange={(value) => setParameters(value)}
              theme="vs-dark"
              options={{
                minimap: { enabled: false },
                padding: { top: 5, bottom: 5 },
                lineNumbers: "off",
                fontSize: 14,
              }}
              height="20vh"
            />
          </div>
          <Button type="submit">Create</Button>
        </Form>
      </DialogContent>
    </Dialog>
  );
}
