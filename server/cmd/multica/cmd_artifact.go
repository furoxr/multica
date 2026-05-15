package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/multica-ai/multica/server/internal/cli"
)

var artifactCmd = &cobra.Command{
	Use:   "artifact",
	Short: "Work with artifacts",
}

var artifactListCmd = &cobra.Command{
	Use:   "list",
	Short: "List artifacts in the workspace",
	RunE:  runArtifactList,
}

var artifactGetCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get artifact details",
	Args:  exactArgs(1),
	RunE:  runArtifactGet,
}

var artifactCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new artifact",
	RunE:  runArtifactCreate,
}

var artifactUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an artifact",
	Args:  exactArgs(1),
	RunE:  runArtifactUpdate,
}

var artifactDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an artifact",
	Args:  exactArgs(1),
	RunE:  runArtifactDelete,
}

func init() {
	artifactCmd.AddCommand(artifactListCmd)
	artifactCmd.AddCommand(artifactGetCmd)
	artifactCmd.AddCommand(artifactCreateCmd)
	artifactCmd.AddCommand(artifactUpdateCmd)
	artifactCmd.AddCommand(artifactDeleteCmd)

	artifactListCmd.Flags().String("origin-issue", "", "Filter by origin issue ID")
	artifactListCmd.Flags().String("output", "table", "Output format: table or json")

	artifactGetCmd.Flags().String("output", "json", "Output format: table or json")

	artifactCreateCmd.Flags().String("title", "", "Artifact title (required)")
	artifactCreateCmd.Flags().String("summary", "", "Short artifact summary")
	artifactCreateCmd.Flags().String("content", "", "Artifact content")
	artifactCreateCmd.Flags().String("content-file", "", "Read artifact content from a UTF-8 file")
	artifactCreateCmd.Flags().String("content-type", "text/markdown", "Content type: text/markdown, text/plain, or application/json")
	artifactCreateCmd.Flags().String("project", "", "Project ID")
	artifactCreateCmd.Flags().String("origin-issue", "", "Origin issue ID")
	artifactCreateCmd.Flags().String("origin-task", "", "Origin task ID")
	artifactCreateCmd.Flags().String("output", "json", "Output format: table or json")

	artifactUpdateCmd.Flags().String("title", "", "New title")
	artifactUpdateCmd.Flags().String("summary", "", "New summary")
	artifactUpdateCmd.Flags().String("content", "", "New content")
	artifactUpdateCmd.Flags().String("content-file", "", "Read new content from a UTF-8 file")
	artifactUpdateCmd.Flags().String("content-type", "", "New content type")
	artifactUpdateCmd.Flags().String("project", "", "Project ID")
	artifactUpdateCmd.Flags().String("origin-issue", "", "Origin issue ID")
	artifactUpdateCmd.Flags().String("origin-task", "", "Origin task ID")
	artifactUpdateCmd.Flags().String("output", "json", "Output format: table or json")

	artifactDeleteCmd.Flags().Bool("yes", false, "Skip confirmation prompt")
}

func readArtifactContent(cmd *cobra.Command) (string, bool, error) {
	if cmd.Flags().Changed("content-file") {
		path, _ := cmd.Flags().GetString("content-file")
		b, err := os.ReadFile(path)
		if err != nil {
			return "", false, fmt.Errorf("read --content-file: %w", err)
		}
		return string(b), true, nil
	}
	if cmd.Flags().Changed("content") {
		v, _ := cmd.Flags().GetString("content")
		return v, true, nil
	}
	return "", false, nil
}

