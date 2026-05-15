import { useParams } from "react-router-dom";
import { useQuery } from "@tanstack/react-query";
import { ArtifactDetailPage as SharedArtifactDetailPage } from "@multica/views/artifacts";
import { artifactDetailOptions } from "@multica/core/artifacts";
import { useWorkspaceId } from "@multica/core/hooks";
import { useDocumentTitle } from "@/hooks/use-document-title";

export function ArtifactDetailPage() {
  const { id } = useParams<{ id: string }>();
  const wsId = useWorkspaceId();
  const { data: artifact } = useQuery(artifactDetailOptions(wsId, id ?? ""));

  useDocumentTitle(artifact?.title ?? "Artifact");

  if (!id) return null;
  return <SharedArtifactDetailPage artifactId={id} />;
}
