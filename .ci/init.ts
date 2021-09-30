import { GithubCheckRunPlugin } from "pkg/buildy/github@1/plugins.ts";
import { registerPlugin, Workspace } from "runtime/core.ts";
import { BlueprintSyncPlugin } from "runtime/blueprint.ts";

export async function setup(ws: Workspace) {
  registerPlugin(
    new GithubCheckRunPlugin({
      repositorySlug: "b1naryth1ef/lardoon",
    }),
  );
  registerPlugin(
    new BlueprintSyncPlugin({}),
  );
}
