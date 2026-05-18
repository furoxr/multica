"use client";

import { useMemo, useState } from "react";
import { Copy, FileText, Search, Trash2 } from "lucide-react";
import { toast } from "sonner";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { artifactKeys, artifactListOptions } from "@multica/core/artifacts";
import { useWorkspaceId } from "@multica/core/hooks";
import { useWorkspacePaths } from "@multica/core/paths";
import { api } from "@multica/core/api";
import type { ArtifactSummary } from "@multica/core/types";
import { Input } from "@multica/ui/components/ui/input";
import { Skeleton } from "@multica/ui/components/ui/skeleton";
import { Button } from "@multica/ui/components/ui/button";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@multica/ui/components/ui/select";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@multica/ui/components/ui/alert-dialog";
import { PageHeader } from "../layout/page-header";
import { AppLink } from "../navigation";

const CONTENT_TYPE_ALL = "all";

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
  onCopyId,
  onRequestDelete,
}: {
  artifact: ArtifactSummary;
  href: string;
  onCopyId: (artifact: ArtifactSummary) => void;
  onRequestDelete: (artifact: ArtifactSummary) => void;
}) {
  const stopAndRun = (fn: () => void) => (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();
    fn();
  };

  return (
    <AppLink
      href={href}
      className="group grid min-h-16 grid-cols-[minmax(0,1fr)_120px_140px_auto] items-center gap-4 border-b px-5 py-3 transition-colors hover:bg-muted/40 max-md:grid-cols-1 max-md:gap-1"
    >
      <div className="min-w-0">
        <div className="flex items-center gap-2">
          <FileText className="h-4 w-4 shrink-0 text-muted-foreground" />
          <span className="shrink-0 font-mono text-xs text-muted-foreground">
            {artifact.identifier}
          </span>
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
      <div className="flex items-center gap-1 opacity-0 transition-opacity group-hover:opacity-100 focus-within:opacity-100">
        <Button
          type="button"
          variant="ghost"
          size="icon"
          className="h-7 w-7"
          title="Copy identifier"
          aria-label="Copy identifier"
          onClick={stopAndRun(() => onCopyId(artifact))}
        >
          <Copy className="h-3.5 w-3.5" />
        </Button>
        <Button
          type="button"
          variant="ghost"
          size="icon"
          className="h-7 w-7 text-muted-foreground hover:text-destructive"
          title="Delete artifact"
          aria-label="Delete artifact"
          onClick={stopAndRun(() => onRequestDelete(artifact))}
        >
          <Trash2 className="h-3.5 w-3.5" />
        </Button>
      </div>
    </AppLink>
  );
}

export function ArtifactsPage() {
  const wsId = useWorkspaceId();
  const paths = useWorkspacePaths();
  const qc = useQueryClient();
  const [search, setSearch] = useState("");
  const [contentTypeFilter, setContentTypeFilter] = useState<string>(CONTENT_TYPE_ALL);
  const [pendingDelete, setPendingDelete] = useState<ArtifactSummary | null>(null);

  const { data: artifacts = [], isLoading, error } = useQuery(
    artifactListOptions(wsId),
  );

  const contentTypes = useMemo(() => {
    const set = new Set<string>();
    for (const a of artifacts) {
      if (a.content_type) set.add(a.content_type);
    }
    return Array.from(set).sort();
  }, [artifacts]);

  const filtered = useMemo(() => {
    const q = search.trim().toLowerCase();
    return artifacts.filter((artifact) => {
      if (
        contentTypeFilter !== CONTENT_TYPE_ALL &&
        artifact.content_type !== contentTypeFilter
      ) {
        return false;
      }
      if (!q) return true;
      return [artifact.identifier, artifact.title, artifact.summary, artifact.content_type]
        .join(" ")
        .toLowerCase()
        .includes(q);
    });
  }, [artifacts, search, contentTypeFilter]);

  const deleteMutation = useMutation({
    mutationFn: (id: string) => api.deleteArtifact(id),
    onSuccess: (_data, id) => {
      qc.setQueryData<ArtifactSummary[]>(artifactKeys.all(wsId), (prev) =>
        prev ? prev.filter((a) => a.id !== id) : prev,
      );
      qc.removeQueries({ queryKey: artifactKeys.detail(wsId, id) });
      toast.success("Artifact deleted");
    },
    onError: (err) => {
      toast.error(err instanceof Error ? err.message : "Failed to delete artifact");
    },
  });

  const handleCopyId = (artifact: ArtifactSummary) => {
    void navigator.clipboard
      .writeText(artifact.identifier)
      .then(() => toast.success(`Copied ${artifact.identifier}`))
      .catch(() => toast.error("Failed to copy identifier"));
  };

  const handleConfirmDelete = () => {
    if (!pendingDelete) return;
    const target = pendingDelete;
    setPendingDelete(null);
    deleteMutation.mutate(target.id);
  };

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

      <div className="flex h-12 shrink-0 items-center gap-2 border-b px-4">
        <div className="relative">
          <Search className="pointer-events-none absolute left-2.5 top-1/2 h-3.5 w-3.5 -translate-y-1/2 text-muted-foreground" />
          <Input
            value={search}
            onChange={(event) => setSearch(event.target.value)}
            placeholder="Search artifacts"
            className="h-8 w-72 pl-8 text-sm max-sm:w-[calc(100vw-2rem)]"
          />
        </div>
        <Select
          value={contentTypeFilter}
          onValueChange={(v) => setContentTypeFilter(v ?? CONTENT_TYPE_ALL)}
        >
          <SelectTrigger size="sm" className="w-44">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value={CONTENT_TYPE_ALL}>All types</SelectItem>
            {contentTypes.map((ct) => (
              <SelectItem key={ct} value={ct}>
                {ct}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
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
                onCopyId={handleCopyId}
                onRequestDelete={setPendingDelete}
              />
            ))}
          </div>
        )}
      </div>

      <AlertDialog
        open={!!pendingDelete}
        onOpenChange={(v) => { if (!v) setPendingDelete(null); }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete artifact?</AlertDialogTitle>
            <AlertDialogDescription>
              {pendingDelete
                ? `${pendingDelete.identifier} · ${pendingDelete.title} will be permanently deleted.`
                : ""}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>Cancel</AlertDialogCancel>
            <AlertDialogAction variant="destructive" onClick={handleConfirmDelete}>
              Delete
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </div>
  );
}
