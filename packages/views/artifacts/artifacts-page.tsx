"use client";

import { useMemo, useState } from "react";
import { FileText, Search } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import { artifactListOptions } from "@multica/core/artifacts";
import { useWorkspaceId } from "@multica/core/hooks";
import { useWorkspacePaths } from "@multica/core/paths";
import type { ArtifactSummary } from "@multica/core/types";
import { Input } from "@multica/ui/components/ui/input";
import { Skeleton } from "@multica/ui/components/ui/skeleton";
import { PageHeader } from "../layout/page-header";
import { AppLink } from "../navigation";

function formatDate(value: string) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return new Intl.DateTimeFormat(undefined, {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  }).format(date);
}

function ArtifactRow({
  artifact,
  href,
}: {
  artifact: ArtifactSummary;
  href: string;
}) {
  return (
    <AppLink
      href={href}
      className="grid min-h-16 grid-cols-[minmax(0,1fr)_120px_140px] items-center gap-4 border-b px-5 py-3 transition-colors hover:bg-muted/40 max-md:grid-cols-1 max-md:gap-1"
    >
      <div className="min-w-0">
        <div className="flex items-center gap-2">
          <FileText className="h-4 w-4 shrink-0 text-muted-foreground" />
          <div className="truncate text-sm font-medium">{artifact.title}</div>
        </div>
        {artifact.summary && (
          <div className="mt-1 truncate text-xs text-muted-foreground">
            {artifact.summary}
          </div>
        )}
      </div>
      <div className="font-mono text-xs text-muted-foreground">
        {artifact.content_type}
      </div>
      <div className="text-xs text-muted-foreground md:text-right">
        {formatDate(artifact.updated_at)}
      </div>
    </AppLink>
  );
}

export function ArtifactsPage() {
  const wsId = useWorkspaceId();
  const paths = useWorkspacePaths();
  const [search, setSearch] = useState("");
  const { data: artifacts = [], isLoading, error } = useQuery(
    artifactListOptions(wsId),
  );

  const filtered = useMemo(() => {
    const q = search.trim().toLowerCase();
    if (!q) return artifacts;
    return artifacts.filter((artifact) =>
      [artifact.title, artifact.summary, artifact.content_type]
        .join(" ")
        .toLowerCase()
        .includes(q),
    );
  }, [artifacts, search]);

  return (
    <div className="flex h-full min-h-0 flex-col bg-background">
      <PageHeader className="justify-between px-5">
        <div className="flex items-center gap-2">
          <FileText className="h-4 w-4 text-muted-foreground" />
          <h1 className="text-sm font-medium">Artifacts</h1>
          {artifacts.length > 0 && (
            <span className="font-mono text-xs tabular-nums text-muted-foreground/70">
              {artifacts.length}
            </span>
          )}
        </div>
      </PageHeader>

      <div className="flex h-12 shrink-0 items-center border-b px-4">
        <div className="relative">
          <Search className="pointer-events-none absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted-foreground" />
          <Input
            value={search}
            onChange={(event) => setSearch(event.target.value)}
            placeholder="Search artifacts"
            className="h-8 w-72 pl-8 text-sm max-sm:w-[calc(100vw-2rem)]"
          />
        </div>
      </div>

      <div className="min-h-0 flex-1 overflow-auto">
        {isLoading ? (
          <div className="space-y-0">
            {Array.from({ length: 6 }).map((_, index) => (
              <div key={index} className="border-b px-5 py-4">
                <Skeleton className="h-4 w-64" />
                <Skeleton className="mt-2 h-3 w-96 max-w-full" />
              </div>
            ))}
          </div>
        ) : error ? (
          <div className="px-5 py-8 text-sm text-destructive">
            Failed to load artifacts.
          </div>
        ) : filtered.length === 0 ? (
          <div className="flex h-full flex-col items-center justify-center px-6 text-center">
            <div className="flex h-12 w-12 items-center justify-center rounded-full bg-muted">
              <FileText className="h-6 w-6 text-muted-foreground" />
            </div>
            <h2 className="mt-4 text-base font-semibold">No artifacts</h2>
            <p className="mt-1 max-w-md text-sm text-muted-foreground">
              Reusable content created by members or agents will appear here.
            </p>
          </div>
        ) : (
          <div>
            {filtered.map((artifact) => (
              <ArtifactRow
                key={artifact.id}
                artifact={artifact}
                href={paths.artifactDetail(artifact.id)}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
