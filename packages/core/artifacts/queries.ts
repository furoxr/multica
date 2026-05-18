import { queryOptions } from "@tanstack/react-query";
import { api } from "../api";

export const artifactKeys = {
  all: (wsId: string) => ["workspaces", wsId, "artifacts"] as const,
  detail: (wsId: string, artifactId: string) =>
    [...artifactKeys.all(wsId), artifactId] as const,
};

export function artifactListOptions(wsId: string) {
  return queryOptions({
    queryKey: artifactKeys.all(wsId),
    queryFn: () => api.listArtifacts(),
    enabled: !!wsId,
  });
}

export function artifactDetailOptions(wsId: string, artifactId: string) {
  return queryOptions({
    queryKey: artifactKeys.detail(wsId, artifactId),
    queryFn: () => api.getArtifact(artifactId),
    enabled: !!wsId && !!artifactId,
  });
}
