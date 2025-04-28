// import type { Route } from "./+types/dashboard";
import Dashboard from "../dashboard/dashboard";

export function meta() {
  return [{ title: "Dashboard" }, { name: "description", content: "Welcome to React Router!" }];
}

export default function DashboardPage() {
  return <Dashboard />;
}
