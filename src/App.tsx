import React from "react";
import { Link, Route, Switch } from "react-router-dom";
import useFetch, { CachePolicies } from "use-http";
import ReplayDetails, { Replay } from "./ReplayDetails";

function ReplayItem({ replay }: { replay: Replay }) {
  // TODO: Render local time + information
  return (
    <Link
      to={`/replay/${replay.id}`}
      className="p-2 border border-blue-400 rounded-sm hover:bg-blue-200 hover:border-blue-500 w-full"
    >
      {replay.recording_time}
    </Link>
  );
}

function ReplayList() {
  const { data } = useFetch<Array<Replay>>("/api/replay", {
    cachePolicy: CachePolicies.NO_CACHE,
  }, []);

  if (!data) return <></>;
  return (
    <div
      className="border border-gray-300 bg-gray-200 rounded-sm p-2 shadow-sm flex flex-col items-center gap-2"
    >
      {data.map((it) => <ReplayItem key={it.id} replay={it} />)}
    </div>
  );
}

function App() {
  return (
    <div
      className="flex h-screen w-screen"
    >
      <div className="m-auto">
        <Switch>
          <Route exact path="/" component={ReplayList} />
          <Route
            exact
            path="/replay/:replayId"
            render={({ match: { params: { replayId } } }) => {
              return <ReplayDetails replayId={parseInt(replayId)} />;
            }}
          />
        </Switch>
      </div>
    </div>
  );
}

export default App;
