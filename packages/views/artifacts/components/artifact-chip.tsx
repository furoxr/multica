"use client";

import { useQuery } from "@tanstack/react-query";
import { FileText } from "lucide-react";
import { artifactDetailOptions, artifactListOptions } from "@multica/core/artifacts";
import { useWorkspaceId } from "@multica/core/hooks";

export interface ArtifactChipProps {
  artifactId: string;
  fallbackLabel?: string;
  className?: string;
}

const BASE_CLASS =
  "artifact-mention inline-flex items-center gap-1.5 rounded-md border mx-0.5 px-2 py-0.5 text-xs max-w-72";

export function ArtifactChip({ artifactId, fallbackLabel, className }: ArtifactChipProps) {
  const wsId = useWorkspaceId();
  const { data: artifacts = [] } = useQuery(artifactListOptions(wsId));
  const listArtifact = artifacts.find((a) => a.id === artifactId);

  const { data: detailArtifact } = useQuery({
    ...artifactDetailOptions(wsId, artifactId),
    enabled: !listArtifact,
  });

  const artifact = listArtifact ?? detailArtifact;
  const cls = className ? `${BASE_CLASS} ${className}` : BASE_CLASS;

  if (!artifact) {
    return (
      <span className={cls}>
        <FileText className="h-3.5 w-3.5 shrink-0 text-muted-foreground" />
        <span className="font-medium text-muted-foreground">
          {fallbackLabel ?? artifactId.slice(0, 8)}
        </span>
      </span>
    );
  }

  return (
    <span className={cls}>
      <FileText className="h-3.5 w-3.5 shrink-0 text-muted-foreground" />
      <span className="font-medium text-muted-foreground shrink-0">
        {artifact.identifier}
      </span>
      <span className="text-foreground truncate">{artifact.title}</span>
    </span>
  );
}
