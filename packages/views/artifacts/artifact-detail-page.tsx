"use client";

import { ArrowLeft, Copy, FileText } from "lucide-react";
import { toast } from "sonner";
import { useQuery } from "@tanstack/react-query";
import { artifactDetailOptions } from "@multica/core/artifacts";
import { useWorkspaceId } from "@multica/core/hooks";
import { useWorkspacePaths } from "@multica/core/paths";
import { Button } from "@multica/ui/components/ui/button";
import { Skeleton } from "@multica/ui/components/ui/skeleton";
import { PageHeader } from "../layout/page-header";
import { AppLink } from "../navigation";
import { Markdown } from "../common/markdown";

function formatDate(value: string) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return new Intl.DateTimeFormat(undefined, {
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(date);
}

function MetaRow({ label, value }: { label: string; value?: string | null }) {
  if (!value) return null;
  return (
    <div className="grid grid-cols-[120px_minmax(0,1fr)] gap-3 text-sm">
      <div className="text-muted-foreground">{label}</div>
      <div className="min-w-0 break-all font-mono text-xs">{value}</div>
    </div>
  );
}

export function ArtifactDetailPage({ artifactId }: { artifactId: string }) {
  const wsId = useWorkspaceId();
  const paths = useWorkspacePaths();
  const { data: artifact, isLoading, error } = useQuery(
    artifactDetailOptions(wsId, artifactId),
  );

  return (
    <div className="flex h-full min-h-0 flex-col bg-background">
      <PageHeader className="justify-between px-5">
        <div className="flex min-w-0 items-center gap-2">
          <Button
            type="button"
            size="icon"
            variant="ghost"
            className="h-8 w-8"
            render={<AppLink href={paths.artifacts()} />}
          >
            <ArrowLeft className="h-4 w-4" />
          </Button>
          <FileText className="h-4 w-4 text-muted-foreground" />
          {artifact?.identifier && (
            <button
              type="button"
              onClick={() => {
                void navigator.clipboard.writeText(artifact.identifier);
                toast.success(`Copied ${artifact.identifier}`);
              }}
              className="flex shrink-0 items-center gap-1 rounded px-1 py-0.5 font-mono text-xs text-muted-foreground hover:bg-accent hover:text-foreground"
              title="Copy identifier"
            >
              <span>{artifact.identifier}</span>
              <Copy className="h-3 w-3" />
            </button>
          )}
          <h1 className="truncate text-sm font-medium">
            {artifact?.title ?? "Artifact"}
          </h1>
        </div>
      </PageHeader>

      <div className="min-h-0 flex-1 overflow-auto">
        {isLoading ? (
          <div className="mx-auto max-w-5xl px-6 py-6">
            <Skeleton className="h-7 w-72" />
            <Skeleton className="mt-3 h-4 w-96 max-w-full" />
            <Skeleton className="mt-8 h-64 w-full" />
          </div>
        ) : error || !artifact ? (
          <div className="px-5 py-8 text-sm text-destructive">
            Artifact not found.
          </div>
        ) : (
          <main className="mx-auto grid max-w-6xl grid-cols-[minmax(0,1fr)_280px] gap-8 px-6 py-6 max-lg:grid-cols-1">
            <article className="min-w-0">
              <div className="font-mono text-xs text-muted-foreground">
                {artifact.identifier}
              </div>
              <h2 className="mt-1 text-2xl font-semibold tracking-normal">
                {artifact.title}
              </h2>
              {artifact.summary && (
                <p className="mt-2 text-sm text-muted-foreground">
                  {artifact.summary}
                </p>
              )}
              <div className="mt-6 rounded-md border bg-card p-5">
                {artifact.content_type === "text/markdown" ? (
                  <Markdown>{artifact.content || "_Empty artifact_"}</Markdown>
                ) : artifact.content_type === "application/json" ? (
                  <pre className="overflow-auto whitespace-pre-wrap text-sm">
                    {artifact.content}
                  </pre>
                ) : (
                  <pre className="overflow-auto whitespace-pre-wrap text-sm">
                    {artifact.content}
                  </pre>
                )}
              </div>
            </article>

            <aside className="space-y-3 border-l pl-6 max-lg:border-l-0 max-lg:border-t max-lg:pl-0 max-lg:pt-6">
              <MetaRow label="Type" value={artifact.content_type} />
              <MetaRow label="Creator" value={`${artifact.creator_type}:${artifact.creator_id}`} />
              <MetaRow label="Created" value={formatDate(artifact.created_at)} />
              <MetaRow label="Updated" value={formatDate(artifact.updated_at)} />
              <MetaRow label="Project" value={artifact.project_id} />
              <MetaRow label="Origin issue" value={artifact.origin_issue_id} />
              <MetaRow label="Origin task" value={artifact.origin_task_id} />
            </aside>
          </main>
        )}
      </div>
    </div>
  );
}
