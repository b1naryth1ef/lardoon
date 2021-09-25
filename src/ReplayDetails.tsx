import React from "react";
import useFetch, { CachePolicies } from "use-http";

export type ReplayObject = {
  id: number;
  replay_id: number;
  name: string;
  pilot: string;
  created_offset: number;
  deleted_offset: number;
};

export type Replay = {
  id: number;
  name: string;
  reference_time: string;
  recording_time: string;
  title: string;
  data_source: string;
  data_recorder: string;
};

function ReplayObject({ object }: { object: ReplayObject }) {
  return (
    <div>
      {object.pilot} {object.name}{" "}
      ({object.deleted_offset - object.created_offset}s)
    </div>
  );
}

export default function ReplayDetails({ replayId }: { replayId: number }) {
  const { data } = useFetch<Replay & { objects: Array<ReplayObject> }>(
    `/api/replay/${replayId}`,
    {
      cachePolicy: CachePolicies.NO_CACHE,
    },
    [],
  );

  if (!data) return <></>;
  return (
    <div className="flex flex-col gap-2 p-2">
      {data.objects.map((it) => <ReplayObject object={it} />)}
    </div>
  );
}
