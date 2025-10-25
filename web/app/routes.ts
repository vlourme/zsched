import {
  type RouteConfig,
  index,
  layout,
  route,
} from "@react-router/dev/routes";

export default [
  layout("routes/_layout.tsx", [
    index("routes/home.tsx"),
    route("tasks", "routes/tasks.tsx"),
    route("tasks/:name", "routes/task.tsx"),
    route("logs/:task_id", "routes/logs.tsx"),
  ]),
] satisfies RouteConfig;
