"use client";

import * as React from "react";
import {
  Markdown as MarkdownBase,
  type MarkdownProps as MarkdownBaseProps,
  type RenderMode,
} from "@multica/ui/markdown";
import { useConfigStore } from "@multica/core/config";
import { useWorkspacePaths } from "@multica/core/paths";
import { IssueMentionCard } from "../issues/components/issue-mention-card";
import { ArtifactChip } from "../artifacts/components/artifact-chip";
import { AppLink } from "../navigation";

export type { RenderMode };

export type MarkdownProps = MarkdownBaseProps;

/**
 * Default renderMention that delegates to IssueMentionCard for issue mentions
 * and renders a styled span for other mention types.
 */
function defaultRenderMention({
  type,
  id,
  artifactHref,
}: {
  type: string;
  id: string;
  artifactHref: string;
}): React.ReactNode {
  if (type === "issue") {
    return <IssueMentionCard issueId={id} />;
  }
  if (type === "artifact") {
    return (
      <AppLink href={artifactHref} className="inline-flex">
        <ArtifactChip
          artifactId={id}
          className="cursor-pointer hover:bg-accent transition-colors"
        />
      </AppLink>
    );
  }
  return null;
}

/**
 * App-level Markdown wrapper that injects IssueMentionCard via renderMention
 * and cdnDomain from the config store for file card rendering.
 */
export function Markdown(props: MarkdownProps): React.JSX.Element {
  const cdnDomain = useConfigStore((s) => s.cdnDomain);
  const paths = useWorkspacePaths();
  return (
    <MarkdownBase
      renderMention={({ type, id }) =>
        defaultRenderMention({ type, id, artifactHref: paths.artifactDetail(id) })
      }
      cdnDomain={cdnDomain}
      {...props}
    />
  );
}

export const MemoizedMarkdown = React.memo(Markdown);
MemoizedMarkdown.displayName = "MemoizedMarkdown";