func runArtifactList(cmd *cobra.Command, _ []string) error {
	client, err := newAPIClient(cmd)
	if err != nil {
		return err
	}
	path := "/api/artifacts"
	if originIssue, _ := cmd.Flags().GetString("origin-issue"); originIssue != "" {
		path += "?origin_issue_id=" + originIssue
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var artifacts []map[string]any
	if err := client.GetJSON(ctx, path, &artifacts); err != nil {
		return fmt.Errorf("list artifacts: %w", err)
	}

	output, _ := cmd.Flags().GetString("output")
	if output == "json" {
		return cli.PrintJSON(os.Stdout, artifacts)
	}

	headers := []string{"ID", "TITLE", "TYPE", "ORIGIN_ISSUE", "UPDATED_AT"}
	rows := make([][]string, 0, len(artifacts))
	for _, a := range artifacts {
		rows = append(rows, []string{
			strVal(a, "id"),
			strVal(a, "title"),
			strVal(a, "content_type"),
			strVal(a, "origin_issue_id"),
			strVal(a, "updated_at"),
		})
	}
	cli.PrintTable(os.Stdout, headers, rows)
	return nil
}

func runArtifactGet(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient(cmd)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var artifact map[string]any
	if err := client.GetJSON(ctx, "/api/artifacts/"+args[0], &artifact); err != nil {
		return fmt.Errorf("get artifact: %w", err)
	}

	output, _ := cmd.Flags().GetString("output")
	if output == "json" {
		return cli.PrintJSON(os.Stdout, artifact)
	}

	headers := []string{"ID", "TITLE", "TYPE", "ORIGIN_ISSUE", "UPDATED_AT"}
	rows := [][]string{{
		strVal(artifact, "id"),
		strVal(artifact, "title"),
		strVal(artifact, "content_type"),
		strVal(artifact, "origin_issue_id"),
		strVal(artifact, "updated_at"),
	}}
	cli.PrintTable(os.Stdout, headers, rows)
	return nil
}

func runArtifactCreate(cmd *cobra.Command, _ []string) error {
	client, err := newAPIClient(cmd)
	if err != nil {
		return err
	}
	title, _ := cmd.Flags().GetString("title")
	if title == "" {
		return fmt.Errorf("--title is required")
	}
	body := map[string]any{"title": title}
	if v, _ := cmd.Flags().GetString("summary"); v != "" {
		body["summary"] = v
	}
	if v, ok, err := readArtifactContent(cmd); err != nil {
		return err
	} else if ok {
		body["content"] = v
	}
	if v, _ := cmd.Flags().GetString("content-type"); v != "" {
		body["content_type"] = v
	}
	if v, _ := cmd.Flags().GetString("project"); v != "" {
		body["project_id"] = v
	}
	if v, _ := cmd.Flags().GetString("origin-issue"); v != "" {
		body["origin_issue_id"] = v
	}
	if v, _ := cmd.Flags().GetString("origin-task"); v != "" {
		body["origin_task_id"] = v
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var result map[string]any
	if err := client.PostJSON(ctx, "/api/artifacts", body, &result); err != nil {
		return fmt.Errorf("create artifact: %w", err)
	}
	output, _ := cmd.Flags().GetString("output")
	if output == "json" {
		return cli.PrintJSON(os.Stdout, result)
	}
	fmt.Printf("Artifact created: %s (%s)\n", strVal(result, "title"), strVal(result, "id"))
	return nil
}

func runArtifactUpdate(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient(cmd)
	if err != nil {
		return err
	}
	body := map[string]any{}
	for _, name := range []string{"title", "summary"} {
		if cmd.Flags().Changed(name) {
			v, _ := cmd.Flags().GetString(name)
			body[name] = v
		}
	}
	if v, ok, err := readArtifactContent(cmd); err != nil {
		return err
	} else if ok {
		body["content"] = v
	}
	if cmd.Flags().Changed("content-type") {
		v, _ := cmd.Flags().GetString("content-type")
		body["content_type"] = v
	}
	if cmd.Flags().Changed("project") {
		v, _ := cmd.Flags().GetString("project")
		body["project_id"] = v
	}
	if cmd.Flags().Changed("origin-issue") {
		v, _ := cmd.Flags().GetString("origin-issue")
		body["origin_issue_id"] = v
	}
	if cmd.Flags().Changed("origin-task") {
		v, _ := cmd.Flags().GetString("origin-task")
		body["origin_task_id"] = v
	}
	if len(body) == 0 {
		return fmt.Errorf("no fields to update")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var result map[string]any
	if err := client.PutJSON(ctx, "/api/artifacts/"+args[0], body, &result); err != nil {
		return fmt.Errorf("update artifact: %w", err)
	}
	output, _ := cmd.Flags().GetString("output")
	if output == "json" {
		return cli.PrintJSON(os.Stdout, result)
	}
	fmt.Printf("Artifact updated: %s (%s)\n", strVal(result, "title"), strVal(result, "id"))
	return nil
}

func runArtifactDelete(cmd *cobra.Command, args []string) error {
	client, err := newAPIClient(cmd)
	if err != nil {
		return err
	}
	yes, _ := cmd.Flags().GetBool("yes")
	if !yes {
		fmt.Printf("Are you sure you want to delete artifact %s? This cannot be undone. [y/N] ", args[0])
		var answer string
		fmt.Scanln(&answer)
		if answer != "y" && answer != "Y" {
			fmt.Println("Cancelled")
			return nil
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := client.DeleteJSON(ctx, "/api/artifacts/"+args[0]); err != nil {
		return fmt.Errorf("delete artifact: %w", err)
	}
	fmt.Println("Artifact deleted")
	return nil
}
