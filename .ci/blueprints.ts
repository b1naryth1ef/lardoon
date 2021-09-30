import { Blueprint } from "runtime/blueprint.ts";
import { print, pushStep, spawnChildJob } from "runtime/core.ts";
import * as Docker from "pkg/buildy/docker@1.0/mod.ts";
import { uploadArtifact } from "runtime/artifacts.ts";

export const buildBinary = Blueprint<
  { os?: string; arch?: string; version?: string }
>(
  "Build Binary",
  async ({ ws, os, arch, version }) => {
    pushStep("Build UI");
    const yarnRes = await Docker.run(
      "yarn && yarn build",
      {
        image: "node:16",
        copy: [
          "dist/**",
          "src/**",
          "*.js",
          "*.json",
          "yarn.lock",
        ],
      },
    );

    await yarnRes.copy("dist/");

    pushStep("Build Lardoon Binary");
    const res = await Docker.run(
      "mkdir /tmp/build && mv cmd dist *.go go.mod go.sum /tmp/build && cd /tmp/build && go build -o lardoon cmd/lardoon/main.go && mv lardoon /",
      {
        image: `golang:1.17`,
        copy: [
          "cmd/**",
          "dist/**",
          "server/**",
          "*.go",
          "go.mod",
          "go.sum",
        ],
        env: [`GOOS=${os || "linux"}`, `GOARCH=${arch || "amd64"}`],
      },
    );

    if (version !== undefined) {
      await res.copy("/lardoon");

      pushStep("Upload Lardoon Binary");
      const uploadRes = await uploadArtifact("lardoon", {
        name: `lardoon-${os}-${arch}-${version}`,
        published: true,
        labels: [
          "lardoon",
          `arch:${arch}`,
          `os:${os}`,
          `version:${version}`,
        ],
      });
      print(
        `Uploaded binary to ${
          uploadRes.generatePublicURL(
            ws.org,
            ws.repository,
          )
        }`,
      );
    }
  },
);

const semVerRe =
  /v([0-9]+)\.([0-9]+)\.([0-9]+)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+[0-9A-Za-z-]+)?/;

export const githubPush = Blueprint(
  "GitHub Push",
  async ({ ws }) => {
    let version;

    const versionTags = ws.commit.tags.filter((tag) => semVerRe.test(tag));
    if (versionTags.length == 1) {
      print(
        `Found version tag ${versionTags[0]}, will build release artifacts.`,
      );
      version = versionTags[0];
    } else if (versionTags.length > 1) {
      throw new Error(`Found too many version tags: ${versionTags}`);
    }

    await spawnChildJob("blueprint:build-binary", {
      alias: "Build Linux amd64",
      args: { os: "linux", arch: "amd64", version: version },
    });

    await spawnChildJob("blueprint:build-binary", {
      alias: "Build Windows amd64",
      args: { os: "windows", arch: "amd64", version: version },
    });
  },
);
