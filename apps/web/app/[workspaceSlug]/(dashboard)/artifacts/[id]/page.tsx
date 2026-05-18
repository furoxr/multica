"use client";

import { use } from "react";
import { ArtifactDetailPage } from "@multica/views/artifacts";

export default function ArtifactDetailRoute({
  params,
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = use(params);
  return <ArtifactDetailPage artifactId={id} />;
}
